package data

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type NewUser struct {
	ID           int64
	AccessToken  string
	RefreshToken string
}

type User struct {
	ID           int64  `db:"user_id"`
	AccessToken  string `db:"access_token"`
	RefreshToken string `db:"refresh_token"`
}

type UserRepo interface {
	Create(ctx context.Context, user NewUser) (err error)
	GetByID(ctx context.Context, id int64) (user User, err error)
}

type PostgresUserRepo struct {
	db *sqlx.DB
}

func (p PostgresUserRepo) Create(ctx context.Context, user NewUser) (err error) {
	stmt, err := p.db.PreparexContext(ctx, `
		INSERT INTO bee_schema.users (id, access_token, refresh_token)
		VALUES ($1, $2, $3)
	`)
	if err != nil {
		return fmt.Errorf("preparing query: %v", err)
	}

	_, err = stmt.ExecContext(ctx, user.ID, user.AccessToken, user.RefreshToken)
	if err != nil {
		return fmt.Errorf("executing INSERT query: %v", err)
	}

	return nil
}

func (p PostgresUserRepo) GetByID(ctx context.Context, id int64) (user User, err error) {
	stmt, err := p.db.PreparexContext(ctx, `
		SELECT id, access_token, refresh_token
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

var _ UserRepo = &PostgresUserRepo{}

func NewPostgresUserRepo(db *sqlx.DB) *PostgresUserRepo {
	return &PostgresUserRepo{db: db}
}
