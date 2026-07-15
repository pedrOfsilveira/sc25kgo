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

func (app *App) activeUser(w http.ResponseWriter) (User, bool) {
	user, err := app.DB.GetActiveUser()
	switch {
	case errors.Is(err, sql.ErrNoRows):
		respondError(w, http.StatusNotFound, "active user not found")
		return User{}, false
	case err != nil:
		log.Printf("get active user: %v", err)
		respondError(w, http.StatusInternalServerError, "failed to get active user")
		return User{}, false
	}

	return user, true
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

	if len(req.PhotoURL) > 2048 {
		respondError(w, http.StatusBadRequest, "photoUrl is too long")
		return
	}

	user, ok := app.activeUser(w)
	if !ok {
		return
	}

	completion, err := app.DB.CompleteStage(
		r.Context(),
		user.ID,
		stage,
		req.PhotoURL,
	)

	switch {
	case errors.Is(err, sql.ErrNoRows):
		respondError(w, http.StatusNotFound, "user not found")
		return
	case err != nil:
		log.Printf("complete stage %d for user %d: %v", stageID, user.ID, err)
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
	user, ok := app.activeUser(w)
	if !ok {
		return
	}

	stages, err := app.DB.GetCompletedStages(user.ID)
	if err != nil {
		log.Printf("get completed stages for user %d: %v", user.ID, err)
		respondError(w, http.StatusInternalServerError, "failed to get completed stages")
		return
	}

	respondJSON(w, http.StatusOK, stages)
}

func (app *App) GetRunHistoryHandler(w http.ResponseWriter, r *http.Request) {
	user, ok := app.activeUser(w)
	if !ok {
		return
	}

	runs, err := app.DB.GetRunHistory(user.ID)
	if err != nil {
		log.Printf("get run history for user %d: %v", user.ID, err)
		respondError(w, http.StatusInternalServerError, "failed to get run history")
		return
	}

	respondJSON(w, http.StatusOK, runs)
}

func (app *App) GetProgressHandler(w http.ResponseWriter, r *http.Request) {
	user, ok := app.activeUser(w)
	if !ok {
		return
	}

	completedStageCount, err := app.DB.GetCompletedStageCount(user.ID)
	if err != nil {
		log.Printf("get completed stage count for user %d: %v", user.ID, err)
		respondError(w, http.StatusInternalServerError, "failed to get progress")
		return
	}

	totalStageCount, err := app.DB.GetTotalStageCount()
	if err != nil {
		log.Printf("get total stage count: %v", err)
		respondError(w, http.StatusInternalServerError, "failed to get progress")
		return
	}

	var nextStage *StageSummary
	stage, err := app.DB.GetNextStage(user.ID)
	switch {
	case errors.Is(err, sql.ErrNoRows):
	case err != nil:
		log.Printf("get next stage for user %d: %v", user.ID, err)
		respondError(w, http.StatusInternalServerError, "failed to get progress")
		return
	default:
		nextStage = &stage
	}

	respondJSON(w, http.StatusOK, ProgressResponse{
		User:                user,
		CompletedStageCount: completedStageCount,
		TotalStageCount:     totalStageCount,
		NextStage:           nextStage,
	})
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
