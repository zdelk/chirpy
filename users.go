package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
)

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Email string `json:"email"`
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

	user, err := cfg.DB.CreateUser(context.Background(), params.Email)
	if err != nil {
		log.Printf("Error creating User: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Couldn't create user", nil)
		return
	}

	respondWithJson(w, http.StatusCreated, response{
		User: User{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
		},
	})

}
