package service

import (
	"github.com/romm80/gophermart.git/internal/app"
	"github.com/romm80/gophermart.git/internal/app/models"
	"github.com/romm80/gophermart.git/internal/app/storage"
	"strconv"
)

type Orders struct {
	store storage.OrdersStore
}

func NewOrders(store storage.OrdersStore) *Orders {
	return &Orders{store: store}
}

func (o *Orders) UploadOrder(userID int, number string) error {
	_, err := strconv.Atoi(number)
	if err != nil {
		return app.ErrInvalidRequestFormat
	}
	if !checkLuhn(number) {
		return app.ErrInvalidOrderFormat
	}

	order := models.Order{
		Number: number,
		Status: models.OrderStatusNew,
		UserID: userID,
	}

	return o.store.AddOrder(order)
}

func (o *Orders) GetOrders(userID int) ([]models.Order, error) {
	return o.store.GetOrders(userID)
}

func (o *Orders) UpdateOrder(order models.AccrualOrder) error {
	return o.store.UpdateOrder(order)
}
