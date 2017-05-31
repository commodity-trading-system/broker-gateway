package executor

import (
	"broker-gateway/entities"
	"github.com/shopspring/decimal"
	"broker-gateway/enum"
	"container/heap"
)

type OrderBook struct {
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


	// Limit order book
	buyBook *MaxHeap
	sellBook *MinHeap

	// Market order book
	marketBuyBook []*entities.Consignation
	marketSellBook []*entities.Consignation
}

func NewOrderBook(d DB) *OrderBook  {
	return &OrderBook{
		db:d,
		buyBook: NewMaxHeap(),
		sellBook: NewMinHeap(),
		triggerBuyPoint: NewMinHeap(),
		triggerSellPoint: NewMaxHeap(),
		marketBuyBook: make([]*entities.Consignation,0,5),
		marketSellBook: make([]*entities.Consignation,0,5),
		topBuy: decimal.Zero,
		topSell: decimal.Zero,
		lastPrice: decimal.Zero,
		totalBuy: 0,
		totalSell: 0,
	}
}

// At most one consignation will cause matching
// If consignation is matched, return true
func (book *OrderBook) matchLimit(consignation *entities.Consignation, p HeapInterface) bool  {
	currentLevel := p.Top()
	for currentLevel != nil &&
		( (consignation.Direction == enum.OrderDirection_BUY && consignation.Price.Cmp(currentLevel.Price) >= 0) ||
			(consignation.Direction == enum.OrderDirection_SELL && consignation.Price.Cmp(currentLevel.Price) <= 0) ) {

		finished, isRest := book.matchOneLevel(consignation,currentLevel)
		if !isRest {
			book.sellBook.Pop()
		}

		if finished {
			return true
		} else {
			currentLevel = book.buyBook.Top()
		}
	}

	book.updateTopBuyAndSell()

	return consignation.OpenQuantity == 0
}


/**
 *  First return value represent whether consignation.OpenQuantity==0
 *  Second return value represent whether currentLevel Quantity==0
 */
func (book *OrderBook) matchOneLevel(consignation *entities.Consignation, currentLevel *Level) (bool, bool)  {
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
				consignation.Status = enum.ConsignationStatus_FINISHED
				book.db.Save(consignation)
				break
			} else if res == enum.MatchCreatOrder_RESULT_EQUAL {
				consignation.Status = enum.ConsignationStatus_FINISHED
				currentLevel.Consignations[i].Status = enum.ConsignationStatus_FINISHED
				book.db.Save(consignation)
				book.db.Save(currentLevel.Consignations[i])
				break
			} else {
				currentLevel.Consignations[i].Status = enum.ConsignationStatus_FINISHED
				book.db.Save(currentLevel.Consignations[i])
			}
		}
		// Match sell order
	} else {
		for ; i< length;i++  {

			_, res := book.matchAndCreatOrder(currentLevel.Consignations[i], currentLevel.Consignations[i],currentLevel.Price)

			if res == enum.MatchCreatOrder_RESULT_BUY_MORE {
				consignation.Status = enum.ConsignationStatus_FINISHED
				book.db.Save(consignation)
				break
			} else if res == enum.MatchCreatOrder_RESULT_EQUAL {
				consignation.Status = enum.ConsignationStatus_FINISHED
				currentLevel.Consignations[i].Status = enum.ConsignationStatus_FINISHED
				book.db.Save(consignation)
				book.db.Save(currentLevel.Consignations[i])
				break
			} else {
				currentLevel.Consignations[i].Status = enum.ConsignationStatus_FINISHED
				book.db.Save(currentLevel.Consignations[i])
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

	return consignation.OpenQuantity==0, currentLevel.Consignations[length-1].OpenQuantity == 0
}


func (book *OrderBook) scanTriggerPoint(start, end decimal.Decimal)  {
	// Price increased, only stop buy can be triggered
	if end.Cmp(start) > 0 {
		currentLevel := book.triggerBuyPoint.Top()
		for currentLevel.Price.Cmp(end) < 0  {
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
func (book *OrderBook) addLimit(consignation *entities.Consignation)  {
	// TODO Insert consignation according to it's priority
	lastDealPrice := book.lastPrice
	if consignation.Direction == enum.OrderDirection_BUY {

		// Won't cause match
		if consignation.Price.Cmp(book.topSell) < 0 {
			heap.Push(book.buyBook, Level{
				Price: consignation.Price,
				Consignations:[]*entities.Consignation{consignation},
			})
			return
		}

		matchFinish := book.matchLimit(consignation,book.sellBook)
		if !matchFinish {
			book.db.Save(consignation)
			heap.Push(book.buyBook, Level{
				Price: consignation.Price,
				Consignations:[]*entities.Consignation{consignation},
			})
		}
		if book.lastPrice.Cmp(lastDealPrice) > 0 {
			book.scanTriggerPoint(lastDealPrice,book.lastPrice)
		}

	} else if consignation.Direction == enum.OrderDirection_SELL {
		// Won't cause match
		if consignation.Price.Cmp(book.topSell) > 0 {
			heap.Push(book.sellBook, Level{
				Price: consignation.Price,
				Consignations:[]*entities.Consignation{consignation},
			})
			return
		}

		matchFinish :=book.matchLimit(consignation,book.buyBook)
		if !matchFinish {
			book.db.Save(consignation)
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

// Process market order
func (book *OrderBook) addMarket(consignation *entities.Consignation)  {
	if consignation.Direction == enum.OrderDirection_BUY {
		if (book.sellBook.Len()) == 0 {
			// TODO invalid market order, should be cancelled
		} else {
			consignation.Price = enum.MAX_PRICE
			finished := book.matchLimit(consignation, book.sellBook)
			if !finished {
				book.marketBuyBook = append(book.marketBuyBook, consignation)
			}
			book.db.Save(consignation)

		}
	} else if consignation.Direction == enum.OrderDirection_SELL {
		if (book.buyBook.Len()) == 0 {
			// TODO invalid market order, should be cancelled
		} else {
			consignation.Price = enum.MIN_PRICE
			finished := book.matchLimit(consignation, book.buyBook)
			if !finished {
				book.marketBuyBook = append(book.marketBuyBook, consignation)
			}
			book.db.Save(consignation)
		}
	}
}

func (book *OrderBook) addStop(consignation *entities.Consignation)  {
	if consignation.Direction == enum.OrderDirection_BUY {
		if consignation.Price.Cmp(book.topSell) <= 0 {
			// TODO invalid stop buy order
		} else {
			heap.Push(book.triggerBuyPoint,Level{
				Price: consignation.Price,
				Consignations: []*entities.Consignation{consignation},
			})
		}
	} else if consignation.Direction == enum.OrderDirection_SELL {
		if consignation.Price.Cmp(book.topBuy) >= 0 {
			// TODO invalid stop buy order
		} else {
			heap.Push(book.triggerSellPoint,Level{
				Price: consignation.Price,
				Consignations: []*entities.Consignation{consignation},
			})
		}
	}
}

func (book *OrderBook) cancel(consignation *entities.Consignation)  {
	// TODO
}

// Return 1 if buyConsignation's quantity > 0 after mactching
// Return 0 if buyConsignation's quantity = 0 and sellConsignation's quantity = 0 after mactching
// Return -1 if buyConsignation's quantity < 0 after mactching
func (book *OrderBook) matchAndCreatOrder(buyConsignation *entities.Consignation, sellConsignation *entities.Consignation, price decimal.Decimal) (*entities.Order, int)  {
	var quantity int
	var res int
	if buyConsignation.OpenQuantity > sellConsignation.OpenQuantity {
		quantity = sellConsignation.OpenQuantity
		quantity = sellConsignation.OpenQuantity
		sellConsignation.OpenQuantity = 0
		buyConsignation.OpenQuantity -= sellConsignation.OpenQuantity

		res = enum.MatchCreatOrder_RESULT_BUY_MORE
	} else if buyConsignation.OpenQuantity == sellConsignation.OpenQuantity {
		quantity = sellConsignation.OpenQuantity
		sellConsignation.OpenQuantity = 0
		buyConsignation.OpenQuantity = 0
		res = enum.MatchCreatOrder_RESULT_EQUAL
	} else {
		quantity = buyConsignation.OpenQuantity
		buyConsignation.OpenQuantity = 0
		sellConsignation.OpenQuantity -= buyConsignation.OpenQuantity
		res = enum.MatchCreatOrder_RESULT_SELL_MORE
	}
	order := &entities.Order{
		BuyerId: buyConsignation.FirmId,
		SellerId: sellConsignation.FirmId,
		FutureId: buyConsignation.FutureId,
		BuyerConsignationId: buyConsignation.ID,
		SellerConsignationId: sellConsignation.ID,
		Quantity: quantity,
		Price: price,
	}
	book.db.Save(order)
	return order,res
}

func (book *OrderBook) updateTopBuyAndSell()  {
	if book.buyBook.Top() != nil {
		book.topBuy = book.buyBook.Top().Price
	}
	if book.sellBook.Top() != nil {
		book.topSell = book.sellBook.Top().Price
	}
}
