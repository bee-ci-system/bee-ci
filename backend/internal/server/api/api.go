package api

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	l "github.com/bee-ci/bee-ci-system/internal/common/logger"
	"github.com/bee-ci/bee-ci-system/internal/common/middleware"
	"github.com/bee-ci/bee-ci-system/internal/common/userid"
	"github.com/bee-ci/bee-ci-system/internal/data"
)

type App struct {
	BuildRepo data.BuildRepo
	LogsRepo  data.LogsRepo
	RepoRepo  data.RepoRepo
	UserRepo  data.UserRepo
	jwtSecret []byte
}

func NewApp(buildRepo data.BuildRepo, logsRepo data.LogsRepo, repoRepo data.RepoRepo, userRepo data.UserRepo, jwtSecret []byte) *App {
	return &App{
		BuildRepo: buildRepo,
		LogsRepo:  logsRepo,
		RepoRepo:  repoRepo,
		UserRepo:  userRepo,
		jwtSecret: jwtSecret,
	}
}

func (a *App) Mux() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /user", a.getUser)

	mux.HandleFunc("GET /repos/", a.getRepos)

	mux.HandleFunc("GET /builds/", a.getBuilds)

	mux.HandleFunc("GET /builds/{build_id}", a.getBuild)

	mux.HandleFunc("GET /builds/{build_id}/logs", a.getBuildLogs)

	authMux := middleware.WithJWT(mux, a.jwtSecret)
	return authMux
}

type GetUserDTO struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

func (a *App) getUser(w http.ResponseWriter, r *http.Request) {
	logger, _ := l.FromContext(r.Context())

	userID, ok := userid.FromContext(r.Context())
	if !ok {
		msg := "invalid user ID"
		logger.Debug(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	user, err := a.UserRepo.Get(r.Context(), userID)
	if err != nil {
		msg := fmt.Sprintf("failed to get user with id: %d", userID)
		logger.Debug(msg, slog.Any("error", err))
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	response := GetUserDTO{
		ID:   user.ID,
		Name: user.Username,
	}

	responseBodyBytes, err := json.Marshal(response)
	if err != nil {
		msg := "failed to marshal user ID into json"
		logger.Error(msg, slog.Any("error", err))
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(responseBodyBytes)
}

func (a *App) getRepos(w http.ResponseWriter, r *http.Request) {
	logger, _ := l.FromContext(r.Context())

	userID, ok := userid.FromContext(r.Context())
	if !ok {
		msg := "invalid user ID"
		logger.Debug(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	repos, err := a.RepoRepo.GetAll(r.Context(), userID)
	if err != nil {
		msg := "failed to get repositories"
		logger.Debug(msg, slog.Any("error", err))
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	decoder := json.NewEncoder(w)
	err = decoder.Encode(repos)
	if err != nil {
		msg := "failed to encode repositories into json"
		logger.Debug(msg, slog.Any("error", err))
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
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
		result, err = a.BuildRepo.GetAllByUserID(r.Context(), userID)
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

func (a *App) getBuild(w http.ResponseWriter, r *http.Request) {
	logger, _ := l.FromContext(r.Context())

	userID, ok := userid.FromContext(r.Context())
	if !ok {
		msg := "invalid user ID"
		logger.Debug(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	buildID, err := strconv.ParseInt(r.URL.Query().Get("build_id"), 10, 64)
	if err != nil {
		msg := "invalid build ID"
		logger.Debug(msg, slog.Any("error", err))
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	result, err := a.BuildRepo.Get(r.Context(), userID)
	if err != nil {
		msg := fmt.Sprintf("failed to get build with id %d from repo", buildID)
		logger.Debug(msg, slog.Any("error", err))
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	// TODO: Authorization (check if buildID belongs to userID)

	responseBodyBytes, err := json.Marshal(result)
	if err != nil {
		msg := "failed to marshal build into json"
		logger.Error(msg, slog.Any("error", err))
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(responseBodyBytes)
}

func (a *App) getBuildLogs(w http.ResponseWriter, r *http.Request) {
	logger, _ := l.FromContext(r.Context())

	_, ok := userid.FromContext(r.Context())
	if !ok {
		msg := "invalid user ID"
		logger.Debug(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	buildID, err := strconv.ParseInt(r.URL.Query().Get("build_id"), 10, 64)
	if err != nil {
		msg := "invalid build ID"
		logger.Debug(msg, slog.Any("error", err))
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	// TODO: Authorization (check if buildID belongs to userID)

	logs, err := a.LogsRepo.Get(r.Context(), buildID)
	if err != nil {
		msg := fmt.Sprintf("failed to get logs for build with id %d", buildID)
		logger.Debug(msg, slog.Any("error", err))
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	for _, logLine := range logs {
		_, _ = w.Write([]byte(logLine))
	}
}
