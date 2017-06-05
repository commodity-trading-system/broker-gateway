package queier

import (
	"broker-gateway/executor"
	"broker-gateway/entities"
	"github.com/jinzhu/gorm"
)

type Querier interface {
	AllOrders(limit, offset int) []entities.Order
	Orders(firmId, limit, offset int) []entities.Order
	OrderById(firmId int, id string)	entities.Order

	AllConsignations(limit, offset int) []entities.Consignation
	Consignations(firmId,limit, offset int) []entities.Consignation
	ConsignationById(firmId int, id string) entities.Consignation

	Futures()	[]entities.Future
	Quotations(futureId,limit,offset int)[]entities.Quotation


	Update(entity interface{})
	Save(entity interface{})

	Query() *gorm.DB
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

func (q querier) Query() *gorm.DB {
	return q.db.Query()
}

func (q querier) AllOrders(limit, offset int) []entities.Order {
	var orders []entities.Order
	q.db.Query().
		Limit(limit).
		Offset(offset).
		Order("created_at desc").
		Find(&orders)
	return orders
}

func (q querier) Orders(firmId, limit, offset int) []entities.Order {
	var orders []entities.Order
	q.db.Query().
		Limit(limit).
		Offset(offset).
		Order("created_at desc").
		Where("buyer_id = ? ",firmId).
		Or("seller_id = ? ", firmId).
		Find(&orders)
	return orders
}

func (q querier) OrderById(firmId int, id string) entities.Order {
	var order entities.Order
	q.db.Query().Where("id = ? ", id).First(&order)
	return order
}

func (q querier) AllConsignations(limit, offset int) []entities.Consignation {
	var consignations []entities.Consignation
	q.db.Query().
		Limit(limit).
		Offset(offset).
		Order("created_at desc").
		//Where("firm_id = ? ",firmId).
		Find(&consignations)
	return consignations
}

func (q querier) Consignations(firmId,limit, offset int) []entities.Consignation {
	var consignations []entities.Consignation
	q.db.Query().
		Limit(limit).
		Offset(offset).
		Order("created_at desc").
		Where("firm_id = ? ",firmId).
		Find(&consignations)
	return consignations
}

func (q querier) ConsignationById(firmId int, id string) entities.Consignation {
	var consignation entities.Consignation
	q.db.Query().Where("id = ? ", id).First(&consignation)
	return consignation
}

func (q querier) Update(entity interface{})  {
	q.db.Save(entity)
}

func (q querier) Save(entity interface{})  {
	q.db.Save(entity)
}



func (q querier) Futures() []entities.Future {
	var futures  []entities.Future
	q.db.Query().Find(&futures)
	return futures
}

func (q querier) Quotations(futureId int,limit, offset int)[]entities.Quotation  {
	var entity  []entities.Quotation
	q.db.Query().
		Limit(limit).
		Offset(offset).
		Order("created_at desc").
		Where("future_id = ?",futureId).
		//Where("created_at > ? AND created_at < ?",start, end).
		Find(&entity)
	return entity
}



