package data

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/jmoiron/sqlx"
)

type Repo struct {
	ID                   int64     `db:"id"`
	Name                 string    `db:"name"`
	UserID               int64     `db:"user_id"`
	LatestCommitSHA      string    `db:"latest_commit_sha"`
	LatestCommitPushedAt time.Time `db:"latest_commit_pushed_at"`
	Description          string    `db:"description"`
}

func (r Repo) LogValue() slog.Value {
	return slog.GroupValue(
		slog.Int64("id", r.ID),
		slog.String("name", r.Name),
		slog.Int64("user_id", r.UserID),
	)
}

type RepoRepo interface {
	Create(ctx context.Context, repo []Repo) (err error)
	Delete(ctx context.Context, id []int64) (err error)
	Get(ctx context.Context, id int64) (repo *Repo, err error)

	// GetAll retrieves all repositories for a given user and whose names are substrings of searchRepo.
	//
	// If searchRepo is empty, all repositories are considered.
	GetAll(ctx context.Context, searchRepo string, userID int64) (repos []Repo, err error)

	UpdateLatestCommit(id int64, sha string, pushedAt time.Time) (err error)

	UpdateDescription(id int64, newDescription string) (err error)
}

type PostgresRepoRepo struct {
	db *sqlx.DB
}

func (p PostgresRepoRepo) Create(ctx context.Context, repos []Repo) (err error) {
	_, err = p.db.NamedExecContext(
		ctx,
		`INSERT INTO bee_schema.repos (id, name, user_id, latest_commit_sha, latest_commit_pushed_at)
		VALUES (:id, :name, :user_id, :latest_commit, :latest_commit_pushed_at)`,
		repos,
	)
	if err != nil {
		return fmt.Errorf("executing INSERT query: %v", err)
	}

	return nil
}

func (p PostgresRepoRepo) Delete(ctx context.Context, ids []int64) (err error) {
	idsInStruct := make([]interface{}, 0, len(ids))
	for _, i := range ids {
		idsInStruct = append(idsInStruct, struct {
			ID int64 `db:"id"`
		}{
			ID: i,
		})
	}

	stmt, err := p.db.PrepareNamedContext(ctx, `
		DELETE FROM bee_schema.repos
		WHERE id = :id
	`)
	if err != nil {
		return fmt.Errorf("preparing query: %v", err)
	}

	_, err = stmt.ExecContext(ctx, idsInStruct)
	if err != nil {
		return fmt.Errorf("executing DELETE query: %v", err)
	}

	return nil
}

func (p PostgresRepoRepo) Get(ctx context.Context, id int64) (repo *Repo, err error) {
	repo = &Repo{}
	err = p.db.GetContext(ctx, repo, `
		SELECT id, name, user_id
		FROM bee_schema.repos
		WHERE id = $1
	`, id)
	if err != nil {
		return nil, fmt.Errorf("selecting from repos: %v", err)
	}

	return repo, nil
}

func (p PostgresRepoRepo) GetAll(ctx context.Context, searchRepo string, userID int64) (repos []Repo, err error) {
	repos = make([]Repo, 0)

	query := `
		SELECT id, name, user_id
		FROM bee_schema.repos
		WHERE user_id = $1`
	args := []interface{}{userID}
	if searchRepo != "" {
		query += " AND name ILIKE $2"
		args = append(args, "%"+searchRepo+"%")
	}
	err = p.db.SelectContext(ctx, &repos, query, args...)
	if err != nil {
		return nil, fmt.Errorf("selecting from repos: %v", err)
	}
	return repos, nil
}

func (p PostgresRepoRepo) UpdateLatestCommit(id int64, sha string, pushedAt time.Time) (err error) {
	_, err = p.db.NamedExecContext(
		context.Background(),
		`UPDATE bee_schema.repos
		SET latest_commit_sha = :latest_commit, latest_commit_pushed_at = :latest_commit_pushed_at
		WHERE id = :id`,
		map[string]interface{}{
			"id":                      id,
			"latest_commit":           sha,
			"latest_commit_pushed_at": pushedAt,
		})
	if err != nil {
		return fmt.Errorf("executing UPDATE query: %v", err)
	}

	return nil
}

func (p PostgresRepoRepo) UpdateDescription(id int64, newDescription string) (err error) {
	_, err = p.db.NamedExecContext(
		context.Background(),
		`UPDATE bee_schema.repos
		SET description = :description
		WHERE id = :id`,
		map[string]interface{}{
			"id":          id,
			"description": newDescription,
		})
	if err != nil {
		return fmt.Errorf("executing UPDATE query: %v", err)
	}

	return nil
}

var _ RepoRepo = &PostgresRepoRepo{}

func NewPostgresRepoRepo(db *sqlx.DB) *PostgresRepoRepo {
	return &PostgresRepoRepo{db: db}
}
