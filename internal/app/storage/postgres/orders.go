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

	rows, err := tx.Query(ctx, `SELECT user_id, "number" FROM orders WHERE "number" = ($1)`, order.Number)
	if err != nil {
		return err
	}
	if rows.Next() {
		var usr int
		var num string
		if err := rows.Scan(&usr, &num); err != nil {
			return err
		}
		if num == order.Number {
			if usr == order.UserID {
				return app.ErrOrderUploaded
			}
			return app.ErrOrderUploadedAnotherUser
		}
	}

	if _, err := tx.Exec(ctx, `INSERT INTO orders (user_id, "number", status) VALUES($1, $2, $3)`, order.UserID, order.Number, order.Status); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (o *OrdersDB) GetOrders(userID int) ([]models.Order, error) {
	ctx := context.Background()
	conn, err := o.pool.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	rows, err := conn.Query(ctx, `SELECT "number", status, accrual, uploaded_at FROM orders WHERE user_id = ($1)`, userID)
	if err != nil {
		return nil, err
	}
	orders := make([]models.Order, 0)
	for rows.Next() {
		order := models.Order{}
		if err := rows.Scan(&order.Number, &order.Status, &order.Accrual, &order.UploadedAt.Time); err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	return orders, nil
}

func (o *OrdersDB) UpdateOrder(order models.AccrualOrder) error {
	ctx := context.Background()
	conn, err := o.pool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	_, err = conn.Exec(ctx, `UPDATE orders SET status=($2), accrual=($3) WHERE "number"=($1)`, order.Order, order.Status, order.Accrual)

	return err
}
