package data

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type NewUser struct {
	ID           int64
	Username     string
	AccessToken  string
	RefreshToken string
}

type User struct {
	ID           int64  `db:"id"`
	Username     string `db:"username"`
	AccessToken  string `db:"access_token"`
	RefreshToken string `db:"refresh_token"`
}

type UserRepo interface {
	Upsert(ctx context.Context, user NewUser) (err error)
	GetByID(ctx context.Context, id int64) (user User, err error)
	DeleteByID(ctx context.Context, id int64) (err error)
}

type PostgresUserRepo struct {
	db *sqlx.DB
}

func (p PostgresUserRepo) Upsert(ctx context.Context, user NewUser) (err error) {
	stmt, err := p.db.PreparexContext(ctx, `
		INSERT INTO bee_schema.users (id, username, access_token, refresh_token)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (id) DO UPDATE
		SET username = EXCLUDED.username,
			access_token = EXCLUDED.access_token,
			refresh_token = EXCLUDED.refresh_token
	`)
	if err != nil {
		return fmt.Errorf("preparing query: %v", err)
	}

	_, err = stmt.ExecContext(ctx, user.ID, user.Username, user.AccessToken, user.RefreshToken)
	if err != nil {
		return fmt.Errorf("executing INSERT query: %v", err)
	}

	return nil
}

func (p PostgresUserRepo) GetByID(ctx context.Context, id int64) (user User, err error) {
	stmt, err := p.db.PreparexContext(ctx, `
		SELECT id, username, access_token, refresh_token
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

func (p PostgresUserRepo) DeleteByID(ctx context.Context, id int64) (err error) {
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
