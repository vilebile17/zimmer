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

func (cfg *apiConfig) handInAssignmentHandler(response http.ResponseWriter, request *http.Request) {
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

	students, err := cfg.dbQueries.GetStudentsForClass(request.Context(), classID)
	if err != nil {
		respondWithError(response, request, "couldn't get students for the class", err, http.StatusUnauthorized)
		return
	}
	inClass := false
	for _, student := range students {
		if student.ID == userID {
			inClass = true
			break
		}
	}
	if !inClass {
		respondWithError(response, request, "you can only hand in to an assignment which you are a student in", err, http.StatusUnauthorized)
		return
	}

	type Params struct {
		Answers string `json:"answers"`
	}

	var params Params
	decoder := json.NewDecoder(request.Body)
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(response, request, "There was an error decoding the parameters (optional {answers:string})", err, http.StatusBadRequest)
		return
	}

	assignmentID, err := uuid.Parse(request.PathValue("assignmentID"))
	if err != nil {
		respondWithError(response, request, "There was an error parsing that assignmentID", err, http.StatusBadRequest)
		return
	}
	assignment, err := cfg.dbQueries.GetAssignmentFromID(request.Context(), assignmentID)
	if err != nil {
		respondWithError(response, request, "couldn't get the assignment from the database", err, http.StatusBadRequest)
		return
	}

	if !assignment.AllowLate && assignment.DueAt.Valid && time.Now().After(assignment.DueAt.Time) {
		respondWithError(response, request, "you can't hand in an assignment after its due date", err, http.StatusBadRequest)
		return
	}

	submission, err := cfg.dbQueries.CreateSubmission(request.Context(), database.CreateSubmissionParams{
		AssignmentID: assignmentID,
		UserID:       userID,
		Answers: sql.NullString{
			String: params.Answers,
			Valid:  params.Answers != "",
		},
	})
	respondWithJSON(response, request, submission, http.StatusCreated)
	fmt.Printf("User %v just submitted their work for assignment %v", userID, assignmentID)
}

func (cfg *apiConfig) getSubmissionsHandler(response http.ResponseWriter, request *http.Request) {
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

	class, err := cfg.dbQueries.GetClassFromClassID(request.Context(), classID)
	if err != nil {
		respondWithError(response, request, "couldn't get the class from the database", err, http.StatusUnauthorized)
		return
	}

	if class.TeacherID != userID {
		respondWithError(response, request, "you can't get the submissions for the assignment if you're not the teacher", nil, http.StatusUnauthorized)
		return
	}

	assignmentID, err := uuid.Parse(request.PathValue("assignmentID"))
	if err != nil {
		respondWithError(response, request, "There was an error parsing that assignmentID", err, http.StatusBadRequest)
		return
	}

	submissions, err := cfg.dbQueries.GetSubmissionsForAssignment(request.Context(), assignmentID)
	if err != nil {
		respondWithError(response, request, "couldn't get the submissions for that assignment (it may not exist)", err, http.StatusBadRequest)
		return
	}

	respondWithJSON(response, request, submissions, http.StatusOK)
}
