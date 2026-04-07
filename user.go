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

type UserAndJWT struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	Token     string    `json:"token"`
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
		Email            string `json:"email"`
		Password         string `json:"password"`
		ExpiresInSeconds int    `json:"expires_in_seconds"` // optional
	}

	var params Params
	decoder := json.NewDecoder(request.Body)
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(response, request, "There was an error decoding the parameters (required {email:string, password:string})", err, http.StatusBadRequest)
		return
	}

	if params.Email == "" || params.Password == "" {
		respondWithError(response, request, "Email and password can't be empty (required {email:string, password:string})", nil, http.StatusBadRequest)
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

	expiresIn := time.Hour
	if params.ExpiresInSeconds != 0 {
		expiresIn = time.Second * time.Duration(params.ExpiresInSeconds)
	}
	token, err := auth.MakeJWT(dbUser.ID, cfg.JWTSecret, expiresIn)
	if err != nil {
		respondWithError(response, request, "Failed to make the JWT token", nil, http.StatusUnauthorized)
		return
	}

	respondWithJSON(response, request, UserAndJWT{
		ID:        dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Email:     dbUser.Email,
		Name:      dbUser.Name,
		Token:     token,
	}, http.StatusOK)
	fmt.Printf("User %v (%v) just logged in\n", dbUser.Name, dbUser.Email)
}

func (cfg *apiConfig) updateUserHandler(response http.ResponseWriter, request *http.Request) {
	userID, err := cfg.getUserIDFromHeader(request.Header)
	if err != nil {
		respondWithError(response, request, "unable to get the userID from the auth header", err, http.StatusBadRequest)
		return
	}

	dbUser, err := cfg.dbQueries.GetUserFromID(request.Context(), userID)
	if err != nil {
		respondWithError(response, request, "couldn't get a user from the database with that ID", err, http.StatusBadRequest)
		return
	}

	type Params struct {
		OldPassword string `json:"old_password"`
		NewName     string `json:"new_name"`     //optional
		NewEmail    string `json:"new_email"`    //optional
		NewPassword string `json:"new_password"` //optional
	}

	var params Params
	decoder := json.NewDecoder(request.Body)
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(response, request, "There was an error decoding the parameters (required {old_password:string)", err, http.StatusBadRequest)
		return
	}

	if params.OldPassword == "" {
		respondWithError(response, request, "You must supply the old_password", nil, http.StatusUnauthorized)
		return
	}
	passwordMatches, err := auth.CheckPasswordHash(params.OldPassword, dbUser.HashedPassword)
	if err != nil {
		respondWithError(response, request, "There was an error validating the password", err, http.StatusBadRequest)
		return
	}
	if !passwordMatches {
		respondWithError(response, request, "Old password doesn't match!", nil, http.StatusUnauthorized)
		return
	}

	if params.NewPassword != "" && len(params.NewPassword) < 8 {
		respondWithError(response, request, "The password must be at least 8 characters long", nil, http.StatusBadRequest)
		return
	}

	hashedPassword, err := auth.HashPassword(params.NewPassword)
	if err != nil {
		respondWithError(response, request, "There was an error hashing the new password", err, http.StatusBadRequest)
		return
	}

	var newPassword string
	var newName string
	var newEmail string
	if params.NewPassword != "" {
		newPassword = hashedPassword
	} else {
		newPassword = dbUser.HashedPassword
	}
	if params.NewName != "" {
		newName = params.NewName
	} else {
		newName = dbUser.Name
	}
	if params.NewEmail != "" {
		newEmail = params.NewEmail
	} else {
		newEmail = dbUser.Email
	}

	newUser, err := cfg.dbQueries.UpdateUser(request.Context(), database.UpdateUserParams{
		ID:             userID,
		Name:           newName,
		Email:          newEmail,
		HashedPassword: newPassword,
	})
	if err != nil {
		respondWithError(response, request, "couldn't update user's information in the database. it's possible that you chose the same name or email as another user", err, http.StatusBadRequest)
		return
	}

	respondWithJSON(response, request, User{
		ID:        newUser.ID,
		CreatedAt: newUser.CreatedAt,
		UpdatedAt: newUser.UpdatedAt,
		Name:      newUser.Name,
		Email:     newUser.Email,
	}, http.StatusCreated)
	fmt.Printf("Just updated user %v to have a name of %v, an email of %v and a hashed password of... Wait, that's a secret :)\n", userID, newName, newEmail)
}

func (cfg *apiConfig) getUserHandler(response http.ResponseWriter, request *http.Request) {
	notLoggedIn := false
	requesterID, err := cfg.getUserIDFromHeader(request.Header)
	if err != nil {
		notLoggedIn = true
	}

	userID, err := uuid.Parse(request.PathValue("userID"))
	if err != nil {
		respondWithError(response, request, "There was an error getting and parsing the user ID from the URL", err, http.StatusUnauthorized)
		return
	}

	user, err := cfg.dbQueries.GetUserFromID(request.Context(), userID)
	if err != nil {
		respondWithError(response, request, "Couldn't get the user from the database", err, http.StatusUnauthorized)
		return
	}

	type UserNotEqualToRequester struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		Name      string    `json:"name"`
	}
	if notLoggedIn || userID != requesterID {
		respondWithJSON(response, request, UserNotEqualToRequester{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			Name:      user.Name,
		}, http.StatusOK)
	} else {
		respondWithJSON(response, request, User{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
			Name:      user.Name,
		}, http.StatusOK)
	}
}

func (cfg *apiConfig) deleteUserHandler(response http.ResponseWriter, request *http.Request) {
	userID, err := cfg.getUserIDFromHeader(request.Header)
	if err != nil {
		respondWithError(response, request, "unable to get the userID from the auth header", err, http.StatusBadRequest)
		return
	}

	user, err := cfg.dbQueries.GetUserFromID(request.Context(), userID)
	if err != nil {
		respondWithError(response, request, "unable to retrieve the user", err, http.StatusBadRequest)
		return
	}

	err = cfg.dbQueries.DeleteUser(request.Context(), userID)
	if err != nil {
		respondWithError(response, request, "unable to delete the user", err, http.StatusBadRequest)
		return
	}

	type Farewell struct {
		Message string `json:"message"`
	}
	respondWithJSON(response, request, Farewell{fmt.Sprintf("Goodbye %v, it's a shame to see you go :(", user.Name)}, http.StatusOK)
}
