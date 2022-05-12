package postgres

import (
	"context"
	"errors"
	"github.com/golang-migrate/migrate/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/romm80/gophermart.git/internal/app/server"
)

func NewPostgresDB() (*pgxpool.Pool, error) {
	if err := migrateDB(); err != nil {
		return nil, err
	}

	pool, err := pgxpool.Connect(context.Background(), server.CFG.DB)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(context.Background()); err != nil {
		return nil, err
	}

	return pool, nil
}

func migrateDB() error {

	m, err := migrate.New(
		"file://db/migrations",
		server.CFG.DB)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}
	return nil
}
