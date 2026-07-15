package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (app *App) Routes() http.Handler {
	r := chi.NewRouter()

	r.Get("/api/stages", app.GetStagesHandler)
	r.Get("/api/stages/completed", app.GetCompletedStagesHandler)
	r.Get("/api/stages/{id}", app.GetStageHandler)

	r.Post("/api/stages/{id}/complete", app.CompleteStageHandler)
	r.Get("/api/runs", app.GetRunHistoryHandler)
	r.Get("/api/progress", app.GetProgressHandler)

	r.Get("/api/users", app.GetUsersHandler)
	r.Get("/api/users/{id}", app.GetUserHandler)

	return r
}
