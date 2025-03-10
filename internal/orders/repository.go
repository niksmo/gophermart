package orders

import (
	"context"
	"errors"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/niksmo/gophermart/internal/errs"
)

type OrdersRepository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) OrdersRepository {
	return OrdersRepository{db: db}
}

func (r OrdersRepository) Create(ctx context.Context, userID int64, orderNumber int64) error {
	stmt := `
	INSERT INTO orders (user_id, number) VALUES ($1, $2);
	`

	_, err := r.db.Exec(ctx, stmt, userID, orderNumber)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return errs.ErrOrderUploaded
		}
		return err
	}

	return nil
}

func (r OrdersRepository) ReadByOrderNumber(ctx context.Context, orderNumber int64) error {
	//user_id, orderNumber, status, accrual, uploadedAt
	return nil
}
