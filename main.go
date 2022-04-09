package main

import (
	"io"
	"math"
	_ "net/http/pprof"
)

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

type PriceSide struct {
	Price float32
	Side  OrderSide
}

type OrderBook struct {
	BuyOrders  *OrderLevel
	SellOrders *OrderLevel
	OutFile    io.Writer
	// mapping from id to price
	IdToPrice map[int]PriceSide
}

type Order struct {
	Price float32
	Side  OrderSide
	Type  OrderType
	Size  int
	Id    int
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
          OrderCount: order.Size,
				}
				level.Orders = append(level.Orders, order)
				book.IdToPrice[order.Id] = PriceSide{order.Price, BUY}
				book.BuyOrders = &level
				return
			}
		}
		order = book.MatchOrderBuy(order)
		// NO one is selling lower than the least
		// selling price
		if order.Type == LIMIT && order.Size > 0 {
			newLevel := book.BuyOrders.Insert(order)
			book.IdToPrice[order.Id] = PriceSide{order.Price, BUY}
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
          OrderCount: order.Size,
				}
				level.Orders = append(level.Orders, order)
				book.IdToPrice[order.Id] = PriceSide{order.Price, SELL}
				book.SellOrders = &level
				return
			}
		}
		order = book.MatchOrderSell(order)
		// NO one is buying at higher price
		// than the highest
		if order.Type == LIMIT && order.Size > 0 {
			newLevel := book.SellOrders.Insert(order)
			book.IdToPrice[order.Id] = PriceSide{order.Price, SELL}
			if newLevel != nil && newLevel.Price < book.BestSell() {
				book.SellOrders = newLevel
			}
			return
		}
	}
}

func (book *OrderBook) Delete(id int) bool{
  price, present := book.IdToPrice[id]
  if present == false {
    return false
  }
  switch price.Side {
    case BUY:
      return book.BuyOrders.Delete(id, price.Price)
    case SELL:
      return book.SellOrders.Delete(id, price.Price)
  }
  return false
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
		order = book.SellOrders.MatchOrder(order, book.OutFile)
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
		order = book.BuyOrders.MatchOrder(order, book.OutFile)
		if book.BuyOrders.OrderCount == 0 {
			book.BuyOrders = book.BuyOrders.LesserLevel
		}
	}
	return order
}

type WriterStub struct{}

func (w WriterStub) Write(b []byte) (int, error) {
	return len(b), nil
}

func main() {
	someComp := OrderBook{}
	someComp.IdToPrice = make(map[int]PriceSide)
	someComp.OutFile = WriterStub{}
	i := 0
	someComp.Insert(Order{10.0, BUY, LIMIT, 10, i})
	i++
	someComp.Insert(Order{10.0, BUY, LIMIT, 10, i})
	i++
	someComp.Insert(Order{10.0, BUY, LIMIT, 10, i})
	i++
	someComp.Insert(Order{11.0, BUY, LIMIT, 10, i})
	i++
	someComp.Insert(Order{12.0, BUY, LIMIT, 10, i})
	i++
	someComp.Insert(Order{9.0, BUY, LIMIT, 10, i})
	i++
	someComp.Insert(Order{8.0, BUY, LIMIT, 10, i})
	i++
	someComp.Insert(Order{10.0, SELL, LIMIT, 10, i})
	i++
	someComp.Insert(Order{10.0, SELL, LIMIT, 15, i})
	i++
	someComp.Insert(Order{11.0, SELL, LIMIT, 10, i})
	i++
	someComp.Insert(Order{13.0, SELL, LIMIT, 10, i})
	i++
	someComp.Insert(Order{9.0, SELL, LIMIT, 10, i})
	i++
	someComp.Insert(Order{8.0, SELL, LIMIT, 10, i})
	i++
}
