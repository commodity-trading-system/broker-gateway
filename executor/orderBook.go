package executor

import (
	"broker-gateway/entities"
	"github.com/shopspring/decimal"
	"broker-gateway/enum"
	"container/heap"
	"github.com/satori/go.uuid"
	"strconv"
)

type OrderBook interface {
	Reset()
	AddLimit(consignation *entities.Consignation)
	AddMarket(consignation *entities.Consignation)
	AddStop(consignation *entities.Consignation)
	AddCancel(consignation *entities.Consignation)
}

type PublishCallback  func(futureId int,buy,sell map[decimal.Decimal]int)

type orderBook struct {
	// db
	db DB

	// For stop order
	triggerBuyPoint *MinHeap
	triggerSellPoint *MaxHeap

	topSell decimal.Decimal
	topBuy decimal.Decimal
	lastPrice decimal.Decimal
	totalBuy uint
	totalSell uint

	futureId int


	// Limit order book
	buyBook *MaxHeap
	sellBook *MinHeap

	// Market order book
	marketBuyBook []*entities.Consignation
	marketSellBook []*entities.Consignation

	publisher Publisher
}

func NewOrderBook(futureId int, d DB, publisher Publisher) OrderBook  {
	return &orderBook{
		db:d,
		buyBook: NewMaxHeap(),
		sellBook: NewMinHeap(),
		triggerBuyPoint: NewMinHeap(),
		triggerSellPoint: NewMaxHeap(),
		marketBuyBook: make([]*entities.Consignation,0,5),
		marketSellBook: make([]*entities.Consignation,0,5),
		topBuy: enum.MIN_PRICE,
		topSell: enum.MAX_PRICE,
		lastPrice: decimal.Zero,
		totalBuy: 0,
		totalSell: 0,
		publisher: publisher,
		futureId: futureId,
	}
}

// At most one consignation will cause matching
// If consignation is matched, return true
func (book *orderBook) matchLimit(consignation *entities.Consignation, p HeapInterface) bool  {
	currentLevel := p.Top()
	for currentLevel != nil &&
		( (consignation.Direction == enum.OrderDirection_BUY && consignation.Price.Cmp(currentLevel.Price) >= 0) ||
			(consignation.Direction == enum.OrderDirection_SELL && consignation.Price.Cmp(currentLevel.Price) <= 0) ) {

		finished, isRest := book.matchOneLevel(consignation,currentLevel)
		if !isRest {
			heap.Pop(p)
		}

		if finished {
			return true
		} else {
			currentLevel = p.Top()
		}
	}

	book.updateTopBuyAndSell()

	return consignation.OpenQuantity == 0
}


/**
 *  First return value represent whether consignation.OpenQuantity==0
 *  Second return value represent whether currentLevel Quantity==0
 */
func (book *orderBook) matchOneLevel(consignation *entities.Consignation, currentLevel *Level) (bool, bool)  {
	length := len(currentLevel.Consignations)
	i := 0
	if length == 0 {
		return false, true
	}
	//var buyConsignation *entities.Consignation
	//var sellConsignation *entities.Consignation
	// Match buy order
	if consignation.Direction == enum.OrderDirection_BUY {
		for ; i< length;i++  {

			_, res := book.matchAndCreatOrder(consignation, currentLevel.Consignations[i],currentLevel.Price)
			if res == enum.MatchCreatOrder_RESULT_SELL_MORE {
				//consignation.Status = enum.ConsignationStatus_FINISHED
				//book.db.Save(consignation)
				//book.publishStatus(consignation)
				break
			} else if res == enum.MatchCreatOrder_RESULT_EQUAL {
				//consignation.Status = enum.ConsignationStatus_FINISHED
				//currentLevel.Consignations[i].Status = enum.ConsignationStatus_FINISHED
				//book.db.Save(consignation)
				//book.db.Save(currentLevel.Consignations[i])
				//book.publishStatus(consignation)
				//book.publishStatus(currentLevel.Consignations[i])
				break
			} else {
				//currentLevel.Consignations[i].Status = enum.ConsignationStatus_FINISHED
				//book.db.Save(currentLevel.Consignations[i])
				//book.publishStatus(currentLevel.Consignations[i])
			}

		}
		// Match sell order
	} else {
		for ; i< length;i++  {

			_, res := book.matchAndCreatOrder(currentLevel.Consignations[i], consignation,currentLevel.Price)

			if res == enum.MatchCreatOrder_RESULT_BUY_MORE {
				//currentLevel.Consignations[i].Status = enum.ConsignationStatus_FINISHED
				//book.db.Save(currentLevel.Consignations[i])
				//book.publishStatus(currentLevel.Consignations[i])
				break
			} else if res == enum.MatchCreatOrder_RESULT_EQUAL {
				//consignation.Status = enum.ConsignationStatus_FINISHED
				//currentLevel.Consignations[i].Status = enum.ConsignationStatus_FINISHED
				//book.db.Save(consignation)
				//book.db.Save(currentLevel.Consignations[i])
				//book.publishStatus(consignation)
				break
			} else {
				//currentLevel.Consignations[i].Status = enum.ConsignationStatus_FINISHED
				//book.db.Save(currentLevel.Consignations[i])
				//book.publishStatus(currentLevel.Consignations[i])
			}
		}
	}

	// Remove finished consignations from current level
	for ; i< length ; i++ {
		if currentLevel.Consignations[i].OpenQuantity > 0 {
			currentLevel.Consignations = currentLevel.Consignations[i:]
		}
	}

	book.lastPrice = currentLevel.Price

	return consignation.OpenQuantity==0, currentLevel.Consignations[length-1].OpenQuantity > 0
}


func (book *orderBook) scanTriggerPoint(start, end decimal.Decimal)  {
	return
	// Price increased, only stop buy can be triggered
	if end.Cmp(start) > 0 {
		currentLevel := book.triggerBuyPoint.Top()
		for currentLevel != nil && currentLevel.Price.Cmp(end) < 0  {
			heap.Push(book.buyBook,currentLevel)
			book.triggerBuyPoint.Pop()
			currentLevel = book.triggerBuyPoint.Top()
		}
		// Match stop order with limit order which have the same price
		if currentLevel.Price.Cmp(end) == 0 {

			// Try to match
			i := 0
			for ; i<len(currentLevel.Consignations) ;i++  {
				finish := book.matchLimit(currentLevel.Consignations[i],book.sellBook)
				if finish {
					break
				}
			}
			// Push extra stop order into buy order book
			for ; i< len(currentLevel.Consignations); i++ {
				if currentLevel.Consignations[i].OpenQuantity > 0 {
					heap.Push(book.buyBook,Level{
						Price: currentLevel.Price,
						Consignations: currentLevel.Consignations[i:],
					})
				}
			}
		}

		// Price decreased, only stop sell can be triggered
	} else if end.Cmp(start) < 0 {
		currentLevel := book.triggerSellPoint.Top()
		for currentLevel.Price.Cmp(end) > 0  {
			heap.Push(book.sellBook,currentLevel)
			book.triggerSellPoint.Pop()
			currentLevel = book.triggerSellPoint.Top()
		}
		// Match stop order with limit order which have the same price
		if currentLevel.Price.Cmp(end) == 0 {

			// Try to match
			i := 0
			for ; i<len(currentLevel.Consignations) ;i++  {
				finish := book.matchLimit(currentLevel.Consignations[i],book.buyBook)
				if finish {
					break
				}
			}
			// Push extra stop order into sell order book
			for ; i< len(currentLevel.Consignations); i++ {
				if currentLevel.Consignations[i].OpenQuantity > 0 {
					heap.Push(book.sellBook,Level{
						Price: currentLevel.Price,
						Consignations: currentLevel.Consignations[i:],
					})
				}
			}
		}

	}
}

// Porcess limit order
func (book *orderBook) AddLimit(consignation *entities.Consignation)  {
	// TODO Insert consignation according to it's priority
	lastDealPrice := book.lastPrice
	if consignation.Direction == enum.OrderDirection_BUY {

		// Won't cause match
		if book.sellBook.Len() == 0 || consignation.Price.Cmp(book.topSell) < 0 {
			heap.Push(book.buyBook, Level{
				Price: consignation.Price,
				Consignations:[]*entities.Consignation{consignation},
			})

			book.db.Save(&consignation)
			book.publishStatus(consignation)
		} else {
			matchFinish := book.matchLimit(consignation,book.sellBook)
			book.db.Save(consignation)
			book.publishStatus(consignation)
			if !matchFinish {
				heap.Push(book.buyBook, Level{
					Price: consignation.Price,
					Consignations:[]*entities.Consignation{consignation},
				})
			}
			if book.lastPrice.Cmp(lastDealPrice) > 0 {
				book.scanTriggerPoint(lastDealPrice,book.lastPrice)
			}
		}

	} else if consignation.Direction == enum.OrderDirection_SELL {
		// Won't cause match
		if book.buyBook.Len() == 0 || consignation.Price.Cmp(book.topSell) > 0 {
			heap.Push(book.sellBook, Level{
				Price: consignation.Price,
				Consignations:[]*entities.Consignation{consignation},
			})
			book.db.Save(consignation)
			book.publishStatus(consignation)
			return
		} else {
			matchFinish :=book.matchLimit(consignation,book.buyBook)
			book.db.Save(consignation)
			book.publishStatus(consignation)
			if !matchFinish {
				heap.Push(book.sellBook, Level{
					Price: consignation.Price,
					Consignations:[]*entities.Consignation{consignation},
				})
			}

			if book.lastPrice.Cmp(lastDealPrice) < 0 {
				book.scanTriggerPoint(book.lastPrice,lastDealPrice)
			}
		}
	}

	book.updateTopBuyAndSell()
}

// Process market order
func (book *orderBook) AddMarket(consignation *entities.Consignation)  {
	if consignation.Direction == enum.OrderDirection_BUY {
		if (book.sellBook.Len()) == 0 {
			consignation.Status = enum.ConsignationStatus_INVALID
			book.db.Save(consignation)
			book.publishStatus(consignation)
		} else {

			consignation.OpenQuantity = consignation.Quantity
			consignation.Price = enum.MAX_PRICE
			finished := book.matchLimit(consignation, book.sellBook)
			if !finished {
				book.marketBuyBook = append(book.marketBuyBook, consignation)
			}
			book.db.Save(consignation)
			book.publishStatus(consignation)

		}
	} else if consignation.Direction == enum.OrderDirection_SELL {
		if (book.buyBook.Len()) == 0 {
			consignation.Status = enum.ConsignationStatus_INVALID
			book.db.Save(consignation)
			book.publishStatus(consignation)
		} else {
			consignation.OpenQuantity = consignation.Quantity
			consignation.Price = enum.MIN_PRICE
			finished := book.matchLimit(consignation, book.buyBook)
			if !finished {
				book.marketBuyBook = append(book.marketBuyBook, consignation)
			}
			book.db.Save(consignation)
			book.publishStatus(consignation)
		}
	}
}

func (book *orderBook) AddStop(consignation *entities.Consignation)  {
	if consignation.Direction == enum.OrderDirection_BUY {
		if consignation.Price.Cmp(book.topSell) <= 0 {
			consignation.Status = enum.ConsignationStatus_INVALID
		} else {
			heap.Push(book.triggerBuyPoint,Level{
				Price: consignation.Price,
				Consignations: []*entities.Consignation{consignation},
			})
		}
	} else if consignation.Direction == enum.OrderDirection_SELL {
		if consignation.Price.Cmp(book.topBuy) >= 0 {
			consignation.Status = enum.ConsignationStatus_INVALID
		} else {
			heap.Push(book.triggerSellPoint,Level{
				Price: consignation.Price,
				Consignations: []*entities.Consignation{consignation},
			})
		}
	}

	book.db.Save(consignation)
	book.publishStatus(consignation)
}

func (book *orderBook) cancelLevel(consignation *entities.Consignation,consignations []*entities.Consignation) bool {
	for i:=0;i<len(consignations);i++ {
		if consignations[i].ID == consignation.ID {
			consignations[i].Status = enum.ConsignationStatus_CANCELLED
			book.publishStatus(consignations[i])
			book.db.Save(consignations[i])
			consignations = append(consignations[:i], consignations[i+1:]...)
			return true
		}
	}
	return false
}

func (book *orderBook) AddCancel(consignation *entities.Consignation)  {
	if consignation.Direction == enum.OrderDirection_BUY {
		book.buyBook.Travel(func(level *Level) bool {
			if book.cancelLevel(consignation, level.Consignations) {
				return false
			}
			return true
		})
		book.cancelLevel(consignation,book.marketBuyBook)
		book.triggerBuyPoint.Travel(func(level *Level) bool {
			if book.cancelLevel(consignation, level.Consignations) {
				return false
			}
			return true
		})
	} else {
		book.sellBook.Travel(func(level *Level) bool {
			if book.cancelLevel(consignation, level.Consignations) {
				return false
			}
			return true
		})
		book.cancelLevel(consignation,book.marketSellBook)
		book.triggerBuyPoint.Travel(func(level *Level) bool {
			if book.cancelLevel(consignation, level.Consignations) {
				return false
			}
			return true
		})
		//
		//if consignation.Type == enum.OrdType_LIMIT {
		//
		//} else if consignation.Type == enum.OrdType_MARKET {
		//
		//} else if consignation.Type == enum.OrdType_STOP {
		//
		//}
	}
}

func (book *orderBook) Reset()  {
	book.db.Empty()
	book.db.Migrate()
	book = &orderBook{
		db:book.db,
		buyBook: NewMaxHeap(),
		sellBook: NewMinHeap(),
		triggerBuyPoint: NewMinHeap(),
		triggerSellPoint: NewMaxHeap(),
		marketBuyBook: make([]*entities.Consignation,0,5),
		marketSellBook: make([]*entities.Consignation,0,5),
		topBuy: enum.MIN_PRICE,
		topSell: enum.MAX_PRICE,
		lastPrice: decimal.Zero,
		totalBuy: 0,
		totalSell: 0,
	}
}

// Return 1 if buyConsignation's quantity > 0 after mactching
// Return 0 if buyConsignation's quantity = 0 and sellConsignation's quantity = 0 after mactching
// Return -1 if buyConsignation's quantity < 0 after mactching
func (book *orderBook) matchAndCreatOrder(buyConsignation *entities.Consignation, sellConsignation *entities.Consignation, price decimal.Decimal) (*entities.Order, int)  {
	var quantity int
	var res int
	if buyConsignation.OpenQuantity > sellConsignation.OpenQuantity {
		quantity = sellConsignation.OpenQuantity
		sellConsignation.OpenQuantity = 0
		buyConsignation.OpenQuantity -= quantity
		sellConsignation.Status = enum.ConsignationStatus_FINISHED
		if buyConsignation.OpenQuantity > 0 {
			buyConsignation.Status = enum.ConsignationStatus_PARTIAL
		} else {
			buyConsignation.Status = enum.ConsignationStatus_FINISHED
		}
		res = enum.MatchCreatOrder_RESULT_BUY_MORE
	} else if buyConsignation.OpenQuantity == sellConsignation.OpenQuantity {
		quantity = sellConsignation.OpenQuantity
		sellConsignation.OpenQuantity = 0
		buyConsignation.OpenQuantity = 0
		res = enum.MatchCreatOrder_RESULT_EQUAL
		sellConsignation.Status = enum.ConsignationStatus_FINISHED
		buyConsignation.Status = enum.ConsignationStatus_FINISHED
	} else {
		quantity = buyConsignation.OpenQuantity
		buyConsignation.OpenQuantity = 0
		sellConsignation.OpenQuantity -= quantity
		res = enum.MatchCreatOrder_RESULT_SELL_MORE
		buyConsignation.Status = enum.ConsignationStatus_FINISHED
		if sellConsignation.OpenQuantity > 0 {
			sellConsignation.Status = enum.ConsignationStatus_PARTIAL
		} else {
			sellConsignation.Status = enum.ConsignationStatus_FINISHED
		}
	}
	order := &entities.Order{
		BuyerId: buyConsignation.FirmId,
		SellerId: sellConsignation.FirmId,
		FutureId: buyConsignation.FutureId,
		BuyerConsignationId: buyConsignation.ID,
		SellerConsignationId: sellConsignation.ID,
		Quantity: quantity,
		Price: price,
		Commission: decimal.Zero,
		Status:1,
		ID: uuid.NewV1(),
	}
	quotation := &entities.Quotation{
		FutureId: buyConsignation.FutureId,
		Price: price,
	}

	book.publisher.PublishLatestPrice(strconv.Itoa(order.FutureId),price)
	book.db.Save(order)
	book.db.Save(quotation)
	book.db.Save(buyConsignation)
	book.db.Save(sellConsignation)
	book.publishStatus(buyConsignation)
	book.publishStatus(sellConsignation)
	return order,res
}

func (book *orderBook) updateTopBuyAndSell()  {
	if book.buyBook.Top() != nil {
		book.topBuy = book.buyBook.Top().Price
	} else {
		book.topBuy = enum.MIN_PRICE
	}
	if book.sellBook.Top() != nil {
		book.topSell = book.sellBook.Top().Price
	} else {
		book.topSell = enum.MAX_PRICE
	}
}

func (book *orderBook) publishDepth()  {
	if book.publisher != nil {
		buy := map[decimal.Decimal]int{}
		sell := map[decimal.Decimal]int{}
		book.buyBook.Travel(func(level *Level) bool {
			quantity := 0
			for i:=0; i<len(level.Consignations) ;i++  {
				quantity += level.Consignations[i].OpenQuantity
			}
			buy[level.Price] = quantity
			return true
		})
		book.sellBook.Travel(func(level *Level) bool {
			quantity := 0
			for i:=0; i<len(level.Consignations) ;i++  {
				quantity += level.Consignations[i].OpenQuantity
			}
			sell[level.Price] = quantity
			return  true
		})
		book.publisher.Publish(book.futureId, buy, sell)
	}

}

func(book *orderBook) publishStatus (consignation *entities.Consignation)  {
	if book.publisher != nil {
		book.publisher.PublishStatus(consignation.ID.String(),consignation.Status)
	}

}
