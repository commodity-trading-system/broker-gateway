package entities

import (
	"github.com/shopspring/decimal"
)

type Consignation struct {
	Type int
	Quantity int
	FutureId int
	Price decimal.Decimal
	OpenQuantity int
	Direction int
	FirmId int
}
