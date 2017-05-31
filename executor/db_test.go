package executor

import (
	"testing"
	"strconv"
	"os"
	"fmt"
	"broker-gateway/entities"
	"github.com/satori/go.uuid"
	"github.com/shopspring/decimal"
)

var d DB

func TestNewDB(t *testing.T) {
	port,_ := strconv.ParseInt(os.Getenv("MYSQL_PORT"),10,32)
	config := DBConfig{
		Host: os.Getenv("MYSQL_HOST"),
		Port: int(port),
		Password: os.Getenv("MYSQL_PASSWORD"),
		DBName: os.Getenv("MYSQL_DB"),
		User: os.Getenv("MYSQL_USER"),
	}
	fmt.Println(config)

	db,err := NewDB(config)
	d = db
	if err != nil {
		t.Error(err)
	}
}

func TestDb_Migrate(t *testing.T) {
	d.Migrate()
}

func TestDb_Create(t *testing.T) {
	order := entities.Order{
		Price: decimal.New(2300,-2),
		BuyerId: 1,
		SellerId: 1,
		ID: uuid.NewV1(),
		BuyerConsignationId: uuid.NewV1(),
		SellerConsignationId: uuid.NewV1(),
		Quantity:10,
		FutureId: 1,
	}
	fmt.Println(order)
	d.Create(&order)
}