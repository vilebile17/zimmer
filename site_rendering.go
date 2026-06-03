package main

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) renderClass(response http.ResponseWriter, request *http.Request) {
	classID, err := uuid.Parse(request.PathValue("classID"))
	if err != nil {
		respondWithError(response, request, "There was an error parsing that classID", err, http.StatusBadRequest)
		return
	}

	class, err := cfg.dbQueries.GetClassFromClassID(request.Context(), classID)
	if err != nil {
		respondWithError(response, request, "There was an error getting the class data from the database", err, http.StatusBadRequest)
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
        	This is the classes page for %v
        </body>
</html>
	`, class.Name))
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
