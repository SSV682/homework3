package sql

import (
	domain "billing-service/internal/domain/models"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/elgris/sqrl"
	"github.com/google/uuid"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
	"strings"
)

const (
	buildQuery   = "build query: %v"
	executeQuery = "execute query: %v"
)

type AccountFields string
type OutboxFields string

func (s AccountFields) String() string {
	return string(s)
}
func (s OutboxFields) String() string {
	return string(s)
}

const (
	idColumn            AccountFields = "id"
	amountColumn        AccountFields = "amount"
	idOutboxColumn      OutboxFields  = "id"
	topicOutboxColumn   OutboxFields  = "topic"
	messageOutboxColumn OutboxFields  = "message"
)

func allAccountColumns() []AccountFields {
	return []AccountFields{
		idColumn,
		amountColumn,
	}
}

func allOutboxColumns() []OutboxFields {
	return []OutboxFields{
		idOutboxColumn,
		topicOutboxColumn,
		messageOutboxColumn,
	}
}

func createOutboxColumns() []OutboxFields {
	return []OutboxFields{
		topicOutboxColumn,
		messageOutboxColumn,
	}
}

func createAccountColumns() []AccountFields {
	return []AccountFields{
		idColumn,
		amountColumn,
	}
}

func accountColumns(fn func() []AccountFields) []string {
	fs := fn()
	result := make([]string, 0, len(fs))

	for _, v := range fs {
		result = append(result, v.String())
	}

	return result
}

func outboxColumns(fn func() []OutboxFields) []string {
	fs := fn()
	result := make([]string, 0, len(fs))

	for _, v := range fs {
		result = append(result, v.String())
	}

	return result
}

const (
	defaultSchema = "user_service"
	billingTable  = defaultSchema + "." + "account"
	outboxTable   = defaultSchema + "." + "outbox"
)

type DBClient interface {
	Get(dest interface{}, query string, args ...interface{}) error
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	NamedExec(query string, arg interface{}) (sql.Result, error)
	NamedExecContext(ctx context.Context, query string, arg interface{}) (sql.Result, error)
	QueryRowx(query string, args ...interface{}) *sqlx.Row
	QueryRowxContext(ctx context.Context, query string, args ...interface{}) *sqlx.Row
}

type sqlBillingProvider struct {
	pool                *sqlx.DB
	maxDeliveriesPerDay int
	messageTopic        string
	commandTopic        string
}

func NewSQLBillingProvider(pool *sqlx.DB, messageTopic, commandTopic string) *sqlBillingProvider {
	return &sqlBillingProvider{
		pool:         pool,
		messageTopic: messageTopic,
		commandTopic: commandTopic,
	}
}

var (
	queryBuilder       = sqrl.NewSelectBuilder(sqrl.StatementBuilder).PlaceholderFormat(sqrl.Dollar)
	queryInsertBuilder = sqrl.NewInsertBuilder(sqrl.StatementBuilder).PlaceholderFormat(sqrl.Dollar)
	queryDeleteBuilder = sqrl.NewDeleteBuilder(sqrl.StatementBuilder).PlaceholderFormat(sqrl.Dollar)
)

func (s *sqlBillingProvider) CheckPossiblePayment(ctx context.Context, order domain.Order) error {
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
		s.CreateOutboxCommand(ctx, domain.ResponseCommand{
			Topic: s.commandTopic,
			Command: domain.Command{
				OrderID: order.ID,
				Status:  domain.PaymentRejected,
			},
		})

		return fmt.Errorf("approve order: %v", err)
	}

	var enough bool
	if err = tx.QueryRowContext(ctx, "SELECT amount>=$1 FROM user_service.account WHERE id=$2", order.TotalPrice, order.UserID).Scan(&enough); err != nil {
		return fail(err)
	}

	if enough {
		_, err = tx.ExecContext(ctx, "UPDATE user_service.account SET amount = amount -$1 WHERE id = $2", order.TotalPrice, order.UserID)
		if err != nil {
			return fail(err)
		}
	} else {
		return fail(errors.New("there are not enough money"))
	}

	if err = tx.Commit(); err != nil {
		return fail(fmt.Errorf("commit transaction: %w", err))
	}

	s.CreateOutboxCommand(ctx, domain.ResponseCommand{
		Topic: s.commandTopic,
		Command: domain.Command{
			OrderID: order.ID,
			Status:  domain.PaymentApproved,
		},
	})

	return nil
}

func (s *sqlBillingProvider) CreateOutboxCommand(ctx context.Context, command domain.ResponseCommand) (int64, error) {
	message, err := json.Marshal(NewResponseCommand(command.Command))
	if err != nil {
		//TODO:
	}

	q := queryInsertBuilder.
		Insert(outboxTable).
		Columns(strings.Join(outboxColumns(createOutboxColumns), ", ")).
		Values(command.Topic, message)

	query, args, err := q.ToSql()
	if err != nil {
		return 0, fmt.Errorf(buildQuery, err)
	}

	var id int64

	err = s.pool.QueryRowxContext(ctx, query, args...).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf(executeQuery, err)
	}

	return id, nil
}

func (s *sqlBillingProvider) GetNextOutboxCommand(ctx context.Context) (*domain.ReadyResponseCommand, error) {
	q := queryBuilder.
		Select(strings.Join(outboxColumns(allOutboxColumns), ", ")).
		From(outboxTable)

	query, args, err := q.ToSql()
	if err != nil {
		return nil, fmt.Errorf(buildQuery, err)
	}

	var c CommandRow

	err = s.pool.GetContext(ctx, &c, query, args...)
	if err != nil {
		return nil, fmt.Errorf(executeQuery, err)
	}

	return &domain.ReadyResponseCommand{
		ID:      c.ID,
		Topic:   c.Topic,
		Command: c.Message,
	}, nil
}

func (s *sqlBillingProvider) DeleteOutboxCommand(ctx context.Context, id int64) error {
	q := queryDeleteBuilder.
		Delete(outboxTable).
		Where(sqrl.Eq{idOutboxColumn.String(): id})

	query, args, err := q.ToSql()
	if err != nil {
		return fmt.Errorf(buildQuery, err)
	}

	if _, err = s.pool.ExecContext(ctx, query, args...); err != nil {
		return fmt.Errorf(executeQuery, err)
	}

	return nil
}

func (s *sqlBillingProvider) DetailAccount(ctx context.Context, id uuid.UUID) (*domain.Account, error) {
	q := queryBuilder.
		Select(strings.Join(accountColumns(allAccountColumns), ", ")).
		From(billingTable).
		Where(sqrl.Eq{idColumn.String(): id})

	query, args, err := q.ToSql()
	if err != nil {
		return nil, fmt.Errorf(buildQuery, err)
	}

	var ar AccountRow

	if err = s.pool.GetContext(ctx, &ar, query, args...); err != nil {
		return nil, fmt.Errorf(executeQuery, err)
	}

	return ar.ToModel(), nil
}

func (s *sqlBillingProvider) FillAccount(ctx context.Context, id uuid.UUID, amount float64) error {
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
		return fmt.Errorf("fill account: %v", err)
	}

	_, err = tx.ExecContext(ctx, "UPDATE user_service.account SET amount = amount +$1 WHERE id = $2", amount, id.String())
	if err != nil {
		return fail(err)
	}

	if err = tx.Commit(); err != nil {
		return fail(fmt.Errorf("commit transaction: %w", err))
	}

	return nil
}

func (s *sqlBillingProvider) RejectPayment(ctx context.Context, order domain.Order) error {
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
		return fmt.Errorf("fill account: %v", err)
	}

	_, err = tx.ExecContext(ctx, "UPDATE user_service.account SET amount = amount +$1 WHERE id = $2", order.TotalPrice, order.UserID.String())
	if err != nil {
		return fail(err)
	}

	if err = tx.Commit(); err != nil {
		return fail(fmt.Errorf("commit transaction: %w", err))
	}

	s.CreateOutboxCommand(ctx, domain.ResponseCommand{
		Topic: s.commandTopic,
		Command: domain.Command{
			OrderID: order.ID,
			Status:  domain.PaymentRejected,
		},
	})

	return nil
}

func (s *sqlBillingProvider) CreateAccount(ctx context.Context, id uuid.UUID) error {
	q := queryInsertBuilder.
		Insert(billingTable).
		Columns(strings.Join(accountColumns(createAccountColumns), ", ")).
		Values(id.String(), 0)

	query, args, err := q.ToSql()
	if err != nil {
		return fmt.Errorf(buildQuery, err)
	}

	err = s.pool.QueryRowxContext(ctx, query, args...).Err()
	if err != nil {
		return fmt.Errorf(executeQuery, err)
	}

	return nil
}
