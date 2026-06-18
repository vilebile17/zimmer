package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/vilebile17/zimmer/internal/auth"
	"github.com/vilebile17/zimmer/internal/database"
)

func (cfg *apiConfig) middlewareIncServerHits(next http.Handler) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		cfg.homePageViews.Add(1)
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
		return uuid.Nil, errors.New("there was an error getting the bearer token from the headers")
	}
	userID, err := auth.ValidateJWT(bearer, cfg.JWTSecret)
	if err != nil {
		return uuid.Nil, errors.New("couldn't find a user with that JWT token")
	}
	return userID, nil
}

func (cfg *apiConfig) getUserIDFromCookie(request *http.Request) (uuid.UUID, error) {
	cookie, err := request.Cookie("token")
	if err != nil {
		return uuid.Nil, errors.New("there was an error getting the token from the cookie")
	}
	bearer := cookie.Value

	userID, err := auth.ValidateJWT(bearer, cfg.JWTSecret)
	if err != nil {
		return uuid.Nil, errors.New("couldn't find a user with that JWT token")
	}
	return userID, nil
}

func (cfg *apiConfig) getUserID(request *http.Request) (uuid.UUID, error) {
	cookieID, cookieErr := cfg.getUserIDFromCookie(request)
	headerID, headerErr := cfg.getUserIDFromHeader(request.Header)

	if cookieErr == nil {
		return cookieID, nil
	} else if headerErr == nil {
		return headerID, nil
	} else {
		return uuid.Nil, fmt.Errorf("Couldn't get the token: %v, %v", cookieErr, headerErr)
	}
}

func (cfg *apiConfig) isUserInThisClass(context context.Context, userID, classID uuid.UUID) (bool, error) {
	var classes []database.Class

	classesAsStudent, err := cfg.dbQueries.GetClassesAsStudent(context, userID)
	if err != nil {
		return false, err
	}
	classes = append(classes, classesAsStudent...)

	classesAsTeacher, err := cfg.dbQueries.GetClassesAsTeacher(context, userID)
	if err != nil {
		return false, err
	}
	classes = append(classes, classesAsTeacher...)

	for _, class := range classes {
		if class.ID == classID {
			return true, nil
		}
	}
	return false, nil
}

func (cfg *apiConfig) teacherActions(request *http.Request) (int, error) {
	classID, err := uuid.Parse(request.PathValue("classID"))
	if err != nil {
		return http.StatusBadRequest, errors.New("There was an error parsing that classID")
	}

	userID, err := cfg.getUserID(request)
	if err != nil {
		return http.StatusUnauthorized, errors.New("couldn't get userID from the auth header")
	}

	class, err := cfg.dbQueries.GetClass(request.Context(), classID)
	if err != nil {
		return http.StatusUnauthorized, errors.New("couldn't get the class from the database")
	}

	if class.TeacherID != userID {
		return http.StatusUnauthorized, errors.New("you can't get the access this action without being the class teacher")
	}

	return 0, nil
}
