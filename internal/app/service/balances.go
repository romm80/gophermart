package service

import (
	"github.com/romm80/gophermart.git/internal/app"
	"github.com/romm80/gophermart.git/internal/app/models"
	"github.com/romm80/gophermart.git/internal/app/storage"
)

type Balances struct {
	store storage.BalancesStore
}

func NewBalances(store storage.BalancesStore) *Balances {
	return &Balances{store: store}
}

func (b *Balances) CurrentBalance(userID int) (*models.CurrentBalance, error) {
	return b.store.CurrentBalance(userID)
}

func (b *Balances) Withdraw(userID int, order models.OrderBalance) error {
	if !checkLuhn(order.Order) {
		return app.ErrInvalidOrderFormat
	}
	return b.store.Withdraw(userID, order)
}

func (b *Balances) Withdrawals(userID int) ([]models.OrderBalance, error) {
	return b.store.Withdrawals(userID)
}

func (b *Balances) Accrual(userID int, order models.AccrualOrder) error {
	return b.store.Accrual(userID, order)
}
