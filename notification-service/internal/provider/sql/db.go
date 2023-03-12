package sql

import (
	"context"
	"database/sql"
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
	userIDColumn  NotificationField = "user_id"
	messageColumn NotificationField = "message"
)

func notificationColumnsForCreate() []NotificationField {
	return []NotificationField{
		userIDColumn,
		messageColumn,
	}
}

func allNotificationColumns() []NotificationField {
	return []NotificationField{
		idColumn,
		userIDColumn,
		messageColumn,
	}
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
)

func (s *sqlProvider) Create(ctx context.Context, p domain.Order) error {
	q := queryInsertBuilder.
		Insert(notificationTable).
		Columns(strings.Join(allColumns(notificationColumnsForCreate), ", ")).
		Values(p.UserID, fmt.Sprintf("Order %d %s", p.ID, p.Status))

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
			UserID:  v.UserID,
			Message: v.Message,
		})
	}

	return notifications, nil
}
