package main

import (
	"context"
	"encoding/json"
	"net/http"
)

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, res *http.Request) {

	chirps, err := cfg.DB.GetAllChirps(context.Background())
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "error returning chirps", err)
		return
	}

	goodChirps := []Chirp{}
	for _, chirp := range chirps {
		nextChirp := Chirp{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		}
		goodChirps = append(goodChirps, nextChirp)
	}
	jsonData, err := json.Marshal(goodChirps)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error marshaling chirps", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}
