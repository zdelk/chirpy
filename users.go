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

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	type response struct {
		User
	}

	params := parameters{}
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error Decoding Json: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode user", err)
		return
	}
	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Couldn't hash password", err)
		return
	}

	user, err := cfg.DB.CreateUser(context.Background(), database.CreateUserParams{
		Email:          params.Email,
		HashedPassword: hashedPassword,
	})

	if err != nil {
		log.Printf("Error creating User: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Couldn't create user", err)
		return
	}

	respondWithJson(w, http.StatusCreated, response{
		User: User{
			ID:             user.ID,
			CreatedAt:      user.CreatedAt,
			UpdatedAt:      user.UpdatedAt,
			Email:          user.Email,
			HashedPassword: user.HashedPassword,
			IsChirpyRed:    user.IsChirpyRed,
		},
	})

}

func (cfg *apiConfig) handlerUpdateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}
	type response struct {
		ID          uuid.UUID `json:"id"`
		CreatedAt   time.Time `json:"created_at"`
		UpdatedAt   time.Time `json:"updated_at"`
		Email       string    `json:"email"`
		IsChirpyRed bool      `json:"is_chirpy_red"`
	}
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Printf("no token in header: %v", err)
		respondWithError(w, http.StatusUnauthorized, "missing token", err)
	}
	userID, err := auth.ValidateJWT(token, cfg.secret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "invalid jwt", err)
		return
	}
	params := parameters{}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error decoding", err)
		return
	}

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error hashing password", err)
		return
	}
	user, err := cfg.DB.UpdateUser(context.Background(), database.UpdateUserParams{
		ID:             userID,
		Email:          params.Email,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error updating user", err)
		return
	}

	respondWithJson(w, http.StatusOK, response{
		ID:          user.ID,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
		Email:       user.Email,
		IsChirpyRed: user.IsChirpyRed,
	})

}
