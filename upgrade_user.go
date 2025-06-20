package main

import (
	"context"
	"encoding/json"
	"net/http"

	"workspace/github.com/zdelk/chirpy/internal/auth"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerUpgradeUser(w http.ResponseWriter, r *http.Request) {
	type data struct {
		UserID uuid.UUID `json:"user_id"`
	}
	type parameters struct {
		Event string `json:"event"`
		Data  data   `json:"data"`
	}

	apiKey, err := auth.GetAPIKey(r.Header)
	if apiKey != cfg.apiKey {
		respondWithError(w, http.StatusUnauthorized, "Request not from Polka", err)
		return
	}

	params := parameters{}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	if params.Event != "user.upgraded" {
		respondWithJson(w, http.StatusNoContent, struct{}{})
		return
	}
	err = cfg.DB.UpgradeUser(context.Background(), params.Data.UserID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "User not Found", err)
		return
	}
	respondWithJson(w, http.StatusNoContent, struct{}{})

}
