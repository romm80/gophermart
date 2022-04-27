package api

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/romm80/gophermart.git/internal/app"
	"github.com/romm80/gophermart.git/internal/app/models"
	"log"
	"net/http"
)

func (a *API) getBalance(c *gin.Context) {
	balance, err := a.Services.CurrentBalance(c.GetString("user"))
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	type bal struct {
		Current   float64 `json:"current"`
		Withdrawn float64 `json:"withdrawn"`
	}
	b := bal{}
	b.Current, _ = balance.Current.Float64()
	b.Withdrawn, _ = balance.Withdrawn.Float64()

	log.Println("getBalance", b)
	c.JSON(http.StatusOK, b)
}

func (a *API) withdraw(c *gin.Context) {
	var order models.OrderBalance
	if err := json.NewDecoder(c.Request.Body).Decode(&order); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	err := a.Services.BalancesService.Withdraw(c.GetString("user"), order)

	if err != nil {
		if errors.Is(err, app.ErrNotEnoughFunds) {
			c.AbortWithStatus(http.StatusPaymentRequired)
			return
		}
		if errors.Is(err, app.ErrInvalidOrderFormat) {
			c.AbortWithStatus(http.StatusUnprocessableEntity)
			return
		}
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.Status(http.StatusOK)
}

func (a *API) withdrawals(c *gin.Context) {
	orders, err := a.Services.Withdrawals(c.GetString("user"))
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	if len(orders) == 0 {
		c.Status(http.StatusNoContent)
		return
	}
	c.JSON(http.StatusOK, orders)
}
