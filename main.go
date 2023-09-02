package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type apiConfig struct {
	fileseverHits int
}

func main() {
	const port = "8000"
	const filePathRoot = "."

	apiCfg := apiConfig{
		fileseverHits: 0,
	}

	router := chi.NewRouter()
	apiRouter := chi.NewRouter()
	metricRouter := chi.NewRouter()
	fsHandler := apiCfg.middlewareMetric(http.StripPrefix("/app", http.FileServer(http.Dir(filePathRoot))))
	router.Handle("/app", fsHandler)
	router.Handle("/app/*", fsHandler)

	apiRouter.Get("/healthz", handlerReadiness)
	apiRouter.Post("/validate_chirp", ChirpsValidate)

	metricRouter.Get("/metrics", apiCfg.handlerMetric)

	router.Mount("/api", apiRouter)
	router.Mount("/admin", metricRouter)
	corsMux := middlewareCors(router)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: corsMux,
	}

	log.Printf("Serving files from %s on port: %s\n", filePathRoot, port)
	log.Fatal(server.ListenAndServe())
}
