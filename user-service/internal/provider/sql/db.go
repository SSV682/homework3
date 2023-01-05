package sql

import (
	"context"
	"fmt"
	"github.com/elgris/sqrl"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
	"user-service/internal/domain/errors"
	"user-service/internal/domain/models"
)

const (
	buildQuery   = "build query: %v"
	executeQuery = "execute query: %v , %s, %v"
)

const (
	idColumn           = "id"
	usernameColumn     = "username"
	firstNameColumn    = "firstname"
	lastNameColumn     = "lastname"
	emailColumn        = "email"
	phoneColumn        = "phone"
	passwordColumn     = "password"
	allColumns         = idColumn + ", " + usernameColumn + ", " + firstNameColumn + ", " + lastNameColumn + ", " + emailColumn + ", " + phoneColumn + ", " + passwordColumn
	allExceptIDColumns = usernameColumn + ", " + firstNameColumn + ", " + lastNameColumn + ", " + emailColumn + ", " + phoneColumn + ", " + passwordColumn
)

const (
	defaultSchema = "user_service"
	usersTable    = defaultSchema + "." + "users"
)

var (
	queryBuilder       = sqrl.NewSelectBuilder(sqrl.StatementBuilder).PlaceholderFormat(sqrl.Dollar)
	queryInsertBuilder = sqrl.NewInsertBuilder(sqrl.StatementBuilder).PlaceholderFormat(sqrl.Dollar)
	queryUpdateBuilder = sqrl.NewUpdateBuilder(sqrl.StatementBuilder).PlaceholderFormat(sqrl.Dollar)
	queryDeleteBuilder = sqrl.NewUpdateBuilder(sqrl.StatementBuilder).PlaceholderFormat(sqrl.Dollar)
)

type sqlProvider struct {
	pool *sqlx.DB
}

func NewSQLProvider(pool *sqlx.DB) *sqlProvider {
	return &sqlProvider{pool: pool}
}

func (s *sqlProvider) CreateUser(ctx context.Context, user *models.User) (string, error) {
	q := queryInsertBuilder.
		Insert(usersTable).
		Columns(allExceptIDColumns).
		Values(user.Username, user.Firstname, user.Lastname, user.Email, user.Phone, user.Password).
		Returning(idColumn)

	query, args, err := q.ToSql()
	if err != nil {
		return "", fmt.Errorf(buildQuery, err)
	}

	var id string

	err = s.pool.QueryRowContext(ctx, query, args...).Scan(&id)
	if err != nil {
		//return "", fmt.Errorf(executeQuery, err, query, s.pool)
		return "", errors.ErrDuplicateUser
	}

	return id, err
}

func (s *sqlProvider) GetUser(ctx context.Context, id string) (models.User, error) {
	var user models.User

	q := queryBuilder.
		Select(allColumns).
		From(usersTable).
		Where(sqrl.Eq{idColumn: id})

	query, args, err := q.ToSql()
	if err != nil {
		return user, fmt.Errorf(buildQuery, err)
	}

	err = s.pool.GetContext(ctx, &user, query, args...)
	if err != nil {
		return user, fmt.Errorf(executeQuery, err, query, s.pool)
	}

	return user, nil
}

func (s *sqlProvider) DeleteUser(ctx context.Context, id string) error {
	q := queryDeleteBuilder.
		From(usersTable).
		Where(sqrl.Eq{idColumn: id})

	query, args, err := q.ToSql()
	if err != nil {
		return fmt.Errorf(buildQuery, err)
	}

	_, err = s.pool.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf(executeQuery, err, query, s.pool)
	}

	return nil
}

func (s *sqlProvider) UpdateUser(ctx context.Context, id string, user *models.User) error {

	q := queryUpdateBuilder.
		Update(usersTable).
		Set(usernameColumn, user.Username).
		Set(firstNameColumn, user.Firstname).
		Set(lastNameColumn, user.Lastname).
		Set(emailColumn, user.Email).
		Set(phoneColumn, user.Phone).
		Set(passwordColumn, user.Password).
		Where(sqrl.Eq{idColumn: id})

	query, args, err := q.ToSql()
	if err != nil {
		return fmt.Errorf(buildQuery, err)
	}

	_, err = s.pool.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf(executeQuery, err, query, s.pool)
	}

	return nil
}
