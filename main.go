package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/appointments_api/db"
	"github.com/joho/godotenv"
)

func main() {

	// load env file
	envErr := godotenv.Load()
	if envErr != nil {
		log.Fatalf("error loading .env file: %v", envErr)
	}

	// initialize database
	dbErr := db.InitDatabase()
	if dbErr != nil {
		log.Fatalf("error initializing database: %v", dbErr)
	}
	fmt.Println("database connected successfully...")
	defer db.Db.Close()

	// init routes
	port := ":8080"
	router := http.NewServeMux()
	initRoutes(router)

	// start listening
	fmt.Println("server running on port", port)
	if err := http.ListenAndServe(port, router); err != nil {
		fmt.Println(err.Error())
	}
}
