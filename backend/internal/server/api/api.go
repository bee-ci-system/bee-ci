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

	mux.HandleFunc("GET /user/", a.getUser)
	mux.HandleFunc("GET /my-repositories/", a.getMyRepositories)
	mux.HandleFunc("GET /dashboard/", a.getDashboard)

	mux.HandleFunc("GET /repos/", a.getRepos)

	mux.HandleFunc("GET /builds/", a.getBuilds)

	mux.HandleFunc("GET /builds/{build_id}", a.getBuild)

	mux.HandleFunc("GET /builds/{build_id}/logs", a.getBuildLogs)

	authMux := middleware.WithJWT(mux, a.jwtSecret)
	return authMux
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

	response := getUserDTO{
		ID:   user.ID,
		Name: user.Username,
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		msg := "failed to encode user into json"
		logger.Error(msg, slog.Any("error", err))
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}
}

func (a *App) getMyRepositories(w http.ResponseWriter, r *http.Request) {
	logger, _ := l.FromContext(r.Context())

	currentPage, err := strconv.Atoi(r.URL.Query().Get("currentPage"))
	if err != nil {
		currentPage = 0
	}

	params := getMyRepositoriesParams{
		CurrentPage: currentPage,
		Search:      r.URL.Query().Get("search"),
	}

	userID, ok := userid.FromContext(r.Context())
	if !ok {
		msg := "invalid user ID"
		logger.Debug(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	repos, err := a.RepoRepo.GetAll(r.Context(), params.Search, userID)
	if err != nil {
		msg := "failed to get my repositories"
		logger.Debug(msg, slog.Any("error", err))
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	response := getMyRepositoriesDTO{
		Repositories:      toRepositories(repos),
		TotalRepositories: len(repos),
		TotalPages:        2137, // TODO: use correct values
		CurrentPage:       69,
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		msg := "failed to encode my repositories into json"
		logger.Debug(msg, slog.Any("error", err))
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}
}

func (a *App) getDashboard(w http.ResponseWriter, r *http.Request) {
	logger, _ := l.FromContext(r.Context())

	userID, ok := userid.FromContext(r.Context())

	if !ok {
		msg := "invalid user ID"
		logger.Debug(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	builds, err := a.BuildRepo.GetAllByUserID(r.Context(), userID)
	if err != nil {
		msg := "failed to get all builds"
		logger.Debug(msg, slog.Any("error", err))
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}
	successfulBuilds := 0
	for _, build := range builds {
		if build.Status == "success" {
			successfulBuilds++
		}
	}

	stats := statsDTO{
		TotalPipelines:        len(builds),
		SuccessfulPipelines:   successfulBuilds,
		UnsuccessfulPipelines: len(builds) - successfulBuilds,
	}

	repos, err := a.RepoRepo.GetAll(r.Context(), "", userID)
	if err != nil {
		msg := "failed to get my repositories"
		logger.Debug(msg, slog.Any("error", err))
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	pipelines := make([]pipelineDashboardData, 0)
	// Get the latest build for every repo
	for _, repo := range repos {
		build, err := a.BuildRepo.GetLatestByRepoID(r.Context(), userID, repo.ID)
		if err != nil {
			msg := "failed to get latest build for repo"
			logger.Error(msg, slog.Any("error", err))
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}
		if build == nil {
			continue
		}

		pipeline := pipelineDashboardData{
			ID:             strconv.FormatInt(build.ID, 10),
			RepositoryName: repo.Name,
			CommitName:     build.CommitMsg,
			Status:         build.Status,
		}
		pipelines = append(pipelines, pipeline)
	}

	response := getDashboardDataDTO{
		Stats:        stats,
		Repositories: toRepositories(repos),
		Pipelines:    pipelines,
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		msg := "failed to encode dashboard data into json"
		logger.Debug(msg, slog.Any("error", err))
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}
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

	repos, err := a.RepoRepo.GetAll(r.Context(), "", userID)
	if err != nil {
		msg := "failed to get repositories"
		logger.Debug(msg, slog.Any("error", err))
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(repos)
	if err != nil {
		msg := "failed to encode repositories into json"
		logger.Debug(msg, slog.Any("error", err))
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}
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

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(result)
	if err != nil {
		msg := "failed to encode builds into json"
		logger.Error(msg, slog.Any("error", err))
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}
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

	result, err := a.BuildRepo.Get(r.Context(), userID, buildID)
	if err != nil {
		msg := fmt.Sprintf("failed to get build with id %d from repo", buildID)
		logger.Debug(msg, slog.Any("error", err))
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		msg := "failed to encode build into json"
		logger.Error(msg, slog.Any("error", err))
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}
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
