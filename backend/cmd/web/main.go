package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {

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
