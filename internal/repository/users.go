package repository

import (
	"context"

	"github.com/niksmo/gophermart/pkg/sqldb"
)

type UsersRepository struct {
	dbService sqldb.DBService
}

func Users(dbService sqldb.DBService) UsersRepository {
	return UsersRepository{dbService: dbService}
}

func (repo UsersRepository) Create(
	ctx context.Context, login string, password string,
) {
}

func (repo UsersRepository) ReadByLogin(ctx context.Context, login string) {}

func (repo UsersRepository) ReadByID(ctx context.Context, userID int64) {}
