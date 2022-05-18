package main

import (
	"log"
	"net/http"

	"github.com/pmaterer/junkboy"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	dbDSN := "junkboy.db"
	db, err := junkboy.NewSQLiteDB(dbDSN)

	if err != nil {
		log.Fatalf("failed to open db %s: %v", dbDSN, err)
	}

	anchorRepo := junkboy.NewAnchorSQLiteRepository(db)
	anchorService := junkboy.NewAnchorService(anchorRepo)
	anchorHandler := junkboy.NewAnchorHTTPHandler(anchorService)

	router := junkboy.NewRouter("/v1")
	anchorHandler.RegisterRoutes(router)

	mw := junkboy.NewCorsMiddleware(junkboy.NewLoggingMiddleware(router))

	srv := &http.Server{
		Addr:    ":8080",
		Handler: mw,
	}

	log.Printf("Starting server on 127.0.0.1:8080")
	log.Fatal(srv.ListenAndServe())
}
