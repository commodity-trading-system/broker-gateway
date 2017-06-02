package entities

import (
	"github.com/shopspring/decimal"
	"time"
)

type Quotation struct {
	FutureId int	`gorm:"primary_key"`
	Price decimal.Decimal
	CreatedAt time.Time
}
