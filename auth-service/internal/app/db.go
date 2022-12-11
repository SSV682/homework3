package app

import (
	"fmt"

	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
	"user-service/internal/config"
)

const (
	driverName   = "pgx"
	databaseName = "postgres"
)

func initDBPool(cfg config.SQLConfig) (*sqlx.DB, error) {
	dataSource := fmt.Sprintf(
		"user=%s password=%s host=%s port=%s database=%s",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		databaseName,
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
