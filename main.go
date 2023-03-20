package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
)

type ApiConfig struct {
	fileServerHits int
}

func (cfg *ApiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileServerHits++
		next.ServeHTTP(w, r)
	})
}

func main() {
	config := ApiConfig{}

	router := chi.NewRouter()
	mux := http.NewServeMux()
	mux.Handle("/", config.middlewareMetricsInc(http.FileServer(http.Dir("."))))
	mux.Handle("assets/logo.png", http.FileServer(http.Dir("./assets/logo.png")))
	mux.HandleFunc("/healthz", healthzHandler)
	mux.HandleFunc("/metrics", config.numberOfHits)
	corsMux := middlewareCors(router)
	server := &http.Server{Handler: corsMux, Addr: "localhost:8080"}
	server.ListenAndServe()

}

func middlewareCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func healthzHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func (cfg *ApiConfig) numberOfHits(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	result := fmt.Sprintf("Hits: %d", cfg.fileServerHits)
	w.Write([]byte(result))

}
