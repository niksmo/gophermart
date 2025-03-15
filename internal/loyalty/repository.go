package loyalty

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/niksmo/gophermart/internal/errs"
	"github.com/niksmo/gophermart/pkg/logger"
	"github.com/rs/zerolog/log"
)

const (
	T_ADD      = "A"
	T_WITHDRAW = "W"
)

type LoyaltyRepository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) LoyaltyRepository {
	return LoyaltyRepository{db: db}
}

func (r LoyaltyRepository) CreateAccount(ctx context.Context, userID int32) error {
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

func (r LoyaltyRepository) CreateAddTransactions(
	ctx context.Context, transactions []TransactionScheme,
) error {
	logger.Instance.With().Caller().Logger()

	tx, err := r.db.Begin(ctx)
	if err != nil {
		log.Error().Err(err).Msg("beginning tx")
		return err
	}

	defer func() {
		if err != nil {
			if err := tx.Rollback(ctx); err != nil {
				log.Error().Err(err).Msg("rollback tx")
			}
			return
		}
		err = tx.Commit(ctx)
		if err != nil {
			log.Error().Err(err).Msg("committing tx")
		}
	}()

	batch := &pgx.Batch{}
	for _, t := range transactions {
		stmt := `
		WITH transaction AS (
			INSERT INTO bonus_transactions (
				user_id, order_number, transaction_type, transaction_amount
			)
			VALUES ($1, $2, $3, $4)
			RETURNING user_id
		)
		UPDATE bonus_accounts
		SET
			balance = balance + $4, last_update = CURRENT_TIMESTAMP
		WHERE user_id = (SELECT user_id FROM transaction);
		`
		batch.Queue(stmt, t.UserID, t.OrderNumber, T_ADD, t.Amount)
	}
	err = tx.SendBatch(ctx, batch).Close()

	return err
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
	ctx context.Context, userID int32, orderNumber string, amount float64,
) error {
	log := logger.Instance.With().
		Caller().
		Int32("userID", userID).
		Str("orderNumber", orderNumber).
		Float64("amount", amount).
		Logger()

	tx, err := r.db.Begin(ctx)
	if err != nil {
		log.Error().Err(err).Msg("beginning tx")
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

	stmt := `
	SELECT balance
	FROM bonus_accounts 
	WHERE user_id=$1
	FOR UPDATE;
	`
	var current float64
	err = tx.QueryRow(ctx, stmt, userID).Scan(&current)
	if err != nil {
		log.Error().Err(err).Msg("selecting current balance")
		return err
	}

	if current < amount {
		return errs.ErrLoyaltyNotEnoughPoints
	}

	stmt = `
	INSERT INTO bonus_transactions (
	user_id, order_number, transaction_type, transaction_amount
	)
	VALUES (
	$1, $2, $3, $4
	);
	`
	_, err = tx.Exec(ctx, stmt, userID, orderNumber, T_WITHDRAW, amount)
	if err != nil {
		log.Error().Err(err).Msg("inserting bonus transaction")
		return err
	}

	stmt = `
	UPDATE bonus_accounts
	SET
        balance=$2,
		withdraw=withdraw+$3,
		last_update=CURRENT_TIMESTAMP
	WHERE user_id=$1;
	`
	_, err = tx.Exec(ctx, stmt, userID, current-amount, amount)
	if err != nil {
		log.Error().Err(err).Msg("updating user account")
		return err
	}

	return nil
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
	rows, err := r.db.Query(ctx, stmt, userID, T_WITHDRAW)
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
