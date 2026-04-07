package main

import (
	"fmt"
	"net/http"
)

func healthzHandler(response http.ResponseWriter, _ *http.Request) {
	response.Header().Set("Content-Type", "text/plain; charset=utf-8")
	response.WriteHeader(http.StatusOK)
	response.Write([]byte("OK"))
}

func (cfg *apiConfig) metricsHandler(response http.ResponseWriter, _ *http.Request) {
	response.Header().Set("Content-Type", "text/html; charset=utf-8")
	response.WriteHeader(http.StatusOK)
	response.Write([]byte(fmt.Sprintf(`
<html>
<body>
    <h1>Welcome, Bester Zimmer Admin</h1>
    <p>The home page has been visited %d times</p>
    <p>There are currently %v active users</p>
  </body>
</html>`, cfg.homePageViews.Load(), cfg.activeUsers.Load())))
}

func (cfg *apiConfig) resetHandler(response http.ResponseWriter, request *http.Request) {
	if cfg.platform != "dev" {
		respondWithError(response, request, "You must be a developer to reset the user database", nil, http.StatusForbidden)
		fmt.Println(cfg.platform)
		return
	}

	err := cfg.dbQueries.ResetUsers(request.Context())
	if err != nil {
		respondWithError(response, request, "There was an error resetting the users table", err, http.StatusBadRequest)
		return
	}

	fmt.Println("Users table successfully reset")
	response.WriteHeader(http.StatusOK)
	response.Write([]byte("reset users table\n"))
}
