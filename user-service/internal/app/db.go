package app

import (
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
	"user-service/internal/config"
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

	log.Infof("App connected with params: %s", dataSource)

	_, err = pool.Exec("set search_path to user_service")
	if err != nil {
		return nil, fmt.Errorf("failed to set search path: %w", err)
	}

	pool.SetMaxOpenConns(cfg.MaxOpenConns)
	pool.SetMaxIdleConns(cfg.MaxIdleConns)
	pool.SetConnMaxIdleTime(cfg.ConnMaxIdleTime)
	pool.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	if err = pool.Ping(); err != nil {
		return nil, fmt.Errorf("failed to connect to the database: %w", err)
	}

	m, err := migrate.New(
		"file:///migrations",
		fmt.Sprintf(
			"postgres://%s:%s@%s:%s/%s?sslmode=disable",
			cfg.Username,
			cfg.Password,
			cfg.Host,
			cfg.Port,
			cfg.Name,
		))
	if err != nil {
		log.Warnf("couldn't find url: %s", err)
	}
	if err = m.Up(); err != nil {
		log.Warnf("couldn't up: %s", err)
	}

	return pool, nil
}
