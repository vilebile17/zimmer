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
	if err != nil {
		respondWithError(response, request, "there was an error adding the submission to the database", err, http.StatusBadRequest)
	}

	respondWithJSON(response, request, submission, http.StatusCreated)
	fmt.Printf("User %v just submitted their work for assignment %v", userID, assignmentID)
}

func (cfg *apiConfig) getSubmissionsHandler(response http.ResponseWriter, request *http.Request) {
	if statusCode, err := cfg.teacherActions(request); err != nil {
		respondWithError(response, request, "An error occured while trying to validate your identity", err, statusCode)
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

func (cfg *apiConfig) gradeSubmissionsHandler(response http.ResponseWriter, request *http.Request) {
	statusCode, err := cfg.teacherActions(request)
	if err != nil {
		respondWithError(response, request, "couldn't verify that the user is the teacher", err, statusCode)
		return
	}

	submissionID, err := uuid.Parse(request.PathValue("submissionID"))
	if err != nil {
		respondWithError(response, request, "unable to parse submissionID from the URL", err, http.StatusBadRequest)
		return
	}

	type Params struct {
		Score int `json:"score"`
	}

	var params Params
	decoder := json.NewDecoder(request.Body)
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(response, request, "There was an error decoding the parameters (required {score:int})", err, http.StatusBadRequest)
		return
	}

	if params.Score < 0 || params.Score > 100 {
		respondWithError(response, request, fmt.Sprintf("Invalid score: %v/100", params.Score), nil, http.StatusBadRequest)
		return
	}

	submission, err := cfg.dbQueries.GradeSubmission(request.Context(), database.GradeSubmissionParams{
		ID: submissionID,
		Score: sql.NullInt32{
			Int32: int32(params.Score),
			Valid: true,
		},
	})
	respondWithJSON(response, request, submission, http.StatusOK)
}
