package entities

import (
	"github.com/shopspring/decimal"
	"github.com/satori/go.uuid"
	"github.com/jinzhu/gorm"
	"time"
)

type Order struct {
	ID uuid.UUID	`gorm:"primary_key;type:varchar(255);unique_index"`
	Quantity int
	Price decimal.Decimal  `gorm:"type:decimal(10,2)"`
	BuyerId int
	SellerId int
	FutureId int
	Status int
	SellerCommission decimal.Decimal	`gorm:"type:decimal(10,2)"`
	BuyerCommission decimal.Decimal		`gorm:"type:decimal(10,2)"`

	BuyerConsignationId uuid.UUID	`gorm:"type:varchar(255);index"`
	SellerConsignationId uuid.UUID	`gorm:"type:varchar(255);index"`


	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`

}



func TransformForFirm(o interface{}, firmId int) Order  {
	order := o.(Order)
	if order.BuyerId == firmId {
		order.SellerId = -1
		order.SellerConsignationId = uuid.FromStringOrNil("")
	} else if order.SellerId == firmId {
		order.BuyerId = -1
		order.BuyerConsignationId = uuid.FromStringOrNil("")
	} else {
		order = Order{}
	}

	return order
}



func (order *Order) BeforeCreate(scope *gorm.Scope) error {
	scope.SetColumn("id", uuid.NewV1().String())
	return nil
}
