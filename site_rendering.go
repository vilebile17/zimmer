package main

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) renderClass(response http.ResponseWriter, request *http.Request) {
	classID, err := uuid.Parse(request.PathValue("classID"))
	if err != nil {
		respondWithError(response, request, "There was an error parsing that classID", err, http.StatusBadRequest)
		return
	}

	class, err := cfg.dbQueries.GetClass(request.Context(), classID)
	if err != nil {
		respondWithError(response, request, "There was an error getting the class data from the database", err, http.StatusBadRequest)
		return
	}

	userID, err := cfg.getUserID(request)
	if err != nil {
		respondWithError(response, request, "Unable to get the userID from cookies and from headers", err, http.StatusUnauthorized)
		return
	}

	isInClass, err := cfg.isUserInThisClass(request.Context(), userID, classID)
	if err != nil {
		respondWithError(response, request, "Unable to check if user is in the class", err, http.StatusUnauthorized)
		return
	}
	if !isInClass {
		respondWithError(response, request, "You can only view a class if you are in it", err, http.StatusUnauthorized)
		return
	}

	tmpl, err := template.ParseFiles("./app/classes/index.html")
	if err != nil {
		respondWithError(response, request, "Failed to retrieve the html file from the /app folder", err, http.StatusInternalServerError)
		return
	}

	response.Header().Set("Content-Type", "text/html; charset=utf-8")
	response.WriteHeader(http.StatusOK)

	err = tmpl.Execute(response, struct {
		Name        string
		TeacherName string
		TeacherID   string
	}{
		Name:        class.Name,
		TeacherName: class.TeacherName,
		TeacherID:   class.TeacherID.String(),
	})
	if err != nil {
		fmt.Println(err)
	}
}

func (cfg *apiConfig) renderUser(response http.ResponseWriter, request *http.Request) {
	userID, err := uuid.Parse(request.PathValue("userID"))
	if err != nil {
		respondWithError(response, request, "There was an error parsing that userID", err, http.StatusBadRequest)
		return
	}

	user, err := cfg.dbQueries.GetUserFromID(request.Context(), userID)
	if err != nil {
		respondWithError(response, request, "There was an error getting user data from the database", err, http.StatusBadRequest)
		return
	}

	response.Header().Set("Content-Type", "text/html")
	response.WriteHeader(http.StatusOK)
	var data []byte
	response.Write(fmt.Appendf(data, `
	<!doctype html>
<html>
        <head>
                <meta charset="utf-8" />
                <meta
                        name="viewport"
                        content="width=device-width, initial-scale=1"
                />
                <title>Bester Zimmer</title>
                <link
                        rel="stylesheet"
                        href="https://fonts.googleapis.com/css2?family=Inter:wght@400;600;700&display=swap"
                />
                <link href="/default.css" rel="stylesheet" />
        </head>
        <body>
        	This is the user profile for the one and only <b>%v</b>
        </body>
</html>
	`, user.Name))
}
