package main

import (
	"context"
	"log"
	"net/http"
	"time"
	"workspace/github.com/zdelk/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {

	type response struct {
		Token string `json:"token"`
	}

	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Printf("error getting token: %v", err)
		respondWithError(w, http.StatusInternalServerError, "error getting refresh token", err)
		return
	}

	token, err := cfg.DB.GetRefreshToken(context.Background(), refreshToken)
	if err != nil {
		log.Printf("token not in database: %v", err)
		respondWithError(w, http.StatusUnauthorized, "unauthorized access", err)
		return
	}
	if token.RevokedAt.Valid {
		log.Printf("token revoked: %v", err)
		respondWithError(w, http.StatusUnauthorized, "token revoked", err)
		return
	}

	accessToken, err := auth.MakeJWT(token.UserID, cfg.secret, (1 * time.Hour))
	if err != nil {
		log.Printf("error making jwt: %v", err)
		respondWithError(w, http.StatusInternalServerError, "error making jwt", err)
		return
	}

	respondWithJson(w, http.StatusOK, response{
		Token: accessToken,
	})
}
