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
		Title        string `json:"title"`
		DueAt        string `json:"due_at"`
		Instructions string `json:"instructions"`
		AllowLate    bool   `json:"allow_late"`
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

	var dueAt time.Time
	if params.DueAt != "" {
		fmt.Println(params.DueAt)
		dueAt, err = time.Parse(time.RFC3339, params.DueAt)
		if err != nil {
			respondWithError(response, request, "couldn't parse the due_at parameter (it uses RFC3339)", err, http.StatusBadRequest)
			return
		}
	}

	assignment, err := cfg.dbQueries.CreateAssignment(request.Context(), database.CreateAssignmentParams{
		ClassID: classID,
		Title:   params.Title,
		DueAt: sql.NullTime{
			Time:  dueAt,
			Valid: !dueAt.IsZero(),
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

func (cfg *apiConfig) getAssignmentsForAClassHandler(response http.ResponseWriter, request *http.Request) {
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

	if isInClass, err := cfg.isUserInThisClass(request.Context(), userID, classID); err != nil {
		respondWithError(response, request, "couldn't find all users for this class", err, http.StatusUnauthorized)
		return
	} else if !isInClass {
		respondWithError(response, request, "you can only view assignments for classes which you are in", err, http.StatusUnauthorized)
		return
	}

	assignments, err := cfg.dbQueries.GetAssignmentsForAClass(request.Context(), classID)
	if err != nil {
		respondWithError(response, request, "couldn't retrieve the assignments from the database", err, http.StatusBadRequest)
		return
	}

	respondWithJSON(response, request, assignments, http.StatusOK)
	fmt.Printf("user %v, just got their assignments for class %v\n", userID, classID)
}

func (cfg *apiConfig) getAssignmentHandler(response http.ResponseWriter, request *http.Request) {
	classID, err := uuid.Parse(request.PathValue("classID"))
	if err != nil {
		respondWithError(response, request, "There was an error parsing that classID", err, http.StatusBadRequest)
		return
	}
	AssignmentID, err := uuid.Parse(request.PathValue("assignmentID"))
	if err != nil {
		respondWithError(response, request, "There was an error parsing that assignmentID", err, http.StatusBadRequest)
		return
	}

	userID, err := cfg.getUserIDFromHeader(request.Header)
	if err != nil {
		respondWithError(response, request, "couldn't get userID from the auth header", err, http.StatusUnauthorized)
		return
	}

	if isInClass, err := cfg.isUserInThisClass(request.Context(), userID, classID); err != nil {
		respondWithError(response, request, "couldn't find all users for this class", err, http.StatusUnauthorized)
		return
	} else if !isInClass {
		respondWithError(response, request, "you can only view assignments for classes which you are in", err, http.StatusUnauthorized)
		return
	}

	Assignment, err := cfg.dbQueries.GetAssignmentFromID(request.Context(), AssignmentID)
	if err != nil {
		respondWithError(response, request, "couldn't get that assignment from the database", err, http.StatusBadRequest)
		return
	}

	respondWithJSON(response, request, Assignment, http.StatusOK)
	fmt.Printf("Just got info about assignment %v: %v\n", AssignmentID, Assignment.Title)
}

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
