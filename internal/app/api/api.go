package api

import (
	"github.com/gin-gonic/gin"
	"github.com/romm80/gophermart.git/internal/app/service"
	"net/http"
)

const authHeader = "Authorization"

type API struct {
	Services *service.Services
}

func NewAPI(services *service.Services) *API {
	return &API{Services: services}
}

func (a *API) mainHandler(c *gin.Context) {
	c.Status(http.StatusOK)
}

func (a *API) RoutingInit() *gin.Engine {
	api := gin.Default()

	api.GET("/", a.mainHandler)

	user := api.Group("/api/user")
	user.POST("/register", a.registerUser)
	user.POST("/login", a.loginUser)

	user.Use(a.authMiddleware)

	orders := user.Group("/orders")
	orders.POST("", a.uploadOrder)
	orders.GET("", a.getOrders)

	balances := user.Group("/balance")
	balances.GET("/", a.getBalance)
	balances.POST("/withdraw", a.withdraw)
	balances.GET("/withdrawals", a.withdrawals)

	return api
}
