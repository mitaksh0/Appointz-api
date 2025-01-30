package utils

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/appointments_api/models"
	"golang.org/x/crypto/bcrypt"
)

// generate API response
func GenerateResponse(w http.ResponseWriter, status int, data interface{}) {

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "http://127.0.0.1:5500")
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	r := models.Response{
		StatusCode: status,
		Message:    data,
	}

	res, err := json.Marshal(r)
	if err != nil {
		http.Error(w, "data parsing error", http.StatusInternalServerError)
	}

	w.WriteHeader(status)
	w.Write(res)

}

func GeneratePreflightRequest(w http.ResponseWriter, status int, data interface{}) {
	// Handle preflight requests
	w.Header().Set("Access-Control-Allow-Origin", "http://127.0.0.1:5500")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.WriteHeader(http.StatusNoContent) // 204 No Content
}

// HashPassword hashes the given password and returns the hashed password or an error.
func HashPassword(password string) (string, error) {
	// Generate the hashed password with a cost of bcrypt.DefaultCost
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}

	return string(hashedBytes), nil
}

// ComparePassword compares a plaintext password with a hashed password and returns true if they match.
func ComparePassword(hashedPassword, password string) bool {
	// Compare the hashed password with the plaintext password
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
