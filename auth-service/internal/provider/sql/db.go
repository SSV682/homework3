package sql

import (
	"auth-service/internal/domain/models"
	"auth-service/internal/provider"
	"context"
	"errors"
	"fmt"
	"github.com/elgris/sqrl"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
)

const (
	buildQuery   = "build query: %v"
	executeQuery = "execute query: %v"
)

const (
	idColumn        = "id"
	usernameColumn  = "username"
	firstNameColumn = "firstname"
	lastNameColumn  = "lastname"
	emailColumn     = "email"
	phoneColumn     = "phone"
	passwordColumn  = "password"
)

const (
	defaultSchema = "user_service"
	usersTable    = defaultSchema + "." + "users"
)

type sqlUserProvider struct {
	pool *sqlx.DB
}

var _ provider.UserProvider = &sqlUserProvider{}

func NewSQLBusinessRulesProvider(pool *sqlx.DB) *sqlUserProvider {
	return &sqlUserProvider{
		pool: pool,
	}
}

var (
	queryBuilder = sqrl.NewSelectBuilder(sqrl.StatementBuilder).PlaceholderFormat(sqrl.Dollar)
)

func (s *sqlUserProvider) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	q := queryBuilder.
		Select(idColumn,
			usernameColumn,
			firstNameColumn,
			lastNameColumn,
			emailColumn,
			phoneColumn,
			passwordColumn).
		From(usersTable).
		Where(sqrl.Eq{usernameColumn: username})

	query, args, err := q.ToSql()
	if err != nil {
		return nil, fmt.Errorf(buildQuery, err)
	}

	var rows []models.User

	if err = s.pool.SelectContext(ctx, &rows, query, args...); err != nil {
		return nil, fmt.Errorf(executeQuery, err)
	}

	switch len(rows) {
	case 0:
		return nil, errors.New("not found username")
	case 1:
		return &rows[0], nil
	default:
		return nil, errors.New("username isn't unique")
	}
}

func (s *sqlUserProvider) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	q := queryBuilder.
		Select(idColumn,
			usernameColumn,
			firstNameColumn,
			lastNameColumn,
			emailColumn,
			phoneColumn,
			passwordColumn).
		From(usersTable).
		Where(sqrl.Eq{idColumn: id})

	query, args, err := q.ToSql()
	if err != nil {
		return nil, fmt.Errorf(buildQuery, err)
	}

	var rows []models.User

	if err = s.pool.SelectContext(ctx, &rows, query, args...); err != nil {
		return nil, fmt.Errorf(executeQuery, err)
	}

	switch len(rows) {
	case 0:
		return nil, errors.New("not found username")
	case 1:
		return &rows[0], nil
	default:
		return nil, errors.New("username isn't unique")
	}
}
