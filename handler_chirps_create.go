package main

import (
	"chirpy/internal/database"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Struct for chirp to be stored in database
type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	UserID    uuid.UUID `json:"user_id"`
	Body      string    `json:"body"`
}

// Handler function to validate and create chirps
func (cfg *apiConfig) handlerChirpsCreate(w http.ResponseWriter, r *http.Request) {

	// Setup struct for expected JSON parameters
	type parameters struct {
		Body   string    `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}

	// Decode JSON and gather parameters
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	// Call function to validate chirp body
	cleaned, err := validateChirp(params.Body)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error(), err)
		return
	}

	// Create chirp in database
	chirp, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:   cleaned,
		UserID: params.UserID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create chirp", err)
		return
	}

	// If chirp is valid, respond with 201 status code and full chirp resource
	respondWithJSON(w, http.StatusCreated, Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	})
}

// Function to validate chirp length and content
func validateChirp(body string) (string, error) {

	// Validate chirp length is within limit - if not, respond with error
	const maxChirpLength = 140
	if len(body) > maxChirpLength {
		return "", errors.New("Chirp is too long")

	}

	// Create a map of bad words to clean up
	badWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}

	// call function to clean up any bad words and replace with "****"
	cleaned := getCleanedBody(body, badWords)

	// Return cleaned body string
	return cleaned, nil
}

// Function to clean up chirp body for designated bodywords
func getCleanedBody(body string, badWords map[string]struct{}) string {

	// Take param body and split into word slices
	words := strings.Split(body, " ")

	// Loop through each word and check against slice of bad words. Replace if found
	for i, word := range words {
		loweredWord := strings.ToLower(word)
		if _, ok := badWords[loweredWord]; ok {
			words[i] = "****"
		}
	}

	// Join the slice back into a complete string
	cleaned := strings.Join(words, " ")

	// Return cleaned body string
	return cleaned
}
