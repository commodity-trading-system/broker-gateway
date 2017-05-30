package executor

import (
	"database/sql"
	"strings"
	"github.com/jinzhu/gorm"
	"strconv"
	"broker-gateway/entities"
)

type DB interface {
	Migrate()
}

type DBConfig struct {
	Host string
	Port int
	User string
	Password string
	DBName string
}

type db struct {
	client *gorm.DB
}


func NewDB(config DBConfig) (DB, error)  {
	db, err := gorm.Open("mysql",config.User+":"+
		config.Password + "@tcp" +
		config.Host + ":" +
		strconv.Itoa(config.Port) + ")/"+
		config.DBName+"?charset=utf8")

	if err != nil {
		return nil, err
	}
	return &db{
		client: db,
	},nil
}


func (d *db) Migrate()  {
	d.client.AutoMigrate(&entities.Order{}, &entities.Consignation{}, &entities.Future{}, &entities.Firm{})
}
