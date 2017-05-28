package entities

import (
	"github.com/shopspring/decimal"
	"strconv"
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
