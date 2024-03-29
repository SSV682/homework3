package app

import (
	"fmt"

	"auth-service/internal/config"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
)

const (
	driverName = "pgx"
)

func initDBPool(cfg config.SQLConfig) (*sqlx.DB, error) {
	dataSource := fmt.Sprintf(
		"user=%s password=%s host=%s port=%s database=%s sslmode=disable",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Name,
	)

	pool, err := sqlx.Open(driverName, dataSource)
	if err != nil {
		return nil, fmt.Errorf("failed to create pool of connections: %w", err)
	}

	pool.SetMaxOpenConns(cfg.MaxOpenConns)
	pool.SetMaxIdleConns(cfg.MaxIdleConns)
	pool.SetConnMaxIdleTime(cfg.ConnMaxIdleTime)
	pool.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	if err = pool.Ping(); err != nil {
		return nil, fmt.Errorf("failed to connect to the database: %w", err)
	}

	return pool, nil
}
