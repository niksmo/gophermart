package orders

import (
	"context"
	"errors"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/niksmo/gophermart/internal/errs"
	"github.com/niksmo/gophermart/pkg/logger"
)

type OrdersRepository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) OrdersRepository {
	return OrdersRepository{db: db}
}

func (r OrdersRepository) Create(
	ctx context.Context, userID int32, orderNumber string,
) (order OrderScheme, err error) {
	stmt := `
	INSERT INTO orders (user_id, number) VALUES ($1, $2)
	RETURNING id, user_id, number, status, accrual, uploaded_at;
	`
	row := r.db.QueryRow(ctx, stmt, userID, orderNumber)
	err = order.ScanRow(row)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			err = errs.ErrOrderUploaded
			return
		}
		logger.Instance.Error().
			Err(err).
			Caller().
			Int32("userID", userID).
			Str("orderNumber", orderNumber).
			Msg("creating order")
		return
	}

	return
}

func (r OrdersRepository) ReadByOrderNumber(
	ctx context.Context, orderNumber string,
) (OrderScheme, error) {
	stmt := `
	SELECT id, user_id, number, status, accrual, uploaded_at
	FROM orders
	WHERE number = $1;
	`

	var order OrderScheme
	err := order.ScanRow(r.db.QueryRow(ctx, stmt, orderNumber))
	if err != nil {
		logger.Instance.Error().Err(err).Caller().Msg("reading order by number")
		return order, err
	}
	return order, nil
}

func (r OrdersRepository) ReadListByUser(
	ctx context.Context, userID int32,
) (OrderListScheme, error) {
	log := logger.Instance.With().Caller().Logger()
	var orderList OrderListScheme
	stmt := `
	SELECT id, user_id, number, status, accrual, uploaded_at
	FROM orders
	WHERE user_id = $1
	ORDER BY uploaded_at DESC;
	`

	rows, err := r.db.Query(ctx, stmt, userID)
	if err != nil {
		log.Error().Err(err).Msg("query order list")
		return orderList, err
	}
	defer rows.Close()

	for rows.Next() {
		err = orderList.ScanRow(rows)
		if err != nil {
			log.Error().Err(err).Msg("scanning row")
			return orderList, err
		}
	}
	return orderList, rows.Err()
}

func (r OrdersRepository) UpdateAccrual(ctx context.Context, orders []OrderScheme) error {
	log := logger.Instance.With().Caller().Logger()
	tx, err := r.db.Begin(ctx)
	if err != nil {
		log.Error().Err(err).Msg("begin tx")
		return err
	}

	defer func() {
		if err != nil {
			if err := tx.Rollback(ctx); err != nil {
				log.Error().Err(err).Msg("rollback tx")
			}
			return
		}
		if err := tx.Commit(ctx); err != nil {
			log.Error().Err(err).Msg("commit tx")
		}
	}()

	batch := &pgx.Batch{}

	var withAccrual []OrderScheme

	for _, order := range orders {
		stmt := `
		UPDATE orders
		SET
			status = $2,
			accrual = $3,
			last_update = CURRENT_TIMESTAMP
		WHERE number = $1;
		`
		batch.Queue(stmt, order.Number, order.Status, order.Accrual)
		if order.Accrual != 0 {
			withAccrual = append(withAccrual, order)
		}
	}

	for _, order := range withAccrual {
		stmt := `
		WITH transaction AS (
			INSERT INTO bonus_transactions (
				user_id, order_number, transaction_type, transaction_amount
			)
			VALUES ($1, $2, 'A', $3)
			RETURNING user_id
		)
		UPDATE bonus_accounts
		SET
			balance = balance + $3, last_update = CURRENT_TIMESTAMP
		WHERE user_id = (SELECT user_id FROM transaction);
		`
		batch.Queue(stmt, order.OwnerID, order.Number, order.Accrual)
	}

	err = tx.SendBatch(ctx, batch).Close()
	return err
}
