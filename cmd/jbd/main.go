package main

import (
	"log"
	"net/http"
)

func main() {
	srv := &http.Server{
		Addr: ":8080",
	}

	log.Printf("Starting server on 127.0.0.1:8080")
	log.Fatal(srv.ListenAndServe())
}
