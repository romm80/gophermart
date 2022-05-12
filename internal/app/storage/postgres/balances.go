package postgres

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/romm80/gophermart.git/internal/app"
	"github.com/romm80/gophermart.git/internal/app/models"
	"github.com/shopspring/decimal"
	"time"
)

type BalancesDB struct {
	pool *pgxpool.Pool
}

var sqlInsertBalance = `INSERT INTO balances(processed_at, user_id, order_id, sum)
						SELECT ($1) AS processed_at, ($2) AS user_id, id AS order_id, ($4) AS sum FROM orders WHERE "number" = ($3);`

func NewBalancesDB(pool *pgxpool.Pool) *BalancesDB {
	return &BalancesDB{pool: pool}
}

func (b *BalancesDB) CurrentBalance(userID int) (*models.CurrentBalance, error) {
	sqlBalance := `SELECT SUM(current) AS current, SUM(withdrawn) AS withdrawn FROM
						(SELECT SUM(sum) AS current, 0 AS withdrawn
							FROM balances
						WHERE user_id = ($1)	
						UNION ALL
						SELECT 0, SUM(ABS(sum))
							FROM balances
						WHERE user_id = ($1) AND sum < 0) AS tmp;`

	ctx := context.Background()
	conn, err := b.pool.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	balance := &models.CurrentBalance{}
	err = conn.QueryRow(ctx, sqlBalance, userID).Scan(&balance.Current, &balance.Withdrawn)
	if err != nil {
		return nil, err
	}

	return balance, nil
}

func (b *BalancesDB) Withdraw(userID int, order models.OrderBalance) error {
	balance, err := b.CurrentBalance(userID)
	if err != nil {
		return err
	}

	if balance.Current.Sub(order.Sum.Decimal).LessThan(decimal.Zero) {
		return app.ErrNotEnoughFunds
	}

	ctx := context.Background()
	conn, err := b.pool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	_, err = conn.Exec(ctx, sqlInsertBalance, time.Now(), userID, order.Order, order.Sum.Neg())
	if err != nil {
		return err
	}
	return nil
}

func (b *BalancesDB) Withdrawals(userID int) ([]models.OrderBalance, error) {
	ctx := context.Background()
	conn, err := b.pool.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	sql := `SELECT orders.number, ABS(sum), processed_at FROM balances
				INNER JOIN orders ON balances.order_id = orders.id 
			WHERE balances.user_id = ($1) AND sum < 0;`
	rows, err := conn.Query(ctx, sql, userID)
	if err != nil {
		return nil, err
	}
	orders := make([]models.OrderBalance, 0)
	for rows.Next() {
		order := models.OrderBalance{}
		if err := rows.Scan(&order.Order, &order.Sum, &order.ProcessedAt.Time); err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	return orders, nil
}

func (b *BalancesDB) Accrual(userID int, order models.AccrualOrder) error {
	ctx := context.Background()
	conn, err := b.pool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	_, err = conn.Exec(ctx, sqlInsertBalance, time.Now(), userID, order.Order, order.Accrual)
	if err != nil {
		return err
	}
	return nil
}
