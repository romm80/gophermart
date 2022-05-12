package postgres

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v4"
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

func (a *AuthDB) CreateUser(user models.User) (int, error) {
	var userID int
	ctx := context.Background()
	conn, err := a.pool.Acquire(ctx)
	if err != nil {
		return userID, err
	}
	defer conn.Release()

	err = conn.QueryRow(ctx, `INSERT INTO users (login, password) VALUES($1, $2) ON CONFLICT DO NOTHING RETURNING id`, user.Login, user.Password).Scan(&userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return userID, app.ErrLoginIsUsed
		}
		return userID, err
	}
	return userID, nil
}

func (a *AuthDB) GetUserID(user models.User) (int, error) {
	var userID int
	ctx := context.Background()
	conn, err := a.pool.Acquire(ctx)
	if err != nil {
		return 0, err
	}
	defer conn.Release()

	err = conn.QueryRow(ctx, `SELECT id FROM users WHERE login = ($1) AND password = ($2)`, user.Login, user.Password).Scan(&userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, app.ErrInvalidLoginOrPassword
		}
		return 0, err
	}
	return userID, nil
}
