package db

import (
	"context"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
	"homework2/internal/domain/models"
)

type PgUserRepository struct {
	db *sqlx.DB
}

func NewPgUserRepository(dsn string) (*PgUserRepository, error) {
	db, err := sqlx.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return &PgUserRepository{db: db}, nil
}

func (pgus *PgUserRepository) CreateUser(ctx context.Context, user *models.User) error {
	query := `
		INSERT INTO users(username, firstname, lastname, email, phone)
		VALUES (:username, :firstname, :lastname, :email, :phone)
		RETURNING id;
	`
	rows, err := pgus.db.NamedQueryContext(ctx, query, map[string]interface{}{
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

func (pgus *PgUserRepository) GetUser(ctx context.Context, id int64) (user models.User, err error) {
	user = models.User{}
	err = pgus.db.QueryRowContext(ctx, `SELECT id, username, firstname, lastname, email, phone FROM users WHERE id = $1;`, id).Scan(&user.Id, &user.Username, &user.Firstname, &user.Lastname, &user.Email, &user.Phone)
	return
}

func (pgus *PgUserRepository) DeleteUser(ctx context.Context, id int64) (err error) {
	_, err = pgus.db.ExecContext(ctx, `DELETE FROM users WHERE id = $1;`, id)
	return

}

func (pgus *PgUserRepository) UpdateUser(ctx context.Context, id int64, user *models.User) error {
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
	rows, err := pgus.db.NamedQueryContext(ctx, query, map[string]interface{}{
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
