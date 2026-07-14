package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func respondJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}

func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, map[string]string{
		"error": message,
	})
}

func (app *App) GetStagesHandler(w http.ResponseWriter, r *http.Request) {
	stages, err := app.DB.GetStages()
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to get stages")
		return
	}
	respondJSON(w, http.StatusOK, stages)
}

func (app *App) GetStageHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))

	if err != nil || id <= 0 {
		respondError(w, http.StatusBadRequest, "invalid stage id")
		return
	}

	stage, err := app.DB.GetStage(id)

	switch {
	case errors.Is(err, sql.ErrNoRows):
		respondError(w, http.StatusNotFound, "stage not found")
		return
	case err != nil:
		log.Printf("get stage %d: %v", id, err)
		respondError(w, http.StatusInternalServerError, "failed to get stage")
		return
	}

	respondJSON(w, http.StatusOK, stage)
}

func (app *App) CompleteStageHandler(w http.ResponseWriter, r *http.Request) {
	stageID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil || stageID <= 0 {
		respondError(w, http.StatusBadRequest, "invalid stage id")
		return
	}

	stage, err := app.DB.GetStage(stageID)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		respondError(w, http.StatusNotFound, "stage not found")
		return
	case err != nil:
		log.Printf("get stage %d: %v", stageID, err)
		respondError(w, http.StatusInternalServerError, "failed to get stage")
		return
	}

	userID := 1                //TODO: MAKE IT DYNAMIC, GET FROM JWT OR SESSION
	photoURL := "path/to/file" //TODO: MAKE IT DYNAMIC, GET FROM FILE UPLOAD
	pointsEarned := 100        //TODO: MAKE IT DYNAMIC, CALCULATE BASED ON STAGE AND USER PERFORMANCE

	completion, err := NewCompletion(userID, stage.ID, pointsEarned, photoURL)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to create completion")
		return
	}

	err = app.DB.CompleteStage(completion)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to complete stage")
		return
	}

	respondJSON(w, http.StatusOK, map[string]any{
		"message":      "stage completed",
		"stage":        stage,
		"pointsEarned": pointsEarned,
		"photoUrl":     photoURL,
	})
}

func (app *App) GetCompletedStagesHandler(w http.ResponseWriter, r *http.Request) {
	stages, err := app.DB.GetCompletedStages()
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to get completed stages")
		return
	}
	respondJSON(w, http.StatusOK, stages)
}

func (app *App) GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	users, err := app.DB.GetUsers()
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to get users")
		return
	}
	respondJSON(w, http.StatusOK, users)
}

func (app *App) GetUserHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil || id <= 0 {
		respondError(w, http.StatusBadRequest, "invalid id")
		return
	}

	user, err := app.DB.GetUser(id)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		respondError(w, http.StatusNotFound, "user not found")
		return
	case err != nil:
		log.Printf("get user %d: %v", id, err)
		respondError(w, http.StatusInternalServerError, "failed to get user")
		return
	}

	respondJSON(w, http.StatusOK, user)
}
