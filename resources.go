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

const ResourceType = "resource"
const AnnouncementType = "announcement"

type ClassContent struct {
	ID          uuid.UUID `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	ContentType string    `json:"content_type"`
	ClassID     uuid.UUID `json:"class_id"`
	Title       string    `json:"title"`
	Content     string    `json:"content"`
}

func (cfg *apiConfig) createResourceHandler(response http.ResponseWriter, request *http.Request) {
	cfg.createClassContentHandler(response, request, ResourceType)
}

func (cfg *apiConfig) createAnnouncementHandler(response http.ResponseWriter, request *http.Request) {
	cfg.createClassContentHandler(response, request, AnnouncementType)
}

func (cfg *apiConfig) createClassContentHandler(response http.ResponseWriter, request *http.Request, contentType string) {
	classID, err := uuid.Parse(request.PathValue("classID"))
	if err != nil {
		respondWithError(response, request, "There was an error parsing that classID", err, http.StatusBadRequest)
		return
	}

	statusCode, err := cfg.teacherActions(request)
	if err != nil {
		respondWithError(response, request, err.Error(), err, statusCode)
		return
	}

	type Params struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}

	var params Params
	decoder := json.NewDecoder(request.Body)
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(response, request, "There was an error decoding the parameters {title:string}, (Optional: {content:string})", err, http.StatusBadRequest)
		return
	}

	if params.Title == "" {
		respondWithError(response, request, "'title' parameter cannot be empty", err, http.StatusBadRequest)
		return
	}

	cc, err := cfg.dbQueries.CreateClassContent(request.Context(), database.CreateClassContentParams{
		ClassID:     classID,
		Title:       params.Title,
		ContentType: contentType,
		Content: sql.NullString{
			String: params.Content,
			Valid:  params.Content != "",
		},
	})
	if err != nil {
		respondWithError(response, request, "there was an error adding the resource to the database", err, http.StatusBadRequest)
		return
	}

	respondWithJSON(response, request, ClassContent{
		ID:          cc.ID,
		CreatedAt:   cc.CreatedAt,
		UpdatedAt:   cc.UpdatedAt,
		ContentType: cc.ContentType,
		ClassID:     cc.ClassID,
		Title:       cc.Title,
		Content:     cc.Content.String,
	}, http.StatusCreated)
	fmt.Printf("a class content was created for class %v\n", classID)
}

func (cfg *apiConfig) getResourceHandler(response http.ResponseWriter, request *http.Request) {
	cfg.getClassContentHandler(response, request, ResourceType)
}
func (cfg *apiConfig) getAnnouncementHandler(response http.ResponseWriter, request *http.Request) {
	cfg.getClassContentHandler(response, request, AnnouncementType)
}

func (cfg *apiConfig) getClassContentHandler(response http.ResponseWriter, request *http.Request, contentType string) {
	classID, err := uuid.Parse(request.PathValue("classID"))
	if err != nil {
		respondWithError(response, request, "There was an error parsing that classID", err, http.StatusBadRequest)
		return
	}

	contentID, err := uuid.Parse(request.PathValue("contentID"))
	if err != nil {
		respondWithError(response, request, "There was an error parsing that resourceID", err, http.StatusBadRequest)
		return
	}

	userID, err := cfg.getUserID(request)
	if err != nil {
		respondWithError(response, request, "Couldn't get the userID from the cookies or headers", err, http.StatusUnauthorized)
		return
	}

	isInClass, err := cfg.isUserInThisClass(request.Context(), userID, classID)
	if err != nil {
		respondWithError(response, request, "Couldn't verify if user is in the class", err, http.StatusUnauthorized)
		return
	}
	if !isInClass {
		respondWithError(response, request, "You must be in the class to view the resource", err, http.StatusForbidden)
		return
	}

	cc, err := cfg.dbQueries.GetClassContent(request.Context(), contentID)
	if err != nil {
		respondWithError(response, request, "There was an error getting the resource from the database", err, http.StatusBadRequest)
		return
	}

	respondWithJSON(response, request, ClassContent{
		ID:        cc.ID,
		CreatedAt: cc.CreatedAt,
		UpdatedAt: cc.UpdatedAt,
		ClassID:   cc.ClassID,
		Title:     cc.Title,
		Content:   cc.Content.String,
	}, http.StatusOK)
}

func (cfg *apiConfig) getResourcesForClassHandler(response http.ResponseWriter, request *http.Request) {
	cfg.getClassContentForClassHandler(response, request, ResourceType)
}
func (cfg *apiConfig) getAnnouncementsForClassHandler(response http.ResponseWriter, request *http.Request) {
	cfg.getClassContentForClassHandler(response, request, AnnouncementType)
}

func (cfg *apiConfig) getClassContentForClassHandler(response http.ResponseWriter, request *http.Request, contentType string) {
	classID, err := uuid.Parse(request.PathValue("classID"))
	if err != nil {
		respondWithError(response, request, "There was an error parsing that classID", err, http.StatusBadRequest)
		return
	}

	userID, err := cfg.getUserID(request)
	if err != nil {
		respondWithError(response, request, "Couldn't get the userID from the cookies or headers", err, http.StatusUnauthorized)
		return
	}

	isInClass, err := cfg.isUserInThisClass(request.Context(), userID, classID)
	if err != nil {
		respondWithError(response, request, "Couldn't verify if user is in the class", err, http.StatusUnauthorized)
		return
	}
	if !isInClass {
		respondWithError(response, request, "You must be in the class to view the resource", err, http.StatusForbidden)
		return
	}

	resources, err := cfg.dbQueries.GetClassContentForClass(request.Context(), database.GetClassContentForClassParams{
		ContentType: contentType,
		ClassID:     classID,
	})
	if err != nil {
		respondWithError(response, request, "Couldn't get the content from the database", err, http.StatusBadRequest)
		return
	}

	type ResponsePayload struct {
		ID        uuid.UUID `json:"id"`
		Title     string    `json:"title"`
		CreatedAt time.Time `json:"created_at"`
	}
	var payload []ResponsePayload

	for _, resource := range resources {
		payload = append(payload, ResponsePayload{
			ID:        resource.ID,
			Title:     resource.Title,
			CreatedAt: resource.CreatedAt,
		})
	}
	respondWithJSON(response, request, payload, http.StatusOK)
	fmt.Printf("Just got some content from the class %v\n", classID)
}

func (cfg *apiConfig) updateClassContentHandler(response http.ResponseWriter, request *http.Request) {
	contentID, err := uuid.Parse(request.PathValue("contentID"))
	if err != nil {
		respondWithError(response, request, "There was an error parsing that contentID", err, http.StatusBadRequest)
		return
	}

	statusCode, err := cfg.teacherActions(request)
	if err != nil {
		respondWithError(response, request, err.Error(), err, statusCode)
		return
	}

	type Params struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}

	var params Params
	decoder := json.NewDecoder(request.Body)
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(response, request, "There was an error decoding the parameters (Optional: {title:string, content:string})", err, http.StatusBadRequest)
		return
	}

	oldClassContent, err := cfg.dbQueries.GetClassContent(request.Context(), contentID)
	if err != nil {
		respondWithError(response, request, "There was an error getting the old resource data", err, http.StatusBadRequest)
		return
	}

	var newTitle string
	var newContent string
	if params.Title == "" {
		newTitle = oldClassContent.Title
	} else {
		newTitle = params.Title
	}
	if params.Content == "" {
		newContent = oldClassContent.Content.String
	} else {
		newContent = params.Content
	}

	cc, err := cfg.dbQueries.UpdateClassContent(request.Context(), database.UpdateClassContentParams{
		ID:    contentID,
		Title: newTitle,
		Content: sql.NullString{
			String: newContent,
			Valid:  newContent != "",
		},
	})
	if err != nil {
		respondWithError(response, request, "There was an error updating the SQL table", err, http.StatusBadRequest)
		return
	}

	respondWithJSON(response, request, ClassContent{
		ID:        cc.ID,
		CreatedAt: cc.CreatedAt,
		UpdatedAt: cc.UpdatedAt,
		ClassID:   cc.ClassID,
		Title:     cc.Title,
		Content:   cc.Content.String,
	}, http.StatusCreated)
	fmt.Println("Just updated some class content of id ", contentID)
}
