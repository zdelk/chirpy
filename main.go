package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"workspace/github.com/zdelk/chirpy/internal/database"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error with db: %v", err)
	}

	dbQueries := database.New(db)
	const filepathRoot = "."
	const port = "8080"

	sMux := http.NewServeMux()
	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
		DB:             dbQueries,
	}

	sMux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))))

	sMux.HandleFunc("GET /api/healthz", handlerReadiness)
	sMux.HandleFunc("GET /admin/metrics", apiCfg.handlerCount)
	sMux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)
	sMux.HandleFunc("POST /api/validate_chirp", handlerValidate)

	localServer := &http.Server{
		Addr:    ":" + port,
		Handler: sMux,
	}

	log.Printf("Serving file from %s on port %s\n", filepathRoot, port)
	log.Fatal(localServer.ListenAndServe())

}

type apiConfig struct {
	fileserverHits atomic.Int32
	DB             *database.Queries
}
