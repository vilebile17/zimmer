package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/vilebile17/zimmer/internal/auth"
	"github.com/vilebile17/zimmer/internal/database"
)

func (cfg *apiConfig) createClassHandler(response http.ResponseWriter, request *http.Request) {
	type Params struct {
		Name string `json:"name"`
	}

	var params Params
	decoder := json.NewDecoder(request.Body)
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(response, request, "There was an error decoding the parameters (required {name:string})", err, http.StatusBadRequest)
		return
	}

	if params.Name == "" {
		respondWithError(response, request, "The name parameter can't be empty", nil, http.StatusBadRequest)
		return
	}

	bearer, err := auth.GetBearerToken(request.Header)
	if err != nil {
		respondWithError(response, request, "There was an error finding the Authorization header", err, http.StatusUnauthorized)
		return
	}
	userID, err := auth.ValidateJWT(bearer, cfg.JWTSecret)
	if err != nil {
		respondWithError(response, request, "There was an error finding the user ID with that JWT", err, http.StatusUnauthorized)
		return
	}

	dbClass, err := cfg.dbQueries.CreateClass(request.Context(), database.CreateClassParams{
		Name:      params.Name,
		TeacherID: userID,
	})
	if err != nil {
		respondWithError(response, request, "There was an error adding a class to the database", err, http.StatusBadRequest)
		return
	}

	fmt.Printf("Just made a class called '%v' with the teacher ID '%v'\n", dbClass.Name, dbClass.TeacherID)
	respondWithJSON(response, request, dbClass, http.StatusCreated)
}
