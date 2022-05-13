package book

import (
	"log"
	"math"
)

type OrderSide int

const (
	BUY OrderSide = iota
	SELL
)

type OrderType int

const (
	LIMIT OrderType = iota
	MARKET
)

type PriceSide struct {
	Price float32
	Side  OrderSide
}

type OrderBook struct {
	BuyOrders  *OrderLevel
	SellOrders *OrderLevel
	Reporter   TradeReporter
}

type Order struct {
	Price float32
	Side  OrderSide
	Type  OrderType
	Size  int
	Id    int
}

type TradeReporter func(to int, from int, price float32, size int)

func reporterStub(to int, from int, price float32, size int) {
	log.Printf("Traded at %0.2f from %d to %d of size %d", price, to, from, size)
}

func (book *OrderBook) bestBuy() float32 {
	if book.BuyOrders == nil {
		return math.SmallestNonzeroFloat32
	}
	return book.BuyOrders.Price
}

func (book *OrderBook) bestSell() float32 {
	if book.SellOrders == nil {
		return math.MaxFloat32
	}
	return book.SellOrders.Price
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func (book *OrderBook) insertBuy(order Order) {
	if book.SellOrders == nil {
		if order.Type == MARKET {
			return // TODO: Some sort of error ??
		}
		if book.BuyOrders == nil {
			level := NewLevel(order)
			book.BuyOrders = level
			return
		}
	}
	order = book.matchOrderBuy(order)
	// NO one is selling lower than the least
	// selling price
	if order.Type == LIMIT && order.Size > 0 {
		newLevel := book.BuyOrders.Insert(order)
		if newLevel != nil && newLevel.Price > book.bestBuy() {
			book.BuyOrders = newLevel
		}
		return
	}
}

func (book *OrderBook) insertSell(order Order) {
	// IF there are no buy orders
	if book.BuyOrders == nil { // unlikely
		if order.Type == MARKET {
			return // TODO: Some sort of error ??
		}
		if book.SellOrders == nil {
			level := NewLevel(order)
			book.SellOrders = level
			return
		}
	}
	order = book.matchOrderSell(order)
	// NO one is buying at higher price
	// than the highest
	if order.Type == LIMIT && order.Size > 0 {
		newLevel := book.SellOrders.Insert(order)
		if newLevel != nil && newLevel.Price < book.bestSell() {
			book.SellOrders = newLevel
		}
		return
	}
}

func (book *OrderBook) Insert(order Order) {
	switch order.Side {
	case BUY:
		book.insertBuy(order)
	case SELL:
		book.insertSell(order)
	}
}

func (book *OrderBook) matchOrderBuy(order Order) Order {
	for order.Size > 0 {
		// If the buying price is greater than
		// the highest selling price then we can't
		// do anything about it
		if order.Type == LIMIT && order.Price < book.bestSell() {
			break
		}
		order = book.SellOrders.Match(order, book.Reporter)

		if book.SellOrders.OrderCount == 0 {
			book.SellOrders = book.SellOrders.GreaterLevel
		}
	}
	return order
}

func (book *OrderBook) matchOrderSell(order Order) Order {
	for order.Size > 0 {
		// If the buying price is lesser than
		// the lowest buying price then we can't
		// do anything about it
		if order.Type == LIMIT && order.Price > book.bestBuy() {
			break
		}
		order = book.BuyOrders.Match(order, book.Reporter)
		if book.BuyOrders.OrderCount == 0 {
			book.BuyOrders = book.BuyOrders.LesserLevel
		}
	}
	return order
}
