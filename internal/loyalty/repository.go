package loyalty

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/niksmo/gophermart/internal/errs"
	"github.com/niksmo/gophermart/pkg/database"
	"github.com/niksmo/gophermart/pkg/logger"
)

const (
	tWithdraw = "W"
)

type LoyaltyRepository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) LoyaltyRepository {
	return LoyaltyRepository{db: db}
}

func (r LoyaltyRepository) CreateAccount(
	ctx context.Context, userID int32,
) error {
	stmt := `
	INSERT INTO bonus_accounts (user_id) VALUES ($1);
	`
	_, err := r.db.Exec(ctx, stmt, userID)
	if err != nil {
		logger.Instance.Error().
			Err(err).
			Caller().
			Msg("creating bonus account")
		return err
	}
	return nil
}

func (r LoyaltyRepository) ReadBalance(
	ctx context.Context, userID int32,
) (BalanceScheme, error) {
	stmt := `
	SELECT id, user_id, balance, withdraw, last_update
	FROM bonus_accounts
	WHERE user_id=$1;
	`
	var balance BalanceScheme
	err := balance.ScanRow(r.db.QueryRow(ctx, stmt, userID))
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

func (r LoyaltyRepository) ReduceBalance(
	ctx context.Context, userID int32, orderNumber string, amount float32,
) error {
	log := logger.Instance.With().
		Caller().
		Int32("userID", userID).
		Str("orderNumber", orderNumber).
		Float32("amount", amount).
		Logger()

	tx, err := r.db.Begin(ctx)
	if err != nil {
		log.Error().Err(err).Msg("beginning tx")
	}

	currentBalance, err := selectCurrentBalance(ctx, tx, userID)
	if err != nil {
		log.Error().Err(err).Msg("selecting current balance")
		return database.CloseTX(ctx, tx, err, log)
	}

	if currentBalance < amount {
		return database.CloseTX(
			ctx, tx, errs.ErrLoyaltyNotEnoughPoints, log,
		)
	}

	err = insertBonusTransaction(ctx, tx, userID, orderNumber, amount)
	if err != nil {
		log.Error().Err(err).Msg("inserting bonus transaction")
		return database.CloseTX(ctx, tx, err, log)
	}

	err = updateCurrentBalance(ctx, tx, userID, amount)
	if err != nil {
		log.Error().Err(err).Msg("updating user account")
		return database.CloseTX(ctx, tx, err, log)
	}

	return database.CloseTX(ctx, tx, nil, log)
}

func (r LoyaltyRepository) ReadWithdrawals(
	ctx context.Context, userID int32,
) (WithdrawalsScheme, error) {
	log := logger.Instance.With().
		Caller().
		Int32("userID", userID).
		Logger()

	stmt := `
	SELECT order_number, transaction_amount, processed_at
	FROM bonus_transactions
	WHERE user_id=$1 AND transaction_type=$2
	ORDER BY processed_at DESC;
	`
	rows, err := r.db.Query(ctx, stmt, userID, tWithdraw)
	if err != nil {
		log.Error().Err(err).Msg("selecting loyalty account transactions")
	}
	defer rows.Close()

	var withdrawals WithdrawalsScheme
	for rows.Next() {
		err = withdrawals.ScanRow(rows)
		if err != nil {
			log.Error().Err(err).Msg("scanning row")
			return withdrawals, err
		}
	}

	return withdrawals, rows.Err()
}

func selectCurrentBalance(
	ctx context.Context, tx pgx.Tx, userID int32,
) (float32, error) {
	stmt := `
	SELECT balance
	FROM bonus_accounts 
	WHERE user_id=$1
	FOR UPDATE;
	`
	var currentBalance float32
	err := tx.QueryRow(ctx, stmt, userID).Scan(&currentBalance)
	return currentBalance, err
}

func insertBonusTransaction(
	ctx context.Context, tx pgx.Tx, userID int32, orderNumber string, amount float32,
) error {
	stmt := `
    INSERT INTO bonus_transactions (
        user_id, order_number, transaction_type, transaction_amount
    )
    VALUES (
        $1, $2, $3, $4
    );
    `
	_, err := tx.Exec(ctx, stmt, userID, orderNumber, tWithdraw, amount)
	return err
}

func updateCurrentBalance(
	ctx context.Context, tx pgx.Tx, userID int32, amount float32,
) error {
	stmt := `
	UPDATE bonus_accounts
	SET
        balance=balance-$2,
		withdraw=withdraw+$2,
		last_update=CURRENT_TIMESTAMP
	WHERE user_id=$1;
	`
	_, err := tx.Exec(ctx, stmt, userID, amount)
	return err
}
