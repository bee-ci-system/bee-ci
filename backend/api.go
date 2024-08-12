package main

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

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
	mux.HandleFunc("GET /builds/", func(w http.ResponseWriter, r *http.Request) {
		// Q: Where do we get userID from?
		// A: From the JWT that was issued by our backend

		// For now, let's assume that the userID is in header
		userID, err := strconv.ParseInt(r.Header.Get("X-User-ID"), 10, 64)
		if err != nil {
			http.Error(w, "invalid user ID", http.StatusBadRequest)
			return
		}

		// var bartekpaciaID int64 = 40357511
		// userID := bartekpaciaID

		builds, err := a.BuildRepo.GetAll(r.Context(), userID)
		if err != nil {
			msg := "failed to get builds"
			slog.Error(msg, slog.Any("error", err))
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}

		responseBodyBytes, err := json.Marshal(builds)
		if err != nil {
			msg := "failed to marshal builds"
			slog.Error(msg, slog.Any("error", err))
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(responseBodyBytes)
	})

	mux.HandleFunc("GET /builds/{buildID}", func(w http.ResponseWriter, r *http.Request) {
	})

	mux.HandleFunc("POST /auth/", func(w http.ResponseWriter, r *http.Request) {
	})

	return mux
}
