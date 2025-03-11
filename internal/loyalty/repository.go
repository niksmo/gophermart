package loyalty

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/niksmo/gophermart/pkg/logger"
)

const (
	tAdd      = "A"
	tWithdraw = "W"
)

type LoyaltyRepository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) LoyaltyRepository {
	return LoyaltyRepository{db: db}
}

func (r LoyaltyRepository) Create(ctx context.Context, userID int32) error {
	stmt := `
	INSERT INTO bonus_accounts (user_id) VALUES ($1);
	`
	_, err := r.db.Exec(ctx, stmt, userID)
	if err != nil {
		logger.Instance.Error().Err(err).Caller().Msg("creating bonus account")
		return err
	}
	return nil
}

func (r LoyaltyRepository) Read(ctx context.Context, userID int32) (BalanceScheme, error) {
	stmt := `
	SELECT id, user_id, balance, withdraw, last_update
	FROM bonus_accounts
	WHERE user_id=$1;
	`
	var balance BalanceScheme
	err := r.db.QueryRow(ctx, stmt, userID).Scan(
		&balance.ID,
		&balance.OwnerID,
		&balance.Balance,
		&balance.Withdraw,
		&balance.LastUpdate,
	)
	if err != nil {
		logger.Instance.Error().
			Err(err).
			Caller().
			Int32("userID", userID).
			Msg("reading user balance")
		return balance, err
	}
	return balance, nil
}

func (r LoyaltyRepository) Update(ctx context.Context) error {
	return nil
}
