package main

import (
	"context"
	"encoding/json"
	"net/http"
	"sort"
	"workspace/github.com/zdelk/chirpy/internal/database"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	var chirps []database.Chirp
	var err error
	userID := r.URL.Query().Get("author_id")
	sortKey := r.URL.Query().Get("sort")

	if userID == "" {
		chirps, err = cfg.DB.GetAllChirps(context.Background())
	} else {
		userUUID, err := uuid.Parse(userID)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "error parsing userID", err)
			return
		}
		chirps, err = cfg.DB.GetChirpsAuthor(context.Background(), userUUID)
	}

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

	if sortKey == "desc" {
		sort.Slice(goodChirps, func(i, j int) bool {
			return goodChirps[i].CreatedAt.After(goodChirps[j].CreatedAt)
		})
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
