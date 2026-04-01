package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/vilebile17/zimmer/internal/auth"
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
		Email    string `json:"email"`
		Name     string `json:"name"`
		Password string `json:"password"`
	}

	var params Params
	decoder := json.NewDecoder(request.Body)
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(response, request, "There was an error decoding the parameters (required {email:string, name:string, password:string})", err, http.StatusBadRequest)
		return
	}

	if params.Email == "" || params.Name == "" || params.Password == "" {
		respondWithError(response, request, "The email, name and password parameters can't be empty", nil, http.StatusBadRequest)
		return
	}

	if len(params.Password) < 8 {
		respondWithError(response, request, "The password must be at least 8 characters long", nil, http.StatusBadRequest)
		return
	}

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(response, request, "There was an error hashing the password", err, http.StatusBadRequest)
		return
	}

	dbUser, err := cfg.dbQueries.CreateUser(request.Context(), database.CreateUserParams{
		Name:           params.Name,
		Email:          params.Email,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		handleErrorsFromCreatingUser(err, response, request)
		return
	}

	fmt.Printf("Just made a user called '%v' with the email '%v'\n", dbUser.Name, dbUser.Email)
	fmt.Printf("The user's hashed password is %v\n", dbUser.HashedPassword)
	respondWithJSON(response, request, User{
		ID:        dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Email:     dbUser.Email,
		Name:      dbUser.Name,
	}, http.StatusCreated)
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

func (cfg *apiConfig) loginHandler(response http.ResponseWriter, request *http.Request) {
	type Params struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var params Params
	decoder := json.NewDecoder(request.Body)
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(response, request, "There was an error decoding the parameters (required {email:string, password:string})", err, http.StatusBadRequest)
		return
	}

	dbUser, err := cfg.dbQueries.GetUserFromEmail(request.Context(), params.Email)
	if err != nil {
		respondWithError(response, request, "couldn't find an account with that email", err, http.StatusBadRequest)
		return
	}

	passwordMatches, err := auth.CheckPasswordHash(params.Password, dbUser.HashedPassword)
	if err != nil {
		respondWithError(response, request, "There was an error validating the password", err, http.StatusBadRequest)
		return
	}

	if !passwordMatches {
		respondWithError(response, request, "Password doesn't match!", nil, http.StatusUnauthorized)
		return
	}

	respondWithJSON(response, request, User{
		ID:        dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Email:     dbUser.Email,
		Name:      dbUser.Name,
	}, http.StatusOK)
	fmt.Printf("User %v (%v) just logged in\n", dbUser.Name, dbUser.Email)
}
