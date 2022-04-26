package api

import (
	"encoding/json"
	"errors"
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

	err := a.Services.CreateUser(user)
	if err != nil {
		if errors.Is(err, app.ErrLoginIsUsed) {
			c.AbortWithStatus(http.StatusConflict)
			return
		}
		if errors.Is(err, app.ErrInvalidLoginOrPassword) {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	token, err := a.Services.GenerateToken(user)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
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

	token, err := a.Services.GenerateToken(user)
	if err != nil {
		if errors.Is(err, app.ErrInvalidLoginOrPassword) {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.Header(authHeader, token)
	c.Status(http.StatusOK)
}
