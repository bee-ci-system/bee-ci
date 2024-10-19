package data

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/jmoiron/sqlx"
)

type NewUser struct {
	ID       int64
	Username string
}

type User struct {
	ID       int64  `db:"id"`
	Username string `db:"username"`
}

func (u User) LogValue() slog.Value {
	return slog.GroupValue(
		slog.Int64("id", u.ID),
		slog.String("username", u.Username),
	)
}

type UserRepo interface {
	Upsert(ctx context.Context, user NewUser) (err error)
	Get(ctx context.Context, id int64) (user User, err error)
	Delete(ctx context.Context, id int64) (err error)
}

type PostgresUserRepo struct {
	db *sqlx.DB
}

func (p PostgresUserRepo) Upsert(ctx context.Context, user NewUser) (err error) {
	stmt, err := p.db.PreparexContext(ctx, `
		INSERT INTO bee_schema.users (id, username)
		VALUES ($1, $2)
		ON CONFLICT (id) DO UPDATE
		SET username = EXCLUDED.username
	`)
	if err != nil {
		return fmt.Errorf("preparing query: %v", err)
	}

	_, err = stmt.ExecContext(ctx, user.ID, user.Username)
	if err != nil {
		return fmt.Errorf("executing INSERT query: %v", err)
	}

	return nil
}

func (p PostgresUserRepo) Get(ctx context.Context, id int64) (user User, err error) {
	stmt, err := p.db.PreparexContext(ctx, `
		SELECT id, username
		FROM bee_schema.users
		WHERE id = $1
	`)
	if err != nil {
		return User{}, fmt.Errorf("preparing query: %v", err)
	}

	err = stmt.GetContext(ctx, &user, id)
	if err != nil {
		return User{}, fmt.Errorf("executing SELECT query: %v", err)
	}

	return user, nil
}

func (p PostgresUserRepo) Delete(ctx context.Context, id int64) (err error) {
	stmt, err := p.db.PreparexContext(ctx, `
		DELETE FROM bee_schema.users
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

var _ UserRepo = &PostgresUserRepo{}

func NewPostgresUserRepo(db *sqlx.DB) *PostgresUserRepo {
	return &PostgresUserRepo{db: db}
}
