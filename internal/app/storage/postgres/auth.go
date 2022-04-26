package postgres

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/romm80/gophermart.git/internal/app"
	"github.com/romm80/gophermart.git/internal/app/models"
)

type AuthDB struct {
	pool *pgxpool.Pool
}

func NewAuthDB(pool *pgxpool.Pool) *AuthDB {
	return &AuthDB{pool: pool}
}

func (a *AuthDB) CreateUser(user models.User) error {
	ctx := context.Background()
	conn, err := a.pool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	tx, err := conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	rows, err := tx.Query(ctx, `SELECT login FROM users WHERE login = ($1)`, user.Login)
	if err != nil {
		return err
	}
	if rows.Next() {
		return app.ErrLoginIsUsed
	}

	if _, err := tx.Exec(ctx, `INSERT INTO users (login, password) VALUES($1, $2)`, user.Login, user.Password); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (a *AuthDB) GetUser(user models.User) error {
	ctx := context.Background()
	conn, err := a.pool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	rows, err := conn.Query(ctx, `SELECT login FROM users WHERE login = ($1) AND password = ($2)`, user.Login, user.Password)
	if err != nil {
		return err
	}
	if !rows.Next() {
		return app.ErrInvalidLoginOrPassword
	}
	return nil
}
