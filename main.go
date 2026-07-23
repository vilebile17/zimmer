package main

import (
	"context"
	"crypto/tls"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	dotenv "github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/vilebile17/zimmer/internal/database"
	"golang.org/x/crypto/acme/autocert"
)

type apiConfig struct {
	homePageViews atomic.Int32
	activeUsers   atomic.Int32
	dbQueries     *database.Queries
	JWTSecret     string
	port          string
	domain        string
	devMode       bool
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
	cfg.JWTSecret = os.Getenv("JWT_SECRET")
	cfg.port = ":" + os.Getenv("PORT")
	httpsMode := false

	if cfg.port == ":" {
		cfg.port = ":8080"
		fmt.Println("Defaulting PORT env variable to ':8080'")
	}
	if cfg.port == ":443" {
		httpsMode = true
		fmt.Println("HTTPS mode initiated")
	}

	cfg.domain = os.Getenv("DOMAIN")
	if cfg.domain == "" {
		cfg.domain = "localhost"
		fmt.Println("Defaulting DOMAIN env variable to 'localhost'")
	}

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

	// Determine if we are running in local development mode (domain is localhost)
	cfg.devMode = cfg.domain == "localhost"
	if cfg.devMode {
		fmt.Println("devMode initiated")
	}

	if httpsMode && !cfg.devMode {
		// Production HTTPS using Let's Encrypt
		whitelist := []string{cfg.domain, "www." + cfg.domain, "zimmer." + cfg.domain}

		certManager := autocert.Manager{
			Prompt:     autocert.AcceptTOS,
			HostPolicy: autocert.HostWhitelist(whitelist...),
			Cache:      autocert.DirCache("certs"),
		}

		// HTTP → HTTPS redirect and ACME challenge handling.
		go func() {
			http.ListenAndServe(":80", certManager.HTTPHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				target := "https://" + r.Host + r.URL.RequestURI()
				http.Redirect(w, r, target, http.StatusMovedPermanently)
			})))
		}()

		server.TLSConfig = &tls.Config{GetCertificate: certManager.GetCertificate}
	} else if httpsMode && cfg.devMode {
		// Development mode – load a self‑signed certificate for localhost.
		cert, err := tls.LoadX509KeyPair("dev.crt", "dev.key")
		if err != nil {
			log.Fatalf("failed to load dev TLS certificate: %v", err)
		}
		server.TLSConfig = &tls.Config{Certificates: []tls.Certificate{cert}}
	}

	// ======= API ENDPOINTS ========
	// general utility stuff
	mux.Handle("/", cfg.middlewareIncServerHits(http.FileServer(http.Dir("./app"))))
	mux.HandleFunc("/healthz", http.HandlerFunc(healthzHandler))
	mux.HandleFunc("/metrics", cfg.metricsHandler)
	mux.HandleFunc("POST /api/reset", cfg.resetHandler)

	// site stuff
	mux.HandleFunc("GET /c/{classID}", cfg.renderClass)
	mux.HandleFunc("GET /u/{userID}", cfg.renderUser)
	mux.HandleFunc("GET /a/{assignmentID}", cfg.renderAssignment)
	mux.HandleFunc("GET /s/{submissionID}", cfg.renderSubmission)

	// user stuff
	mux.HandleFunc("POST /api/users", cfg.createUserHandler)
	mux.HandleFunc("PUT /api/users", cfg.updateUserHandler)
	mux.HandleFunc("PUT /api/users/profile", cfg.updateUserProfileHandler)
	mux.HandleFunc("GET /api/users", cfg.getUserFromCookie)
	mux.HandleFunc("GET /api/users/{userID}", cfg.getUserHandler)
	mux.HandleFunc("DELETE /api/users", cfg.deleteUserHandler)
	mux.HandleFunc("POST /api/login", cfg.loginHandler)
	mux.HandleFunc("POST /api/logout", cfg.logoutHandler)

	// classes stuff
	mux.HandleFunc("POST /api/classes", cfg.createClassHandler)
	mux.HandleFunc("GET /api/classes", cfg.getClassesForUserHandler)
	mux.HandleFunc("GET /api/classes/{classID}", cfg.getClassHandler)
	mux.HandleFunc("DELETE /api/classes/{classID}", cfg.deleteClass)
	mux.HandleFunc("POST /api/classes/{classID}/members", cfg.joinClassHandler)
	mux.HandleFunc("DELETE /api/classes/{classID}/members/{userID}", cfg.removeFromClass)
	mux.HandleFunc("GET /api/classes/{classID}/members", cfg.getUsersForClassHandler)

	// assignments stuff
	mux.HandleFunc("POST /api/classes/{classID}/assignments", cfg.createAssignmentHandler)
	mux.HandleFunc("PUT /api/classes/{classID}/assignments/{assignmentID}", cfg.updateAssignmentHandler)
	mux.HandleFunc("GET /api/classes/{classID}/assignments", cfg.getAssignmentsForAClassHandler)
	mux.HandleFunc("GET /api/classes/{classID}/assignments/{assignmentID}", cfg.getAssignmentHandler)
	mux.HandleFunc("DELETE /api/classes/{classID}/assignments/{assignmentID}", cfg.deleteAssignmentHandler)
	mux.HandleFunc("GET /api/numAssignmentsDue", cfg.getNumAssignmentsHandler)

	// submissions stuff
	mux.HandleFunc("POST /api/classes/{classID}/assignments/{assignmentID}/submissions", cfg.handInAssignmentHandler)
	mux.HandleFunc("PUT /api/classes/{classID}/assignments/{assignmentID}/submissions", cfg.updateSubmissionHandler)
	mux.HandleFunc("GET /api/classes/{classID}/assignments/{assignmentID}/submissions", cfg.getSubmissionsHandler)
	mux.HandleFunc("PUT /api/classes/{classID}/assignments/{assignmentID}/submissions/{submissionID}", cfg.gradeSubmissionsHandler)

	// resource stuff
	mux.HandleFunc("POST /api/classes/{classID}/resources", cfg.createResourceHandler)
	mux.HandleFunc("GET /api/classes/{classID}/resources", cfg.getResourcesForClassHandler)
	mux.HandleFunc("GET /api/classes/{classID}/resources/{contentID}", cfg.getResourceHandler)
	mux.HandleFunc("PUT /api/classes/{classID}/resources/{contentID}", cfg.updateClassContentHandler)

	// announcement stuff
	mux.HandleFunc("POST /api/classes/{classID}/announcements", cfg.createAnnouncementHandler)
	mux.HandleFunc("GET /api/classes/{classID}/announcements", cfg.getAnnouncementsForClassHandler)
	mux.HandleFunc("GET /api/classes/{classID}/announcements/{contentID}", cfg.getAnnouncementHandler)
	mux.HandleFunc("PUT /api/classes/{classID}/announcements/{contentID}", cfg.updateClassContentHandler)

	if httpsMode {
		fmt.Printf("Hosting Zimmer at https://%s\n", cfg.domain)
		if err := server.ListenAndServeTLS("", ""); err != nil {
			fmt.Println(err)
		}
	} else {
		fmt.Printf("Hosting Zimmer at http://%s%s\n", cfg.domain, cfg.port)
		if err := server.ListenAndServe(); err != nil {
			fmt.Println(err)
		}
	}
}
