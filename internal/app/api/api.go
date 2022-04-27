package api

import (
	"github.com/gin-gonic/gin"
	"github.com/romm80/gophermart.git/internal/app/service"
	"github.com/romm80/gophermart.git/internal/app/service/workers"
	"net/http"
)

const authHeader = "Authorization"

type API struct {
	Services      *service.Services
	AccrualWorker *workers.AccrualWorker
}

func NewAPI(services *service.Services) *API {
	worker := workers.NewAccrualWorker(1000)
	worker.Run(services)

	return &API{
		Services:      services,
		AccrualWorker: worker,
	}
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
	user.Use(a.gzipMiddleware)

	orders := user.Group("/orders")
	orders.POST("", a.uploadOrder)
	orders.GET("", a.getOrders)

	balances := user.Group("/balance")
	balances.GET("/", a.getBalance)
	balances.POST("/withdraw", a.withdraw)
	balances.GET("/withdrawals", a.withdrawals)

	return api
}
