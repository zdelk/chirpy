package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, req *http.Request) {
	godotenv.Load()
	if os.Getenv("PLATFORM") != "dev" {
		respondWithError(w, http.StatusForbidden, "forbidden", nil)
	}
	cfg.DB.DelUsers(context.Background())
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	cfg.fileserverHits.Swap(0)
	w.Write(fmt.Appendf([]byte{}, "Hits: %d", cfg.fileserverHits.Load()))
}
