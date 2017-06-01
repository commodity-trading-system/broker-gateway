package executor

import (
	"github.com/shopspring/decimal"
	"broker-gateway/entities"
	"container/heap"
)

type Level struct {
	Price decimal.Decimal
	Consignations []*entities.Consignation
}

type Heap []Level

type HeapInterface interface {
	heap.Interface
	Top() *Level
	Travel(func(level *Level))
}

func (h Heap) Len() int { return len(h) }
func (h Heap) Less(i, j int) bool { return h[i].Price.Cmp(h[j].Price)< 0 }
func (h Heap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

// Append new consignation to the level with same price
func (h *Heap) Push(x interface{}) {
	level,_ := x.(Level)
	for i:= len(*h) - 1; i>=0 ;i--  {
		if (*h)[i].Price.Equal(level.Price) {
			(*h)[i].Consignations = append( (*h)[i].Consignations,level.Consignations[0])
			return
		}
	}
	*h = append(*h, x.(Level))
}


func (h *Heap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

// Return the top element without removing it
func (h *Heap) Top() *Level {
	if len(*h) > 0 {
		return &(*h)[0]
	}
	return nil
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

func NewMaxHeap() *MaxHeap {
	return &MaxHeap{
		Heap{},
	}
}

func NewMinHeap() *MinHeap {
	return &MinHeap{
		Heap{},
	}
}

func (h *Heap)Travel(callback func(level *Level))  {
	for i:=0; i<len(*h) ;i++  {
		callback(&(*h)[i])
	}
}