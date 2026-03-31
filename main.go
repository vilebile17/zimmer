package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	dotenv "github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/vilebile17/beste_zimmer/internal/database"
)

type apiConfig struct {
	server_hits atomic.Int32
	dbQueries   *database.Queries
	platform    string
}

func main() {
	const port = "8080"
	cfg := apiConfig{}

	dotenv.Load()
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	cfg.dbQueries = database.New(db)
	cfg.platform = os.Getenv("PLATFORM")

	mux := http.NewServeMux()
	server := http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	mux.Handle("/", cfg.middlewareIncServerHits(http.FileServer(HandlerGetEndpointPath(""))))
	mux.HandleFunc("/healthz", http.HandlerFunc(healthzHandler))
	mux.HandleFunc("/metrics", cfg.metricsHandler)
	mux.HandleFunc("POST /api/users", cfg.createUserHandler)
	mux.HandleFunc("POST /api/classes", cfg.createClassHandler)
	mux.HandleFunc("POST /api/reset", cfg.resetHandler)

	fmt.Printf("Hosting Beste Zimmer at http://localhost:%s\n", port)
	if err := server.ListenAndServe(); err != nil {
		fmt.Println(err)
	}
}
