package entities

import (
	"github.com/shopspring/decimal"
	"time"
)

type Quotation struct {
	FutureId int	`gorm:"primary_key"`
	Price decimal.Decimal	`gorm:"type:decimal(10,2)"`
	CreatedAt time.Time	`sql:"index"`
}
