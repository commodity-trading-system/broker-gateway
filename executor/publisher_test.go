package executor

import (
	"testing"
	"github.com/coreos/etcd/client"
	"strings"
	"os"
	"broker-gateway/enum"
	"github.com/shopspring/decimal"
)

var p Publisher

func TestNewPublisher(t *testing.T) {
	endpoints := strings.FieldsFunc(os.Getenv("ETCD_ENDPOINTS"), func(s rune) bool {
		return s==enum.MARSHAL_DELI
	})
	p = NewPublisher(client.Config{
		Endpoints: endpoints,
	})
	p.Publish(1,map[decimal.Decimal]int{
		decimal.New(5000,-2):100,
	},map[decimal.Decimal]int{
		decimal.New(5100,-2):90,
	})
}
