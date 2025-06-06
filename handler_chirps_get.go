package main

import (
	"net/http"

	"github.com/google/uuid"
)

// Handler function to retreive one specified chirp from database
func (cfg *apiConfig) handlerChirpsGet(w http.ResponseWriter, r *http.Request) {

	// Get specified chirp ID
	chirpIDString := r.PathValue("chirpID")

	// Validate chirp ID is found
	chirpID, err := uuid.Parse(chirpIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID", err)
		return
	}

	// Retreive chirp from database via specified ID
	dbChirp, err := cfg.db.GetChirp(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't get chirp", err)
		return
	}

	// Call function to respond with JSON containing specified chirp data
	respondWithJSON(w, http.StatusOK, Chirp{
		ID:        dbChirp.ID,
		CreatedAt: dbChirp.CreatedAt,
		UpdatedAt: dbChirp.UpdatedAt,
		UserID:    dbChirp.UserID,
		Body:      dbChirp.Body,
	})
}

// Handler function to retrieve all chirps from database
func (cfg *apiConfig) handlerChirpsRetrieve(w http.ResponseWriter, r *http.Request) {

	// Retreive all chirps from database
	dbChirps, err := cfg.db.GetChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't retreive chirps", err)
		return
	}

	// Create an array for chirps
	chirps := []Chirp{}

	// Loop through each chirp and append to chirps array for JSON response
	for _, dbChirp := range dbChirps {
		chirps = append(chirps, Chirp{
			ID:        dbChirp.ID,
			CreatedAt: dbChirp.CreatedAt,
			UpdatedAt: dbChirp.UpdatedAt,
			UserID:    dbChirp.UserID,
			Body:      dbChirp.Body,
		})
	}

	// Call function to respond with JSON containing array for chirps
	respondWithJSON(w, http.StatusOK, chirps)
}
