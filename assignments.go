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
	var err error
	if statusCode, err := cfg.teacherActions(request); err != nil {
		respondWithError(response, request, "There was an error validating your identity", err, statusCode)
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
	} else {
		params.AllowLate = true
	}

	// It cannot error as it would have done that earlier
	classID, _ := uuid.Parse(request.PathValue("classID"))

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
	if err != nil {
		respondWithError(response, request, "There was an error adding the assignment to the database", err, http.StatusBadRequest)
	}

	respondWithJSON(response, request, assignment, http.StatusCreated)
	fmt.Printf("just made an assignment for class %v, with the title '%v'. It has a due date of %v, instructions of '%v' and the statement 'it allows late submissions' is %v\n", classID, params.Title, params.DueAt, params.Instructions, assignment.AllowLate)
}

func (cfg *apiConfig) getAssignmentsForAClassHandler(response http.ResponseWriter, request *http.Request) {
	classID, err := uuid.Parse(request.PathValue("classID"))
	if err != nil {
		respondWithError(response, request, "There was an error parsing that classID", err, http.StatusBadRequest)
		return
	}

	userID, err := cfg.getUserID(request)
	if err != nil {
		respondWithError(response, request, "couldn't get userID", err, http.StatusUnauthorized)
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

func (cfg *apiConfig) getNumAssignmentsHandler(response http.ResponseWriter, request *http.Request) {
	userID, err := cfg.getUserID(request)
	if err != nil {
		respondWithError(response, request, "couldn't get the user ID from headers nor cookies", err, http.StatusUnauthorized)
		return
	}

	num, err := cfg.dbQueries.GetNumAssignmentsToDo(request.Context(), userID)
	if err != nil {
		respondWithError(response, request, "couldn't get the number of assignments due from the database", err, http.StatusBadRequest)
		return
	}

	type Response struct {
		Num int64 `json:"num"`
	}
	r := Response{Num: num}
	respondWithJSON(response, request, r, http.StatusOK)
}
