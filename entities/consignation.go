package entities

import (
	"github.com/shopspring/decimal"
	"strconv"
	"github.com/satori/go.uuid"
	"time"
)

type Consignation struct {
	ID uuid.UUID	`gorm:"type:varchar(255);unique_index"`
	Type int
	Quantity int
	FutureId int
	Price decimal.Decimal	`gorm:"type:decimal(10,2)"`
	OpenQuantity int
	Direction int
	FirmId int

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`


	// Relations
	Orders []Order	`gorm:"-"`
}

func (c Consignation) MarshalBinary() (data []byte, err error) {
	return []byte(strconv.Itoa(c.Type) + "," +
		strconv.Itoa(c.Quantity) + "," +
		strconv.Itoa(c.FutureId) + "," +
		strconv.Itoa(c.FutureId) + "," +
		strconv.Itoa(c.OpenQuantity) + "," +
		strconv.Itoa(c.Direction) + "," +
		strconv.Itoa(c.FirmId) + "," +
		c.Price.String() + ","),nil
}

func (c Consignation) UnmarshalBinary(data []byte) error  {
	return nil
}
