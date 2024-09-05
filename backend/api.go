package main

import (
	"encoding/json"
	"github.com/bartekpacia/ghapp/internal/userid"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/bartekpacia/ghapp/data"
	l "github.com/bartekpacia/ghapp/internal/logger"
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

	mux.HandleFunc("GET /builds/", a.getBuilds)

	mux.HandleFunc("GET /builds/{buildID}", func(w http.ResponseWriter, r *http.Request) {
	})

	mux.HandleFunc("POST /auth/", func(w http.ResponseWriter, r *http.Request) {
	})

	authMux := WithJWT(mux)
	return authMux
}

func (a *App) getBuilds(w http.ResponseWriter, r *http.Request) {
	logger, _ := l.FromContext(r.Context())

	userID, ok := userid.FromContext(r.Context())
	if !ok {
		msg := "invalid user ID"
		logger.Debug(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	var result []data.FatBuild

	if r.URL.Query().Get("repo_id") == "" {
		var err error
		result, err = a.BuildRepo.GetAll(r.Context(), userID)
		if err != nil {
			msg := "failed to get all builds"
			logger.Debug(msg, slog.Any("error", err))
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}
	} else {
		repoID, err := strconv.ParseInt(r.URL.Query().Get("repo_id"), 10, 64)
		if err != nil {
			msg := "invalid repo ID"
			logger.Debug(msg, slog.Any("error", err))
			http.Error(w, msg, http.StatusBadRequest)
			return
		}

		result, err = a.BuildRepo.GetAllByRepoID(r.Context(), userID, repoID)
		if err != nil {
			msg := "failed to get builds by repo id"
			logger.Debug(msg, slog.Any("error", err))
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}
	}

	responseBodyBytes, err := json.Marshal(result)
	if err != nil {
		msg := "failed to marshal builds"
		logger.Error(msg, slog.Any("error", err))
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(responseBodyBytes)
}
