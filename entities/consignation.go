package entities

import (
	"github.com/shopspring/decimal"
	"strconv"
	"github.com/satori/go.uuid"
	"time"
	"strings"
	"broker-gateway/enum"
)

type Consignation struct {
	ID uuid.UUID	`gorm:"primary_key;type:varchar(255);unique_index"`
	Type int
	Quantity int
	FutureId int
	Price decimal.Decimal	`gorm:"type:decimal(10,2)"`
	OpenQuantity int
	Direction int
	FirmId int
	Status int

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`


	// Relations
	Orders []Order	`gorm:"-"`
}

func (c Consignation) MarshalBinary() (data []byte, err error) {
	return []byte(strconv.Itoa(c.Type) + string(enum.MARSHAL_DELI) +
		strconv.Itoa(c.Quantity) + string(enum.MARSHAL_DELI) +
		strconv.Itoa(c.FutureId) + string(enum.MARSHAL_DELI) +
		strconv.Itoa(c.OpenQuantity) + string(enum.MARSHAL_DELI) +
		strconv.Itoa(c.Direction) + string(enum.MARSHAL_DELI) +
		strconv.Itoa(c.FirmId) + string(enum.MARSHAL_DELI) +
		c.Price.String() + string(enum.MARSHAL_DELI) +
		c.ID.String()+ string(enum.MARSHAL_DELI) +
		strconv.Itoa(c.Status)),nil
}


func split(s rune) bool {
	if s == enum.MARSHAL_DELI {
		return true
	}
	return false
}

func (c *Consignation) UnmarshalBinary(data []byte) error  {
	res := strings.FieldsFunc(string(data),split)
	c.Type,_ 		= strconv.Atoi(res[0])
	c.Quantity,_ 		= strconv.Atoi(res[1])
	c.FutureId,_ 		= strconv.Atoi(res[2])
	c.OpenQuantity,_ 	= strconv.Atoi(res[3])
	c.Direction,_ 		= strconv.Atoi(res[4])
	c.FirmId,_ 		= strconv.Atoi(res[5])
	price,_			:= strconv.ParseFloat(res[6],32)
	c.Price 		= decimal.NewFromFloat(price)
	id, _ 			:= uuid.FromString(res[7])
	c.ID = id
	c.Status, _ 		= strconv.Atoi(res[8])
	return nil
}



func WapperUnmarshalBinary(consignation *Consignation, data[]byte)  {
	consignation.UnmarshalBinary(data)
}
