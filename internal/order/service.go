package order

import "github.com/niksmo/gophermart/internal/repository"

type OrderService struct {
	repository repository.OrdersRepository
}

func NewService(repository repository.OrdersRepository) OrderService {
	return OrderService{repository: repository}
}
