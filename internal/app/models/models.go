package models

import (
	"fmt"
	"github.com/shopspring/decimal"
	"time"
)

const TimeFormat = time.RFC3339

type CustomTime struct {
	time.Time
}

func (c *CustomTime) MarshalJSON() ([]byte, error) {
	if c.Time.IsZero() {
		return nil, nil
	}
	return []byte(fmt.Sprintf(`"%s"`, c.Time.Format(TimeFormat))), nil
}

type User struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type Order struct {
	User       string     `json:"-"`
	Number     string     `json:"number"`
	Status     string     `json:"status"`
	Accrual    float64    `json:"accrual,omitempty"`
	UploadedAt CustomTime `json:"uploaded_at"`
}

type CurrentBalance struct {
	Current   decimal.Decimal `json:"current"`
	Withdrawn decimal.Decimal `json:"withdrawn"`
}

type OrderBalance struct {
	Order       string          `json:"order"`
	Sum         decimal.Decimal `json:"sum"`
	ProcessedAt CustomTime      `json:"processed_at,omitempty"`
}
