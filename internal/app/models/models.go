package models

import (
	"encoding/json"
	"fmt"
	"github.com/shopspring/decimal"
	"time"
)

type OrderStatus string

const (
	TimeFormat                       = time.RFC3339
	OrderStatusNew       OrderStatus = "NEW"
	OrderStatusInvalid   OrderStatus = "INVALID"
	OrderStatusProcessed OrderStatus = "PROCESSED"
)

type CustomTime struct {
	time.Time
}

func (c *CustomTime) MarshalJSON() ([]byte, error) {
	if c.Time.IsZero() {
		return nil, nil
	}
	return []byte(fmt.Sprintf(`"%s"`, c.Time.Format(TimeFormat))), nil
}

type CustomDecimal struct {
	decimal.Decimal
}

func (c *CustomDecimal) MarshalJSON() ([]byte, error) {
	n, _ := c.Float64()
	return json.Marshal(n)
}

type User struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type Order struct {
	UserID     int           `json:"-"`
	Number     string        `json:"number"`
	Status     OrderStatus   `json:"status"`
	Accrual    CustomDecimal `json:"accrual,omitempty"`
	UploadedAt CustomTime    `json:"uploaded_at"`
}

type CurrentBalance struct {
	Current   CustomDecimal `json:"current"`
	Withdrawn CustomDecimal `json:"withdrawn"`
}

type OrderBalance struct {
	Order       string        `json:"order"`
	Sum         CustomDecimal `json:"sum"`
	ProcessedAt CustomTime    `json:"processed_at,omitempty"`
}

type AccrualOrder struct {
	Order   string        `json:"order"`
	Status  OrderStatus   `json:"status"`
	Accrual CustomDecimal `json:"accrual"`
}
