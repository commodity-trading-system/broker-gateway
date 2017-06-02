package executor

import (
	"golang.org/x/net/context"
	"github.com/coreos/etcd/client"
	"strings"
	"strconv"
	"github.com/shopspring/decimal"
)

type Publisher interface {
	Publish(futureId int, buy, sell map[decimal.Decimal]int)
	PublishStatus(id string, status int)
}

type publisher struct {
	etcd client.Client
	kapi client.KeysAPI
}

const PublishKeyBuy  	= "futures/future_id/buy"
const PublishKeySell	= "futures/future_id/sell"
const PublishKeyStatus  = "consignations/id/status"

func (p *publisher) Publish(futureId int,buy,sell map[decimal.Decimal]int)  {

	p.kapi.Set(context.Background(),
		strings.Replace(PublishKeyBuy,"future_id",strconv.Itoa(futureId),1),
		getPricesString(buy),nil)
	p.kapi.Set(context.Background(),
		strings.Replace(PublishKeySell,"future_id",strconv.Itoa(futureId),1),
		getPricesString(sell),nil)
}

func (p *publisher) PublishStatus(id string, status int) ()  {
	p.kapi.Set(context.Background(),
		strings.Replace(PublishKeyStatus,"id",id,1),
		strconv.Itoa(status),nil)
}

func NewPublisher(config client.Config) Publisher  {
	c, err := client.New(config)
	if err != nil {
		return nil
	}
	return &publisher{
		etcd: c,
		kapi: client.NewKeysAPI(c),
	}
}

func getPricesString(data map[decimal.Decimal]int) string {
	res := ""
	for price, quantity := range data {
		res += price.String()+"="+strconv.Itoa(quantity)+ ","
	}
	if res != "" {
		res = res[:len(res)-1]
	}
	return res
}
