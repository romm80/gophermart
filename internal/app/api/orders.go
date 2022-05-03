package api

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/romm80/gophermart.git/internal/app"
	"github.com/romm80/gophermart.git/internal/app/service/workers"
	"io/ioutil"
	"net/http"
)

func (a *API) uploadOrder(c *gin.Context) {
	order, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	userID := c.GetInt("user_id")
	if err := a.Services.AuthService.ValidUserID(userID); err != nil {
		c.AbortWithStatus(app.ErrStatusCode(err))
		return
	}

	err = a.Services.OrdersService.UploadOrder(userID, string(order))
	if err != nil {
		if errors.Is(err, app.ErrOrderUploaded) {
			c.Status(http.StatusOK)
			return
		}
		c.AbortWithStatus(app.ErrStatusCode(err))
		return
	}

	a.AccrualWorker.QueueTask(workers.Task{
		UserID: userID,
		Order:  string(order),
	})
	c.Status(http.StatusAccepted)
}

func (a *API) getOrders(c *gin.Context) {
	userID := c.GetInt("user_id")
	if err := a.Services.AuthService.ValidUserID(userID); err != nil {
		c.AbortWithStatus(app.ErrStatusCode(err))
		return
	}

	orders, err := a.Services.GetOrders(userID)
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
