package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
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

func (cfg *apiConfig) joinClassHandler(response http.ResponseWriter, request *http.Request) {
	classID, err := uuid.Parse(request.PathValue("classID"))
	if err != nil {
		respondWithError(response, request, "There was an error getting and parsing the class ID from the URL", err, http.StatusUnauthorized)
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

	users_classes, err := cfg.dbQueries.JoinClass(request.Context(), database.JoinClassParams{
		UserID:  userID,
		ClassID: classID,
	})
	if err != nil {
		respondWithError(response, request, "There was an error finding that class...", err, http.StatusBadRequest)
		return
	}

	fmt.Printf("User %v just joined a class %v\n", users_classes.UserID, classID)
	respondWithJSON(response, request, users_classes, http.StatusCreated)
}

func (cfg *apiConfig) getClassesForUserHandler(response http.ResponseWriter, request *http.Request) {
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

	classes, err := cfg.dbQueries.GetClassesForUserID(request.Context(), userID)
	if err != nil {
		respondWithError(response, request, "There was an error finding any classes for that user", err, http.StatusBadRequest)
		return
	}

	fmt.Printf("User %v just got their classes\n", userID)
	if len(classes) == 0 {
		type NoClasses struct {
			Message string `json:"message"`
		}
		respondWithJSON(response, request, NoClasses{Message: "Oof, it doesn't look like you're signed up for any classes yet"}, http.StatusOK)
		return
	}

	respondWithJSON(response, request, classes, http.StatusOK)
}
