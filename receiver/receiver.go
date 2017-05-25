package receiver

import (
	"github.com/quickfixgo/quickfix"
	"github.com/go-redis/redis"
	"github.com/quickfixgo/quickfix/fix44/newordersingle"
	"github.com/quickfixgo/quickfix/enum"
	"github.com/quickfixgo/quickfix/tag"
	"broker-gateway/entities"
)

type Receiver struct {
	redisClient redis.Client
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
			Addr:     config.RedisHost+":"+config.RedisPort,
			Password: config.RedisPwd,
			DB:       config.RedisDB,
		}),
		MessageRouter: quickfix.NewMessageRouter(),
	}
	r.AddRoute(newordersingle.Route(r.OnNewOrderSingle))
	return r
}

func (r*Receiver) OnCreate(sessionID quickfix.SessionID)()  { return }
func (r*Receiver) OnLogon(sessionID quickfix.SessionID)()  { return }
func (r*Receiver) OnLogout(sessionID quickfix.SessionID)()  { return }
func (r*Receiver) ToAdmin(message quickfix.Message, sessionID quickfix.SessionID)()  { return}
func (r*Receiver) ToApp(message quickfix.Message, sessionID quickfix.SessionID) error { return nil}
func (r*Receiver) FromAdmin(message quickfix.Message, sessionID quickfix.SessionID) quickfix.MessageRejectError { return nil }

func (r*Receiver) FromApp(message quickfix.Message, sessionID quickfix.SessionID) quickfix.MessageRejectError  {
	return r.Route(message, sessionID)
}

func (r*Receiver) OnNewOrderSingle(msg newordersingle.NewOrderSingle, session quickfix.SessionID)  quickfix.MessageRejectError {
	var futureId quickfix.FIXInt

	// Tag40
	ordType, err := msg.GetOrdType()
	if err != nil {
		return err
	}

	if ordType != enum.OrdType_LIMIT &&
		ordType != enum.OrdType_MARKET &&
		ordType != enum.OrdType_STOP &&
		ordType != enum.OrdType_STOP_LIMIT {
		return quickfix.ValueIsIncorrect(tag.OrdType)
	}

	// Tag 38
	orderQty, err := msg.GetOrderQty()
	if err != nil {
		return err
	}

	// Tag 44
	price, err := msg.GetPrice()
	if err != nil {
		return err
	}

	err = msg.GetField(quickfix.Tag(40),&futureId)
	if err != nil {
		return err
	}

	order := entities.Consignation{
		Price: price,
		Quantity: orderQty,
	}

	r.redisClient.RPush("future_"+futureId, order)
	return nil
}