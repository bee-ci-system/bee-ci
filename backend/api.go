package main

import (
	"net/http"

	"github.com/bartekpacia/ghapp/data"
)

type App struct {
	BuildRepo data.BuildRepo
}

func NewApp(buildRepo data.BuildRepo) *App {
	return &App{
		BuildRepo: buildRepo,
	}
}

func (a *App) Mux() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /builds", func(w http.ResponseWriter, r *http.Request) {
	})

	mux.HandleFunc("POST /auth", func(w http.ResponseWriter, r *http.Request) {
	})

	return mux
}
