package main

import (
	"flag"
	"os"
	"github.com/quickfixgo/quickfix"
	"broker-gateway/receiver"
	"github.com/joho/godotenv"

	"log"
	"strconv"
)

func main()  {
	env, err := godotenv.Read()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	port,_ := strconv.ParseInt(env["REDIS_PORT"],10,32)
	db,_ := strconv.ParseInt(env["REDIS_DB"],10,32)
	config := receiver.ReceiverConfig{
		RedisHost: env["REDIS_HOST"],
		RedisPort: int(port),
		RedisPwd: env["REDIS_PASSWORD"],
		RedisDB: int(db),
	}

	app := receiver.NewReceiver(config)
	flag.Parse()
	fileName := flag.Arg(0)

	cfg, _ := os.Open(fileName)
	appSettings, _ := quickfix.ParseSettings(cfg)
	storeFactory := quickfix.NewMemoryStoreFactory()
	logFactory, _ := quickfix.NewFileLogFactory(appSettings)
	acceptor, _ := quickfix.NewAcceptor(app, storeFactory, appSettings, logFactory)

	_ = acceptor.Start()
	//for condition == true { do something }
	acceptor.Stop()
}
