package executor

import (
	"github.com/shopspring/decimal"
	"broker-gateway/entities"
)

type Level struct {
	Price decimal.Decimal
	Consignations []*entities.Consignation
}

type Heap []Level

func (h Heap) Len() int { return len(h) }
func (h Heap) Less(i, j int) bool { return h[i].Price.Cmp(h[j].Price)< 0 }
func (h Heap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *Heap) Push(x interface{}) {
	*h = append(*h, x.(Level))
}

func (h *Heap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

type BuyHeap struct {
	Heap
}

func (h *BuyHeap) Less(i, j int) bool {
	return h.Heap[i].Price.Cmp(h.Heap[j].Price)< 0
}

type SellHeap struct {
	Heap
}

func (h *SellHeap) Less(i, j int) bool {
	return h.Heap[i].Price.Cmp(h.Heap[j].Price)> 0
}

func NewBuyHeap(direction int) *BuyHeap {
	return &BuyHeap{
		Heap{},
	}
}

func NewSellHeap(direction int) *SellHeap {
	return &SellHeap{
		Heap{},
	}
}