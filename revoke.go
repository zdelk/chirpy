package main

import (
	"context"
	"log"
	"net/http"
	"workspace/github.com/zdelk/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Printf("error with getting bearer token: %v", err)
		respondWithError(w, http.StatusUnauthorized, "error getting bearer", err)
		return
	}

	err = cfg.DB.RevokeToken(context.Background(), token)
	if err != nil {
		log.Printf("error revoking token: %v", err)
		respondWithError(w, http.StatusBadRequest, "error revoking token", err)
		return
	}
	respondWithJson(w, http.StatusNoContent, struct{}{})
}
