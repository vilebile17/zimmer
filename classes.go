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

	userID, err := cfg.getUserID(request)
	if err != nil {
		respondWithError(response, request, "couldn't authorize user", err, http.StatusUnauthorized)
		return
	}

	dbClass, err := cfg.dbQueries.CreateClass(request.Context(), database.CreateClassParams{
		Name:      params.Name,
		TeacherID: userID,
	})
	if err != nil {
		respondWithError(response, request, "couldn't add class to database, it might already exist", err, http.StatusBadRequest)
		return
	}

	fmt.Printf("Just made a class called '%v' with the teacher ID '%v'\n", dbClass.Name, dbClass.TeacherID)
	respondWithJSON(response, request, dbClass, http.StatusCreated)
}

func (cfg *apiConfig) joinClassHandler(response http.ResponseWriter, request *http.Request) {
	classID, err := uuid.Parse(request.PathValue("classID"))
	if err != nil {
		respondWithError(response, request, "Invalid classID", err, http.StatusBadRequest)
		return
	}

	userID, err := cfg.getUserID(request)
	if err != nil {
		respondWithError(response, request, "couldn't authorize user", err, http.StatusUnauthorized)
		return
	}

	class, err := cfg.dbQueries.GetClass(request.Context(), classID)
	if err != nil {
		respondWithError(response, request, "couldn't retrieve the class from the database", err, http.StatusBadRequest)
		return
	}

	if class.TeacherID == userID {
		respondWithError(response, request, "you're the teacher mate, why are ya joining your own class as a student", err, http.StatusUnauthorized)
		return
	}

	if !class.AllowJoining {
		respondWithError(response, request, "the class doesn't allow new users to join :(", err, http.StatusUnauthorized)
		return
	}

	students_classes, err := cfg.dbQueries.JoinClass(request.Context(), database.JoinClassParams{
		StudentID: userID,
		ClassID:   classID,
	})
	if err != nil {
		respondWithError(response, request, `couldn't join, perhaps you're already there ¯\_(ツ)_/¯`, err, http.StatusBadRequest)
		return
	}

	fmt.Printf("User %v just joined a class %v\n", students_classes.StudentID, classID)
	respondWithJSON(response, request, students_classes, http.StatusCreated)
}

func (cfg *apiConfig) getClassesForUserHandler(response http.ResponseWriter, request *http.Request) {
	userID, err := cfg.getUserID(request)
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
		ClassesAsStudent []database.Class `json:"classesAsStudent"`
		ClassesAsTeacher []database.Class `json:"classesAsTeacher"`
	}

	respondWithJSON(response, request, Classes{
		ClassesAsStudent: classesAsStudent,
		ClassesAsTeacher: classesAsTeacher,
	}, http.StatusOK)
}

func (cfg *apiConfig) getClassHandler(response http.ResponseWriter, request *http.Request) {
	userID, err := cfg.getUserIDFromHeader(request.Header)
	if err != nil {
		respondWithError(response, request, "couldn't authorize user", err, http.StatusUnauthorized)
		return
	}

	classID, err := uuid.Parse(request.PathValue("classID"))
	if err != nil {
		respondWithError(response, request, "There was an error getting and parsing the class ID from the URL", err, http.StatusBadRequest)
		return
	}

	b, err := cfg.isUserInThisClass(request.Context(), userID, classID)
	if err != nil {
		respondWithError(response, request, "couldn't validate that you're in this class", err, http.StatusUnauthorized)
		return
	}
	if !b {
		respondWithError(response, request, "you can't get a class' info if you aren't in it :)", err, http.StatusUnauthorized)
		return
	}

	class, err := cfg.dbQueries.GetClass(request.Context(), classID)
	if err != nil {
		respondWithError(response, request, "couldn't get the class data from the database", err, http.StatusBadRequest)
		return
	}
	respondWithJSON(response, request, class, http.StatusOK)
}

func (cfg *apiConfig) getUsersForClassHandler(response http.ResponseWriter, request *http.Request) {
	classID, err := uuid.Parse(request.PathValue("classID"))
	if err != nil {
		respondWithError(response, request, "There was an error parsing that classID", err, http.StatusBadRequest)
		return
	}

	userID, err := cfg.getUserID(request)
	if err != nil {
		respondWithError(response, request, "couldn't get userID from the auth header", err, http.StatusUnauthorized)
		return
	}

	if inClass, err := cfg.isUserInThisClass(request.Context(), userID, classID); err != nil {
		respondWithError(response, request, "couldn't get figure out if you are in the class", err, http.StatusUnauthorized)
		return
	} else if !inClass {
		respondWithError(response, request, "you can only view the users for a class that you are in", nil, http.StatusUnauthorized)
		return
	}

	class, err := cfg.dbQueries.GetClass(request.Context(), classID)
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
		ID:        teacherFull.ID,
		Name:      teacherFull.Name,
		CreatedAt: teacherFull.CreatedAt,
	}

	students, err := cfg.dbQueries.GetStudentsForClass(request.Context(), classID)
	if err != nil {
		respondWithError(response, request, "There was an error finding users for that class", err, http.StatusBadRequest)
		return
	}

	type responseJSON struct {
		Teacher  database.GetStudentsForClassRow   `json:"teacher"`
		Students []database.GetStudentsForClassRow `json:"students"`
	}
	respondWithJSON(response, request, responseJSON{
		Teacher:  teacher,
		Students: students,
	}, http.StatusOK)
	fmt.Printf("Just got all of the users in class %v\n", classID)
}

func (cfg *apiConfig) removeFromClass(response http.ResponseWriter, request *http.Request) {
	classID, err := uuid.Parse(request.PathValue("classID"))
	if err != nil {
		respondWithError(response, request, "There was an error parsing that classID", err, http.StatusBadRequest)
		return
	}
	userToRemove, err := uuid.Parse(request.PathValue("userID"))
	if err != nil {
		respondWithError(response, request, "There was an error parsing that userID", err, http.StatusBadRequest)
		return
	}

	userID, err := cfg.getUserIDFromHeader(request.Header)
	if err != nil {
		respondWithError(response, request, "couldn't get userID from the auth header", err, http.StatusUnauthorized)
		return
	}

	class, err := cfg.dbQueries.GetClass(request.Context(), classID)
	if err != nil {
		respondWithError(response, request, "there was an error retrieving the class data from the database", err, http.StatusBadRequest)
		return
	}

	if userID != userToRemove && userID != class.TeacherID {
		respondWithError(response, request, "You cannot remove another user from a class if you aren't the teacher.", err, http.StatusUnauthorized)
		return
	}

	if err = cfg.dbQueries.RemoveUserFromClass(request.Context(), database.RemoveUserFromClassParams{
		StudentID: userToRemove,
		ClassID:   classID,
	}); err != nil {
		respondWithError(response, request, "there was an error removing the user from the class", err, http.StatusBadRequest)
		return
	}

	fmt.Printf("Successfully removed user %v from the class %d\n", userToRemove, classID)
	response.WriteHeader(http.StatusOK)
}
