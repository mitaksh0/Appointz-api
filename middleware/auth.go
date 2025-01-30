package middleware

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/appointments_api/models"
	"github.com/appointments_api/utils"
	"github.com/golang-jwt/jwt/v4"
)

func JWTAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// handle preflight response
		if r.Method == http.MethodOptions {
			allowed := os.Getenv("ALLOWED_URL")
			w.Header().Set("Access-Control-Allow-Origin", allowed)
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.WriteHeader(http.StatusNoContent) // 204 No Content
			return
		}

		cookie, _ := r.Cookie("token")

		if cookie == nil {
			utils.GenerateResponse(w, http.StatusUnauthorized, "unauthorized")
			return
		}

		tokenString := cookie.Value
		if tokenString == "" /* || !strings.HasPrefix(tokenString, "Bearer ")*/ {
			utils.GenerateResponse(w, http.StatusForbidden, "forbidden")
			return
		}

		// tokenString = strings.TrimPrefix(tokenString, tokenString)
		secretKey := os.Getenv("SECRET_KEY")

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, http.ErrNotSupported
			}
			return []byte(secretKey), nil
		})

		if err != nil || !token.Valid {
			fmt.Println(err)
			fmt.Println("err")
			utils.GenerateResponse(w, http.StatusUnauthorized, "unauthorized")
			return
		}

		claimsMap := token.Claims.(jwt.MapClaims)

		ctx := context.WithValue(r.Context(), models.PayloadContextKey, claimsMap)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func PreflightResponse(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			allowed := os.Getenv("ALLOWED_URL")
			w.Header().Set("Access-Control-Allow-Origin", allowed)
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.WriteHeader(http.StatusNoContent) // 204 No Content
			return
		}

		next.ServeHTTP(w, r)
	})
}
