package sql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/elgris/sqrl"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
	domain "notification-service/internal/domain/models"
	"strings"
)

const (
	buildQuery   = "build query: %v"
	executeQuery = "execute query: %v"
)

type NotificationField string

func (s NotificationField) String() string {
	return string(s)
}

const (
	idColumn      NotificationField = "id"
	mailColumn    NotificationField = "mail"
	messageColumn NotificationField = "message"
)

type UserInfoField string

func (s UserInfoField) String() string {
	return string(s)
}

const (
	userIDUserInfoColumn UserInfoField = "user_id"
	mailUserInfoColumn   UserInfoField = "mail"
)

func notificationColumnsForCreate() []NotificationField {
	return []NotificationField{
		mailColumn,
		messageColumn,
	}
}

func allNotificationColumns() []NotificationField {
	return []NotificationField{
		idColumn,
		mailColumn,
		messageColumn,
	}
}

func allUserInfoColumns() string {
	m := []UserInfoField{
		userIDUserInfoColumn,
		mailUserInfoColumn,
	}

	result := make([]string, 0, len(m))

	for _, v := range m {
		result = append(result, v.String())
	}

	return strings.Join(result, ",")
}

func allColumns(fn func() []NotificationField) []string {
	fs := fn()
	result := make([]string, 0, len(fs))

	for _, v := range fs {
		result = append(result, v.String())
	}

	return result
}

const (
	defaultSchema     = "user_service"
	notificationTable = defaultSchema + "." + "notification"
	userInfoTable     = defaultSchema + "." + "user_info"
)

type DBClient interface {
	Get(dest interface{}, query string, args ...interface{}) error
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	NamedExec(query string, arg interface{}) (sql.Result, error)
	NamedExecContext(ctx context.Context, query string, arg interface{}) (sql.Result, error)
	QueryRowx(query string, args ...interface{}) *sqlx.Row
	QueryRowxContext(ctx context.Context, query string, args ...interface{}) *sqlx.Row
}

type sqlProvider struct {
	pool *sqlx.DB
}

func NewSQLProvider(pool *sqlx.DB) *sqlProvider {
	return &sqlProvider{
		pool: pool,
	}
}

var (
	queryBuilder       = sqrl.NewSelectBuilder(sqrl.StatementBuilder).PlaceholderFormat(sqrl.Dollar)
	queryInsertBuilder = sqrl.NewInsertBuilder(sqrl.StatementBuilder).PlaceholderFormat(sqrl.Dollar)
	queryUpdateBuilder = sqrl.NewUpdateBuilder(sqrl.StatementBuilder).PlaceholderFormat(sqrl.Dollar)
)

func (s *sqlProvider) Create(ctx context.Context, p domain.Notification) error {

	q := queryInsertBuilder.
		Insert(notificationTable).
		Columns(strings.Join(allColumns(notificationColumnsForCreate), ", ")).
		Values(p.Mail, p.Message)

	query, args, err := q.ToSql()
	if err != nil {
		return fmt.Errorf(buildQuery, err)
	}

	_, err = s.pool.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf(executeQuery, err)
	}

	return nil
}

func (s *sqlProvider) List(ctx context.Context) ([]*domain.Notification, error) {
	q := queryBuilder.
		Select(strings.Join(allColumns(allNotificationColumns), ", ")).
		From(notificationTable)

	query, args, err := q.ToSql()
	if err != nil {
		return nil, fmt.Errorf(buildQuery, err)
	}

	var rows []NotificationRow

	if err = s.pool.SelectContext(ctx, &rows, query, args...); err != nil {
		return nil, fmt.Errorf(executeQuery, err)
	}

	notifications := make([]*domain.Notification, 0, len(rows))

	for _, v := range rows {
		notifications = append(notifications, &domain.Notification{
			ID:      v.ID,
			Mail:    v.Mail,
			Message: v.Message,
		})
	}

	return notifications, nil
}

func (s *sqlProvider) UpdateUserInfo(ctx context.Context, user domain.User) error {
	tx, err := s.pool.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}

	defer func() {
		if txErr := tx.Rollback(); txErr != nil && !errors.Is(txErr, sql.ErrTxDone) {
			//log.GetLoggerFromContext(ctx).Errorf("Failed rollback transaction: %v", txErr)
		}
	}()

	existingUser, err := s.getUserByID(ctx, tx, user.ID)
	if err != nil {
		return fmt.Errorf("get user by id for update: %w", err)
	}

	if existingUser != nil {
		s.updateUser(ctx, existingUser.ID, user)
	} else {
		s.createUser(ctx, user)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}

func (s *sqlProvider) GetUserByID(ctx context.Context, id string) (*domain.User, error) {
	return s.getUserByID(ctx, s.pool, id)
}

func (s *sqlProvider) getUserByID(ctx context.Context, db DBClient, id string) (*domain.User, error) {
	q := queryBuilder.
		Select(allUserInfoColumns()).
		From(userInfoTable).
		Where(sqrl.Eq{userIDUserInfoColumn.String(): id}).Limit(1)

	query, args, err := q.ToSql()
	if err != nil {
		return nil, fmt.Errorf(buildQuery, err)
	}

	var row UserInfoRow

	if err = db.GetContext(ctx, &row, query, args...); err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, nil
		default:
			return nil, fmt.Errorf(executeQuery, err)
		}
	}

	return row.ToModel(), nil
}

func (s *sqlProvider) updateUser(ctx context.Context, id string, user domain.User) error {
	q := queryUpdateBuilder.
		Update(userInfoTable).
		Set(mailUserInfoColumn.String(), user.Mail).
		Where(sqrl.Eq{userIDUserInfoColumn.String(): id})

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

func (s *sqlProvider) createUser(ctx context.Context, p domain.User) error {
	q := queryInsertBuilder.
		Insert(userInfoTable).
		Columns(allUserInfoColumns()).
		Values(p.ID, p.Mail)

	query, args, err := q.ToSql()
	if err != nil {
		return fmt.Errorf(buildQuery, err)
	}

	_, err = s.pool.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf(executeQuery, err)
	}

	return nil
}
