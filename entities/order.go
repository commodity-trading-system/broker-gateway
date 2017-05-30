package entities

import (
	"github.com/shopspring/decimal"
	"github.com/satori/go.uuid"
	"github.com/jinzhu/gorm"
)

type Order struct {
	gorm.Model
	Quantity uint32
	Price decimal.Decimal
	Id uuid.UUID
	BuyerId uint32
	SellerId uint32
	FutureId uint32
	Status uint32
	Commission decimal.Decimal

	BuyerConsignationId uuid.UUID
	SellerConsignationId uuid.UUID

}
