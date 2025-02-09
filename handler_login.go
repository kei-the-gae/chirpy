package main

import (
	"encoding/json"
	"net/http"

	"github.com/kei-the-gae/chirpy/internal/auth"
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
	if params.Email == "" {
		respondWithError(w, http.StatusBadRequest, "Email is required", nil)
		return
	}
	if params.Password == "" {
		respondWithError(w, http.StatusBadRequest, "Password is required", nil)
		return
	}
	user, err := cfg.db.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid email or password", nil)
		return
	}
	if err := auth.CheckPasswordHash(params.Password, user.HashedPassword); err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid email or password", nil)
		return
	}
	respondWithJSON(w, http.StatusOK, mapToUser(user))
}
