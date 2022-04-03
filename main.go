package main

import (
	"math"
)

type OrderBook struct {
	BuyOrders  *OrderLevel
	SellOrders *OrderLevel
}

type OrderSide int

const (
	BUY  OrderSide = iota
	SELL OrderSide = iota
)

type OrderType int

const (
	LIMIT  OrderType = iota
	MARKET OrderType = iota
)

type Order struct {
	Price float32
	Side  OrderSide
	Type  OrderType
	Size  int
}

func (book *OrderBook) BestBuy() float32 {
	if book.BuyOrders == nil {
		return math.SmallestNonzeroFloat32
	}
	return book.BuyOrders.Price
}
func (book *OrderBook) BestSell() float32 {
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

func (book *OrderBook) Insert(order Order) {
	switch order.Side {
	case BUY:
		if book.SellOrders == nil {
			if order.Type == MARKET {
				return // TODO: Some sort of error ??
			}
			if book.BuyOrders == nil {
				level := OrderLevel{
					Price:  order.Price,
					Orders: make([]Order, 0, 10),
				}
				level.Orders = append(level.Orders, order)
				book.BuyOrders = &level
				return
			}
		}
		order = book.MatchOrderBuy(order)
		// NO one is selling lower than the least
		// selling price
		if order.Type == LIMIT && order.Size > 0 {
			newLevel := book.BuyOrders.Insert(order)
			if newLevel != nil && newLevel.Price > book.BestBuy() {
				book.BuyOrders = newLevel
			}
			return
		}
	case SELL:
		// IF there are no buy orders
		if book.BuyOrders == nil { // unlikely
			if order.Type == MARKET {
				return // TODO: Some sort of error ??
			}
			if book.SellOrders == nil {
				level := OrderLevel{
					Price:  order.Price,
					Orders: make([]Order, 0, 10),
				}
				level.Orders = append(level.Orders, order)
				book.SellOrders = &level
				return
			}
		}
		order = book.MatchOrderSell(order)
		// NO one is buying at higher price
		// than the highest
		if order.Type == LIMIT && order.Size > 0 {
			newLevel := book.SellOrders.Insert(order)
			if newLevel != nil && newLevel.Price < book.BestSell() {
				book.SellOrders = newLevel
			}
			return
		}
	}
}

func (book *OrderBook) MatchOrderBuy(order Order) Order {
BuyOrderLoop:
	for order.Size > 0 {
		// If the buying price is greater than
		// the highest selling price then we can't
		// do anything about it
		if order.Type == LIMIT && order.Price < book.BestSell() {
			break BuyOrderLoop
		}
		order = book.SellOrders.MatchOrder(order)
		if book.SellOrders.OrderCount == 0 {
			book.SellOrders = book.SellOrders.GreaterLevel
		}
	}
	return order
}

func (book *OrderBook) MatchOrderSell(order Order) Order {
SellOrderLoop:
	for order.Size > 0 {
		// If the buying price is lesser than
		// the lowest buying price then we can't
		// do anything about it
		if order.Type == LIMIT && order.Price > book.BestBuy() {
			break SellOrderLoop
		}
		order = book.BuyOrders.MatchOrder(order)
		if book.BuyOrders.OrderCount == 0 {
			book.BuyOrders = book.BuyOrders.LesserLevel
		}
	}
	return order
}

func main() {
	{
		someComp := OrderBook{}
		someComp.Insert(Order{10.0, BUY, LIMIT, 10})
		someComp.Insert(Order{10.0, BUY, LIMIT, 10})
		someComp.Insert(Order{10.0, BUY, LIMIT, 10})
		someComp.Insert(Order{11.0, BUY, LIMIT, 10})
		someComp.Insert(Order{12.0, BUY, LIMIT, 10})
		someComp.Insert(Order{9.0, BUY, LIMIT, 10})
		someComp.Insert(Order{8.0, BUY, LIMIT, 10})
		someComp.Insert(Order{10.0, SELL, LIMIT, 10})
		someComp.Insert(Order{10.0, SELL, LIMIT, 15})
		someComp.Insert(Order{11.0, SELL, LIMIT, 10})
		someComp.Insert(Order{13.0, SELL, LIMIT, 10})
		someComp.Insert(Order{9.0, SELL, LIMIT, 10})
		someComp.Insert(Order{8.0, SELL, LIMIT, 10})
	}
}
