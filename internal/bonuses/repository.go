package bonuses

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/niksmo/gophermart/pkg/logger"
)

const (
	addT      = "A"
	withdrawT = "W"
)

type BonusesRepository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) BonusesRepository {
	return BonusesRepository{db: db}
}

func (r BonusesRepository) CreateAccount(ctx context.Context, userID int32) error {
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
