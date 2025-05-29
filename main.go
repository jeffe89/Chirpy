package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
	"time"
)

// Struct for in-memory data
type apiConfig struct {
	fileserverHits atomic.Int32
}

func main() {

	// Initialize an apiConfig struct
	apiCfg := apiConfig{}

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
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Register a handler function for the /metrics path to display hit count
	mux.HandleFunc("GET /metrics/", apiCfg.hitCounter)

	// Register a handler function for the /reset path to reset hit count
	mux.HandleFunc("POST /reset", apiCfg.resetCounter)

	// Setup file server with the /app/ path
	fileServer := http.FileServer(http.Dir("./"))
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", fileServer)))

	// Use servers ListenAndServe method to start server
	log.Fatal(s.ListenAndServe())
}

// Middleware method to increment server hits
func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

// Handler method for apiConfig struct to display hit count
func (a *apiConfig) hitCounter(w http.ResponseWriter, r *http.Request) {
	hitCount := fmt.Sprintf("Hits: %v", a.fileserverHits.Load())
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte(hitCount))
}

// Handler method for apiConfig struct to reset hit count
func (a *apiConfig) resetCounter(w http.ResponseWriter, r *http.Request) {
	a.fileserverHits.Store(0)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte("Hit Count Reset!"))
}
