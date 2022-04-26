package api

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/romm80/gophermart.git/internal/app"
	"io/ioutil"
	"net/http"
)

func (a *API) uploadOrder(c *gin.Context) {
	order, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	err = a.Services.OrdersService.UploadOrder(c.GetString("user"), string(order))
	if err != nil {
		if errors.Is(err, app.ErrInvalidRequestFormat) {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		if errors.Is(err, app.ErrInvalidOrderFormat) {
			c.AbortWithStatus(http.StatusUnprocessableEntity)
			return
		}
		if errors.Is(err, app.ErrOrderUploadedAnotherUser) {
			c.AbortWithStatus(http.StatusConflict)
			return
		}
		if errors.Is(err, app.ErrOrderUploaded) {
			c.Status(http.StatusOK)
			return
		}
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.Status(http.StatusAccepted)
}

func (a *API) getOrders(c *gin.Context) {
	orders, err := a.Services.GetOrders(c.GetString("user"))
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
