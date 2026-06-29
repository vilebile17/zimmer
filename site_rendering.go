package main

import (
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func (cfg *apiConfig) renderClass(response http.ResponseWriter, request *http.Request) {
	classID, err := uuid.Parse(request.PathValue("classID"))
	if err != nil {
		respondWithErrorPage(response, request, "Invalid classID", err, http.StatusBadRequest)
		return
	}

	class, err := cfg.dbQueries.GetClass(request.Context(), classID)
	if err != nil {
		respondWithErrorPage(response, request, "No class was found with that ID", err, http.StatusNotFound)
		return
	}

	userID, err := cfg.getUserID(request)
	if err != nil {
		respondWithErrorPage(response, request, "Unable to authenticate user", err, http.StatusUnauthorized)
		return
	}

	isInClass, err := cfg.isUserInThisClass(request.Context(), userID, classID)
	if err != nil {
		respondWithErrorPage(response, request, "Unable to authenticate user", err, http.StatusUnauthorized)
		return
	}
	if !isInClass {
		respondWithErrorPage(response, request, "You cannot view a class which you're not a member in", err, http.StatusForbidden)
		return
	}

	tmpl, err := template.ParseFiles("./app/classes/index.html")
	if err != nil {
		respondWithErrorPage(response, request, "Failed to retrieve the html file from the /app folder", err, http.StatusInternalServerError)
		return
	}

	response.Header().Set("Content-Type", "text/html; charset=utf-8")
	response.WriteHeader(http.StatusOK)

	err = tmpl.Execute(response, struct {
		Name        string
		TeacherName string
		TeacherID   string
		ClassID     string
	}{
		Name:        class.Name,
		TeacherName: class.TeacherName,
		TeacherID:   class.TeacherID.String(),
		ClassID:     classID.String(),
	})
	if err != nil {
		fmt.Println(err)
	}
}

func (cfg *apiConfig) renderUser(response http.ResponseWriter, request *http.Request) {
	userID, err := uuid.Parse(request.PathValue("userID"))
	if err != nil {
		respondWithErrorPage(response, request, "Invalid userID", err, http.StatusBadRequest)
		return
	}

	user, err := cfg.dbQueries.GetUserFromID(request.Context(), userID)
	if err != nil {
		respondWithErrorPage(response, request, "Couldn't find a user with that ID", err, http.StatusNotFound)
		return
	}

	tmpl, err := template.ParseFiles("./app/users/index.html")
	if err != nil {
		respondWithErrorPage(response, request, "Failed to retrieve the html file from the /app folder", err, http.StatusInternalServerError)
		return
	}

	response.Header().Set("Content-Type", "text/html; charset=utf-8")
	response.WriteHeader(http.StatusOK)

	err = tmpl.Execute(response, struct {
		Username string
		JoinedAt string
		Bio      string
	}{
		Username: user.Name,
		JoinedAt: user.CreatedAt.Format("02/01/2006"),
		Bio:      user.Bio,
	})
	if err != nil {
		fmt.Println(err)
	}
}

func (cfg *apiConfig) renderAssignment(response http.ResponseWriter, request *http.Request) {
	assignmentID, err := uuid.Parse(request.PathValue("assignmentID"))
	if err != nil {
		respondWithErrorPage(response, request, "Invalid assignmentID", err, http.StatusBadRequest)
		return
	}

	assignment, err := cfg.dbQueries.GetAssignmentFromID(request.Context(), assignmentID)
	if err != nil {
		respondWithErrorPage(response, request, "Couldn't find an assignment with that ID", err, http.StatusBadRequest)
		return
	}

	userID, err := cfg.getUserID(request)
	if err != nil {
		respondWithErrorPage(response, request, "Couldn't authenticate the user", err, http.StatusUnauthorized)
		return
	}

	isInClass, err := cfg.isUserInThisClass(request.Context(), userID, assignment.ClassID)
	if err != nil {
		respondWithErrorPage(response, request, "Couldn't authenticate the user", err, http.StatusUnauthorized)
		return
	}
	if !isInClass {
		respondWithErrorPage(response, request, "You cannot view the assignment as your not in the class", err, http.StatusForbidden)
		return
	}

	tmpl, err := template.ParseFiles("./app/assignments/index.html")
	if err != nil {
		respondWithErrorPage(response, request, "Failed to retrieve the html file from the /app folder", err, http.StatusInternalServerError)
		return
	}

	response.Header().Set("Content-Type", "text/html; charset=utf-8")
	response.WriteHeader(http.StatusOK)

	var dueAt string
	if assignment.DueAt.Valid {
		dueAt = "Due on " + assignment.DueAt.Time.Format(time.RFC1123)
	} else {
		dueAt = "No due date"
	}

	err = tmpl.Execute(response, struct {
		Title        string
		ID           string
		ClassID      string
		Instructions string
		DueAt        string
	}{
		Title:        assignment.Title,
		ID:           assignmentID.String(),
		ClassID:      assignment.ClassID.String(),
		Instructions: assignment.Instructions.String,
		DueAt:        dueAt,
	})
	if err != nil {
		fmt.Println(err)
	}
}

func (cfg *apiConfig) renderSubmission(response http.ResponseWriter, request *http.Request) {
	submissionID, err := uuid.Parse(request.PathValue("submissionID"))
	if err != nil {
		respondWithErrorPage(response, request, "Invalid submissionID", err, http.StatusBadRequest)
		return
	}

	submission, err := cfg.dbQueries.GetSubmission(request.Context(), submissionID)
	if err != nil {
		respondWithErrorPage(response, request, "Couldn't find a submission with that ID", err, http.StatusBadRequest)
		return
	}

	userID, err := cfg.getUserID(request)
	if err != nil {
		respondWithErrorPage(response, request, "Couldn't authenticate the user", err, http.StatusUnauthorized)
		return
	}

	class, err := cfg.dbQueries.GetClass(request.Context(), submission.ClassID)
	if err != nil {
		respondWithErrorPage(response, request, "Couldn't get class data", err, http.StatusBadRequest)
		return
	}

	if class.TeacherID != userID {
		respondWithError(response, request, "You must be a teacher to access this page", err, http.StatusForbidden)
		return
	}

	tmpl, err := template.ParseFiles("./app/submissions/index.html")
	if err != nil {
		respondWithErrorPage(response, request, "Failed to retrieve the html file from the /app folder", err, http.StatusInternalServerError)
		return
	}

	var grade string
	if submission.Grade.Valid {
		grade = fmt.Sprintf("%v", submission.Grade.Int32)
	} else {
		grade = ""
	}

	response.Header().Set("Content-Type", "text/html; charset=utf-8")
	response.WriteHeader(http.StatusOK)

	err = tmpl.Execute(response, struct {
		AssignmentName string
		AssignmentID   string
		StudentName    string
		OldGrade       string
		ClassID        string
		SubmissionID   string
		Work           string
		UpdatedAt      string
	}{
		AssignmentName: submission.AssignmentTitle,
		AssignmentID:   submission.AssignmentID.String(),
		StudentName:    submission.UserName,
		OldGrade:       grade,
		ClassID:        submission.ClassID.String(),
		SubmissionID:   submissionID.String(),
		Work:           submission.Answers.String,
		UpdatedAt:      submission.UpdatedAt.Format(time.RFC1123),
	})
	if err != nil {
		fmt.Println(err)
	}
}

func respondWithErrorPage(response http.ResponseWriter, request *http.Request, message string, err error, statusCode int) {
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(message)
	}

	tmpl, err := template.ParseFiles("./app/error.html")
	if err != nil {
		respondWithError(response, request, "Failed to retrieve the html file from the /app folder", err, http.StatusInternalServerError)
		return
	}

	response.Header().Set("Content-Type", "text/html; charset=utf-8")
	response.WriteHeader(http.StatusOK)

	err = tmpl.Execute(response, struct {
		ErrorCode    int
		ErrorMessage string
	}{
		ErrorCode:    statusCode,
		ErrorMessage: message,
	})
	if err != nil {
		fmt.Println(err)
	}
}
