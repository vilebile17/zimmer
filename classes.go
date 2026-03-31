package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/vilebile17/beste_zimmer/internal/database"
)

type Class struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
	TeacherID uuid.UUID `json:"teacher_id"`
}

func (cfg *apiConfig) createClassHandler(response http.ResponseWriter, request *http.Request) {
	type Params struct {
		Name      string    `json:"name"`
		TeacherID uuid.UUID `json:"teacher_id"`
	}

	var params Params
	decoder := json.NewDecoder(request.Body)
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(response, request, "There was an error decoding the parameters (required {name:string, teacher_id:UUID})", err, http.StatusBadRequest)
		return
	}

	if params.Name == "" || params.TeacherID.String() == "" {
		respondWithError(response, request, "The teacher_id and name parameters can't be empty", nil, http.StatusBadRequest)
		return
	}

	dbClass, err := cfg.dbQueries.CreateClass(request.Context(), database.CreateClassParams{
		Name:      params.Name,
		TeacherID: params.TeacherID,
	})
	if err != nil {
		respondWithError(response, request, "There was an error adding a class to the database", err, http.StatusBadRequest)
		return
	}

	fmt.Printf("Just made a class called '%v' with the teacher ID '%v'\n", dbClass.Name, dbClass.TeacherID)
	respondWithJSON(response, request, dbClass, http.StatusCreated)
}
