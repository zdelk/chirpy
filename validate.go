package main

import (
	"encoding/json"
	"net/http"
	"slices"
	"strings"
)

func handlerValidate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	type returnVals struct {
		// Valid        bool   `json:"valid"`
		Cleaned_Body string `json:"cleaned_body"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
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

	respondWithJson(w, http.StatusOK, returnVals{
		Cleaned_Body: cleanedText,
	})
}
