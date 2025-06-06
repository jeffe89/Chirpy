package main

import "net/http"

// Handler method for apiConfig struct to reset hit count and users
func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {

	// Check if in dev environment
	if cfg.platform != "dev" {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("Reset is only allowed in dev environment."))
		return
	}

	// Reset hit count to 0
	cfg.fileserverHits.Store(0)

	// Reset database
	err := cfg.db.Reset(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to reset the database: " + err.Error()))
		return
	}

	// Set Status OK and Respond with statement confirming database reset
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits reset to 0 and database reset to initial state."))
}
