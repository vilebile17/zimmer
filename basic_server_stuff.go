package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func HandlerGetEndpointPath(url string) http.FileSystem {
	return http.Dir("./app" + url)
}

func (cfg *apiConfig) middlewareIncServerHits(next http.Handler) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		cfg.server_hits.Add(1)
		next.ServeHTTP(response, request)
	})
}

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
    <h1>Welcome, Beste Zimmer Admin</h1>
    <p>The home page has been visited %d times!</p>
  </body>
</html>`, cfg.server_hits.Load())))
}

func (cfg *apiConfig) resetHandler(response http.ResponseWriter, request *http.Request) {
	if cfg.platform != "dev" {
		respondWithError(response, request, "You must be a developer to reset the user database", nil, http.StatusForbidden)
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

func respondWithError(response http.ResponseWriter, _ *http.Request, message string, err error, statusCode int) {
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(message)
	}

	type ErrorJSON struct {
		Error string `json:"error"`
	}

	errorJSON := ErrorJSON{
		message,
	}

	data, err := json.Marshal(errorJSON)
	if err != nil {
		fmt.Printf("Error encoding error message into json: %s\n", err)
		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(statusCode)
	response.Write(data)
}

func respondWithJSON(response http.ResponseWriter, request *http.Request, payload any, statusCode int) {
	data, err := json.Marshal(payload)
	if err != nil {
		fmt.Println(err)
		respondWithError(response, request, "an error occured while marshalling the payload", err, http.StatusBadRequest)
		return
	}

	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(statusCode)
	response.Write(data)
}
