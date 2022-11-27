package sql

import (
	"context"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
	"homework2/internal/domain/models"
	"homework2/internal/provider"
)

type sqlProvider struct {
	pool *sqlx.DB
}

var _ provider.UserProvider = &sqlProvider{}

func NewSQLProvider(pool *sqlx.DB) *sqlProvider {
	return &sqlProvider{pool: pool}
}

func (s *sqlProvider) CreateUser(ctx context.Context, user *models.User) error {
	query := `
		INSERT INTO users(username, firstname, lastname, email, phone)
		VALUES (:username, :firstname, :lastname, :email, :phone)
		RETURNING id;
	`
	rows, err := s.pool.NamedQueryContext(ctx, query, map[string]interface{}{
		"username":  user.Username,
		"firstname": user.Firstname,
		"lastname":  user.Lastname,
		"email":     user.Email,
		"phone":     user.Phone,
	})
	defer rows.Close()
	if rows.Next() {
		rows.Scan(&user.Id)
	}
	return err
}

func (s *sqlProvider) GetUser(ctx context.Context, id int64) (user models.User, err error) {
	user = models.User{}
	err = s.pool.QueryRowContext(ctx, `SELECT id, username, firstname, lastname, email, phone FROM users WHERE id = $1;`, id).Scan(&user.Id, &user.Username, &user.Firstname, &user.Lastname, &user.Email, &user.Phone)
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
