package main

import (
	"net/http"

	"github.com/appointments_api/controllers"
	"github.com/appointments_api/middleware"
)

func initRoutes(router *http.ServeMux) {

	// routes without JWT middleware
	router.Handle("/login", middleware.PreflightResponse(http.HandlerFunc(controllers.Login)))
	router.Handle("/register", middleware.PreflightResponse(http.HandlerFunc(controllers.Register)))

	// routes with JWT middleware

	// AUTH ROUTES
	router.Handle("/logout", middleware.JWTAuth(http.HandlerFunc(controllers.Logout)))
	// router.Handle("/register", middleware.JWTAuth(http.HandlerFunc(controllers.Register)))

	// USER ROUTES
	router.Handle("/patients", middleware.JWTAuth(http.HandlerFunc(controllers.PatientsHandler)))
	// router.Handle("/patients", http.HandlerFunc(controllers.PatientsHandler))
	router.Handle("/appointments", http.HandlerFunc(controllers.AppointmentHandler))

}
