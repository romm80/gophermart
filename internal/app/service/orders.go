package service

import (
	"github.com/romm80/gophermart.git/internal/app"
	"github.com/romm80/gophermart.git/internal/app/models"
	"github.com/romm80/gophermart.git/internal/app/storage"
	"strconv"
)

const (
	NEW        = "NEW"
	PROCESSING = "PROCESSING"
	INVALID    = "INVALID"
	PROCESSED  = "PROCESSED"
)

type Orders struct {
	store storage.OrdersStore
}

func NewOrders(store storage.OrdersStore) *Orders {
	return &Orders{store: store}
}

func (o *Orders) UploadOrder(user, number string) error {
	_, err := strconv.Atoi(number)
	if err != nil {
		return app.ErrInvalidRequestFormat
	}
	if !checkLuhn(number) {
		return app.ErrInvalidOrderFormat
	}
	order := models.Order{
		Number: number,
		Status: NEW,
		User:   user,
	}

	return o.store.AddOrder(order)
}

func (o *Orders) GetOrders(user string) ([]models.Order, error) {
	return o.store.GetOrders(user)
}
