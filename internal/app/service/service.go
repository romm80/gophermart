package service

import (
	"github.com/romm80/gophermart.git/internal/app/models"
	"strconv"
)

type AuthService interface {
	CreateUser(user models.User) (string, error)
	LoginUser(user models.User) (string, error)
	ParseToken(token string) (int, error)
	ValidUserID(userID int) error
}

type OrdersService interface {
	UploadOrder(userID int, order string) error
	UpdateOrder(order models.AccrualOrder) error
	GetOrders(userID int) ([]models.Order, error)
}

type BalancesService interface {
	CurrentBalance(userID int) (*models.CurrentBalance, error)
	Withdraw(userID int, order models.OrderBalance) error
	Withdrawals(userID int) ([]models.OrderBalance, error)
	Accrual(userID int, order models.AccrualOrder) error
}

type Services struct {
	AuthService
	OrdersService
	BalancesService
}

func checkLuhn(order string) bool {
	sum := 0
	f := false
	for i := len(order) - 1; i >= 0; i-- {
		num, _ := strconv.Atoi(string(order[i]))
		if f {
			num *= 2
			if num > 9 {
				num -= 9
			}
		}
		f = !f
		sum += num
	}
	return sum%10 == 0
}
