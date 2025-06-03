package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

func main() {
	const filepathRoot = "."
	const port = "8080"

	sMux := http.NewServeMux()
	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
	}

	sMux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))))

	sMux.HandleFunc("/healthz", handlerReadiness)
	sMux.HandleFunc("/metrics", apiCfg.handlerCount)
	sMux.HandleFunc("/reset", apiCfg.handlerReset)

	localServer := &http.Server{
		Addr:    ":" + port,
		Handler: sMux,
	}

	log.Printf("Serving file from %s on port %s\n", filepathRoot, port)
	log.Fatal(localServer.ListenAndServe())

}

func handlerReadiness(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

type apiConfig struct {
	fileserverHits atomic.Int32
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	cfg.fileserverHits.Add(1)

	return next
}

func (cfg *apiConfig) handlerCount(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(fmt.Appendf([]byte{}, "Hits: %d", cfg.fileserverHits.Load()))
}

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	cfg.fileserverHits.Swap(0)
	w.Write(fmt.Appendf([]byte{}, "Hits: %d", cfg.fileserverHits.Load()))
}
