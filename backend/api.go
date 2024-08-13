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
			msg := "invalid user ID"
			slog.Debug(msg, slog.Any("error", err))
			http.Error(w, msg, http.StatusBadRequest)
			return
		}

		// var bartekpaciaID int64 = 40357511
		// userID := bartekpaciaID
		var result []data.FatBuild

		if r.URL.Query().Get("repo_id") == "" {
			result, err = a.BuildRepo.GetAll(r.Context(), userID)
			if err != nil {
				msg := "failed to get all builds"
				slog.Debug(msg, slog.Any("error", err))
				http.Error(w, msg, http.StatusInternalServerError)
				return
			}
		} else {
			repoID, err := strconv.ParseInt(r.URL.Query().Get("repo_id"), 10, 64)
			if err != nil {
				msg := "invalid repo ID"
				slog.Debug(msg, slog.Any("error", err))
				http.Error(w, msg, http.StatusBadRequest)
				return
			}

			result, err = a.BuildRepo.GetAllByRepoID(r.Context(), userID, repoID)
			if err != nil {
				msg := "failed to get builds by repo id"
				slog.Debug(msg, slog.Any("error", err))
				http.Error(w, msg, http.StatusInternalServerError)
				return
			}
		}

		responseBodyBytes, err := json.Marshal(result)
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
