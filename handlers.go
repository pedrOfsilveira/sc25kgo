package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"io"
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

type completeStageRequest struct {
	UserID   int    `json:"userId"`
	PhotoURL string `json:"photoUrl"`
}

func calculateCompletionPoints(stage Stage, isRepeat bool) int {
	const basePoints = 100

	points := basePoints
	points += stage.Week * 10 // 10 points per week
	points += stage.Day * 5   // 5 points per day

	if isRepeat {
		points /= 2
	}

	return points
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

	r.Body = http.MaxBytesReader(w, r.Body, 4*1024)

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	var req completeStageRequest
	if err := decoder.Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		respondError(w, http.StatusBadRequest, "request body must only contain a single JSON object")
		return
	}

	if req.UserID <= 0 {
		respondError(w, http.StatusBadRequest, "userId must be positive")
		return
	}

	if len(req.PhotoURL) > 2048 {
		respondError(w, http.StatusBadRequest, "photoUrl is too long")
		return
	}

	completion, err := app.DB.CompleteStage(
		r.Context(),
		req.UserID,
		stage,
		req.PhotoURL,
	)

	switch {
	case errors.Is(err, sql.ErrNoRows):
		respondError(w, http.StatusNotFound, "user not found")
		return
	case err != nil:
		log.Printf("complete stage %d for user %d: %v", stageID, req.UserID, err)
		respondError(w, http.StatusInternalServerError, "failed to complete stage")
		return
	}

	respondJSON(w, http.StatusCreated, map[string]any{
		"message":    "stage completed",
		"stage":      stage,
		"completion": completion,
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
