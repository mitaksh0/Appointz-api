package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/appointments_api/models"
	"github.com/appointments_api/utils"
	"github.com/golang-jwt/jwt/v4"
)

func PatientsHandler(w http.ResponseWriter, r *http.Request) {

	queryParams := r.URL.Query()
	patientId := queryParams.Get("id")

	ctx := r.Context().Value(models.PayloadContextKey)

	payload := ctx.(jwt.MapClaims)
	userId := payload["userId"].(float64)
	role := payload["role"].(string)

	// get user role
	userIdInt := int(userId)
	if userId == 0 || userIdInt == 0 || role == "" {
		utils.GenerateResponse(w, http.StatusUnauthorized, "error in session, please login again")
		return
	}

	// userRoles, err := models.GetRoles(userIdInt)
	// if err != nil {
	// 	utils.GenerateResponse(w, http.StatusInternalServerError, err.Error())
	// 	return
	// }

	switch r.Method {
	case http.MethodGet:
		// get all patient(s)
		// if patient id given, return single object, otherwise return all objects
		if patientId != "" {
			// get single patient
			idInt, _ := strconv.Atoi(patientId)
			res, err := models.GetPatient(idInt)
			if err != nil {
				if err == sql.ErrNoRows {
					utils.GenerateResponse(w, http.StatusNotFound, "no result")
					return
				}
				utils.GenerateResponse(w, http.StatusInternalServerError, err.Error())
				return
			}

			utils.GenerateResponse(w, http.StatusOK, res)
			return
		} else {
			// get all patients info
			res, err := models.GetPatients()
			if err != nil {
				utils.GenerateResponse(w, http.StatusInternalServerError, err.Error())
				return
			}

			// no patients found
			if len(res) == 0 {
				utils.GenerateResponse(w, http.StatusNotFound, "no result")
				return
			}

			utils.GenerateResponse(w, http.StatusOK, res)
			return
		}
	case http.MethodPost:

		if role != "receptionist" {
			utils.GenerateResponse(w, http.StatusUnauthorized, "user not authorised for this request")
			return
		}

		// add patent
		var patient models.Patient
		err := json.NewDecoder(r.Body).Decode(&patient)
		if err != nil {
			fmt.Println(err.Error())
			utils.GenerateResponse(w, http.StatusInternalServerError, "error parsing data")
			return
		}

		if patient.Address == "" || patient.Contact == "" || patient.DOB == "" || patient.FName == "" || patient.Gender == "" || patient.LName == "" {
			utils.GenerateResponse(w, http.StatusBadRequest, "missing data in input field(s)")
			return
		}

		// add patient
		err = models.InsertPatient(patient)
		if err != nil {
			fmt.Println(err.Error())
			utils.GenerateResponse(w, http.StatusInternalServerError, "error inserting")
			return
		}

		utils.GenerateResponse(w, http.StatusOK, "insert success")

	case http.MethodPut:
		// update patient info
		// since we can only update single patient per call
		if patientId == "" {
			utils.GenerateResponse(w, http.StatusBadRequest, "missing patient id")
			return
		}

		var patient models.Patient
		err := json.NewDecoder(r.Body).Decode(&patient)
		if err != nil {
			utils.GenerateResponse(w, http.StatusInternalServerError, "error parsing data")
			return
		}

		// if no value is passed
		if patient == (models.Patient{}) {
			utils.GenerateResponse(w, http.StatusBadRequest, "missing data in input field(s)")
			return
		}
		idInt, _ := strconv.Atoi(patientId)
		err = models.UpdatePatient(patient, idInt)
		if err != nil {
			fmt.Println(err.Error())
			utils.GenerateResponse(w, http.StatusInternalServerError, "error updating")
			return
		}

		utils.GenerateResponse(w, http.StatusOK, "updated successfully")
		return

	case http.MethodDelete:

		// if role is not receptionist, return
		if role != "receptionist" {
			utils.GenerateResponse(w, http.StatusUnauthorized, "user not authorised for this request")
			return
		}

		// delete patient
		// only single patient delete allowed at once, update here in case of batch delete
		if patientId == "" {
			utils.GenerateResponse(w, http.StatusBadRequest, "missing patient id")
			return
		}
		idInt, _ := strconv.Atoi(patientId)
		err := models.DeletePatient(idInt)
		if err != nil {
			utils.GenerateResponse(w, http.StatusInternalServerError, "error updating")
			return
		}

		utils.GenerateResponse(w, http.StatusOK, "deleted successfully")
		return

	default:
		http.Error(w, "invalid request method", http.StatusMethodNotAllowed)
	}
}

// func PatientHandler(w http.ResponseWriter, r *http.Request) {

// 	// get patient id
// 	idStr := strings.TrimPrefix(r.URL.Path, "/patients/")
// 	_, err := strconv.Atoi(idStr)
// 	if err != nil {
// 		http.Error(w, "invalid patient ID", http.StatusBadRequest)
// 		return
// 	}

// 	switch r.Method {
// 	case http.MethodGet:
// 		utils.GenerateResponse(w, http.StatusOK, "got single patient "+idStr)
// 		// get patient
// 	case http.MethodPut:
// 		// update patient
// 	case http.MethodDelete:
// 		// delete patient
// 	default:
// 		http.Error(w, "invalid request method", http.StatusMethodNotAllowed)
// 	}
// }
