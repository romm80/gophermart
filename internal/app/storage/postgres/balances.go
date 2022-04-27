package postgres

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/romm80/gophermart.git/internal/app"
	"github.com/romm80/gophermart.git/internal/app/models"
	"github.com/shopspring/decimal"
	"log"
	"time"
)

type BalancesDB struct {
	pool *pgxpool.Pool
}

func NewBalancesDB(pool *pgxpool.Pool) *BalancesDB {
	return &BalancesDB{pool: pool}
}

func (b *BalancesDB) CurrentBalance(user string) (models.CurrentBalance, error) {
	sqlBalance := `SELECT SUM(current) AS current, SUM(withdrawn) AS withdrawn FROM
						(SELECT SUM(sum) AS current, 0 AS withdrawn
							FROM balances
						WHERE "user" = ($1)	
						UNION ALL
						SELECT 0, SUM(sum)
							FROM balances
						WHERE "user" = ($1) AND sum < 0) AS tmp;`

	balance := models.CurrentBalance{}
	ctx := context.Background()
	conn, err := b.pool.Acquire(ctx)
	if err != nil {
		return balance, err
	}
	defer conn.Release()

	err = conn.QueryRow(ctx, sqlBalance, user).Scan(&balance.Current, &balance.Withdrawn)
	if err != nil {
		return balance, err
	}
	balance.Withdrawn = -balance.Withdrawn

	return balance, nil
}

func (b *BalancesDB) Withdraw(user string, order models.OrderBalance) error {
	balance, err := b.CurrentBalance(user)
	if err != nil {
		return err
	}

	current := decimal.NewFromFloat(balance.Current)
	sum := decimal.NewFromFloat(order.Sum)
	if current.Sub(sum).LessThan(decimal.Zero) {
		return app.ErrNotEnoughFunds
	}

	ctx := context.Background()
	conn, err := b.pool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	_, err = conn.Exec(ctx, `INSERT INTO balances(processed_at, "user", "order", sum) VALUES ($1, $2, $3, $4);`, time.Now(), user, order.Order, sum.Neg())
	if err != nil {
		return err
	}
	return nil
}

func (b *BalancesDB) Withdrawals(user string) ([]models.OrderBalance, error) {
	ctx := context.Background()
	conn, err := b.pool.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	rows, err := conn.Query(ctx, `SELECT "order", sum, processed_at FROM balances WHERE "user" = ($1) AND sum < 0;`, user)
	if err != nil {
		return nil, err
	}
	orders := make([]models.OrderBalance, 0)
	for rows.Next() {
		order := &models.OrderBalance{}
		if err := rows.Scan(&order.Order, &order.Sum, &order.ProcessedAt.Time); err != nil {
			return nil, err
		}
		order.Sum = -order.Sum
		orders = append(orders, *order)
	}
	return orders, nil
}

func (b *BalancesDB) Accrual(user string, order models.AccrualOrder) error {
	ctx := context.Background()
	conn, err := b.pool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	_, err = conn.Exec(ctx, `INSERT INTO balances(processed_at, "user", "order", sum) VALUES ($1, $2, $3, $4);`, time.Now(), user, order.Order, order.Accrual)
	if err != nil {
		return err
	}

	rows, err := conn.Query(ctx, `SELECT processed_at, "user", "order", sum FROM balances`)
	if err != nil {
		log.Println(err)
		return nil
	}

	for rows.Next() {
		var usr string
		order := &models.OrderBalance{}
		if err := rows.Scan(&order.ProcessedAt.Time, &usr, &order.Order, &order.Sum); err != nil {
			log.Println(err)
			return nil
		}
		log.Println(usr, order)
	}

	return nil
}
