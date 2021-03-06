package api

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/romm80/gophermart.git/internal/app"
	"github.com/romm80/gophermart.git/internal/app/models"
	"net/http"
)

func (a *API) getBalance(c *gin.Context) {
	userID := c.GetInt("user_id")
	if userID <= 0 {
		c.AbortWithStatus(app.ErrStatusCode(app.ErrInvalidUserID))
		return
	}

	balance, err := a.Services.CurrentBalance(userID)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, balance)
}

func (a *API) withdraw(c *gin.Context) {
	var order models.OrderBalance
	if err := json.NewDecoder(c.Request.Body).Decode(&order); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	userID := c.GetInt("user_id")
	if userID <= 0 {
		c.AbortWithStatus(app.ErrStatusCode(app.ErrInvalidUserID))
		return
	}

	if err := a.Services.OrdersService.UploadOrder(userID, order.Order); err != nil {
		if errors.Is(err, app.ErrOrderUploaded) {
			c.AbortWithStatus(http.StatusUnprocessableEntity)
			return
		}
		c.AbortWithStatus(app.ErrStatusCode(err))
		return
	}

	if err := a.Services.BalancesService.Withdraw(userID, order); err != nil {
		c.AbortWithStatus(app.ErrStatusCode(err))
		return
	}
	c.Status(http.StatusOK)
}

func (a *API) withdrawals(c *gin.Context) {
	userID := c.GetInt("user_id")
	if userID <= 0 {
		c.AbortWithStatus(app.ErrStatusCode(app.ErrInvalidUserID))
		return
	}

	orders, err := a.Services.Withdrawals(userID)
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
