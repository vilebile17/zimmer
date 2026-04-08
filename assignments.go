package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/vilebile17/zimmer/internal/database"
)

func (cfg *apiConfig) createAssignmentHandler(response http.ResponseWriter, request *http.Request) {
	classID, err := uuid.Parse(request.PathValue("classID"))
	if err != nil {
		respondWithError(response, request, "There was an error parsing that classID", err, http.StatusBadRequest)
		return
	}

	userID, err := cfg.getUserIDFromHeader(request.Header)
	if err != nil {
		respondWithError(response, request, "couldn't get userID from the auth header", err, http.StatusUnauthorized)
		return
	}

	if class, err := cfg.dbQueries.GetClassFromClassID(request.Context(), classID); err != nil {
		respondWithError(response, request, "couldn't get the class from the database", err, http.StatusBadRequest)
		return
	} else if class.TeacherID != userID {
		respondWithError(response, request, "you can only make an assignment for a class which you are a teacher for", nil, http.StatusUnauthorized)
		return
	}

	type Params struct {
		Title        string    `json:"title"`
		DueAt        time.Time `json:"due_at"`
		Instructions string    `json:"instructions"`
		AllowLate    bool      `json:"allow_late"`
	}

	var params Params
	decoder := json.NewDecoder(request.Body)
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(response, request, "There was an error decoding the parameters (required {title:string}, optional {due_at:timestamp, instructions:string, allow_late:boolean})", err, http.StatusBadRequest)
		return
	}

	if params.Title == "" {
		respondWithError(response, request, "The title parameter can't be empty", nil, http.StatusBadRequest)
		return
	}

	assignment, err := cfg.dbQueries.CreateAssignment(request.Context(), database.CreateAssignmentParams{
		ClassID: classID,
		Title:   params.Title,
		DueAt: sql.NullTime{
			Time:  params.DueAt,
			Valid: !params.DueAt.IsZero(),
		},
		Instructions: sql.NullString{
			String: params.Instructions,
			Valid:  params.Instructions != "",
		},
		AllowLate: params.AllowLate,
	})

	respondWithJSON(response, request, assignment, http.StatusCreated)
	fmt.Printf("just made an assignment for class %v, with the title '%v'. It has a due date of %v, instructions of '%v' and the statement 'it allows late submissions' is %v\n", classID, params.Title, params.DueAt, params.Instructions, assignment.AllowLate)
}
