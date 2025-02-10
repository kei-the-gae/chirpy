package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/kei-the-gae/chirpy/internal/auth"
	"github.com/kei-the-gae/chirpy/internal/database"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	params := parameters{}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	// Verify parameters exist
	if params.Email == "" {
		respondWithError(w, http.StatusBadRequest, "Email is required", nil)
		return
	}
	if params.Password == "" {
		respondWithError(w, http.StatusBadRequest, "Password is required", nil)
		return
	}

	// Verify credentials
	user, err := cfg.db.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid email or password", nil)
		return
	}
	if err := auth.CheckPasswordHash(params.Password, user.HashedPassword); err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid email or password", nil)
		return
	}

	// Create JWT Access Token
	expirationTime := time.Hour
	accessToken, err := auth.MakeJWT(
		user.ID,
		cfg.jwtSecret,
		expirationTime,
	)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create access JWT", err)
		return
	}
	// Create Refresh Token
	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create refresh token", err)
		return
	}

	_, err = cfg.db.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		UserID:    user.ID,
		Token:     refreshToken,
		ExpiresAt: time.Now().UTC().Add(time.Hour * 24 * 60),
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't save refresh token", err)
		return
	}

	// Response
	respondWithJSON(w, http.StatusOK, mapToUser(user, accessToken, refreshToken))
}
