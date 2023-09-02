package main

import (
	"log"
	"net/http"
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

	mux := http.NewServeMux()

	mux.Handle("/app/", apiCfg.middlewareMetric(http.StripPrefix("/app", http.FileServer(http.Dir(filePathRoot)))))
	mux.HandleFunc("/healthz", handlerReadiness)
	mux.HandleFunc("/metrics", apiCfg.handlerMetric)
	corsMux := middlewareCors(mux)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: corsMux,
	}

	log.Printf("Serving files from %s on port: %s\n", filePathRoot, port)
	log.Fatal(server.ListenAndServe())
}
