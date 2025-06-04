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

	// Register a handler function for the /api/healthz path to display status of server
	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	// Register a handler function for the /api/validate_chirp to handle message posts
	mux.HandleFunc("POST /api/validate_chirp", handlerChirpsValidate)
	// Register a handler function for the /admin/reset path to reset hit count
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)
	// Register a handler function for the /admin/metrics path to display hit count
	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)

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
