package enum


type OrderType int

const (
	OrdType_MARKET                          = 1
	OrdType_LIMIT                           = 2
	OrdType_STOP                            = 3
	OrdType_CANCEL	                        = 4
)

type OrderDirection string

const (
	OrderDirection_SELL    OrderDirection = "0"
	OrderDirection_BUY     OrderDirection = "1"
)

type TagNum int

const (
	TagNum_FIRMID		TagNum = 11
	TagNum_FUTUREID		TagNum = 12
	TagNum_QUANTITY		TagNum = 13
	TagNum_PRICE		TagNum = 14
	TagNum_DIRECTION	TagNum = 15
	TagNum_OrdType		TagNum = 16
)


