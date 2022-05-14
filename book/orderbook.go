package book

import (
	"math"

	"github.com/keshavchand/orderbook/cti"
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

type Order struct {
	Price float32
	Side  OrderSide
	Type  OrderType
	Size  int
	Id    int
}

type PriceSide struct {
	Price float32
	Side  OrderSide
}

type OrderBook struct {
	BuyOrders  *OrderLevel
	SellOrders *OrderLevel
	M          map[int]PriceSide // Mapping from order id to its price
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

func (book *OrderBook) insertBuy(order Order) []cti.TradedOrder {
	var traded []cti.TradedOrder
	if book.SellOrders == nil {
		if order.Type == MARKET {
			return traded // TODO: Some sort of error ??
		}
		if book.BuyOrders == nil {
			level := newLevel(order)
			book.BuyOrders = level
			return traded
		}
	}
	traded, order = book.matchOrderBuy(order)
	// NO one is selling lower than the least
	// selling price
	if order.Type == LIMIT && order.Size > 0 {
		level := book.BuyOrders.insert(order)
		if level != nil && level.Price > book.bestBuy() {
			book.BuyOrders = level
		}
	}
	return traded
}

func (book *OrderBook) insertSell(order Order) []cti.TradedOrder {
	var traded []cti.TradedOrder
	// IF there are no buy orders
	if book.BuyOrders == nil { // unlikely
		if order.Type == MARKET {
			return traded // TODO: Some sort of error ??
		}
		if book.SellOrders == nil {
			level := newLevel(order)
			book.SellOrders = level
			return traded
		}
	}
	traded, order = book.matchOrderSell(order)
	// NO one is buying at higher price
	// than the highest
	if order.Type == LIMIT && order.Size > 0 {
		newLevel := book.SellOrders.insert(order)
		if newLevel != nil && newLevel.Price < book.bestSell() {
			book.SellOrders = newLevel
		}
	}
	return traded
}

func (book *OrderBook) Insert(order Order) []cti.TradedOrder {
	if order.Type == LIMIT {
		// MARKET orders will not be stored in the books
		if book.M == nil {
			book.M = make(map[int]PriceSide)
		}
		book.M[order.Id] = PriceSide{
			order.Price,
			order.Side,
		}
	}
	switch order.Side {
	case BUY:
		return book.insertBuy(order)
	case SELL:
		return book.insertSell(order)
	}

  return nil
}

func (book *OrderBook) matchOrderBuy(order Order) ([]cti.TradedOrder, Order) {
	var traded []cti.TradedOrder
	for order.Size > 0 {
		// If the buying price is greater than
		// the highest selling price then we can't
		// do anything about it
		if order.Type == LIMIT && order.Price < book.bestSell() {
			break
		}
		t, o := book.SellOrders.match(order)
		order = o
		traded = append(traded, t...)

		if book.SellOrders.OrderCount == 0 {
			book.SellOrders = book.SellOrders.GreaterLevel
		}
	}
	return traded, order
}

func (book *OrderBook) matchOrderSell(order Order) ([]cti.TradedOrder, Order) {
	var traded []cti.TradedOrder
	for order.Size > 0 {
		// If the buying price is lesser than
		// the lowest buying price then we can't
		// do anything about it
		if order.Type == LIMIT && order.Price > book.bestBuy() {
			break
		}
		t, o := book.BuyOrders.match(order)
		order = o
		traded = append(traded, t...)
		if book.BuyOrders.OrderCount == 0 {
			book.BuyOrders = book.BuyOrders.LesserLevel
		}
	}
	return traded, order
}
