package executor

import (
	"fmt"
	"strconv"
	"github.com/quickfixgo/quickfix"
	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
)
 

type Executor struct {
	redisClient *redis.Client
	*quickfix.MessageRouter
	orderBooks map[string] *OrderBook
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
}

func NewExecutor(config ExecutorConfig) (*Executor,error) {

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

	obs := make(map[string] *OrderBook,len(config.Futures))
	for i:= 0; i< len(config.Futures) ;i++  {
		obs[config.Futures[i]] = NewOrderBook(db)
	}

	r := &Executor{
		redisClient: redis.NewClient(&redis.Options{
			Addr:     config.RedisHost+":"+strconv.Itoa(config.RedisPort),
			Password: config.RedisPwd,
			DB:       config.RedisDB,
		}),
		orderBooks: obs,
	}
	return r,nil
}

func (executor *Executor) Execute()  {
	for {
		for futureId, orderBook := range executor.orderBooks {
			result, err := executor.redisClient.RPopLPush("future_"+ futureId,"temp_future_" + futureId).Result()
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println(result)
			fmt.Println(orderBook)

			//orderBook.addLimit()
			// Match successfully, pop the consignation
			executor.redisClient.RPop("temp_future_" + futureId)
		}

	}

}
