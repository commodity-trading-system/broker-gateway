package executor

import (
	"testing"
	"strconv"
	"os"
	"broker-gateway/entities"
	"github.com/satori/go.uuid"
	"github.com/shopspring/decimal"
	"broker-gateway/enum"
)

var book OrderBook

func TestNewOrderBook(t *testing.T) {
	port,_ := strconv.ParseInt(os.Getenv("MYSQL_PORT"),10,32)
	config := DBConfig{
		Host: os.Getenv("MYSQL_HOST"),
		Port: int(port),
		Password: os.Getenv("MYSQL_PASSWORD"),
		DBName: os.Getenv("MYSQL_DB"),
		User: os.Getenv("MYSQL_USER"),
	}

	db,err := NewDB(config)
	d = db
	if err != nil {
		t.Error(err)
	}
	db.Empty()
	db.Migrate()

	book = NewOrderBook(db)
}

func newConsignation(ctype, direction, price, quantity int) *entities.Consignation  {
	return &entities.Consignation{
		ID: uuid.NewV1(),
		Quantity: quantity,
		Price: decimal.New(int64(price),-2),
		FutureId: 1,
		Type: ctype,
		FirmId: 2,
		Direction: direction,
		Status: enum.ConsignationStatus_APPENDING,
		OpenQuantity: quantity,
	}
}
func add(cons []*entities.Consignation)  {
	for i:=0; i<len(cons); i++ {
		if cons[i].Type == enum.OrdType_LIMIT {
			book.AddLimit(cons[i])
		} else if  cons[i].Type == enum.OrdType_MARKET {
			book.AddMarket(cons[i])
		}
	}
}

func TestOrderBook_AddLimit(t *testing.T) {
	c1 := newConsignation(enum.OrdType_LIMIT, enum.OrderDirection_BUY, 5000, 100)
	c2 := newConsignation(enum.OrdType_LIMIT, enum.OrderDirection_SELL, 5100, 100)
	c3 := newConsignation(enum.OrdType_LIMIT, enum.OrderDirection_BUY, 5200, 110)
	c4 := newConsignation(enum.OrdType_LIMIT, enum.OrderDirection_SELL, 5000,110)
	book.AddLimit(c1)
	book.AddLimit(c2)
	book.AddLimit(c3)
	book.AddLimit(c4)
}

func TestOrderBook_AddMarket(t *testing.T) {
	book.Reset()
	cons := []*entities.Consignation{}
	cons = append(cons,newConsignation(enum.OrdType_LIMIT, enum.OrderDirection_BUY, 5000, 100))
	cons = append(cons,newConsignation(enum.OrdType_LIMIT, enum.OrderDirection_BUY, 5100, 100))
	cons = append(cons,newConsignation(enum.OrdType_MARKET, enum.OrderDirection_SELL, 5100, 190))
	add(cons)
}
