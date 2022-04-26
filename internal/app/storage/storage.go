package storage

import (
	"github.com/romm80/gophermart.git/internal/app/models"
)

type AuthStore interface {
	CreateUser(user models.User) error
	GetUser(user models.User) error
}

type OrdersStore interface {
	AddOrder(order models.Order) error
	UpdateOrder(order models.AccrualOrder) error
	GetOrders(user string) ([]models.Order, error)
}

type BalancesStore interface {
	CurrentBalance(user string) (*models.CurrentBalance, error)
	Withdraw(user string, order models.OrderBalance) error
	Withdrawals(user string) ([]models.OrderBalance, error)
	Accrual(user string, order models.AccrualOrder) error
}

type Storage struct {
	AuthStore
	OrdersStore
	BalancesStore
}
