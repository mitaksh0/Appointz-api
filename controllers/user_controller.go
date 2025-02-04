package controllers

import (
	"net/http"

	"github.com/appointments_api/models"
	"github.com/appointments_api/utils"
	"github.com/golang-jwt/jwt/v4"
)

func UsersHandler(w http.ResponseWriter, r *http.Request) {

	queryParams := r.URL.Query()
	userRole := queryParams.Get("role")

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
		// get all user(s)

		if userRole == "" {
			utils.GenerateResponse(w, http.StatusBadRequest, "no role passed")
			return
		}

		// get all patients info
		res, err := models.GetUsers(userRole)
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

	default:
		http.Error(w, "invalid request method", http.StatusMethodNotAllowed)
	}
}
