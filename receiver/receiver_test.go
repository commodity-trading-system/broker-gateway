package receiver

import (
	"testing"
	"os"
	"strconv"
	"github.com/quickfixgo/quickfix"
	"broker-gateway/enum"
	"fmt"
)

var receiver *Receiver

func TestNewReceiver(t *testing.T) {
	port,_ := strconv.ParseInt(os.Getenv("REDIS_PORT"),10,32)
	db,_ := strconv.ParseInt(os.Getenv("REDIS_DB"),10,32)
	config := ReceiverConfig{
		RedisHost: os.Getenv("REDIS_HOST"),
		RedisPort: int(port),
		RedisPwd: os.Getenv("REDIS_PASSWORD"),
		RedisDB: int(db),
	}
	fmt.Println(config)

	receiver = NewReceiver(config)
}

func TestReceiver_FromApp(t *testing.T) {
	msg := quickfix.NewMessage()

	msg.Body.SetInt(quickfix.Tag(enum.TagNum_DIRECTION), enum.OrderDirection_BUY)
	msg.Body.SetInt(quickfix.Tag(enum.TagNum_FIRMID), 100)
	msg.Body.SetInt(quickfix.Tag(enum.TagNum_FUTUREID),1)
	msg.Body.SetInt(quickfix.Tag(enum.TagNum_QUANTITY),100)
	msg.Body.SetString(quickfix.Tag(enum.TagNum_PRICE), "12.00")
	msg.Body.SetInt(quickfix.Tag(enum.TagNum_OrdType), enum.OrdType_LIMIT)
	session := quickfix.SessionID{}
	err := receiver.FromApp(msg, session)
	if err != nil {
		t.Error(err)
	}
}
