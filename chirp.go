package main

import (
	"context"
	"encoding/json"
	"net/http"
	"slices"
	"strings"
	"workspace/github.com/zdelk/chirpy/internal/database"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerCreateChirp(w http.ResponseWriter, r *http.Request) {
	// handlerValidate(w, r)

	type parameters struct {
		Body   string    `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't decode parameters", err)
		return
	}

	const maxChirpLength = 140

	if len(params.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	badWords := []string{"kerfuffle", "sharbert", "fornax"}

	cleanedWords := []string{}
	words := strings.Split(params.Body, " ")

	for _, word := range words {
		if slices.Contains(badWords, strings.ToLower(word)) {
			word = "****"
		}
		cleanedWords = append(cleanedWords, word)
	}

	cleanedText := strings.Join(cleanedWords, " ")

	newChirp, err := cfg.DB.CreateChirp(context.Background(), database.CreateChirpParams{
		Body:   cleanedText,
		UserID: params.UserID,
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't add chirp to db", err)
		return
	}

	respondWithJson(w, http.StatusCreated, newChirp)

}
