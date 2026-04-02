package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/vilebile17/zimmer/internal/auth"
)

func (cfg *apiConfig) middlewareIncServerHits(next http.Handler) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		cfg.server_hits.Add(1)
		next.ServeHTTP(response, request)
	})
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

func (cfg *apiConfig) getUserIDFromHeader(header http.Header) (uuid.UUID, error) {
	bearer, err := auth.GetBearerToken(header)
	if err != nil {
		return uuid.Nil, errors.New("there was an error getting the bearer token")
	}
	userID, err := auth.ValidateJWT(bearer, cfg.JWTSecret)
	if err != nil {
		return uuid.Nil, errors.New("couldn't find a user with that JWT token")
	}
	return userID, nil
}
