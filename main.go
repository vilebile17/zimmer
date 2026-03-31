package main

import (
	"fmt"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	server_hits atomic.Int32
}

func HandlerGetEndpointPath(url string) http.FileSystem {
	return http.Dir("./app" + url)
}

func (cfg *apiConfig) middlewareIncServerHits(next http.Handler) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		cfg.server_hits.Add(1)
		next.ServeHTTP(response, request)
	})
}

func main() {
	const port = "8080"
	cfg := apiConfig{}
	mux := http.NewServeMux()
	server := http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	mux.Handle("/", cfg.middlewareIncServerHits(http.FileServer(HandlerGetEndpointPath(""))))
	mux.HandleFunc("/healthz", http.HandlerFunc(healthzHandler))
	mux.HandleFunc("/metrics", cfg.metricsHandler)

	fmt.Printf("Hosting Beste Zimmer at http://localhost:%s\n", port)
	if err := server.ListenAndServe(); err != nil {
		fmt.Println(err)
	}
}
