package data

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type Repo struct {
	ID     int64  `db:"id"`
	Name   string `db:"name"`
	UserID int64  `db:"user_id"`
}

type RepoRepo interface {
	Create(ctx context.Context, repo Repo) (err error)
	Delete(ctx context.Context, id int64) (err error)
}

type PostgresRepoRepo struct {
	db *sqlx.DB
}

func (p PostgresRepoRepo) Create(ctx context.Context, repo Repo) (err error) {
	stmt, err := p.db.PreparexContext(ctx, `
		INSERT INTO bee_schema.repos (id, name, user_id)
		VALUES ($1, $2, $3)
	`)
	if err != nil {
		return fmt.Errorf("preparing query: %v", err)
	}

	_, err = stmt.ExecContext(ctx, repo.ID, repo.Name, repo.UserID)
	if err != nil {
		return fmt.Errorf("executing INSERT query: %v", err)
	}

	return nil
}

func (p PostgresRepoRepo) Delete(ctx context.Context, id int64) (err error) {
	stmt, err := p.db.PreparexContext(ctx, `
		DELETE FROM bee_schema.repos
		WHERE id = $1
	`)
	if err != nil {
		return fmt.Errorf("preparing query: %v", err)
	}

	_, err = stmt.ExecContext(ctx, id)
	if err != nil {
		return fmt.Errorf("executing DELETE query: %v", err)
	}

	return nil
}

var _ RepoRepo = &PostgresRepoRepo{}

func NewPostgresRepoRepo(db *sqlx.DB) *PostgresRepoRepo {
	return &PostgresRepoRepo{db: db}
}
