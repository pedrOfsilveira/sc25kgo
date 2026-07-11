package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (app *App) Routes() http.Handler {
	r := chi.NewRouter()

	r.Get("/api/stages", app.GetStagesHandler)
	r.Get("/api/stages/{id}", app.GetStageHandler)

	return r
}
