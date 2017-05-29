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

type MaxHeap struct {
	Heap
}

func (h *MaxHeap) Less(i, j int) bool {
	return h.Heap[i].Price.Cmp(h.Heap[j].Price)< 0
}

type MinHeap struct {
	Heap
}

func (h *MinHeap) Less(i, j int) bool {
	return h.Heap[i].Price.Cmp(h.Heap[j].Price)> 0
}

func NewMaxHeap(direction int) *MaxHeap {
	return &MaxHeap{
		Heap{},
	}
}

func NewMinHeap(direction int) *MinHeap {
	return &MinHeap{
		Heap{},
	}
}