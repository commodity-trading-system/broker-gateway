package main

import (
	"github.com/joho/godotenv"
	"strconv"
	"os"
	"broker-gateway/executor"
	"fmt"
	"strings"
	"broker-gateway/enum"
)

func main()  {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	port,_ := strconv.ParseInt(os.Getenv("REDIS_PORT"),10,32)
	db,_ := strconv.ParseInt(os.Getenv("REDIS_DB"),10,32)
	mysqlPort,_ := strconv.ParseInt(os.Getenv("MYSQL_PORT"),10,32)
	futures := strings.FieldsFunc(os.Getenv("FUTURES"), func(s rune) bool {
		return s == enum.MARSHAL_DELI
	})
	config := executor.ExecutorConfig{
		RedisHost: os.Getenv("REDIS_HOST"),
		RedisPort: int(port),
		RedisPwd: os.Getenv("REDIS_PASSWORD"),
		RedisDB: int(db),
		MysqlHost: os.Getenv("MYSQL_HOST"),
		MysqlPort: int(mysqlPort),
		MysqlPwd: os.Getenv("MYSQL_PASSWORD"),
		MysqlDB: os.Getenv("MYSQL_DB"),
		MysqlUser: os.Getenv("MYSQL_USER"),
		Futures: futures,
	}

	exe, err := executor.NewExecutor(config)
	if err != nil {
		fmt.Println(err)
		return
	}
	exe.Execute()
}


