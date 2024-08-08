package data

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type NewBuildRequest struct {
	RepoID    uint64
	CommitSHA string
}

type BuildRepo interface {
	Create(ctx context.Context, build NewBuildRequest) (id uint64, err error)
	Update(ctx context.Context, id uint64, status BuildStatus) (err error)
	GetAll(ctx context.Context, repoID uint64) (builds []NewBuildRequest, err error)
}

type PostgresBuildRepo struct {
	db *sqlx.DB
}

func (p PostgresBuildRepo) Create(ctx context.Context, build NewBuildRequest) (id uint64, err error) {
	stmt, err := p.db.PreparexContext(ctx, `
		INSERT INTO builds (repo_id, commit_sha, status)
		VALUES ($1, $2, 'queued')
		RETURNING id
	`)
	if err != nil {
		return 0, fmt.Errorf("preparing query: %v", err)
	}

	err = stmt.Get(ctx, &id, build.RepoID, build.CommitSHA)
	if err != nil {
		return 0, fmt.Errorf("executing INSERT query: %v", err)
	}

	return id, nil
}

type BuildStatus string

const (
	StatusQueued  BuildStatus = "queued"
	StatusRunning BuildStatus = "running"
	StatusFailed  BuildStatus = "failed"
	StatusSuccess BuildStatus = "success"
)

func (p PostgresBuildRepo) Update(ctx context.Context, id uint64, status BuildStatus) (err error) {
	// UPDATE bee_schema.users SET username = 'dupa' WHERE id = 3;

	stmt, err := p.db.PreparexContext(ctx, `
		UPDATE builds
		SET status = $2
		WHERE id = $1
	`)
	if err != nil {
		return fmt.Errorf("preparing query: %v", err)
	}

	_, err = stmt.ExecContext(ctx, id, status)
	if err != nil {
		return fmt.Errorf("executing UPDATE query: %v", err)
	}

	return nil
}

// TODO: refactor to only get builds for a specific user
func (p PostgresBuildRepo) GetAll(ctx context.Context, repoID uint64) (builds []NewBuildRequest, err error) {
	builds = make([]NewBuildRequest, 0)
	err = p.db.SelectContext(ctx, builds, "SELECT * FROM builds WHERE repo_id = $1", repoID)
	if err != nil {
		return nil, err
	}

	return builds, nil
}

var _ BuildRepo = &PostgresBuildRepo{}

func NewPostgresBuildRepo(db *sqlx.DB) BuildRepo {
	return &PostgresBuildRepo{db: db}
}
