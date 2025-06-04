package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

// Setup struct for expected JSON parameters
type parameters struct {
	Body string `json:"body"`
}

// Handler function to validate chirps
func handlerChirpsValidate(w http.ResponseWriter, r *http.Request) {

	// Setup struct for cleaned JSON response
	type returnVals struct {
		CleanedBody string `json:"cleaned_body"`
	}

	// Decode JSON and gather parameters
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	// Validate chirp length is within limit - if not, respond with error
	const maxChirpLength = 140
	if len(params.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	// Create a map of bad words to clean up
	badWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}

	// Clean up any bad words and replace with "****"
	cleaned := getCleanedBody(params.Body, badWords)

	// If chirp is valid, response with OK (200) status cleaned body text
	respondWithJSON(w, http.StatusOK, returnVals{
		CleanedBody: cleaned,
	})
}

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
