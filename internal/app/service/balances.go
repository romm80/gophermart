package service

import (
	"github.com/romm80/gophermart.git/internal/app"
	"github.com/romm80/gophermart.git/internal/app/models"
	"github.com/romm80/gophermart.git/internal/app/storage"
)

const (
	EXPENSE = "expense"
	RECEIPT = "receipt"
)

type Balances struct {
	store storage.BalancesStore
}

func NewBalances(store storage.BalancesStore) *Balances {
	return &Balances{store: store}
}

func (b *Balances) CurrentBalance(user string) (*models.CurrentBalance, error) {
	return b.store.CurrentBalance(user)
}

func (b *Balances) Withdraw(user string, order models.OrderBalance) error {
	if !checkLuhn(order.Order) {
		return app.ErrInvalidOrderFormat
	}

	return b.store.Withdraw(user, order)
}

func (b *Balances) Withdrawals(user string) ([]models.OrderBalance, error) {
	return b.store.Withdrawals(user)
}
