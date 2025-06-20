package main

import (
	"context"
	"encoding/json"
	"net/http"
	"slices"
	"strings"
	"workspace/github.com/zdelk/chirpy/internal/auth"
	"workspace/github.com/zdelk/chirpy/internal/database"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerCreateChirp(w http.ResponseWriter, r *http.Request) {
	// handlerValidate(w, r)

	type parameters struct {
		Body string `json:"body"`
		// UserID uuid.UUID `json:"user_id"`
	}

	type response struct {
		Chirp
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't decode parameters", err)
		return
	}

	tokenString, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't retrieve bearer token", err)
		return
	}

	userID, err := auth.ValidateJWT(tokenString, cfg.secret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "error with validating jwt", err)
		return
	}

	// if check != params.UserID {
	// 	log.Println("Poster doesnt match acc")
	// 	log.Printf("Poster: %s", params.UserID)
	// 	log.Printf("acc: %s", check)
	// 	respondWithError(w, http.StatusUnauthorized, "Unauthorized to post on acc", nil)
	// 	return
	// }
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
		UserID: userID,
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't add chirp to db", err)
		return
	}

	respondWithJson(w, http.StatusCreated, response{
		Chirp: Chirp{
			ID:        newChirp.ID,
			CreatedAt: newChirp.CreatedAt,
			UpdatedAt: newChirp.UpdatedAt,
			Body:      newChirp.Body,
			UserID:    newChirp.UserID,
		},
	})

}

func (cfg *apiConfig) handlerDeleteChirp(w http.ResponseWriter, r *http.Request) {
	// Parse chirp ID
	chirpID := r.PathValue("chirpID")
	chirpUUID, err := uuid.Parse(chirpID)
	if err != nil {
		respondWithError(w, http.StatusForbidden, "failed to parse chirp id", err)
		return
	}

	// Check existence and retrieve
	chirp, err := cfg.DB.GetChirp(context.Background(), chirpUUID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "chirp doesnt exist", err)
		return
	}

	// Check bearer token
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "no token present in header", err)
		return
	}
	userID, err := auth.ValidateJWT(token, cfg.secret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate JWT", err)
		return
	}

	if userID != chirp.UserID {
		respondWithError(w, http.StatusForbidden, "Request not from Author", err)
		return
	}
	err = cfg.DB.DeleteChirp(context.Background(), chirpUUID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't delete chirp", err)
		return
	}
	respondWithJson(w, http.StatusNoContent, struct{}{})

}
