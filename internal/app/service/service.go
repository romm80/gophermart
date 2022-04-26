package service

import (
	"github.com/romm80/gophermart.git/internal/app/models"
	"strconv"
)

type AuthService interface {
	CreateUser(user models.User) error
	GenerateToken(user models.User) (string, error)
	ParseToken(token string) (string, error)
}

type OrdersService interface {
	UploadOrder(user, order string) error
	GetOrders(user string) ([]models.Order, error)
}

type BalancesService interface {
	CurrentBalance(user string) (*models.CurrentBalance, error)
	Withdraw(user string, order models.OrderBalance) error
	Withdrawals(user string) ([]models.OrderBalance, error)
}

type Services struct {
	AuthService
	OrdersService
	BalancesService
}

func checkLuhn(order string) bool {
	sum := 0
	f := false
	for _, s := range order {
		num, _ := strconv.Atoi(string(s))
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
