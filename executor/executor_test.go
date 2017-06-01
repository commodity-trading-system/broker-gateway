package executor

import (
	"strconv"
	"os"
	"fmt"
	"testing"
	"github.com/go-redis/redis"
	"broker-gateway/enum"
	"strings"
)

var exe *Executor
var red *redis.Client


func TestNewExecutor(t *testing.T) {
	port,_ := strconv.ParseInt(os.Getenv("REDIS_PORT"),10,32)
	db,_ := strconv.ParseInt(os.Getenv("REDIS_DB"),10,32)
	mysqlPort,_ := strconv.ParseInt(os.Getenv("MYSQL_PORT"),10,32)
	etcdEndpoints := strings.FieldsFunc(os.Getenv("ETCD_ENDPOINTS"), func(s rune) bool {
		return s==enum.MARSHAL_DELI
	})
	config := ExecutorConfig{
		RedisHost: os.Getenv("REDIS_HOST"),
		RedisPort: int(port),
		RedisPwd: os.Getenv("REDIS_PASSWORD"),
		RedisDB: int(db),
		MysqlHost: os.Getenv("MYSQL_HOST"),
		MysqlPort: int(mysqlPort),
		MysqlPwd: os.Getenv("MYSQL_PASSWORD"),
		MysqlDB: os.Getenv("MYSQL_DB"),
		MysqlUser: os.Getenv("MYSQL_USER"),
		Futures: []string{"1"},
		EtcdEndpoints: etcdEndpoints,
	}
	fmt.Println(config)
	red = redis.NewClient(&redis.Options{
		Addr:     config.RedisHost+":"+strconv.Itoa(config.RedisPort),
		Password: config.RedisPwd,
		DB:       config.RedisDB,
	})

	intCmd := red.RPush("future_"+"1", newConsignation(enum.OrdType_LIMIT,enum.OrderDirection_BUY,5000, 100))
	if intCmd.Err() != nil {
		t.Error(intCmd.Err())
	}

	_, err := NewExecutor(config)
	if err != nil {
		t.Error(err)
	}

	//exe.Execute()
}

