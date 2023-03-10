package sql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/elgris/sqrl"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
	"stock-service/internal/domain/dto"
	domain "stock-service/internal/domain/models"
	"strings"
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
	idColumn       OrderField = "id"
	nameColumn     OrderField = "name"
	quantityColumn OrderField = "Quantity"
)

func productsColumnsForCreate() []OrderField {
	return []OrderField{
		nameColumn,
		quantityColumn,
	}
}

func allProductsColumns() []OrderField {
	return []OrderField{
		idColumn,
		nameColumn,
		quantityColumn,
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
	productTable  = defaultSchema + "." + "products"
)

type DBClient interface {
	Get(dest interface{}, query string, args ...interface{}) error
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	NamedExec(query string, arg interface{}) (sql.Result, error)
	NamedExecContext(ctx context.Context, query string, arg interface{}) (sql.Result, error)
	QueryRowx(query string, args ...interface{}) *sqlx.Row
	QueryRowxContext(ctx context.Context, query string, args ...interface{}) *sqlx.Row
}

type sqlProductProvider struct {
	pool *sqlx.DB
}

func NewSQLProductProvider(pool *sqlx.DB) *sqlProductProvider {
	return &sqlProductProvider{
		pool: pool,
	}
}

var (
	queryBuilder       = sqrl.NewSelectBuilder(sqrl.StatementBuilder).PlaceholderFormat(sqrl.Dollar)
	queryInsertBuilder = sqrl.NewInsertBuilder(sqrl.StatementBuilder).PlaceholderFormat(sqrl.Dollar)
	queryDeleteBuilder = sqrl.NewDeleteBuilder(sqrl.StatementBuilder).PlaceholderFormat(sqrl.Dollar)
	queryUpdateBuilder = sqrl.NewUpdateBuilder(sqrl.StatementBuilder).PlaceholderFormat(sqrl.Dollar)
)

func (s *sqlProductProvider) CreateProduct(ctx context.Context, p domain.Product) (int64, error) {
	q := queryInsertBuilder.
		Insert(productTable).
		Columns(strings.Join(allColumns(productsColumnsForCreate), ", ")).
		Values(p.Name, p.Quantity).
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

func (s *sqlProductProvider) UpdateProduct(ctx context.Context, id int64, p domain.Product) error {
	q := queryUpdateBuilder.
		Update(productTable).
		Where(sqrl.Eq{idColumn.String(): id})

	q.Set(nameColumn.String(), p.Name)
	q.Set(quantityColumn.String(), p.Quantity)

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

func (s *sqlProductProvider) GetByID(ctx context.Context, id int64) (*domain.Product, error) {
	res, err := getProductByID(ctx, s.pool, id)
	if err != nil {
		return nil, fmt.Errorf("get by id: %w", err)
	}

	return res, nil
}

func (s *sqlProductProvider) RavageStock(ctx context.Context, productsOrder []domain.Product) error {
	tx, err := s.pool.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}

	fail := func(err error) error {
		return fmt.Errorf("approve order: %v", err)
	}

	defer func() {
		if txErr := tx.Rollback(); txErr != nil && !errors.Is(txErr, sql.ErrTxDone) {
			log.Errorf("transaction error: %s", txErr)
		}
	}()

	for _, v := range productsOrder {

		var enough bool
		if err = tx.QueryRowContext(ctx, "SELECT (quantity >= $1) FROM user_service.products WHERE id=$2", v.Quantity, v.ID).Scan(&enough); err != nil {
			if err == sql.ErrNoRows {
				return fail(fmt.Errorf("no such %s", v.Name))
			}
			return fail(err)
		}

		fmt.Println(fmt.Sprintf("product id: %d, quantity: %d, enough %t", v.ID, v.Quantity, enough))
		if !enough {
			return fail(fmt.Errorf("not enough %s", v.Name))
		}

		_, err = tx.ExecContext(ctx, "UPDATE user_service.products SET quantity = quantity- $1 WHERE id = $2", v.Quantity, v.ID)
		if err != nil {
			return fail(err)
		}
	}

	if err = tx.Commit(); err != nil {
		return fail(fmt.Errorf("commit transaction: %w", err))
	}

	return nil
}

func (s *sqlProductProvider) FillStock(ctx context.Context, productsOrder []domain.Product) error {
	tx, err := s.pool.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}

	defer func() {
		if txErr := tx.Rollback(); txErr != nil && !errors.Is(txErr, sql.ErrTxDone) {
		}
	}()

	fail := func(err error) error {
		return fmt.Errorf("approve order: %v", err)
	}

	for _, v := range productsOrder {
		var quantity int64
		exist := true
		if err = tx.QueryRowContext(ctx, "SELECT quantity FROM user_service.products WHERE id=$1", v.ID).Scan(&quantity); err != nil {
			if err == sql.ErrNoRows {
				exist = false
			} else {
				return fail(err)
			}
		}

		if exist {
			_, err = tx.ExecContext(ctx, "UPDATE user_service.products SET quantity = $1 WHERE id = $2", quantity+v.Quantity, v.ID)
			if err != nil {
				return fail(err)
			}
		}
	}

	if err = tx.Commit(); err != nil {
		return fail(fmt.Errorf("commit transaction: %w", err))
	}

	return nil
}

func (s *sqlProductProvider) ListStock(ctx context.Context, filter dto.FilterProductDTO) ([]*domain.Product, error) {
	q := queryBuilder.
		Select(strings.Join(allColumns(allProductsColumns), ", ")).
		From(productTable).
		Offset(filter.Offset).
		Limit(filter.Limit)

	if filter.Quantity != nil {
		q.Where(sqrl.Eq{quantityColumn.String(): *filter.Quantity})
	}

	query, args, err := q.ToSql()
	if err != nil {
		return nil, fmt.Errorf(buildQuery, err)
	}

	rows := make([]ProductRow, 0, filter.Limit)

	if err = s.pool.SelectContext(ctx, &rows, query, args...); err != nil {
		return nil, fmt.Errorf(executeQuery, err)
	}

	products := make([]*domain.Product, 0, len(rows))

	for _, v := range rows {
		products = append(products, &domain.Product{
			ID:       v.ID,
			Quantity: v.Quantity,
			Name:     v.Name,
		})
	}

	return products, nil
}

func getProductByID(ctx context.Context, db DBClient, id int64) (*domain.Product, error) {
	q := queryBuilder.
		Select(strings.Join(allColumns(allProductsColumns), ", ")).
		From(productTable).
		Where(sqrl.Eq{idColumn.String(): id}).Limit(1)

	query, args, err := q.ToSql()
	if err != nil {
		return nil, fmt.Errorf(buildQuery, err)
	}

	var row ProductRow

	if err = db.GetContext(ctx, &row, query, args...); err != nil {
		return nil, fmt.Errorf(executeQuery, err)
	}

	return row.ToModel(), nil
}
