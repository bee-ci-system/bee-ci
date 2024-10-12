// Package api implements the HTTP REST API endpoints.
package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"math"
	"net/http"
	"strconv"
	"time"

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

	mux.HandleFunc("GET /repositories/{$}", a.getRepositories)
	mux.HandleFunc("GET /builds/", a.getBuilds)
	mux.HandleFunc("GET /builds/{build_id}/", a.getBuild)
	mux.HandleFunc("GET /builds/{build_id}/logs", a.getBuildLogs)

	// Actually used by frontend
	mux.HandleFunc("GET /user/", a.getUser)
	mux.HandleFunc("GET /dashboard/", a.getDashboard)
	mux.HandleFunc("GET /my-repositories/", a.getMyRepositories)
	mux.HandleFunc("GET /repositories/{id}/", a.getRepository)
	mux.HandleFunc("GET /pipeline/{id}/", a.getPipeline)

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
	pageSize, err := strconv.Atoi(r.URL.Query().Get("pageSize"))
	if err != nil {
		pageSize = 10
	}

	params := getMyRepositoriesParams{
		CurrentPage: currentPage,
		PageSize:    pageSize,
		Search:      r.URL.Query().Get("search"),
	}

	userID, ok := userid.FromContext(r.Context())
	if !ok {
		msg := "invalid user ID"
		logger.Debug(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	allRepos, err := a.RepoRepo.GetAll(r.Context(), params.Search, userID)
	if err != nil {
		msg := "failed to get my repositories"
		logger.Debug(msg, slog.Any("error", err))
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	totalPages := math.Ceil(float64(len(allRepos)) / float64(params.PageSize))

	// Perform paging
	// TODO: Paging should be done at database-level
	startIndex := (params.CurrentPage) * params.PageSize
	endIndex := startIndex + params.PageSize
	repos := make([]data.Repo, 0)
	for i := startIndex; i < len(allRepos) && i < endIndex; i++ {
		repos = append(repos, allRepos[i])
	}

	response := getMyRepositoriesDTO{
		Repositories:      toRepositories(repos),
		TotalRepositories: len(allRepos),
		TotalPages:        int(totalPages),
		CurrentPage:       params.CurrentPage,
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

func (a *App) getRepository(w http.ResponseWriter, r *http.Request) {
	logger, _ := l.FromContext(r.Context())

	userID, ok := userid.FromContext(r.Context())
	if !ok {
		msg := "invalid user ID"
		logger.Debug(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	repoID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		msg := fmt.Sprintf("invalid repository ID: %s", r.PathValue("id"))
		logger.Debug(msg, slog.Any("error", err))
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	repo, err := a.RepoRepo.Get(r.Context(), repoID)
	if err != nil {
		msg := "failed to get repository"
		logger.Debug(msg, slog.Any("error", err))
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	builds, err := a.BuildRepo.GetAllByRepoID(r.Context(), userID, repoID)
	if err != nil {
		msg := "failed to get all builds for repository"
		logger.Debug(msg, slog.Any("error", err))
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	pipelines := make([]pipeline, 0)
	for _, build := range builds {
		pipeline := pipeline{
			ID:             strconv.FormatInt(build.ID, 10),
			RepositoryName: build.RepoName,
			RepositoryID:   strconv.FormatInt(build.RepoID, 10),
			CommitName:     build.CommitMsg,
			Status:         build.Status,
			StartDate:      build.CreatedAt,
			EndDate:        &build.UpdatedAt,
		}

		pipelines = append(pipelines, pipeline)
	}

	response := getRepositoryDTO{
		ID:               strconv.FormatInt(repo.ID, 10),
		Name:             repo.Name,
		Description:      "Description not available",
		URL:              "URL not available",
		DateOfLastUpdate: time.Time{},
		Pipelines:        pipelines,
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		msg := "failed to encode repositoryDTO into json"
		logger.Debug(msg, slog.Any("error", err))
		http.Error(w, msg, http.StatusInternalServerError)
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
			if errors.Is(err, data.ErrNotFound) {
				continue
			}
			msg := "failed to get latest build for repo"
			logger.Error(msg, slog.Any("error", err))
			http.Error(w, msg, http.StatusInternalServerError)
			return
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

func (a *App) getRepositories(w http.ResponseWriter, r *http.Request) {
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

func (a *App) getPipeline(w http.ResponseWriter, r *http.Request) {
	logger, _ := l.FromContext(r.Context())

	userID, ok := userid.FromContext(r.Context())
	if !ok {
		msg := "invalid user ID"
		logger.Debug(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	buildID, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		msg := fmt.Sprintf("invalid build ID: %s", r.PathValue("id"))
		logger.Debug(msg, slog.Any("error", err))
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	fatBuild, err := a.BuildRepo.Get(r.Context(), userID, buildID)
	if err != nil {
		if errors.Is(err, data.ErrNotFound) {
			msg := fmt.Sprintf("build with id %d not found", buildID)
			http.Error(w, msg, http.StatusNotFound)
			return
		}

		msg := fmt.Sprintf("failed to get build with id %d from repo", buildID)
		logger.Debug(msg, slog.Any("error", err))
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	ppln := pipeline{
		ID:             strconv.FormatInt(fatBuild.ID, 10),
		RepositoryName: fatBuild.RepoName,
		RepositoryID:   strconv.FormatInt(fatBuild.RepoID, 10),
		CommitName:     fatBuild.CommitMsg,
		Status:         fatBuild.Status,
		StartDate:      fatBuild.CreatedAt,
		EndDate:        &fatBuild.UpdatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(ppln)
	if err != nil {
		msg := "failed to encode build into json"
		logger.Error(msg, slog.Any("error", err))
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}
}
