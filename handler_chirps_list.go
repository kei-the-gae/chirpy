package main

import "net/http"

func (cfg *apiConfig) handlerChirpsRetrieve(w http.ResponseWriter, r *http.Request) {
	dbChirps, err := cfg.db.GetChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get chirps", err)
		return
	}

	chirps := []Chirp{}
	for _, chirp := range dbChirps {
		chirps = append(chirps, mapToChirp(chirp))
	}

	respondWithJSON(w, http.StatusOK, chirps)
}
