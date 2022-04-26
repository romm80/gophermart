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
	User  string
	Order string
}

type AccrualWorker struct {
	Tasks chan Task
}

func NewDeleteWorker(size int) *AccrualWorker {
	return &AccrualWorker{
		Tasks: make(chan Task, size),
	}
}

func (r *AccrualWorker) Run(services *service.Services) {
	go func() {
		for {
			for task := range r.Tasks {
				client := http.Client{
					Timeout: 5 * time.Second,
				}
				resp, err := client.Get(fmt.Sprintf("%s/%s/%s", server.CFG.Accrual, "/api/orders/", task.Order))
				if err != nil {
					log.Println(err)
				}

				if resp.StatusCode == http.StatusInternalServerError || resp.StatusCode == http.StatusTooManyRequests {
					time.Sleep(time.Minute)
					r.Add(task.User, task.Order)
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
				if err := services.BalancesService.Accrual(task.User, order); err != nil {
					log.Println(err)
					resp.Body.Close()
					continue
				}
				resp.Body.Close()
				if order.Status != models.PROCESSED && order.Status != models.INVALID {
					r.Add(task.User, task.Order)
				}
			}
		}
	}()
}

func (r *AccrualWorker) Add(user, order string) {
	go func(user, order string) {
		r.Tasks <- Task{
			User:  user,
			Order: order,
		}
	}(user, order)
}
