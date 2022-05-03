package workers

import (
	"encoding/json"
	"fmt"
	"github.com/romm80/gophermart.git/internal/app/models"
	"github.com/romm80/gophermart.git/internal/app/server"
	"github.com/romm80/gophermart.git/internal/app/service"
	"log"
	"net/http"
	"time"
)

type Task struct {
	UserID int
	Order  string
}

type AccrualWorker struct {
	Tasks      chan Task
	httpClient *http.Client
}

func NewAccrualWorker(size int) *AccrualWorker {
	return &AccrualWorker{
		Tasks:      make(chan Task, size),
		httpClient: &http.Client{Timeout: 5 * time.Second},
	}
}

func (r *AccrualWorker) Run(services *service.Services) {
	go func() {
		for task := range r.Tasks {
			resp, err := r.httpClient.Get(fmt.Sprintf("%s/%s/%s", server.CFG.Accrual, "api/orders", task.Order))
			if err != nil {
				log.Println(err)
			}

			if resp.StatusCode == http.StatusInternalServerError || resp.StatusCode == http.StatusTooManyRequests {
				time.Sleep(time.Minute)
				r.QueueTask(task)
				resp.Body.Close()
				continue
			}

			var order models.AccrualOrder
			if err := json.NewDecoder(resp.Body).Decode(&order); err != nil {
				log.Println(err)
				resp.Body.Close()
				continue
			}

			if err := services.OrdersService.UpdateOrder(order); err != nil {
				log.Println(err)
				resp.Body.Close()
				continue
			}
			if err := services.BalancesService.Accrual(task.UserID, order); err != nil {
				log.Println(err)
				resp.Body.Close()
				continue
			}
			resp.Body.Close()
			if order.Status != models.OrderStatusProcessed && order.Status != models.OrderStatusInvalid {
				r.QueueTask(task)
			}
		}
	}()
}

func (r *AccrualWorker) QueueTask(task Task) {
	go func() {
		r.Tasks <- task
	}()
}
