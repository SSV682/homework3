package sql

import (
	"context"
	"fmt"
	"github.com/elgris/sqrl"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
	"strings"
	"user-service/internal/domain/dto"
	"user-service/internal/domain/models"
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
	statusColumn     OrderField = "status"
)

func ordersColumnsForCreate() []OrderField {
	return []OrderField{
		userIDColumn,
		totalPriceColumn,
		createAtColumn,
		statusColumn,
	}
}

func allOrdersColumns() []OrderField {
	return []OrderField{
		idColumn,
		userIDColumn,
		totalPriceColumn,
		createAtColumn,
		statusColumn,
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
		Values(order.UserID(), order.TotalPrice(), order.CreatedAt(), domain.StatusCreated).
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
		Where(sqrl.Eq{idColumn.String(): id}, sqrl.Eq{userIDColumn.String(): userID})

	query, args, err := q.ToSql()
	if err != nil {
		return nil, fmt.Errorf(buildQuery, err)
	}

	var row OrderRow

	if err = s.pool.SelectContext(ctx, &row, query, args...); err != nil {
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

	var rows []OrderRow

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

	_, err = s.pool.ExecContext(ctx, query, args...)
	if err != nil {
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
