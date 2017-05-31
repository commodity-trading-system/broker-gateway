package executor

import (
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
	d, err := gorm.Open("mysql",config.User+":"+
		config.Password + "@tcp(" +
		config.Host + ":" +
		strconv.Itoa(config.Port) + ")/"+
		config.DBName+"?charset=utf8")

	if err != nil {
		return nil, err
	}
	return &db{
		client: d,
	},nil
}


func (d *db) Migrate()  {
	d.client.AutoMigrate(&entities.Future{}, &entities.Firm{}, &entities.Order{}, &entities.Consignation{})
}

func (d *db) Query() *gorm.DB {
	return d.client
}

func (d *db) Create(value interface{}) bool  {
	d.client.Create(value)
	return d.client.NewRecord(value)
}

func (d *db) Save(model interface{})  {
	d.client.Save(model)
}

func (d *db) Update(model interface{}, attrs map[string]string) *gorm.DB {
	return d.client.Model(model).Update(attrs)
}
