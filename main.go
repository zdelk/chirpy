package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"
	"time"

	"workspace/github.com/zdelk/chirpy/internal/database"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL must be set")
	}
	platform := os.Getenv("PLATFORM")
	if platform == "" {
		log.Fatal("PLATFORM must be set")
	}
	jwt_secret := os.Getenv("JWT_SECRET")
	if jwt_secret == "" {
		log.Fatal("JWT_SECRET must be set")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}
	dbQueries := database.New(db)
	const filepathRoot = "."
	const port = "8080"

	sMux := http.NewServeMux()
	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
		DB:             dbQueries,
		platform:       platform,
		secret:         jwt_secret,
	}

	sMux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))))

	sMux.HandleFunc("GET /api/healthz", handlerReadiness)
	sMux.HandleFunc("GET /api/chirps", apiCfg.handlerGetChirps)
	sMux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.handlerGetChirp)
	// sMux.HandleFunc("POST /api/validate_chirp", handlerValidate)
	sMux.HandleFunc("POST /api/users", apiCfg.handlerCreateUser)
	sMux.HandleFunc("POST /api/login", apiCfg.handlerLogin)
	sMux.HandleFunc("POST /api/chirps", apiCfg.handlerCreateChirp)

	sMux.HandleFunc("GET /admin/metrics", apiCfg.handlerCount)
	sMux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)

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
	platform       string
	secret         string
}

type User struct {
	ID             uuid.UUID `json:"id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	Email          string    `json:"email"`
	HashedPassword string    `json:"hashedpassword"`
}

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}
