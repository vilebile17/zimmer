package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	dotenv "github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/vilebile17/zimmer/internal/database"
)

type apiConfig struct {
	homePageViews atomic.Int32
	activeUsers   atomic.Int32
	dbQueries     *database.Queries
	platform      string
	JWTSecret     string
	port          string
}

func main() {
	cfg := apiConfig{}

	err := dotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	cfg.dbQueries = database.New(db)
	cfg.platform = os.Getenv("PLATFORM")
	cfg.JWTSecret = os.Getenv("JWT_SECRET")
	cfg.port = ":" + os.Getenv("PORT")

	activeUsers, err := cfg.dbQueries.GetTotalUserCount(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	cfg.activeUsers.Add(int32(activeUsers))

	mux := http.NewServeMux()
	server := http.Server{
		Addr:    cfg.port,
		Handler: mux,
	}

	// general utility stuff
	mux.Handle("/", cfg.middlewareIncServerHits(http.FileServer(http.Dir("./app"))))
	mux.HandleFunc("/healthz", http.HandlerFunc(healthzHandler))
	mux.HandleFunc("/metrics", cfg.metricsHandler)
	mux.HandleFunc("POST /api/reset", cfg.resetHandler)

	// site stuff
	mux.HandleFunc("GET /c/{classID}", cfg.renderClass)
	mux.HandleFunc("GET /u/{userID}", cfg.renderUser)

	// user stuff
	mux.HandleFunc("POST /api/users", cfg.createUserHandler)
	mux.HandleFunc("PUT /api/users", cfg.updateUserHandler)
	mux.HandleFunc("GET /api/users", cfg.getUserFromCookie)
	mux.HandleFunc("GET /api/users/{userID}", cfg.getUserHandler)
	mux.HandleFunc("DELETE /api/users", cfg.deleteUserHandler)
	mux.HandleFunc("POST /api/login", cfg.loginHandler)

	// classes stuff
	mux.HandleFunc("POST /api/classes", cfg.createClassHandler)
	mux.HandleFunc("GET /api/classes", cfg.getClassesForUserHandler)
	mux.HandleFunc("GET /api/classes/{classID}", cfg.getClassHandler)
	mux.HandleFunc("POST /api/classes/{classID}/members", cfg.joinClassHandler)
	mux.HandleFunc("DELETE /api/classes/{classID}/members/{userID}", cfg.removeFromClass)
	mux.HandleFunc("GET /api/classes/{classID}/members", cfg.getUsersForClassHandler)

	// assignments stuff
	mux.HandleFunc("POST /api/classes/{classID}/assignments", cfg.createAssignmentHandler)
	mux.HandleFunc("GET /api/classes/{classID}/assignments", cfg.getAssignmentsForAClassHandler)
	mux.HandleFunc("GET /api/classes/{classID}/assignments/{assignmentID}", cfg.getAssignmentHandler)

	// submissions stuff
	mux.HandleFunc("POST /api/classes/{classID}/assignments/{assignmentID}/submissions", cfg.handInAssignmentHandler)
	mux.HandleFunc("GET /api/classes/{classID}/assignments/{assignmentID}/submissions", cfg.getSubmissionsHandler)
	mux.HandleFunc("PUT /api/classes/{classID}/assignments/{assignmentID}/submissions/{submissionID}", cfg.gradeSubmissionsHandler)

	fmt.Printf("Hosting Bester Zimmer at http://localhost%s\n", cfg.port)
	if err := server.ListenAndServe(); err != nil {
		fmt.Println(err)
	}
}
