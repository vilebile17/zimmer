package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/vilebile17/zimmer/internal/database"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
}

func (cfg *apiConfig) createUserHandler(response http.ResponseWriter, request *http.Request) {
	type Params struct {
		Email string `json:"email"`
		Name  string `json:"name"`
	}

	var params Params
	decoder := json.NewDecoder(request.Body)
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(response, request, "There was an error decoding the parameters (required {email:string, name:string})", err, http.StatusBadRequest)
		return
	}

	if params.Email == "" || params.Name == "" {
		respondWithError(response, request, "The email and name parameters can't be empty", nil, http.StatusBadRequest)
		return
	}

	dbUser, err := cfg.dbQueries.CreateUser(request.Context(), database.CreateUserParams{
		Name:  params.Name,
		Email: params.Email,
	})
	if err != nil {
		handleErrorsFromCreatingUser(err, response, request)
		return
	}

	fmt.Printf("Just made a user called '%v' with the email '%v'\n", dbUser.Name, dbUser.Email)
	respondWithJSON(response, request, dbUser, http.StatusCreated)
}

func handleErrorsFromCreatingUser(err error, response http.ResponseWriter, request *http.Request) {
	if pqErr, ok := err.(*pq.Error); ok { // Basically, if this error is a postgres error
		if pqErr.Code.Name() == "unique_violation" { // if this error violates the uniqueness requirement of a collumn
			respondWithError(response, request, "You cannot have the same name or email as another user", err, http.StatusConflict)
			return
		}
	}
	respondWithError(response, request, "There was an error adding the user to the database", err, http.StatusBadRequest)
}
