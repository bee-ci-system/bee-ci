package main

import (
	"context"
	"github.com/jmoiron/sqlx"
	"net/http"
)

type App struct {
	BuildService BuildService
}

func NewApp(db *sqlx.DB) *App {
	return &App{
		BuildService: NewPostgresBuildService(db),
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

type Build struct {
	RepoID uint64
	Commit string
}

type BuildService interface {
	Create(ctx context.Context, build Build) (id uint64, err error)
	GetAll(ctx context.Context, repoID uint64) (builds []Build, err error)
}

type PostgresBuildService struct {
	db *sqlx.DB
}

func (p PostgresBuildService) Create(ctx context.Context, build Build) (id uint64, err error) {
	//TODO implement me
	panic("implement me")
}

// TODO: refactor to only get builds for a specific user
func (p PostgresBuildService) GetAll(ctx context.Context, repoID uint64) ([]Build, error) {
	builds := make([]Build, 0)
	err := p.db.SelectContext(ctx, builds, "SELECT * FROM builds WHERE repo_id = $1", repoID)
	if err != nil {
		return nil, err
	}

	return builds, nil
}

var _ BuildService = &PostgresBuildService{}

func NewPostgresBuildService(db *sqlx.DB) BuildService {
	return &PostgresBuildService{db: db}
}
