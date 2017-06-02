package queier

import (
	"broker-gateway/executor"
	"broker-gateway/entities"
	"fmt"
	"github.com/satori/go.uuid"
)

type Querier interface {
	Consignations(int) []entities.Consignation
	ConsignationById(firmId int, id string) entities.Consignation
	Orders(firmId int)	[]entities.Order
	OrderById(firmId int, id string)	entities.Order
	Futures()	[]entities.Future
}

type querier struct {
	db executor.DB
}

func NewQuerier(cfg executor.DBConfig) Querier {
	db, err := executor.NewDB(cfg)
	if err != nil {
		return nil
	}

	return &querier{
		db:db,
	}
}

func (q querier) Orders(firmId int) []entities.Order {
	var orders []entities.Order
	q.db.Query().Where("buyer_id = ? ",firmId).Or("seller_id = ? ", firmId).Find(&orders)
	for i:=0;i<len(orders);i++ {
		if orders[i].BuyerId == firmId {
			orders[i].SellerId = -1
			orders[i].SellerConsignationId = uuid.FromStringOrNil("")
		} else {
			orders[i].BuyerId = -1
			orders[i].BuyerConsignationId = uuid.FromStringOrNil("")
		}
	}
	return orders
}

func (q querier) OrderById(firmId int, id string) entities.Order {
	var order entities.Order
	q.db.Query().Where("id = ? ", id).First(&order)
	if order.BuyerId == firmId {
		order.SellerId = -1
		order.SellerConsignationId = uuid.FromStringOrNil("")
	} else if order.SellerId == firmId {
		order.BuyerId = -1
		order.BuyerConsignationId = uuid.FromStringOrNil("")
	} else {
		order = entities.Order{}
	}
	return order
}

func (q querier) Consignations(firmId int) []entities.Consignation {
	var consignations []entities.Consignation
	q.db.Query().Where("firm_id = ? ",firmId).Find(&consignations)
	return consignations
}

func (q querier) ConsignationById(firmId int, id string) entities.Consignation {
	var consignation entities.Consignation
	q.db.Query().Where("id = ? ", id).First(&consignation)
	return consignation
}



func (q querier) Futures() []entities.Future {
	var futures  []entities.Future
	q.db.Query().Find(&futures)
	fmt.Println(futures)
	return futures
}



