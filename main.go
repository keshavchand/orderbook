package main

import (
	"math"
	"runtime"
)

type OrderBook struct {
	BuyOrders  *OrderLevel
	SellOrders *OrderLevel
}

type OrderLevel struct {
	Price        float32
	Orders       []Order
	GreaterLevel *OrderLevel
	LesserLevel  *OrderLevel
}

func (level *OrderLevel) OrderCount() int {
	levelOrderSize := 0
	for levelOrderSize == 0 {
		for _, buyOrder := range level.Orders {
			levelOrderSize += buyOrder.Size
		}
	}
	return levelOrderSize
}

func (currentLevel *OrderLevel) Insert(order Order) *OrderLevel {
	// If orders are present then iterate
	// to the level least less than the required level
	if currentLevel == nil {
		newLevel := &OrderLevel{
			Price:        order.Price,
			Orders:       make([]Order, 0, 10),
			GreaterLevel: nil,
			LesserLevel:  nil,
		}
		newLevel.Orders = append(newLevel.Orders, order)
		return newLevel
	}
	level := currentLevel
	if level.Price < order.Price {
		for level.GreaterLevel != nil {
			if level.Price == order.Price {
				level.Orders = append(level.Orders, order)
				return nil
			}
			if level.GreaterLevel.Price > order.Price {
				break
			}
			level = level.GreaterLevel
		}
		newLevel := &OrderLevel{
			Price:        order.Price,
			Orders:       make([]Order, 0, 10),
			GreaterLevel: level.GreaterLevel,
			LesserLevel:  level,
		}
		newLevel.Orders = append(newLevel.Orders, order)
		level.GreaterLevel = newLevel
		return newLevel
	} else if level.Price > order.Price {
		for level.LesserLevel != nil {
			if level.Price == order.Price {
				level.Orders = append(level.Orders, order)
				return nil
			}
			if level.LesserLevel.Price < order.Price {
				break
			}
			level = level.LesserLevel
		}
		newLevel := &OrderLevel{
			Price:        order.Price,
			Orders:       make([]Order, 0, 100),
			GreaterLevel: level,
			LesserLevel:  level.LesserLevel,
		}
		newLevel.Orders = append(newLevel.Orders, order)
		level.LesserLevel = newLevel
		return newLevel
	}

	return nil // HOW IT REACHED HERE
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
		{
			order := order
		BuyOrderLoop:
			for {
				if order.Type == LIMIT && order.Price < book.BestSell() {
					break BuyOrderLoop
				}
				for idx, sellOrder := range book.SellOrders.Orders {
					tradeSize := min(order.Size, sellOrder.Size)
					order.Size -= tradeSize
					book.SellOrders.Orders[idx].Size -= tradeSize
					// Write the results to somewhere
					if order.Size == 0 {
						break BuyOrderLoop
					}
				}
				book.SellOrders = book.SellOrders.LesserLevel
			}
		}
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
		{
			order := order
		SellOrderLoop:
			for order.Size > 0 {
				if order.Type == LIMIT && order.Price > book.BestBuy() {
					break SellOrderLoop
				}
				for idx, buyOrder := range book.BuyOrders.Orders {
					tradeSize := min(order.Size, buyOrder.Size)
					order.Size -= tradeSize
					book.BuyOrders.Orders[idx].Size -= tradeSize
					// Write the results to somewhere
					if order.Size == 0 {
						break
					}
				}
				if book.BuyOrders.OrderCount() == 0 {
					book.BuyOrders = book.BuyOrders.GreaterLevel
				}
			}
		}
		// NO one is buying at higher price
		// than the highest
		if order.Type == LIMIT && order.Size > 0 {
			runtime.Breakpoint()
			newLevel := book.SellOrders.Insert(order)
			if newLevel != nil && newLevel.Price < book.BestSell() {
				book.SellOrders = newLevel
			}
			return
		}
	}
}

func main() {
	{
		someComp := OrderBook{}
		someComp.Insert(Order{10.0, BUY, LIMIT, 10})
		someComp.Insert(Order{11.0, BUY, LIMIT, 10})
		someComp.Insert(Order{12.0, BUY, LIMIT, 10})
		someComp.Insert(Order{9.0, BUY, LIMIT, 10})
		someComp.Insert(Order{8.0, BUY, LIMIT, 10})
		// someComp := OrderBook{}
		someComp.Insert(Order{10.0, SELL, LIMIT, 10})
		someComp.Insert(Order{11.0, SELL, LIMIT, 10})
		someComp.Insert(Order{13.0, SELL, LIMIT, 10})
		someComp.Insert(Order{9.0, SELL, LIMIT, 10})
		someComp.Insert(Order{8.0, SELL, LIMIT, 10})
	}
}
