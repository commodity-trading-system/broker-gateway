package entities

import (
	"github.com/shopspring/decimal"
	"time"
)

type Quotation struct {
	FutureId int
	Price decimal.Decimal	`gorm:"type:decimal(10,2)"`
	CreatedAt time.Time	`sql:"index"`
}
