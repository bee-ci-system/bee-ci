package main

import (
	"context"
	"fmt"
	"github.com/bartekpacia/ghapp/data"
	"github.com/jmoiron/sqlx"
	"net/http"
)

type App struct {
	BuildRepo data.BuildRepo
}

func NewApp(db *sqlx.DB) *App {
	return &App{
		BuildRepo: data.NewPostgresBuildRepo(db),
	}
}

func (a *App) Mux() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /builds", func(w http.ResponseWriter, r *http.Request) {

	})

	mux.HandleFunc("POST /auth", func(w http.ResponseWriter, r *http.Request) {

	})

	return mux
})
