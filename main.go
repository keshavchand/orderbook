package main

import "math"

type OrderBook struct {
	BuyOrders  *OrderLevel
	SellOrders *OrderLevel
}

type OrderLevel struct {
	Price        float32
	CurrentOrder []Order
	GreaterLevel *OrderLevel
	LesserLevel  *OrderLevel
}

func (currentLevel *OrderLevel) Insert(order Order) *OrderLevel {
	// If orders are present then iterate
	// to the level least less than the required level
	level := currentLevel
	if level.Price < order.Price {
		for level.GreaterLevel != nil {
			if level.Price == order.Price {
				level.CurrentOrder = append(level.CurrentOrder, order)
				return nil
			}
			if level.GreaterLevel.Price > order.Price {
				break
			}
			level = level.GreaterLevel
		}
		newLevel := &OrderLevel{
			Price:        order.Price,
			CurrentOrder: make([]Order, 0, 10),
			GreaterLevel: level.GreaterLevel,
			LesserLevel:  level,
		}
		level.GreaterLevel = newLevel
		return newLevel
	}
	if level.Price > order.Price {
		for level.LesserLevel != nil {
			if level.Price == order.Price {
				level.CurrentOrder = append(level.CurrentOrder, order)
				return nil
			}
			if level.LesserLevel.Price < order.Price {
				break
			}
			level = level.LesserLevel
		}
		newLevel := &OrderLevel{
			Price:        order.Price,
			CurrentOrder: make([]Order, 0, 100),
			GreaterLevel: level,
			LesserLevel:  level.LesserLevel,
		}
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

func (book *OrderBook) Insert(order Order) {
	switch order.Side {
	case BUY:
		if book.SellOrders == nil {
			if order.Type == MARKET {
				return /* TODO: Some sort of error ?? */
			}
			if book.BuyOrders == nil {
				level := OrderLevel{
					Price:        order.Price,
					CurrentOrder: make([]Order, 0, 10),
				}
				level.CurrentOrder = append(level.CurrentOrder, order)
				book.BuyOrders = &level
				return
			}
		}
		// NO one is selling lower than the least
		// selling price
		if order.Price < book.BestSell() {
			newOrder := book.BuyOrders.Insert(order)
			if newOrder != nil {
				book.BuyOrders = newOrder
			}
			return
		}
	case SELL:
		if book.BuyOrders == nil {
			if order.Type == MARKET {
				return /* TODO: Some sort of error ?? */
			}
			if book.SellOrders == nil {
				level := OrderLevel{
					Price:        order.Price,
					CurrentOrder: make([]Order, 0, 10),
				}
				level.CurrentOrder = append(level.CurrentOrder, order)
				book.SellOrders = &level
				return
			}
		}
		// NO one is buying at higher price
		// than the highest
		if order.Price > book.BestBuy() {
			newOrder := book.SellOrders.Insert(order)
			if newOrder != nil {
				book.SellOrders = newOrder
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
	}

	{
		someComp := OrderBook{}
		someComp.Insert(Order{10.0, SELL, LIMIT, 10})
		someComp.Insert(Order{11.0, SELL, LIMIT, 10})
		someComp.Insert(Order{12.0, SELL, LIMIT, 10})
		someComp.Insert(Order{9.0, SELL, LIMIT, 10})
		someComp.Insert(Order{8.0, SELL, LIMIT, 10})
	}
}
