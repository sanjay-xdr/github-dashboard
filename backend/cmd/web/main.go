package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/sanjay-xdr/github-dashboard/backend/internals/database"
)

func main() {

	database.InitMongo()

	srv := &http.Server{
		Addr:    ":8080",
		Handler: Routes(),
	}
	fmt.Print("Server is starting on port number 8080 ")
	err := srv.ListenAndServe()

	if err != nil {
		log.Fatal(err)
	}
}
