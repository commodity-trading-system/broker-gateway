package executor

import (
	"github.com/jinzhu/gorm"
	"strconv"
	"broker-gateway/entities"
)

type DB interface {
	Migrate()
	// Create a new object
	Create(value interface{})
	// Update completely
	Save(model interface{})
	// Update partially
	Update(model interface{}, attrs map[string]string) *gorm.DB

	Query() *gorm.DB

	Empty()

	Seeder()

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
		config.DBName+"?charset=utf8&parseTime=true")

	if err != nil {
		return nil, err
	}
	return &db{
		client: d,
	},nil
}


func (d *db) Migrate()  {
	d.client.AutoMigrate(&entities.Future{},
		&entities.Firm{},
		&entities.Order{},
		&entities.Consignation{},
		&entities.Quotation{},
		&entities.Commission{})
}

func (d *db) Query() *gorm.DB {
	return d.client
}

func (d *db) Create(value interface{})  {
	d.client.Create(value)
}

func (d *db) Save(model interface{})  {
	d.client.Save(model)
}

func (d *db) Update(model interface{}, attrs map[string]string) *gorm.DB {
	return d.client.Model(model).Update(attrs)
}

func (d *db) Empty() {
	tables := []string{"futures","firms","orders","consignations","quotations","commissions"}
	for i:=0; i<len(tables); i++ {
		d.client.DropTable(tables[i])
	}
}

func (d *db) Seeder()  {
	d.Save(&entities.Future{
		ID: 1,
		Name: "oil",
		Period: "10",
		Description: "-2017.10 oil",
	})

	d.Save(&entities.Future{
		ID: 2,
		Name: "oil",
		Period: "12",
		Description: "-2017.12 oil",
	})

	d.Save(&entities.Future{
		ID: 3,
		Name: "gold",
		Period: "8",
		Description: "-2017.8 gold",
	})

	d.Save(&entities.Future{
		ID: 4,
		Name: "gold",
		Period: "12",
		Description: "-2017.12 gold",
	})

	commissions := [][]int{[]int{1,1,3,1},[]int{1,1,5,2},[]int{1,1,10,3}}

	for i:=0; i<len(commissions) ;i++  {
		d.Save(&entities.Commission{
			ID: i,
			FirmId: commissions[i][0],
			FutureId: commissions[i][1],
			CommissionPercent: commissions[i][2],
			OrderType: commissions[i][3],
		})
	}

}