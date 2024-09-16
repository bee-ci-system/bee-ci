package data

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	l "github.com/bee-ci/bee-ci-system/internal/logger"
	"github.com/jmoiron/sqlx"
)

type NewBuild struct {
	RepoID         int64
	CommitSHA      string
	CommitMsg      string
	InstallationID int64
}

// Build represents a row in the "builds" table.
//
// The JSON struct tags are only to be used when receiving a row from LISTEN/NOTIFY.
type Build struct {
	ID             int64     `db:"id" json:"id"`
	RepoID         int64     `db:"repo_id" json:"repo_id"`
	CommitSHA      string    `db:"commit_sha" json:"commit_sha"`
	CommitMsg      string    `db:"commit_message" json:"commit_message"`
	InstallationID int64     `db:"installation_id" json:"installation_id"`
	Status         string    `db:"status" json:"status"`
	Conclusion     *string   `db:"conclusion" json:"conclusion"`
	CreatedAt      time.Time `db:"created_at" json:"created_at"`
	UpdatedAt      time.Time `db:"updated_at" json:"updated_at"`
}

// FatBuild represents a row in the "builds" table, merged with information from other tables:
// - "repos" table, for repository information (repository name)
// - "users" table, for owner information (userID and user name)
type FatBuild struct {
	Build
	RepoName string `db:"repo_name" json:"repo_name"`
	UserID   int64  `db:"user_id" json:"user_id"`
	UserName string `db:"user_name" json:"user_name"`
}

type BuildRepo interface {
	Create(ctx context.Context, build NewBuild) (id int64, err error)
	UpdateStatus(ctx context.Context, buildID int64, status string) (err error)
	SetConclusion(ctx context.Context, buildID int64, conclusion string) (err error)

	// Get returns a build by its ID.
	Get(ctx context.Context, buildID int64) (build FatBuild, err error)

	// GetAllByUserID returns all builds for all repositories of userID.
	GetAllByUserID(ctx context.Context, userID int64) (builds []FatBuild, err error)

	// GetAllByRepoID returns all builds for the repository of repoID.
	GetAllByRepoID(ctx context.Context, userID, repoID int64) (builds []FatBuild, err error)
}

type PostgresBuildRepo struct {
	db *sqlx.DB
}

func (p PostgresBuildRepo) Create(ctx context.Context, build NewBuild) (id int64, err error) {
	stmt, err := p.db.PreparexContext(ctx, `
		INSERT INTO bee_schema.builds (repo_id, commit_sha, commit_message, installation_id, status)
		VALUES ($1, $2, $3, $4, 'queued')
		RETURNING id
	`)
	if err != nil {
		return 0, fmt.Errorf("preparing query: %v", err)
	}

	err = stmt.GetContext(ctx, &id, build.RepoID, build.CommitSHA, build.CommitMsg, build.InstallationID)
	if err != nil {
		return 0, fmt.Errorf("executing INSERT query: %v", err)
	}

	return id, nil
}

// UpdateStatus sets the status of a build. Available values are: "queued", "in_progress", "completed".
//
// See https://docs.github.com/en/rest/checks/runs?apiVersion=2022-11-28#create-a-check-run
func (p PostgresBuildRepo) UpdateStatus(ctx context.Context, buildID int64, status string) (err error) {
	stmt, err := p.db.PreparexContext(ctx, `
		UPDATE bee_schema.builds
		SET status = $2
		WHERE id = $1
	`)
	if err != nil {
		return fmt.Errorf("preparing query: %v", err)
	}

	_, err = stmt.ExecContext(ctx, buildID, status)
	if err != nil {
		return fmt.Errorf("executing UPDATE query: %v", err)
	}

	return nil
}

// SetConclusion sets the conclusion of a build. Available values are: "canceled", "failure", "success", "timed_out".
//
// See https://docs.github.com/en/rest/checks/runs?apiVersion=2022-11-28#create-a-check-run
func (p PostgresBuildRepo) SetConclusion(ctx context.Context, buildID int64, conclusion string) (err error) {
	stmt, err := p.db.PreparexContext(ctx, `
		UPDATE bee_schema.builds
		SET status = 'completed', conclusion = $2
		WHERE id = $1
	`)
	if err != nil {
		return fmt.Errorf("preparing query: %v", err)
	}

	_, err = stmt.ExecContext(ctx, buildID, conclusion)
	if err != nil {
		return fmt.Errorf("executing UPDATE query: %v", err)
	}

	return nil
}

// TODO: refactor to only get builds for a specific user

func (p PostgresBuildRepo) Get(ctx context.Context, buildID int64) (build FatBuild, err error) {
	logger, _ := l.FromContext(ctx)
	logger.Debug("BuildRepo.Get", slog.Any("buildID", buildID))

	err = p.db.GetContext(ctx, &build, `
		 		SELECT builds.*, repos.name AS repo_name, users.id AS user_id, users.username AS user_name
		 		FROM bee_schema.builds builds
		 		JOIN bee_schema.repos repos ON builds.repo_id = repos.id
		 		JOIN bee_schema.users users ON repos.user_id = users.id
		 		WHERE builds.id = $1
		 	`, buildID)
	if err != nil {
		return FatBuild{}, fmt.Errorf("executing SELECT query for buildID %d: %v", buildID, err)
	}

	return build, nil
}

func (p PostgresBuildRepo) GetAllByUserID(ctx context.Context, userID int64) (builds []FatBuild, err error) {
	logger, _ := l.FromContext(ctx)
	logger.Debug("BuildRepo.GetAllByUserID", slog.Any("userID", userID))

	builds = make([]FatBuild, 0)
	err = p.db.SelectContext(ctx, &builds, `
         		SELECT builds.*, repos.name AS repo_name, users.id AS user_id, users.username AS user_name
         		FROM bee_schema.builds builds
         		JOIN bee_schema.repos repos ON builds.repo_id = repos.id
         		JOIN bee_schema.users users ON repos.user_id = users.id
         		WHERE users.id = $1
		 	`, userID)
	if err != nil {
		return nil, fmt.Errorf("executing SELECT query for userID %d: %v", userID, err)
	}

	return builds, nil
}

func (p PostgresBuildRepo) GetAllByRepoID(ctx context.Context, userID, repoID int64) (builds []FatBuild, err error) {
	logger, _ := l.FromContext(ctx)
	logger.Debug("BuildRepo.GetAllByRepoID", slog.Any("userID", userID), slog.Any("repoID", repoID))

	builds = make([]FatBuild, 0)
	err = p.db.SelectContext(ctx, &builds, `
         		SELECT builds.*, repos.name AS repo_name, users.id AS user_id, users.username AS user_name
         		FROM bee_schema.builds builds
         		JOIN bee_schema.repos repos ON builds.repo_id = repos.id
         		JOIN bee_schema.users users ON repos.user_id = users.id
         		WHERE users.id = $1 AND repos.id = $2
		 	`, userID, repoID)
	if err != nil {
		return nil, fmt.Errorf("executing SELECT query for userID %d and repo ID %d: %v", userID, repoID, err)
	}

	return builds, nil
}

var _ BuildRepo = &PostgresBuildRepo{}

func NewPostgresBuildRepo(db *sqlx.DB) *PostgresBuildRepo {
	return &PostgresBuildRepo{db: db}
}
