package sql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/elgris/sqrl"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
	"order-service/internal/domain/dto"
	"order-service/internal/domain/models"
	"strings"
)

const (
	buildQuery   = "build query: %v"
	executeQuery = "execute query: %v"
	emptyQuery   = "empty query"
)

type OrderField string

func (s OrderField) String() string {
	return string(s)
}

const (
	idColumn         OrderField = "id"
	userIDColumn     OrderField = "user_id"
	totalPriceColumn OrderField = "total_price"
	createAtColumn   OrderField = "created_at"
	deliveryAtColumn OrderField = "delivery_at"
	productsColumn   OrderField = "products"
	addressColumn    OrderField = "address"
	statusColumn     OrderField = "status"
)

func ordersColumnsForCreate() []OrderField {
	return []OrderField{
		userIDColumn,
		totalPriceColumn,
		createAtColumn,
		statusColumn,
		productsColumn,
		deliveryAtColumn,
		addressColumn,
	}
}

func allOrdersColumns() []OrderField {
	return []OrderField{
		idColumn,
		userIDColumn,
		totalPriceColumn,
		createAtColumn,
		statusColumn,
		deliveryAtColumn,
		addressColumn,
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
	ordersTable   = defaultSchema + "." + "orders"
)

type DBClient interface {
	Get(dest interface{}, query string, args ...interface{}) error
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	NamedExec(query string, arg interface{}) (sql.Result, error)
	NamedExecContext(ctx context.Context, query string, arg interface{}) (sql.Result, error)
	QueryRowx(query string, args ...interface{}) *sqlx.Row
	QueryRowxContext(ctx context.Context, query string, args ...interface{}) *sqlx.Row
}

type sqlOrderProvider struct {
	pool *sqlx.DB
}

func NewSQLBusinessRulesProvider(pool *sqlx.DB) *sqlOrderProvider {
	return &sqlOrderProvider{
		pool: pool,
	}
}

var (
	queryBuilder       = sqrl.NewSelectBuilder(sqrl.StatementBuilder).PlaceholderFormat(sqrl.Dollar)
	queryInsertBuilder = sqrl.NewInsertBuilder(sqrl.StatementBuilder).PlaceholderFormat(sqrl.Dollar)
	queryDeleteBuilder = sqrl.NewDeleteBuilder(sqrl.StatementBuilder).PlaceholderFormat(sqrl.Dollar)
	queryUpdateBuilder = sqrl.NewUpdateBuilder(sqrl.StatementBuilder).PlaceholderFormat(sqrl.Dollar)
)

func (s *sqlOrderProvider) CreateOrder(ctx context.Context, order *domain.Order) (int64, error) {
	q := queryInsertBuilder.
		Insert(ordersTable).
		Columns(strings.Join(allColumns(ordersColumnsForCreate), ", ")).
		Values(order.UserID, order.TotalPrice, order.CreatedAt, domain.Created, order.Products, order.DeliveryAt, order.Address).
		Returning(idColumn.String())

	query, args, err := q.ToSql()
	if err != nil {
		return 0, fmt.Errorf(buildQuery, err)
	}

	var id int64

	err = s.pool.QueryRowContext(ctx, query, args...).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf(executeQuery, err)
	}

	return id, nil
}

func (s *sqlOrderProvider) DetailOrder(ctx context.Context, id int64, userID string) (*domain.Order, error) {
	q := queryBuilder.
		Select(strings.Join(allColumns(allOrdersColumns), ", ")).
		From(ordersTable).
		Where(sqrl.Eq{idColumn.String(): id}, sqrl.Eq{userIDColumn.String(): userID}).Limit(1)

	query, args, err := q.ToSql()
	if err != nil {
		return nil, fmt.Errorf(buildQuery, err)
	}

	var row OrderRow

	if err = s.pool.GetContext(ctx, &row, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrOrderNotFound
		}
		return nil, fmt.Errorf(executeQuery, err)
	}

	return domain.RestoreOrderFromDTO(row.ToDTO()), nil
}

func (s *sqlOrderProvider) GetOrderByID(ctx context.Context, id int64) (*domain.Order, error) {
	q := queryBuilder.
		Select(strings.Join(allColumns(allOrdersColumns), ", ")).
		From(ordersTable).
		Where(sqrl.Eq{idColumn.String(): id}).Limit(1)

	query, args, err := q.ToSql()
	if err != nil {
		return nil, fmt.Errorf(buildQuery, err)
	}

	var row OrderRow

	if err = s.pool.GetContext(ctx, &row, query, args...); err != nil {
		return nil, fmt.Errorf(executeQuery, err)
	}

	return domain.RestoreOrderFromDTO(row.ToDTO()), nil
}

func (s *sqlOrderProvider) ListOrders(ctx context.Context, dto *dto.FilterOrderDTO) ([]*domain.Order, error) {
	q := queryBuilder.
		Select(strings.Join(allColumns(allOrdersColumns), ", ")).
		From(ordersTable).
		Where(sqrl.Eq{userIDColumn.String(): dto.UserID}).
		OrderBy(createAtColumn.String() + " DESC").
		Offset(dto.Offset).
		Limit(dto.Limit)

	query, args, err := q.ToSql()
	if err != nil {
		return nil, fmt.Errorf(buildQuery, err)
	}

	rows := make([]OrderRow, 0, dto.Limit)

	if err = s.pool.SelectContext(ctx, &rows, query, args...); err != nil {
		return nil, fmt.Errorf(executeQuery, err)
	}

	orders := make([]*domain.Order, 0, len(rows))

	for _, v := range rows {
		order := domain.RestoreOrderFromDTO(v.ToDTO())
		orders = append(orders, order)
	}

	return orders, nil
}

func (s *sqlOrderProvider) DeleteOrder(ctx context.Context, id int64, userID string) error {
	q := queryDeleteBuilder.
		Delete(ordersTable).
		Where(sqrl.Eq{idColumn.String(): id}, sqrl.Eq{userIDColumn.String(): userID})

	query, args, err := q.ToSql()
	if err != nil {
		return fmt.Errorf(buildQuery, err)
	}

	if _, err = s.pool.ExecContext(ctx, query, args...); err != nil {
		return fmt.Errorf(executeQuery, err)
	}

	return nil
}

func (s *sqlOrderProvider) UpdateOrder(ctx context.Context, id int64, userID string, order *domain.Order) error {
	q := queryUpdateBuilder.
		Update(ordersTable).
		Where(sqrl.Eq{idColumn.String(): id}, sqrl.Eq{userIDColumn.String(): userID})

	if order == nil {
		return fmt.Errorf(emptyQuery)
	}

	q.Set(userIDColumn.String(), order.UserID)
	q.Set(totalPriceColumn.String(), order.TotalPrice)
	q.Set(createAtColumn.String(), order.CreatedAt)
	q.Set(statusColumn.String(), order.Status)

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

func (s *sqlOrderProvider) GetOrderByIDThenUpdate(ctx context.Context, id int64, fn domain.IntermediateOrderFunc) (*domain.Order, error) {
	if fn == nil {
		order, err := getOrderByID(ctx, s.pool, id)
		return order, err
	}

	tx, err := s.pool.BeginTxx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("begin transaction: %w", err)
	}

	defer func() {
		if txErr := tx.Rollback(); txErr != nil && !errors.Is(txErr, sql.ErrTxDone) {
			//log.GetLoggerFromContext(ctx).Errorf("Failed rollback transaction: %v", txErr)
		}
	}()

	order, err := getOrderByID(ctx, tx, id)
	if err != nil {
		return nil, fmt.Errorf("get order: %w", err)
	}

	ok, err := fn(order)
	if err != nil {
		return nil, err
	}

	if !ok {
		return order, nil
	}

	if err = s.update(ctx, tx, order); err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit transaction: %w", err)
	}

	return order, nil
}

func getOrderByID(ctx context.Context, db DBClient, id int64) (*domain.Order, error) {
	q := queryBuilder.
		Select(strings.Join(allColumns(allOrdersColumns), ", ")).
		From(ordersTable).
		Where(sqrl.Eq{idColumn.String(): id}).Limit(1)

	query, args, err := q.ToSql()
	if err != nil {
		return nil, fmt.Errorf(buildQuery, err)
	}

	var row OrderRow

	if err = db.GetContext(ctx, &row, query, args...); err != nil {
		return nil, fmt.Errorf(executeQuery, err)
	}

	return domain.RestoreOrderFromDTO(row.ToDTO()), nil
}

func (s *sqlOrderProvider) update(ctx context.Context, db DBClient, order *domain.Order) error {
	query := `UPDATE user_service.orders 
			SET 
				status=:status,
				created_at=:created_at,
				user_id=:user_id,
				total_price=:total_price
			WHERE id=:id`

	res, err := db.NamedExecContext(ctx, query, FromModel(order))
	if err != nil {
		return fmt.Errorf("execute update order query: %v", err)
	}

	affectedRows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("get affected rows: %v", err)
	}

	if affectedRows == 0 {
		return fmt.Errorf("no rows in result set")
	}

	return nil
}
