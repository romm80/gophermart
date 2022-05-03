package storage

import (
	"github.com/romm80/gophermart.git/internal/app/models"
)

type AuthStore interface {
	CreateUser(user models.User) (int, error)
	GetUserID(user models.User) (int, error)
	ValidUserID(userID int) error
}

type OrdersStore interface {
	AddOrder(order models.Order) error
	UpdateOrder(order models.AccrualOrder) error
	GetOrders(userID int) ([]models.Order, error)
}

type BalancesStore interface {
	CurrentBalance(userID int) (*models.CurrentBalance, error)
	Withdraw(userID int, order models.OrderBalance) error
	Withdrawals(userID int) ([]models.OrderBalance, error)
	Accrual(userID int, order models.AccrualOrder) error
}

type Storage struct {
	AuthStore
	OrdersStore
	BalancesStore
}
