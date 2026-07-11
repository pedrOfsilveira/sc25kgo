package main

import (
	"encoding/json"
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
	idParam := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "invalid id")
		return
	}

	stage, err := app.DB.GetStage(id)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to get stages")
		return
	}

	respondJSON(w, http.StatusOK, stage)
}

func (app *App) CompleteStageHandler(w http.ResponseWriter, r *http.Request) {
	stageIDParam := chi.URLParam(r, "id")

	stageID, err := strconv.Atoi(stageIDParam)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid stage id")
		return
	}

	stage, err := app.DB.GetStage(stageID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to get stage")
		return
	}

	userID := 1
	photoURL := "path/to/file"
	pointsEarned := 100

	completion, err := NewCompletion(userID, stage.ID, pointsEarned, photoURL)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
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
		"photoURL":     photoURL,
	})
}
