package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/vilebile17/zimmer/internal/database"
)

func (cfg *apiConfig) createClassHandler(response http.ResponseWriter, request *http.Request) {
	type Params struct {
		Name string `json:"name"`
	}

	var params Params
	decoder := json.NewDecoder(request.Body)
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(response, request, "There was an error decoding the parameters (required {name:string})", err, http.StatusBadRequest)
		return
	}

	if params.Name == "" {
		respondWithError(response, request, "The name parameter can't be empty", nil, http.StatusBadRequest)
		return
	}

	userID, err := cfg.getUserIDFromHeader(request.Header)
	if err != nil {
		respondWithError(response, request, "couldn't authorize user", err, http.StatusUnauthorized)
		return
	}

	dbClass, err := cfg.dbQueries.CreateClass(request.Context(), database.CreateClassParams{
		Name:      params.Name,
		TeacherID: userID,
	})
	if err != nil {
		respondWithError(response, request, "There was an error adding a class to the database", err, http.StatusBadRequest)
		return
	}

	fmt.Printf("Just made a class called '%v' with the teacher ID '%v'\n", dbClass.Name, dbClass.TeacherID)
	respondWithJSON(response, request, dbClass, http.StatusCreated)
}

func (cfg *apiConfig) joinClassHandler(response http.ResponseWriter, request *http.Request) {
	classID, err := uuid.Parse(request.PathValue("classID"))
	if err != nil {
		respondWithError(response, request, "There was an error getting and parsing the class ID from the URL", err, http.StatusUnauthorized)
		return
	}

	userID, err := cfg.getUserIDFromHeader(request.Header)
	if err != nil {
		respondWithError(response, request, "couldn't authorize user", err, http.StatusUnauthorized)
		return
	}

	students_classes, err := cfg.dbQueries.JoinClass(request.Context(), database.JoinClassParams{
		StudentID: userID,
		ClassID:   classID,
	})
	if err != nil {
		respondWithError(response, request, "There was an error finding that class...", err, http.StatusBadRequest)
		return
	}

	fmt.Printf("User %v just joined a class %v\n", students_classes.StudentID, classID)
	respondWithJSON(response, request, students_classes, http.StatusCreated)
}

func (cfg *apiConfig) getClassesForUserHandler(response http.ResponseWriter, request *http.Request) {
	userID, err := cfg.getUserIDFromHeader(request.Header)
	if err != nil {
		respondWithError(response, request, "couldn't authorize user", err, http.StatusUnauthorized)
		return
	}

	classesAsStudent, err := cfg.dbQueries.GetClassesAsStudent(request.Context(), userID)
	if err != nil {
		respondWithError(response, request, "There was an error finding classes where that user is a student", err, http.StatusBadRequest)
		return
	}

	classesAsTeacher, err := cfg.dbQueries.GetClassesAsTeacher(request.Context(), userID)
	if err != nil {
		respondWithError(response, request, "There was an error finding classes where that user is a student", err, http.StatusBadRequest)
		return
	}

	fmt.Printf("User %v just got their classes\n", userID)
	if len(classesAsStudent) == 0 && len(classesAsTeacher) == 0 {
		type NoClasses struct {
			Message string `json:"message"`
		}
		respondWithJSON(response, request, NoClasses{Message: "Oof, it doesn't look like you're signed up for any classes yet"}, http.StatusOK)
		return
	}

	type Classes struct {
		ClassesAsStudent []database.Class `json:"classes_as_student"`
		ClassesAsTeacher []database.Class `json:"classes_as_teacher"`
	}

	respondWithJSON(response, request, Classes{
		ClassesAsStudent: classesAsStudent,
		ClassesAsTeacher: classesAsTeacher,
	}, http.StatusOK)
}

func (cfg *apiConfig) getUsersForClassHandler(response http.ResponseWriter, request *http.Request) {
	classID, err := uuid.Parse(request.PathValue("classID"))
	if err != nil {
		respondWithError(response, request, "There was an error parsing that classID", err, http.StatusBadRequest)
		return
	}

	class, err := cfg.dbQueries.GetClassFromClassID(request.Context(), classID)
	if err != nil {
		respondWithError(response, request, "There was an error retreiving that class", err, http.StatusBadRequest)
		return
	}
	teacherFull, err := cfg.dbQueries.GetUserFromID(request.Context(), class.TeacherID)
	if err != nil {
		respondWithError(response, request, "There was an error retreiving the class teacher", err, http.StatusBadRequest)
		return
	}
	teacher := database.GetStudentsForClassRow{
		ID:   teacherFull.ID,
		Name: teacherFull.Name,
	}

	students, err := cfg.dbQueries.GetStudentsForClass(request.Context(), classID)
	if err != nil {
		respondWithError(response, request, "There was an error finding users for that class", err, http.StatusBadRequest)
		return
	}

	userID, err := cfg.getUserIDFromHeader(request.Header)
	if err != nil {
		respondWithError(response, request, "couldn't get userID from the auth header", err, http.StatusUnauthorized)
		return
	}
	if userID != teacher.ID {
		isStudent := false
		for _, student := range students {
			if student.ID == userID {
				isStudent = true
			}
		}
		if !isStudent {
			respondWithError(response, request, "you can only get the users for a class that you are in", err, http.StatusUnauthorized)
			return
		}
	}

	type responseJSON struct {
		Teacher  database.GetStudentsForClassRow   `json:"teacher"`
		Students []database.GetStudentsForClassRow `json:"students"`
	}
	respondWithJSON(response, request, responseJSON{
		Teacher:  teacher,
		Students: students,
	}, http.StatusOK)
	fmt.Printf("Just got all of the users in class %v", classID)
}
