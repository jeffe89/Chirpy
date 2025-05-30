package main

import (
	"log"
	"net/http"
	"sync/atomic"
)

// Struct for in-memory data
type apiConfig struct {
	fileserverHits atomic.Int32
}

func main() {

	// Define constants
	const filepathRoot = "."
	const port = "8080"

	// Initialize an apiConfig struct
	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
	}

	// Create a new http.ServeMux
	mux := http.NewServeMux()

	// Setup file server handler with the /app/ path
	fsHandler := apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))
	mux.Handle("/app/", fsHandler)

	// Register a handler function for the /healthz path
	mux.HandleFunc("GET /healthz", handlerReadiness)
	// Register a handler function for the /metrics path to display hit count
	mux.HandleFunc("GET /metrics/", apiCfg.handlerMetrics)
	// Register a handler function for the /reset path to reset hit count
	mux.HandleFunc("POST /reset", apiCfg.handlerReset)

	// Create a new HTTP server struct
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	// Log information on files being served on particular port
	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	// Use servers ListenAndServe method to start server
	log.Fatal(srv.ListenAndServe())
}
