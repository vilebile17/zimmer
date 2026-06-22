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

type Resource struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	ClassID   uuid.UUID `json:"class_id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
}

func (cfg *apiConfig) createResourceHandler(response http.ResponseWriter, request *http.Request) {
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

	resource, err := cfg.dbQueries.CreateResource(request.Context(), database.CreateResourceParams{
		ClassID: classID,
		Title:   params.Title,
		Content: sql.NullString{
			String: params.Content,
			Valid:  params.Content != "",
		},
	})
	if err != nil {
		respondWithError(response, request, "there was an error adding the resource to the database", err, http.StatusBadRequest)
		return
	}

	respondWithJSON(response, request, Resource{
		ID:        resource.ID,
		CreatedAt: resource.CreatedAt,
		UpdatedAt: resource.UpdatedAt,
		ClassID:   resource.ClassID,
		Title:     resource.Title,
		Content:   resource.Content.String,
	}, http.StatusCreated)
	fmt.Printf("a resource was created for class %v", classID)
}

func (cfg *apiConfig) getResourceHandler(response http.ResponseWriter, request *http.Request) {
	classID, err := uuid.Parse(request.PathValue("classID"))
	if err != nil {
		respondWithError(response, request, "There was an error parsing that classID", err, http.StatusBadRequest)
		return
	}

	resourceID, err := uuid.Parse(request.PathValue("resourceID"))
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

	resource, err := cfg.dbQueries.GetResource(request.Context(), resourceID)
	if err != nil {
		respondWithError(response, request, "There was an error getting the resource from the database", err, http.StatusBadRequest)
		return
	}

	respondWithJSON(response, request, Resource{
		ID:        resource.ID,
		CreatedAt: resource.CreatedAt,
		UpdatedAt: resource.UpdatedAt,
		ClassID:   resource.ClassID,
		Title:     resource.Title,
		Content:   resource.Content.String,
	}, http.StatusOK)
}

func (cfg *apiConfig) getResourcesForClassHandler(response http.ResponseWriter, request *http.Request) {
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

	resources, err := cfg.dbQueries.GetResourcesForClass(request.Context(), classID)
	if err != nil {
		respondWithError(response, request, "Couldn't get the classes from the database", err, http.StatusBadRequest)
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
	fmt.Printf("Just got the resources for the class %v", classID)
}

func (cfg *apiConfig) updateResourceHandler(response http.ResponseWriter, request *http.Request) {
	resourceID, err := uuid.Parse(request.PathValue("resourceID"))
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
		respondWithError(response, request, "There was an error decoding the parameters (Optional: {title:string, content:string})", err, http.StatusBadRequest)
		return
	}

	oldResource, err := cfg.dbQueries.GetResource(request.Context(), resourceID)
	if err != nil {
		respondWithError(response, request, "There was an error getting the old resource data", err, http.StatusBadRequest)
		return
	}

	var newTitle string
	var newContent string
	if params.Title == "" {
		newTitle = oldResource.Title
	} else {
		newTitle = params.Title
	}
	if params.Content == "" {
		newContent = oldResource.Content.String
	} else {
		newContent = params.Content
	}

	resource, err := cfg.dbQueries.UpdateResource(request.Context(), database.UpdateResourceParams{
		ID:    resourceID,
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

	respondWithJSON(response, request, Resource{
		ID:        resource.ID,
		CreatedAt: resource.CreatedAt,
		UpdatedAt: resource.UpdatedAt,
		ClassID:   resource.ClassID,
		Title:     resource.Title,
		Content:   resource.Content.String,
	}, http.StatusCreated)
	fmt.Println("Just updated a resource of id ", resourceID)
}
