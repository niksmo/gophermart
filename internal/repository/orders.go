package repository

import "github.com/jackc/pgx/v5/pgxpool"

type OrdersRepository struct {
	db *pgxpool.Pool
}

func Orders(db *pgxpool.Pool) OrdersRepository {
	return OrdersRepository{db: db}
}
