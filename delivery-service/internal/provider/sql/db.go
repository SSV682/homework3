package sql

import (
	"context"
	"database/sql"
	domain "delivery-service/internal/domain/models"
	"errors"
	"fmt"
	"github.com/elgris/sqrl"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
	"strings"
	"time"
)

const (
	buildQuery   = "build query: %v"
	executeQuery = "execute query: %v"
)

type OrderField string

func (s OrderField) String() string {
	return string(s)
}

const (
	idColumn           OrderField = "id"
	orderIDColumn      OrderField = "order_id"
	orderContentColumn OrderField = "order_content"
	addressColumn      OrderField = "address"
	dateColumn         OrderField = "date"
)

func productsColumnsForCreate() []OrderField {
	return []OrderField{
		orderIDColumn,
		orderContentColumn,
		addressColumn,
		dateColumn,
	}
}

func allProductsColumns() []OrderField {
	return []OrderField{
		idColumn,
		orderIDColumn,
		orderContentColumn,
		addressColumn,
		dateColumn,
	}
}

func allColumns(fn func() []OrderField) []string {
	fs := fn()
	result := make([]string, 0, len(fs))

	for _, v := range fs {
		result = append(result, v.String())
	}

	return result
}

const (
	defaultSchema = "user_service"
	deliveryTable = defaultSchema + "." + "delivery"
)

type DBClient interface {
	Get(dest interface{}, query string, args ...interface{}) error
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	NamedExec(query string, arg interface{}) (sql.Result, error)
	NamedExecContext(ctx context.Context, query string, arg interface{}) (sql.Result, error)
	QueryRowx(query string, args ...interface{}) *sqlx.Row
	QueryRowxContext(ctx context.Context, query string, args ...interface{}) *sqlx.Row
}

type sqlDeliveryProvider struct {
	pool                *sqlx.DB
	maxDeliveriesPerDay int
}

func NewSQLProductProvider(pool *sqlx.DB, maxDeliveriesPerDay int) *sqlDeliveryProvider {
	return &sqlDeliveryProvider{
		pool:                pool,
		maxDeliveriesPerDay: maxDeliveriesPerDay,
	}
}

var (
	queryBuilder       = sqrl.NewSelectBuilder(sqrl.StatementBuilder).PlaceholderFormat(sqrl.Dollar)
	queryInsertBuilder = sqrl.NewInsertBuilder(sqrl.StatementBuilder).PlaceholderFormat(sqrl.Dollar)
	queryDeleteBuilder = sqrl.NewDeleteBuilder(sqrl.StatementBuilder).PlaceholderFormat(sqrl.Dollar)
)

func (s *sqlDeliveryProvider) CheckPossibleDelivery(ctx context.Context, entry domain.DeliveryEntry) error {
	tx, err := s.pool.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}

	defer func() {
		if txErr := tx.Rollback(); txErr != nil && !errors.Is(txErr, sql.ErrTxDone) {
			log.Errorf("defer func: %v", txErr)
		}
	}()

	fail := func(err error) error {
		return fmt.Errorf("approve order: %v", err)
	}

	var deliveries int
	if err = tx.QueryRowContext(ctx, "SELECT count(*) FROM user_service.delivery WHERE date=$1", entry.Date).Scan(&deliveries); err != nil {
		if err == sql.ErrNoRows {
			deliveries = 0
		} else {
			return fail(err)
		}
	}

	if deliveries < s.maxDeliveriesPerDay {
		_, err = s.createDelivery(ctx, tx, entry)
		if err != nil {
			return fail(err)
		}
	} else {
		return fail(errors.New("there are no couriers available"))
	}

	if err = tx.Commit(); err != nil {
		return fail(fmt.Errorf("commit transaction: %w", err))
	}

	return nil
}

func (s *sqlDeliveryProvider) ListDelivery(ctx context.Context, date time.Time) ([]*domain.DeliveryEntry, error) {
	q := queryBuilder.
		Select(strings.Join(allColumns(allProductsColumns), ", ")).
		From(deliveryTable).
		Where(sqrl.Eq{dateColumn.String(): date})

	query, args, err := q.ToSql()
	if err != nil {
		return nil, fmt.Errorf(buildQuery, err)
	}

	log.Infof("query: %s, args: %s", query, date)

	rows := make([]DeliveryRow, 0, s.maxDeliveriesPerDay)

	if err = s.pool.SelectContext(ctx, &rows, query, args...); err != nil {
		return nil, fmt.Errorf(executeQuery, err)
	}

	products := make([]*domain.DeliveryEntry, 0, len(rows))

	for _, v := range rows {
		products = append(products, &domain.DeliveryEntry{
			ID:           v.ID,
			OrderID:      v.OrderID,
			OrderContent: v.OrderContent,
			Address:      v.Address,
			Date:         v.Date,
		})
	}

	return products, nil
}

func (s *sqlDeliveryProvider) createDelivery(ctx context.Context, db DBClient, entry domain.DeliveryEntry) (int64, error) {
	q := queryInsertBuilder.
		Insert(deliveryTable).
		Columns(strings.Join(allColumns(productsColumnsForCreate), ", ")).
		Values(entry.OrderID, entry.OrderContent, entry.Address, entry.Date).
		Returning(idColumn.String())

	query, args, err := q.ToSql()
	if err != nil {
		return 0, fmt.Errorf(buildQuery, err)
	}

	var id int64

	err = db.QueryRowxContext(ctx, query, args...).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf(executeQuery, err)
	}

	return id, nil
}

func (s *sqlDeliveryProvider) RejectDelivery(ctx context.Context, orderID int64) error {
	q := queryDeleteBuilder.
		Delete(deliveryTable).
		Where(sqrl.Eq{orderIDColumn.String(): orderID})

	query, args, err := q.ToSql()
	if err != nil {
		return fmt.Errorf(buildQuery, err)
	}

	if _, err = s.pool.ExecContext(ctx, query, args...); err != nil {
		return fmt.Errorf(executeQuery, err)
	}

	return nil
}
