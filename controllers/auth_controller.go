package controllers

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"

	"github.com/appointments_api/models"
	"github.com/appointments_api/utils"
)

func Login(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		utils.GenerateResponse(w, http.StatusMethodNotAllowed, "invalid request method")
		return
	}

	var user models.Credentials
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		utils.GenerateResponse(w, http.StatusBadRequest, "bad request")
		return
	}

	if user.Username == "" || user.Password == "" || user.Role == "" {
		utils.GenerateResponse(w, http.StatusUnauthorized, "missing data in input field(s)")
		return
	}

	dbUser, err := models.GetAuthUser(user.Username, user.Role)
	if err != nil {
		utils.GenerateResponse(w, http.StatusInternalServerError, "server error")
		return
	}

	isMatch := utils.ComparePassword(dbUser.Password, user.Password)
	if !isMatch {
		utils.GenerateResponse(w, http.StatusUnauthorized, "incorrect password/username/role")
		return
	}

	expirationTime := time.Now().Add(15 * time.Minute)

	// Convert expirationTime to NumericDate
	exp := jwt.NewNumericDate(expirationTime)
	// claims := &jwt.RegisteredClaims{
	// 	ExpiresAt: exp,
	// }

	// Combine standard claims with your custom payload
	claims := jwt.MapClaims{
		"exp":    exp,
		"iat":    time.Now().Unix(), // Issued at time (optional)
		"userId": dbUser.ID,
		"role":   user.Role,
	}

	secretKey := os.Getenv("SECRET_KEY")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		utils.GenerateResponse(w, http.StatusInternalServerError, "internal server error")
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    tokenString,
		Expires:  expirationTime,
		Path:     "/",
		Secure:   true, // Set to true if using HTTPS
		HttpOnly: true, // Cookie is only accessible by the server
		SameSite: http.SameSiteNoneMode,
	})

	utils.GenerateResponse(w, http.StatusOK, map[string]string{"token": tokenString})
	// json.NewEncoder(w).Encode(map[string]string{"token": tokenString})

}

func Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.GenerateResponse(w, http.StatusMethodNotAllowed, "invalid request method")
		return
	}
	// register user
	var user models.Credentials
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		utils.GenerateResponse(w, http.StatusInternalServerError, "error parsing data")
		return
	}

	if user.Email == "" || user.Role == "" || user.Username == "" || user.Password == "" || user.Name == "" {
		utils.GenerateResponse(w, http.StatusUnauthorized, "missing data in input field(s)")
		return
	} else if user.Password != user.RePassword {
		utils.GenerateResponse(w, http.StatusUnauthorized, "password does not match re-password")
		return
	}

	// if username is not unique per role, return error
	users, err := models.GetUser(user.Username, user.Email, user.Role)
	if err != nil {
		utils.GenerateResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	// if user already exists with given info, return along with error
	if len(users) != 0 {
		if len(users) >= 2 {
			utils.GenerateResponse(w, http.StatusInternalServerError, "user already exists")
		} else if users[0].Email == user.Email {
			utils.GenerateResponse(w, http.StatusInternalServerError, "email already in use")
		} else {
			utils.GenerateResponse(w, http.StatusInternalServerError, "username already in use")
		}
		return
	}

	// insert user

	// hash password and remove unhashed password
	user.RePassword = ""
	user.Password, err = utils.HashPassword(user.Password)
	if err != nil {
		utils.GenerateResponse(w, http.StatusInternalServerError, "error hashing password")
		return
	}

	err = models.SetUser(user)
	if err != nil {
		utils.GeneratePreflightRequest(w, http.StatusInternalServerError, map[string]string{
			"msg":   "error inserting",
			"error": err.Error(),
		})
		return
	}

	utils.GenerateResponse(w, http.StatusCreated, "new user created")
}

// Delete jwt token by setting expiry to 0 to apply logout functionality
func Logout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		utils.GenerateResponse(w, http.StatusMethodNotAllowed, "invalid request method")
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    "",
		Expires:  time.Unix(0, 0), // Set the expiration time to the past
		Path:     "/",
		Secure:   true, // Set to true if using HTTPS
		HttpOnly: true, // Cookie is only accessible by the server
	})

	utils.GenerateResponse(w, http.StatusOK, "logged out successfully")
}

// Return user info, and their roles for admin page
func AdminPage(w http.ResponseWriter, r *http.Request) {

	if r.Method != "GET" {
		utils.GenerateResponse(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	ctx := r.Context().Value(models.PayloadContextKey)

	payload := ctx.(jwt.MapClaims)
	userId := payload["userId"].(float64)
	role := payload["role"].(string)

	if userId == 0 || role == "" {
		utils.GenerateResponse(w, http.StatusBadRequest, "session timed out")
		return
	}

	user, err := models.GetRoles(int(userId))
	if err != nil {
		utils.GenerateResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	roles, ok := user[role]
	if !ok || len(roles) != 1 {
		utils.GenerateResponse(w, http.StatusInternalServerError, "invalid session")
		return
	}

	userWithRole := roles[0]
	userWithRole.Type = role

	utils.GenerateResponse(w, http.StatusOK, userWithRole)
}
