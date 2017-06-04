package executor

import (
	"fmt"
	"strconv"
	"github.com/quickfixgo/quickfix"
	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
	"broker-gateway/entities"
	"broker-gateway/enum"
	"github.com/coreos/etcd/client"
	"time"
	"github.com/shopspring/decimal"
)

type Executor interface {
	Execute()
}

type executor struct {
	redisClient *redis.Client
	*quickfix.MessageRouter
	orderBooks map[string] OrderBook
	publisher Publisher
	db DB
}

type ExecutorConfig struct {
	RedisHost string
	RedisPort int
	RedisPwd string
	RedisDB int
	MysqlHost string
	MysqlPort int
	MysqlPwd string
	MysqlDB string
	MysqlUser string
	Futures []string
	EtcdEndpoints []string
}

type inspector struct {
	db DB
	//  setting[futureId][firmId][type]
	commissionSetting map[int]map[int]map[int]int
}

func (in *inspector) GetCommission(futureId, firmId, orderType int, amount decimal.Decimal) decimal.Decimal {
	firmAndOrderType, exist := in.commissionSetting[futureId]
	if ! exist {
		return decimal.New(0,0)
	}
	drTy, exist := firmAndOrderType[firmId]
	if ! exist {
		return decimal.New(0,0)
	}

	percent, exist := drTy[orderType]
	if ! exist {
		return decimal.New(0,0)
	}
	return amount.Div(decimal.New(int64(percent),-2))
}


func (in *inspector) MonitorSetting()  {
	for true {
		var commissions []entities.Commission
		in.db.Query().Find(&commissions)
		for i:=0; i< len(commissions) ;i++  {
			in.commissionSetting[commissions[i].FutureId][commissions[i].FirmId][commissions[i].OrderType] = commissions[i].CommissionPercent
		}
		time.Sleep(time.Second * 10)
	}

}

func NewExecutor(config ExecutorConfig) (Executor,error) {

	db, err := NewDB(DBConfig{
		Host: config.MysqlHost,
		Port: config.MysqlPort,
		User: config.MysqlUser,
		Password: config.MysqlPwd,
		DBName: config.MysqlDB,
	})
	if err != nil {
		return nil, err
	}

	db.Empty()
	db.Migrate()
	db.Seeder()

	etcdPublisher := NewPublisher(client.Config{
		Endpoints: config.EtcdEndpoints,
	})

	insp := &inspector{
		commissionSetting: map[int]map[int]map[int]int{},
		db: db,
	}

	go insp.MonitorSetting()

	obs := make(map[string] OrderBook,len(config.Futures))
	for i:= 0; i< len(config.Futures) ;i++  {
		fid,_ :=strconv.Atoi(config.Futures[i])
		obs[config.Futures[i]] = NewOrderBook(fid,db, etcdPublisher,insp)
	}


	r := &executor{
		redisClient: redis.NewClient(&redis.Options{
			Addr:     config.RedisHost+":"+strconv.Itoa(config.RedisPort),
			Password: config.RedisPwd,
			DB:       config.RedisDB,
		}),
		db:db,
		orderBooks: obs,
	}

	go r.MonitorSetting()

	db.Migrate()
	return r,nil
}



func (executor *executor) Execute()  {
	for {
		for futureId, book := range executor.orderBooks {
			result, err := executor.redisClient.RPopLPush("future_"+ futureId,"temp_future_" + futureId).Result()
			if err != nil {
				if err.Error() == "redis: nil" {
					continue
				}
				fmt.Println(err)
				continue

			}

			consignation := entities.Consignation{}
			entities.WapperUnmarshalBinary(&consignation,[]byte(result))
			add(book, &consignation)

			fmt.Println(result)
			fmt.Println(consignation)

			// Match successfully, pop the consignation
			executor.redisClient.RPop("temp_future_" + futureId)
		}
	}
}


func add(book OrderBook, cons *entities.Consignation)  {
	if cons.Type == enum.OrdType_LIMIT {
		book.AddLimit(cons)
	} else if  cons.Type == enum.OrdType_MARKET {
		book.AddMarket(cons)
	} else if cons.Type == enum.OrdType_STOP {
		book.AddStop(cons)
	} else if cons.Type == enum.OrdType_CANCEL {
		book.AddCancel(cons)
	}
}
