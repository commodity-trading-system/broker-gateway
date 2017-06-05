package main

import (
	"flag"
	"os"
	"github.com/quickfixgo/quickfix"
	"broker-gateway/receiver"
	"github.com/joho/godotenv"

	"log"
	"strconv"
	"fmt"
	"os/signal"
)

func main()  {
	env, err := godotenv.Read()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	port,_ := strconv.ParseInt(env["REDIS_PORT"],10,32)
	db,_ := strconv.ParseInt(env["REDIS_DB"],10,32)

	mysqlPort,_ := strconv.ParseInt(os.Getenv("MYSQL_PORT"),10,32)

	config := receiver.ReceiverConfig{
		RedisHost: env["REDIS_HOST"],
		RedisPort: int(port),
		RedisPwd: env["REDIS_PASSWORD"],
		RedisDB: int(db),

		MysqlHost: os.Getenv("MYSQL_HOST"),
		MysqlPort: int(mysqlPort),
		MysqlPwd: os.Getenv("MYSQL_PASSWORD"),
		MysqlDB: os.Getenv("MYSQL_DB"),
		MysqlUser: os.Getenv("MYSQL_USER"),
	}

	app := receiver.NewReceiver(config)
	flag.Parse()
	fileName := flag.Arg(0)


	cfg, err := os.Open(fileName)
	if err != nil {
		fmt.Println("Invalid config file", err)
		return
	}

	appSettings, err := quickfix.ParseSettings(cfg)
	if err != nil {
		fmt.Println("Invalid config file", err)
		return
	}
	storeFactory := quickfix.NewMemoryStoreFactory()
	logFactory, err := quickfix.NewFileLogFactory(appSettings)
	if err != nil {
		fmt.Println("New log factory error", err)
		return
	}
	acceptor, err := quickfix.NewAcceptor(app, storeFactory, appSettings, logFactory)
	if err != nil {
		fmt.Println("New Acceptor error", err)
		return
	}

	err = acceptor.Start()
	if err != nil {
		fmt.Println("Acceptor error", err)
		return
	}
	fmt.Println("after start")

	interrupt := make(chan os.Signal)
	signal.Notify(interrupt, os.Interrupt, os.Kill)
	<-interrupt

	acceptor.Stop()

	//for true {}
	//for condition == true { do something }
	//acceptor.Stop()
}
