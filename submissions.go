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

	userID, err := cfg.getUserID(request)
	if err != nil {
		respondWithError(response, request, "couldn't get userID from the auth header or from cookies :(", err, http.StatusUnauthorized)
		return
	}

	class, err := cfg.dbQueries.GetClass(request.Context(), classID)
	if err != nil {
		respondWithError(response, request, "couldn't get the class info from the database", err, http.StatusUnauthorized)
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

	if class.TeacherID == userID {
		respondWithError(response, request, "ur the teacher bro, y r u handing in??", nil, http.StatusUnauthorized)
		return
	} else if !inClass {
		respondWithError(response, request, "you're not a student in this class!", nil, http.StatusUnauthorized)
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
		return
	}

	respondWithJSON(response, request, submission, http.StatusCreated)
	fmt.Printf("User %v just submitted their work for assignment %v\n", userID, assignmentID)
}

type authStatus int

const (
	notInClass authStatus = iota
	Student
	Teacher
)

func (cfg *apiConfig) getSubmissionsHandler(response http.ResponseWriter, request *http.Request) {
	userID, err := cfg.getUserID(request)
	if err != nil {
		respondWithError(response, request, "couldn't get userID from the auth header or from cookies :(", err, http.StatusUnauthorized)
		return
	}
	classID, err := uuid.Parse(request.PathValue("classID"))
	if err != nil {
		respondWithError(response, request, "There was an error parsing that classID", err, http.StatusBadRequest)
		return
	}

	var status authStatus
	classesAsStudent, err := cfg.dbQueries.GetClassesAsStudent(request.Context(), userID)
	if err != nil {
		respondWithError(response, request, "failed to retrieve the users classes as a student", err, http.StatusUnauthorized)
		return
	}
	for _, class := range classesAsStudent {
		if class.ID == classID {
			status = Student
		}
	}
	class, err := cfg.dbQueries.GetClass(request.Context(), classID)
	if err != nil {
		respondWithError(response, request, "failed to retrieve the class data", err, http.StatusBadRequest)
		return
	}
	if class.TeacherID == userID {
		status = Teacher
	}

	assignmentID, err := uuid.Parse(request.PathValue("assignmentID"))
	if err != nil {
		respondWithError(response, request, "There was an error parsing that assignmentID", err, http.StatusBadRequest)
		return
	}

	switch status {
	case notInClass:
		respondWithError(response, request, "You cannot get a submission for an assignment if you're not even in the class bozo", err, http.StatusUnauthorized)
		return
	case Student:
		submission, err := cfg.dbQueries.GetSubmissionForUser(request.Context(), database.GetSubmissionForUserParams{
			AssignmentID: assignmentID,
			UserID:       userID,
		})
		if err != nil {
			respondWithError(response, request, "There was an error getting the submission", err, http.StatusBadRequest)
			return
		}
		respondWithJSON(response, request, submission, http.StatusOK)
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

func (cfg *apiConfig) updateSubmissionHandler(response http.ResponseWriter, request *http.Request) {
	classID, err := uuid.Parse(request.PathValue("classID"))
	if err != nil {
		respondWithError(response, request, "There was an error parsing that classID", err, http.StatusBadRequest)
		return
	}

	userID, err := cfg.getUserID(request)
	if err != nil {
		respondWithError(response, request, "couldn't get userID from the auth header or from cookies :(", err, http.StatusUnauthorized)
		return
	}

	class, err := cfg.dbQueries.GetClass(request.Context(), classID)
	if err != nil {
		respondWithError(response, request, "couldn't get the class info from the database", err, http.StatusUnauthorized)
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

	if class.TeacherID == userID {
		respondWithError(response, request, "ur the teacher bro, y r u handing in??", nil, http.StatusUnauthorized)
		return
	} else if !inClass {
		respondWithError(response, request, "you're not a student in this class!", nil, http.StatusUnauthorized)
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
		respondWithError(response, request, "you cannot update a submission after its assignment's due date has passed", err, http.StatusBadRequest)
		return
	}

	submission, err := cfg.dbQueries.GetSubmissionForUser(request.Context(), database.GetSubmissionForUserParams{
		AssignmentID: assignmentID,
		UserID:       userID,
	})
	if err != nil {
		respondWithError(response, request, "failed to retrieve old submission", err, http.StatusBadRequest)
		return
	}

	submission, err = cfg.dbQueries.UpdateSubmission(request.Context(), database.UpdateSubmissionParams{
		ID: submission.ID,
		Answers: sql.NullString{
			String: params.Answers,
			Valid:  params.Answers != "",
		},
	})
	if err != nil {
		respondWithError(response, request, "there was an error adding the submission to the database", err, http.StatusBadRequest)
		return
	}

	respondWithJSON(response, request, submission, http.StatusCreated)
	fmt.Printf("User %v just updated their work for assignment %v\n", userID, assignmentID)
}
