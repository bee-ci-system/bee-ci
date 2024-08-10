package data

import (
	"context"
	"fmt"
	
	"github.com/jmoiron/sqlx"
)

type NewUser struct {
	ID           uint64
	AccessToken  string
	RefreshToken string
}

type UserRepo interface {
	Create(ctx context.Context, user NewUser) (err error)
}

type postgresUserRepo struct {
	db *sqlx.DB
}

func (p postgresUserRepo) Create(ctx context.Context, user NewUser) (err error) {
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

var _ UserRepo = &postgresUserRepo{}

func NewPostgresUserRepo(db *sqlx.DB) UserRepo {
	return &postgresUserRepo{db: db}
}
