package executor

import (
	"broker-gateway/entities"
	"github.com/shopspring/decimal"
	"broker-gateway/enum"
	"container/heap"
)

type OrderBook struct {
	// db
	db *DB

	// For stop order
	//triggerPoint map[decimal.Decimal][]*entities.Consignation
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

	//buyBook map[decimal.Decimal][]*entities.Consignation
	//sellBook map[decimal.Decimal][]*entities.Consignation
	// Market order book
	marketBuyBook []*entities.Consignation
	marketSellBook []*entities.Consignation
}

func NewOrderBook(d *DB) *OrderBook  {
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

	return consignation.OpenQuantity == 0
}


/**
 *  First return value represent whether consignation.OpenQuantity==0
 *  Second return value represent whether currentLevel Quantity==0
 */
func (book *OrderBook) matchOneLevel(consignation *entities.Consignation, currentLevel *Level) (bool, bool)  {
	length := len(currentLevel.Consignations)
	for i:=0; i< length;i++  {
		if currentLevel.Consignations[i].OpenQuantity < consignation.OpenQuantity {
			// TODO currentSellLevel.Consignations[i] finish

			consignation.OpenQuantity -= currentLevel.Consignations[i].OpenQuantity
		} else {
			currentLevel.Consignations[i].OpenQuantity -= consignation.OpenQuantity
			// TODO consignation finish
			if currentLevel.Consignations[i].OpenQuantity == 0 {
				// TODO currentSellLevel.Consignations[i] finish

			}
			break
		}
	}
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

	//for price, consignations := range book.triggerPoint {
	//	// Price increased, only stop buy can be triggered
	//	if end.Cmp(start) > 0  && price.Cmp(start) > 0 && price.Cmp(end) <= 0 {
	//		buyConsignations := make([]*entities.Consignation,0,3)
	//		sellConsignations := make([]*entities.Consignation,0,3)
	//		for i:=0; i<len(consignations) ;i++  {
	//			if consignations[i].Direction == enum.OrderDirection_BUY {
	//				buyConsignations = append(buyConsignations,consignations[i])
	//			} else if consignations[i].Direction == enum.OrderDirection_SELL {
	//				sellConsignations = append(sellConsignations,consignations[i])
	//			}
	//		}
	//		if len(buyConsignations) > 0 {
	//			heap.Push(book.buyBook,Level{
	//				Price:price,
	//				Consignations: buyConsignations,
	//			})
	//		}
	//		if len(sellConsignations) > 0 {
	//			heap.Push(book.buyBook,Level{
	//				Price:price,
	//				Consignations: sellConsignations,
	//			})
	//		}
	//	}
	//}
}

// Porcess limit order
func (book *OrderBook) addLimit(consignation *entities.Consignation)  {
	// TODO Insert consignation according to it's priority
	if consignation.Direction == enum.OrderDirection_BUY {

		// Won't cause match
		if consignation.Price.Cmp(book.topSell) < 0 {
			heap.Push(book.buyBook, Level{
				Price: consignation.Price,
				Consignations:[]*entities.Consignation{consignation},
			})
			return
		}

		book.matchLimit(consignation,book.sellBook)


		if consignation.OpenQuantity > 0 {
			heap.Push(book.buyBook, Level{
				Price: consignation.Price,
				Consignations:[]*entities.Consignation{consignation},
			})
			return
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

		book.matchLimit(consignation,book.buyBook)

		if consignation.OpenQuantity > 0 {
			heap.Push(book.sellBook, Level{
				Price: consignation.Price,
				Consignations:[]*entities.Consignation{consignation},
			})
			return
		}
	}
}

// Process market order
func (book *OrderBook) addMarket(consignation *entities.Consignation)  {
	if consignation.Direction == enum.OrderDirection_BUY {
		if (book.sellBook.Len()) == 0 {
			// TODO invalid market order, should be cancelled
		} else {
			currentLevel := book.sellBook.Top()
			for currentLevel != nil {
				for i:=0; i<len(currentLevel.Consignations) ;i++  {
					if currentLevel.Consignations[i].OpenQuantity < consignation.OpenQuantity {
						// TODO currentSellLevel.Consignations[i] finish

						consignation.OpenQuantity -= currentLevel.Consignations[i].OpenQuantity
					} else {
						currentLevel.Consignations[i].OpenQuantity -= consignation.OpenQuantity
						// TODO consignation finish
						if currentLevel.Consignations[i].OpenQuantity == 0 {
							// TODO currentSellLevel.Consignations[i] finish
						}
					}
				}
				book.sellBook.Pop()
				currentLevel = book.buyBook.Top()
				book.topSell = currentLevel.Price
			}
			book.marketBuyBook = append(book.marketBuyBook, consignation)
		}
	} else if consignation.Direction == enum.OrderDirection_SELL {
		if (book.buyBook.Len()) == 0 {
			// TODO invalid market order, should be cancelled
		} else {
			currentLevel := book.buyBook.Top()
			for currentLevel != nil {
				for i:=0; i<len(currentLevel.Consignations) ;i++  {
					if currentLevel.Consignations[i].OpenQuantity < consignation.OpenQuantity {
						// TODO currentSellLevel.Consignations[i] finish

						consignation.OpenQuantity -= currentLevel.Consignations[i].OpenQuantity
					} else {
						currentLevel.Consignations[i].OpenQuantity -= consignation.OpenQuantity
						// TODO consignation finish
						if currentLevel.Consignations[i].OpenQuantity == 0 {
							// TODO currentSellLevel.Consignations[i] finish
						}
					}
				}
				book.buyBook.Pop()
				currentLevel = book.buyBook.Top()
				book.topBuy = currentLevel.Price
			}
			book.marketSellBook = append(book.marketSellBook, consignation)
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
