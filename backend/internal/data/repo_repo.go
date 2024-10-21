package data

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/jmoiron/sqlx"
)

type Repo struct {
	ID     int64  `db:"id"`
	Name   string `db:"name"`
	UserID int64  `db:"user_id"`
}

func (r Repo) LogValue() slog.Value {
	return slog.GroupValue(
		slog.Int64("id", r.ID),
		slog.String("name", r.Name),
		slog.Int64("user_id", r.UserID),
	)
}

type RepoRepo interface {
	Upsert(ctx context.Context, repo []Repo) (err error)
	Delete(ctx context.Context, id []int64) (err error)
	Get(ctx context.Context, id int64) (repo *Repo, err error)

	// GetAll retrieves all repositories for a given user and whose names are substrings of searchRepo.
	//
	// If searchRepo is empty, all repositories are considered.
	GetAll(ctx context.Context, searchRepo string, userID int64) (repos []Repo, err error)
}

type PostgresRepoRepo struct {
	db *sqlx.DB
}

func (p PostgresRepoRepo) Upsert(ctx context.Context, repos []Repo) (err error) {
	_, err = p.db.NamedExecContext(
		ctx,
		`INSERT INTO bee_schema.repos (id, name, user_id)
		VALUES (:id, :name, :user_id)
		ON CONFLICT (id) DO UPDATE SET
		name = EXCLUDED.name,
		user_id = EXCLUDED.user_id
		`,
		repos,
	)
	if err != nil {
		return fmt.Errorf("executing INSERT query: %v", err)
	}

	return nil
}

func (p PostgresRepoRepo) Delete(ctx context.Context, ids []int64) (err error) {
	query, args, err := sqlx.In(`
		DELETE FROM bee_schema.repos
       	WHERE id IN (?)
	`, ids)
	if err != nil {
		return fmt.Errorf("preparing query with IN clause: %v", err)
	}

	query = p.db.Rebind(query)

	_, err = p.db.ExecContext(ctx, query, args...)
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

var _ RepoRepo = &PostgresRepoRepo{}

func NewPostgresRepoRepo(db *sqlx.DB) *PostgresRepoRepo {
	return &PostgresRepoRepo{db: db}
}
