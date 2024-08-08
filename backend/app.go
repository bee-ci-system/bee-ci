package main

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"net/http"
)

type App struct {
	BuildRepo BuildRepo
}

func NewApp(db *sqlx.DB) *App {
	return &App{
		BuildRepo: NewPostgresBuildService(db),
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
	RepoID    uint64
	CommitSHA string
}

type BuildRepo interface {
	Create(ctx context.Context, build Build) (id uint64, err error)
	GetAll(ctx context.Context, repoID uint64) (builds []Build, err error)
}

type PostgresBuildRepo struct {
	db *sqlx.DB
}

func (p PostgresBuildRepo) Create(ctx context.Context, build Build) (id uint64, err error) {
	stmt, err := p.db.Preparex(`
		INSERT INTO builds (repo_id, commit_sha, status)
		VALUES ($1, $2, 'queued')
		RETURNING id
	`)
	if err != nil {
		return 0, fmt.Errorf("preparing query: %v", err)
	}

	err = stmt.Get(ctx, &id, build.RepoID, build.CommitSHA)
	if err != nil {
		return 0, fmt.Errorf("executing query: %v", err)
	}

	return id, nil
}

// TODO: refactor to only get builds for a specific user
func (p PostgresBuildRepo) GetAll(ctx context.Context, repoID uint64) (builds []Build, err error) {
	builds = make([]Build, 0)
	err = p.db.SelectContext(ctx, builds, "SELECT * FROM builds WHERE repo_id = $1", repoID)
	if err != nil {
		return nil, err
	}

	return builds, nil
}

var _ BuildRepo = &PostgresBuildRepo{}

func NewPostgresBuildService(db *sqlx.DB) BuildRepo {
	return &PostgresBuildRepo{db: db}
}
