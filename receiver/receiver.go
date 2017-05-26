package receiver

import (
	"github.com/quickfixgo/quickfix"
	"github.com/go-redis/redis"
	//"github.com/quickfixgo/quickfix/fix44/newordersingle"
	//"github.com/quickfixgo/quickfix/enum"
	"github.com/quickfixgo/quickfix/tag"
	"broker-gateway/entities"
	"broker-gateway/enum"
	"fmt"
)

type Receiver struct {
	redisClient *redis.Client
	*quickfix.MessageRouter
}

type ReceiverConfig struct {
	RedisHost string
	RedisPort int
	RedisPwd string
	RedisDB int
}

func NewReceiver(config ReceiverConfig) *Receiver {
	r := &Receiver{
		redisClient: redis.NewClient(&redis.Options{
			Addr:     config.RedisHost+":"+string(config.RedisPort),
			Password: config.RedisPwd,
			DB:       config.RedisDB,
		}),
		MessageRouter: quickfix.NewMessageRouter(),
	}
	return r
}

func (r*Receiver) OnCreate(sessionID quickfix.SessionID)()  { return }
func (r*Receiver) OnLogon(sessionID quickfix.SessionID)()  { return }
func (r*Receiver) OnLogout(sessionID quickfix.SessionID)()  { return }
func (r*Receiver) ToAdmin(message *quickfix.Message, sessionID quickfix.SessionID)()  { return}
func (r*Receiver) ToApp(message *quickfix.Message, sessionID quickfix.SessionID) error { return nil}
func (r*Receiver) FromAdmin(message *quickfix.Message, sessionID quickfix.SessionID) quickfix.MessageRejectError { return nil }

func (r*Receiver) FromApp(msg *quickfix.Message, sessionID quickfix.SessionID) quickfix.MessageRejectError  {
	var futureId quickfix.FIXString
	var direction quickfix.FIXInt
	var firmId quickfix.FIXInt
	var quantity quickfix.FIXInt
	var price quickfix.FIXDecimal
	var orderType quickfix.FIXInt
	var err quickfix.MessageRejectError

	err = msg.Body.GetField(quickfix.Tag(enum.TagNum_FUTUREID),&futureId)
	if err != nil {
		return err
	}

	err = msg.Body.GetField(quickfix.Tag(enum.TagNum_DIRECTION),&direction)
	if err != nil {
		return err
	}

	err = msg.Body.GetField(quickfix.Tag(enum.TagNum_FIRMID),&firmId)
	if err != nil {
		return err
	}

	err = msg.Body.GetField(quickfix.Tag(enum.TagNum_PRICE),&price)
	if err != nil {
		return err
	}

	err = msg.Body.GetField(quickfix.Tag(enum.TagNum_QUANTITY),&quantity)
	if err != nil {
		return err
	}

	err = msg.Body.GetField(quickfix.Tag(enum.TagNum_OrdType),&orderType)
	if err != nil {
		return err
	}


	if orderType.Int() != enum.OrdType_CANCEL &&
		orderType.Int() != enum.OrdType_MARKET &&
		orderType.Int() != enum.OrdType_STOP &&
		orderType.Int() != enum.OrdType_LIMIT {
		return quickfix.ValueIsIncorrect(tag.OrdType)
	}


	order := entities.Consignation{
		Price: price.Decimal,
		Quantity: quantity.Int(),
		Direction: direction.Int(),
		FirmId: firmId.Int(),
	}


	intCmd := r.redisClient.RPush("future_"+futureId.String(), order)
	if intCmd.Err() != nil {
		return quickfix.NewMessageRejectError(intCmd.String(),0,nil)
	}
	fmt.Println(intCmd.Err())
	return nil
}

//func (r*Receiver) OnNewOrderSingle(msg newordersingle.NewOrderSingle, session quickfix.SessionID)  quickfix.MessageRejectError {
//
//}
