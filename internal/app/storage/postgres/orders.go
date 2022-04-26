package postgres

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/romm80/gophermart.git/internal/app"
	"github.com/romm80/gophermart.git/internal/app/models"
)

type OrdersDB struct {
	pool *pgxpool.Pool
}

func NewOrdersDB(pool *pgxpool.Pool) *OrdersDB {
	return &OrdersDB{pool: pool}
}

func (o *OrdersDB) AddOrder(order models.Order) error {
	ctx := context.Background()
	conn, err := o.pool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	tx, err := conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	rows, err := tx.Query(ctx, `SELECT "user", "number" FROM orders WHERE "number" = ($1)`, order.Number)
	if err != nil {
		return err
	}
	if rows.Next() {
		var usr, num string
		if err := rows.Scan(&usr, &num); err != nil {
			return err
		}
		if num == order.Number {
			if usr == order.User {
				return app.ErrOrderUploaded
			}
			return app.ErrOrderUploadedAnotherUser
		}
	}

	if _, err := tx.Exec(ctx, `INSERT INTO orders ("user", "number", status) VALUES($1, $2, $3)`, order.User, order.Number, order.Status); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (o *OrdersDB) GetOrders(user string) ([]models.Order, error) {
	ctx := context.Background()
	conn, err := o.pool.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	rows, err := conn.Query(ctx, `SELECT "number", status, accrual, uploaded_at FROM orders WHERE "user" = ($1)`, user)
	if err != nil {
		return nil, err
	}
	orders := make([]models.Order, 0)
	for rows.Next() {
		order := &models.Order{}
		if err := rows.Scan(&order.Number, &order.Status, &order.Accrual, &order.UploadedAt.Time); err != nil {
			return nil, err
		}
		orders = append(orders, *order)
	}
	return orders, nil
}
