package data

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

type NewBuild struct {
	RepoID    uint64
	CommitSHA string
}

// Build represents a build in the database.
//
// The JSON struct tags are only to be used when receiving a row from LISTEN/NOTIFY.
type Build struct {
	ID        uint64      `db:"id" json:"id"`
	RepoID    uint64      `db:"repo_id" json:"repo_id"`
	CommitSHA string      `db:"commit_sha" json:"commit_sha"`
	Status    BuildStatus `db:"status" json:"status"`
	CreatedAt time.Time   `db:"created_at" json:"created_at"`
	UpdatedAt time.Time   `db:"updated_at" json:"updated_at"`
}

type BuildRepo interface {
	Create(ctx context.Context, build NewBuild) (id uint64, err error)
	Update(ctx context.Context, id uint64, status BuildStatus) (err error)
	GetAll(ctx context.Context, repoID uint64) (builds []NewBuild, err error)

	ListenUpdates(ctx context.Context, buildID uint64) (<-chan BuildStatus, error)
}

type PostgresBuildRepo struct {
	db *sqlx.DB
}

func (p PostgresBuildRepo) Create(ctx context.Context, build NewBuild) (id uint64, err error) {
	stmt, err := p.db.PreparexContext(ctx, `
		INSERT INTO bee_schema.builds (repo_id, commit_sha, status)
		VALUES ($1, $2, 'queued')
		RETURNING id
	`)
	if err != nil {
		return 0, fmt.Errorf("preparing query: %v", err)
	}

	err = stmt.GetContext(ctx, &id, build.RepoID, build.CommitSHA)
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
		UPDATE bee_schema.builds
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
func (p PostgresBuildRepo) GetAll(ctx context.Context, repoID uint64) (builds []NewBuild, err error) {
	builds = make([]NewBuild, 0)
	err = p.db.SelectContext(ctx, builds, "SELECT * FROM builds WHERE repo_id = $1", repoID)
	if err != nil {
		return nil, err
	}

	return builds, nil
}

func (p PostgresBuildRepo) ListenUpdates(ctx context.Context, buildID uint64) (<-chan BuildStatus, error) {
	updates := make(chan BuildStatus)
	go func() {
		defer close(updates)
		for {
			var status BuildStatus
			err := p.db.GetContext(ctx, &status, "SELECT status FROM builds WHERE id = $1", buildID)
			if err != nil {
				updates <- StatusFailed
				return
			}

			updates <- status
		}
	}()

	// TODO fix
	return nil, nil
}

var _ BuildRepo = &PostgresBuildRepo{}

func NewPostgresBuildRepo(db *sqlx.DB) *PostgresBuildRepo {
	return &PostgresBuildRepo{db: db}
}
