package main

import (
	"net/http"
	"time"

	"chirpy/internal/auth"
)

// Handler function to refresh access token
func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {

	// Struct for JSON request parameters
	type response struct {
		Token string `json:"token"`
	}

	// Get refresh token via bearer token
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't find token", err)
		return
	}

	// Get user data via refresh token
	user, err := cfg.db.GetUserFromRefreshToken(r.Context(), refreshToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't get user for refresh token", err)
		return
	}

	// Make JWT access token
	accessToken, err := auth.MakeJWT(
		user.ID,
		cfg.jwtSecret,
		time.Hour,
	)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate token", err)
		return
	}

	// Respond with access token in JSON format
	respondWithJSON(w, http.StatusOK, response{
		Token: accessToken,
	})
}

// Handler function to revoke access token
func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {

	// Get refresh token via bearer token
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't find token", err)
		return
	}

	// Revoke refresh token and update database
	_, err = cfg.db.RevokeRefreshToken(r.Context(), refreshToken)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't revoke session", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
