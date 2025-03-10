package orders

import (
	"context"
	"errors"
	"time"

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

func (r OrdersRepository) Create(ctx context.Context, userID int32, orderNumber string) error {
	stmt := `
	INSERT INTO orders (user_id, number) VALUES ($1, $2);
	`

	_, err := r.db.Exec(ctx, stmt, userID, orderNumber)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return errs.ErrOrderUploaded
		}
		logger.Instance.Error().Err(err).Caller().Msg("creating order")
		return err
	}

	return nil
}

func (r OrdersRepository) ReadByOrderNumber(
	ctx context.Context, orderNumber string,
) (order OrderScheme, err error) {
	stmt := `
	WITH certain_order AS (
	    SELECT id, user_id, status_id, number, accrual, uploaded_at
		FROM orders
		WHERE number = $1
	)
	SELECT o.id, o.user_id, o.number, s.name AS status, o.accrual, o.uploaded_at
	FROM certain_order AS o
	JOIN order_status AS s ON o.status_id = s.id;
	`

	err = r.db.QueryRow(ctx, stmt, orderNumber).Scan(
		&order.ID,
		&order.OwnerID,
		&order.Number,
		&order.Status,
		&order.Accrual,
		&order.UploadedAt,
	)
	if err != nil {
		logger.Instance.Error().Err(err).Caller().Msg("reading order by number")
		return
	}
	return
}

func (r OrdersRepository) ReadListByUser(
	ctx context.Context, userID int32, orders []OrderScheme,
) ([]OrderScheme, error) {
	stmt := `
	WITH certain_order AS (
	    SELECT id, user_id, status_id, number, accrual, uploaded_at
		FROM orders
		WHERE user_id = $1
	)
	SELECT o.id, o.user_id, o.number, s.name AS status, o.accrual, o.uploaded_at
	FROM certain_order AS o
	JOIN order_status AS s ON o.status_id = s.id
	ORDER BY o.uploaded_at DESC;
	`

	rows, err := r.db.Query(ctx, stmt, userID)
	if err != nil {
		logger.Instance.Error().Err(err).Caller().Msg("reading order list by user")
		return orders, err
	}
	var (
		id         int32
		ownerID    int32
		number     string
		status     string
		accrual    float64
		uploadetAt time.Time
	)

	scanRowFn := func() error {
		orders = append(
			orders,
			OrderScheme{id, ownerID, number, status, accrual, uploadetAt},
		)
		return nil
	}

	_, err = pgx.ForEachRow(
		rows,
		[]any{&id, &ownerID, &number, &status, &accrual, &uploadetAt},
		scanRowFn,
	)

	if err != nil {
		logger.Instance.Error().Err(err).Caller().Msg("reading order list by user")
		return orders, err
	}

	return orders, nil
}
