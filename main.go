package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
)

type apiConfig struct {
	fileserverHits int
}

func main() {
	const filepathRoot = "."
	const port = "8080"

	apiCfg := apiConfig{
		fileserverHits: 0,
	}

	router := chi.NewRouter()
	router.Mount("/", apiCfg.middlewareMetricsInc(http.FileServer(http.Dir(filepathRoot))))

	adminRouter := chi.NewRouter()
	router.Mount("/admin", adminRouter)
	adminRouter.Get("/metrics", apiCfg.handlerMetrics)

	apiRouter := chi.NewRouter()
	router.Mount("/api", apiRouter)
	apiRouter.Get("/healthz", handlerReadiness)
	apiRouter.Post("/validate_chirp", handlerChirpsValidate)

	corsMux := middlewareCors(router)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: corsMux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())
}
