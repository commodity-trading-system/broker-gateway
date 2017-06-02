package receiver

import (
	"github.com/quickfixgo/quickfix"
	"github.com/go-redis/redis"
	"broker-gateway/entities"
	"broker-gateway/enum"
	"fmt"
	"strconv"
	"github.com/satori/go.uuid"
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

	fmt.Println(config.RedisHost+":"+strconv.Itoa(config.RedisPort))
	r := &Receiver{
		redisClient: redis.NewClient(&redis.Options{
			Addr:     config.RedisHost+":"+strconv.Itoa(config.RedisPort),
			Password: config.RedisPwd,
			DB:       config.RedisDB,
		}),
		MessageRouter: quickfix.NewMessageRouter(),
	}
	return r
}

func (r*Receiver) OnCreate(sessionID quickfix.SessionID)()  { return }
func (r*Receiver) OnLogon(sessionID quickfix.SessionID)()  {
	fmt.Println("logon")
	return
}
func (r*Receiver) OnLogout(sessionID quickfix.SessionID)()  {
	fmt.Println("logout")
	return
}
func (r*Receiver) ToAdmin(message *quickfix.Message, sessionID quickfix.SessionID)()  { return}
func (r*Receiver) ToApp(message *quickfix.Message, sessionID quickfix.SessionID) error { return nil}
func (r*Receiver) FromAdmin(message *quickfix.Message, sessionID quickfix.SessionID) quickfix.MessageRejectError { return nil }

func (r*Receiver) FromApp(msg *quickfix.Message, sessionID quickfix.SessionID) quickfix.MessageRejectError  {
	fmt.Println("收到消息",msg)
	var futureId quickfix.FIXInt
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
		return nil
	}


	consignation := entities.Consignation{
		ID: uuid.NewV1(),
		Type: orderType.Int(),
		FutureId: futureId.Int(),
		Price: price.Decimal,
		Quantity: quantity.Int(),
		Direction: direction.Int(),
		FirmId: firmId.Int(),
		Status: enum.ConsignationStatus_APPENDING,
	}

	intCmd := r.redisClient.RPush("future_"+strconv.Itoa(futureId.Int()), consignation)
	if intCmd.Err() != nil {
		return quickfix.NewMessageRejectError(intCmd.String(),0,nil)
	}

	nmsg := quickfix.NewMessage()
	id := quickfix.FIXBytes{}
	id.Read([]byte(consignation.ID.String()))
	nmsg.Body.SetField(quickfix.Tag(enum.TagNum_ID),id)
	quickfix.SendToTarget(nmsg, sessionID)
	return nil
}
