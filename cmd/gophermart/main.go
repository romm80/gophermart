package main

import (
	"flag"
	"github.com/caarlos0/env/v6"
	"github.com/romm80/gophermart.git/internal/app/api"
	"github.com/romm80/gophermart.git/internal/app/server"
	"github.com/romm80/gophermart.git/internal/app/service"
	"github.com/romm80/gophermart.git/internal/app/storage"
	"github.com/romm80/gophermart.git/internal/app/storage/postgres"
	"log"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	if err := env.Parse(&server.CFG); err != nil {
		log.Fatal(err)
	}
	flag.StringVar(&server.CFG.Host, "a", server.CFG.Host, "Server address")
	flag.StringVar(&server.CFG.DB, "d", server.CFG.DB, "Database address")
	flag.StringVar(&server.CFG.Accrual, "r", server.CFG.Accrual, "Accrual address")
	flag.Parse()
	server.CFG.Key = []byte("secret_key")

	dbpool, err := postgres.NewPostgresDB()
	if err != nil {
		log.Fatal(err)
	}

	store := &storage.Storage{
		AuthStore:     postgres.NewAuthDB(dbpool),
		OrdersStore:   postgres.NewOrdersDB(dbpool),
		BalancesStore: postgres.NewBalancesDB(dbpool),
	}
	services := &service.Services{
		AuthService:     service.NewAuth(store.AuthStore),
		OrdersService:   service.NewOrders(store.OrdersStore),
		BalancesService: service.NewBalances(store.BalancesStore),
	}
	handlers := api.NewAPI(services)

	srv, err := server.New(handlers.RoutingInit())
	if err != nil {
		log.Fatal(err)
	}
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
