package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"
	"workspace/github.com/zdelk/chirpy/internal/auth"
	"workspace/github.com/zdelk/chirpy/internal/database"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	type response struct {
		ID           uuid.UUID `json:"id"`
		CreatedAt    time.Time `json:"created_at"`
		UpdatedAt    time.Time `json:"updated_at"`
		Email        string    `json:"email"`
		Token        string    `json:"token"`
		RefreshToken string    `json:"refresh_token"`
		IsChirpyRed  bool      `json:"is_chirpy_red"`
	}

	params := parameters{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Couldn't decode login: %v", err)
		respondWithError(w, http.StatusInternalServerError, "error decoding login", err)
		return
	}

	user, err := cfg.DB.GetEmail(context.Background(), params.Email)
	if err != nil {
		log.Printf("Couldn't return user: %v", err)
		respondWithError(w, http.StatusInternalServerError, "error returning user", err)
		return
	}

	expiresIn := 1 * time.Hour
	tokenString, err := auth.MakeJWT(user.ID, cfg.secret, expiresIn)

	if err != nil {
		log.Printf("error generating jwt: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Error creating jwt", err)
	}

	err = auth.CheckPasswordHash(user.HashedPassword, params.Password)
	if err != nil {
		log.Printf("incorrect password: %v", err)
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}
	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		log.Printf("error making refresh token: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Error creating refresh token", err)
		return
	}

	_, err = cfg.DB.CreateRefreshToken(context.Background(), database.CreateRefreshTokenParams{
		Token:     refreshToken,
		UserID:    user.ID,
		ExpiresAt: time.Now().AddDate(0, 0, 60),
	})
	if err != nil {
		log.Printf("error adding refresh token to database: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Error adding token to database", err)
		return
	}

	respondWithJson(w, http.StatusOK, response{
		ID:           user.ID,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
		Email:        user.Email,
		Token:        tokenString,
		RefreshToken: refreshToken,
		IsChirpyRed:  user.IsChirpyRed,
	})
}
