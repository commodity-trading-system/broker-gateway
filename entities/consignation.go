package entities

import (
	"github.com/shopspring/decimal"
)

type Consignation struct {
	Type int
	Quantity uint32
	FutureId uint32
	Price decimal.Decimal
	OpenQuantity uint32
}
