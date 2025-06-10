package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"workspace/github.com/zdelk/chirpy/internal/auth"
	"workspace/github.com/zdelk/chirpy/internal/database"
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
		},
	})

}
