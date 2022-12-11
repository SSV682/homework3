package sql

import (
	"context"
	"fmt"
	"github.com/elgris/sqrl"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
	"user-service/internal/domain/models"
)

const (
	buildQuery   = "build query: %v"
	executeQuery = "execute query: %v"
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
	usersTable = "users"
)

var (
	queryBuilder       = sqrl.NewSelectBuilder(sqrl.StatementBuilder).PlaceholderFormat(sqrl.Dollar)
	queryInsertBuilder = sqrl.NewInsertBuilder(sqrl.StatementBuilder).PlaceholderFormat(sqrl.Dollar)
)

type sqlProvider struct {
	pool *sqlx.DB
}

func NewSQLProvider(pool *sqlx.DB) *sqlProvider {
	return &sqlProvider{pool: pool}
}

func (s *sqlProvider) CreateUser(ctx context.Context, user *models.User) (int64, error) {
	q := queryInsertBuilder.
		Insert(usersTable).
		Columns(allExceptIDColumns).
		Values(user.Username, user.Firstname, user.Lastname, user.Email, user.Phone, user.Password).
		Returning(idColumn)

	query, args, err := q.ToSql()
	if err != nil {
		return 0, fmt.Errorf(buildQuery, err)
	}

	var id int64

	err = s.pool.QueryRowContext(ctx, query, args...).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf(executeQuery, err)
	}

	return id, err
}

func (s *sqlProvider) GetUser(ctx context.Context, id int64) (user models.User, err error) {
	user = models.User{}
	err = s.pool.QueryRowContext(ctx, `SELECT id, username, firstname, lastname, email, phone, password FROM users WHERE id = $1;`, id).Scan(&user.Id, &user.Username, &user.Firstname, &user.Lastname, &user.Email, &user.Phone, &user.Password)
	return
}

func (s *sqlProvider) DeleteUser(ctx context.Context, id int64) (err error) {
	_, err = s.pool.ExecContext(ctx, `DELETE FROM users WHERE id = $1;`, id)
	return

}

func (s *sqlProvider) UpdateUser(ctx context.Context, id int64, user *models.User) error {
	user.Id = id
	query := `
		UPDATE users
		SET username=:username,
		    firstname=:firstname,
		    lastname=:lastname, 
		    email=:email, 
		    phone=:phone
		WHERE id=:id
	`
	rows, err := s.pool.NamedQueryContext(ctx, query, map[string]interface{}{
		"username":  user.Username,
		"firstname": user.Firstname,
		"lastname":  user.Lastname,
		"email":     user.Email,
		"phone":     user.Phone,
		"id":        id,
	})
	defer rows.Close()
	return err
}
