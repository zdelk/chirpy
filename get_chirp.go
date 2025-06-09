package main

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerGetChirp(w http.ResponseWriter, res *http.Request) {

	chirpID := res.PathValue("chirpID")
	chirpUUID, err := uuid.Parse(chirpID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error parsing id", err)
	}

	chirp, err := cfg.DB.GetChirp(context.Background(), chirpUUID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error getting chirp", err)
	}
	respondWithJson(w, http.StatusOK, Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	})

}
