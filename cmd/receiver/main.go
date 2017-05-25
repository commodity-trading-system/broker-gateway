package main

import (
	"flag"
	"os"
	"github.com/quickfixgo/quickfix"
	"broker-gateway/receiver"
	"github.com/joho/godotenv"

	"log"
)

func main()  {
	env, err := godotenv.Read()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	config := receiver.ReceiverConfig{
		RedisHost: env["REDIS_HOST"],
		RedisPort: env["REDIS_PORT"],
		RedisPwd: env["REDIS_PASSWORD"],
		RedisDB: env["REDIS_DB"],
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
