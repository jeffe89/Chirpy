package main

import (
	"log"
	"net/http"
	"time"
)

func main() {

	// Create a new http.ServeMux
	mux := http.NewServeMux()

	// Create a new HTTP server struct
	s := &http.Server{
		Addr:           ":8080",
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	// Register a handler function for the /healthz path
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Setup file server with the /app/ path
	fileServer := http.FileServer(http.Dir("./"))
	mux.Handle("/app/", http.StripPrefix("/app", fileServer))

	// Use servers ListenAndServe method to start server
	log.Fatal(s.ListenAndServe())
}
