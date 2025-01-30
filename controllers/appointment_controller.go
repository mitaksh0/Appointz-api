package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/appointments_api/models"
	"github.com/appointments_api/utils"
)

// crud operations for appointments table
func AppointmentHandler(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	appointmentId := queryParams.Get("id")
	userId := queryParams.Get("user-id")

	// get user role
	userIdInt, _ := strconv.Atoi(appointmentId)
	if userId == "" || userIdInt == 0 {
		utils.GenerateResponse(w, http.StatusUnauthorized, "user id not found in request")
		return
	}
	userRoles, err := models.GetRoles(userIdInt)
	if err != nil {
		utils.GenerateResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	switch r.Method {
	case http.MethodGet:

		if appointmentId != "" { // get single appointment
			aId, err := strconv.Atoi(appointmentId)
			if err != nil {
				utils.GenerateResponse(w, http.StatusBadRequest, "invalid appointment id")
				return
			}
			res, err := models.GetAppointment(aId)

			if err != nil {
				utils.GenerateResponse(w, http.StatusInternalServerError, err.Error())
				return
			}

			utils.GenerateResponse(w, http.StatusOK, res)
			return

		} else { // get all appointments
			// user date
			var userDate string
			res, err := models.GetAppointments(userDate)
			if err != nil {
				utils.GenerateResponse(w, http.StatusInternalServerError, err.Error())
				return
			} else if len(res) == 0 {
				utils.GenerateResponse(w, http.StatusNotFound, "no appointments found")
				return
			}

			utils.GenerateResponse(w, http.StatusOK, res)
			return

		}

	case http.MethodPost:

		for key := range userRoles {
			if key != "receptionist" {
				utils.GenerateResponse(w, http.StatusUnauthorized, "user not authorized for this request")
				return
			}
		}

		var a models.Appointment

		err := json.NewDecoder(r.Body).Decode(&a)
		if err != nil {
			utils.GenerateResponse(w, http.StatusInternalServerError, "error parsing data")
			return
		}

		if a.AppointmentDate == "" || a.Notes == "" || a.RecepID == 0 || a.PatientID == 0 {
			utils.GenerateResponse(w, http.StatusBadRequest, "missing data in input field(s)")
			return
		}

		err = models.AddAppointment(a)
		if err != nil {
			utils.GenerateResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		utils.GenerateResponse(w, http.StatusOK, "appointment added successfully")
		return

	case http.MethodPut:
		// appointment update function
		var a models.Appointment
		aId, err := strconv.Atoi(appointmentId)
		if err != nil {
			utils.GenerateResponse(w, http.StatusBadRequest, "invalid appointment id")
			return
		}

		err = json.NewDecoder(r.Body).Decode(&a)
		if err != nil {
			utils.GenerateResponse(w, http.StatusInternalServerError, "error parsing data")
			return
		} else if a.AppointmentDate == "" || a.Notes == "" || a.RecepID == 0 || a.PatientID == 0 {
			utils.GenerateResponse(w, http.StatusBadRequest, "missing data in input field(s)")
			return
		}

		err = models.UpdateAppointment(a, aId)
		if err != nil {
			utils.GenerateResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		utils.GenerateResponse(w, http.StatusOK, "updated successfully")
		return

	case http.MethodDelete:

		for key := range userRoles {
			if key != "receptionist" {
				utils.GenerateResponse(w, http.StatusUnauthorized, "user not authorized for this request")
				return
			}
		}

		if appointmentId == "" {
			utils.GenerateResponse(w, http.StatusBadRequest, "missing appointment id")
			return
		}

		aId, err := strconv.Atoi(appointmentId)
		if err != nil {
			utils.GenerateResponse(w, http.StatusBadRequest, "invalid appointment id")
			return
		}

		err = models.DeleteAppointment(aId)
		if err != nil {
			utils.GenerateResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		utils.GenerateResponse(w, http.StatusOK, "delete appointment successfully")
		return

	default:
		utils.GenerateResponse(w, http.StatusMethodNotAllowed, "method not allowed")
	}

}
