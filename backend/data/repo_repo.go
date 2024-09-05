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
	Create(ctx context.Context, repo []Repo) (err error)
	Delete(ctx context.Context, id []int64) (err error)
}

type PostgresRepoRepo struct {
	db *sqlx.DB
}

func (p PostgresRepoRepo) Create(ctx context.Context, repos []Repo) (err error) {
	_, err = p.db.NamedExecContext(
		ctx,
		`INSERT INTO bee_schema.repos (id, name, user_id)
		VALUES (:id, :name, :user_id)`,
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

var _ RepoRepo = &PostgresRepoRepo{}

func NewPostgresRepoRepo(db *sqlx.DB) *PostgresRepoRepo {
	return &PostgresRepoRepo{db: db}
}
