package entities

import (
	"github.com/shopspring/decimal"
	"github.com/satori/go.uuid"
	"github.com/jinzhu/gorm"
	"time"
)

type Order struct {
	ID uuid.UUID	`gorm:"type:varchar(255);unique_index"`
	Quantity int
	Price decimal.Decimal  `gorm:"type:decimal(10,2)"`
	BuyerId int
	SellerId int
	FutureId int
	Status int
	Commission decimal.Decimal	`gorm:"type:decimal(10,2)"`

	BuyerConsignationId uuid.UUID	`gorm:"type:varchar(255);index"`
	SellerConsignationId uuid.UUID	`gorm:"type:varchar(255);index"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`


}

func (order *Order) BeforeCreate(scope *gorm.Scope) error {
	scope.SetColumn("id", order.ID.String())
	return nil
}
