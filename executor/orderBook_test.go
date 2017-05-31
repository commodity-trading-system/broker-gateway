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

func newLimitConsignation(ctype, direction, price, quantity int) *entities.Consignation  {
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

func TestOrderBook_AddLimit(t *testing.T) {
	c1 := newLimitConsignation(enum.OrdType_LIMIT, enum.OrderDirection_BUY, 5000, 100)
	c2 := newLimitConsignation(enum.OrdType_LIMIT, enum.OrderDirection_SELL, 5100, 100)
	c3 := newLimitConsignation(enum.OrdType_LIMIT, enum.OrderDirection_BUY, 5200, 110)
	book.AddLimit(c1)
	book.AddLimit(c2)
	book.AddLimit(c3)
}
