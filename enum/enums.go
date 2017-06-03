package enum

import "github.com/shopspring/decimal"

type OrderType int

const (
	OrdType_MARKET                          = 1
	OrdType_LIMIT                           = 2
	OrdType_STOP                            = 3
	OrdType_CANCEL	                        = 4
)

type OrderDirection string

const (
	OrderDirection_SELL    = 2
	OrderDirection_BUY     = 1
)

type TagNum int

const (
	TagNum_FUTUREID		TagNum = 11
	TagNum_FIRMID		TagNum = 12
	TagNum_QUANTITY		TagNum = 13
	TagNum_PRICE		TagNum = 14
	TagNum_DIRECTION	TagNum = 15
	TagNum_OrdType		TagNum = 16
	TagNum_ID		TagNum = 17
)

const (
	ConsignationStatus_CANCELLED = 0
	ConsignationStatus_APPENDING = 1
	ConsignationStatus_PARTIAL = 2
	ConsignationStatus_FINISHED = 3
)


const (
	MatchCreatOrder_RESULT_BUY_MORE = 1
	MatchCreatOrder_RESULT_EQUAL = 0
	MatchCreatOrder_RESULT_SELL_MORE = -1
)

var (
	MAX_PRICE = decimal.New(999999,-2)
	MIN_PRICE = decimal.Zero
)

const MARSHAL_DELI  =  ','

