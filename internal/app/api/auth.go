package api

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/romm80/gophermart.git/internal/app"
	"github.com/romm80/gophermart.git/internal/app/models"
	"net/http"
)

func (a *API) registerUser(c *gin.Context) {
	var user models.User
	if err := json.NewDecoder(c.Request.Body).Decode(&user); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	token, err := a.Services.CreateUser(user)
	if err != nil {
		c.AbortWithStatus(app.ErrStatusCode(err))
		return
	}

	c.Header(authHeader, token)
	c.Status(http.StatusOK)
}

func (a *API) loginUser(c *gin.Context) {
	var user models.User
	if err := json.NewDecoder(c.Request.Body).Decode(&user); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	token, err := a.Services.LoginUser(user)
	if err != nil {
		c.AbortWithStatus(app.ErrStatusCode(err))
		return
	}

	c.Header(authHeader, token)
	c.Status(http.StatusOK)
}
