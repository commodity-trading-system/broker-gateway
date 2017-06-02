package main

import (
	"fmt"
	"flag"
	"path"
	"os"
	"github.com/quickfixgo/quickfix"
	"broker-gateway/receiver"
	"broker-gateway/enum"
	//"github.com/quickfixgo/quickfix/fix44/newordersingle"
	//"github.com/quickfixgo/quickfix/field"
	//"github.com/quickfixgo/quickfix/fix43/newordercross"
)

func main()  {
	flag.Parse()

	cfgFileName := path.Join("config", "tradeclient.cfg")
	if flag.NArg() > 0 {
		cfgFileName = flag.Arg(0)
	}

	cfg, err := os.Open(cfgFileName)
	if err != nil {
		fmt.Printf("Error opening %v, %v\n", cfgFileName, err)
		return
	}

	appSettings, err := quickfix.ParseSettings(cfg)
	if err != nil {
		fmt.Println("Error reading cfg,", err)
		return
	}

	app := receiver.NewClient()

	fileLogFactory, err := quickfix.NewFileLogFactory(appSettings)

	if err != nil {
		fmt.Println("Error creating file log factory,", err)
		return
	}

	initiator, err := quickfix.NewInitiator(app, quickfix.NewMemoryStoreFactory(), appSettings, fileLogFactory)
	if err != nil {
		fmt.Printf("Unable to create Initiator: %s\n", err)
		return
	}

	initiator.Start()

	for {
		msg := quickfix.NewMessage()

		//temp := newordersingle.New(field.NewClOrdID("test"),
		//	newordercross.NoSides{},
		//	nil,
		//	field.NewOrdType(enum.OrdType_LIMIT))
		msg.Header.SetInt(quickfix.Tag(8),1)
		//msg.Body.SetInt(quickfix.Tag(8),10)

		msg.Body.SetInt(quickfix.Tag(enum.TagNum_DIRECTION), enum.OrderDirection_BUY)
		msg.Body.SetInt(quickfix.Tag(enum.TagNum_FIRMID), 100)
		msg.Body.SetInt(quickfix.Tag(enum.TagNum_FUTUREID),1)
		msg.Body.SetInt(quickfix.Tag(enum.TagNum_QUANTITY),100)
		msg.Body.SetString(quickfix.Tag(enum.TagNum_PRICE), "12.00")
		msg.Body.SetInt(quickfix.Tag(enum.TagNum_OrdType), enum.OrdType_LIMIT)
		msg.Trailer.SetInt(quickfix.Tag(56),1)
		msg.Header.SetInt(quickfix.Tag(56),1)

		err := quickfix.Send(msg)
		if err != nil {
			fmt.Println(err)
		}
	}

	initiator.Stop()
}
