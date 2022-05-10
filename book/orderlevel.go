package book

import (
	"io"
	"log"
)

type OrderLevel struct {
	Price        float32
	Orders       Orders
	GreaterLevel *OrderLevel
	LesserLevel  *OrderLevel
	OrderCount   int
}

/*
// Return true if found false if not found
// TODO: Currently it mainly does a linear search but should
// be upgraded to priority queue on Order.Size bases
func (level *OrderLevel) Delete(id int, price float32, offset int) bool {
	for level.Price < price {
		newLevel := level.GreaterLevel
		if newLevel == nil || newLevel.Price > price {
			return false
		}
		level = newLevel
	}
	for level.Price > price {
		newLevel := level.LesserLevel
		if newLevel == nil || newLevel.Price < price {
			return false
		}
		level = newLevel
	}

	if offset < level.Offset {
		return true
	}

	if level.Orders[offset].Id == id {
		if level.Orders[offset].Size > 0 {
			level.OrderCount -= level.Orders[offset].Size
			level.Orders[offset].Size = 0
			return true
		}
	}
	return false
}
*/

func (level *OrderLevel) Match(order Order, outFile io.Writer) Order {
	for {
		thisOrder, err := level.Orders.Pop()
		if err != nil {
      switch err {
			case ErrNoOrder:
				return order
			default:
				// XXX: PANIC OR SOMETHING
			}
		}
		tradeSize := min(order.Size, thisOrder.Size)
		if tradeSize != 0 {
			side := ""
			switch order.Side {
			case BUY:
				side = "BUYING"
			case SELL:
				side = "SELLING"
			}
			log.Printf("%s at %0.2f from %d to %d of size %d",
				side, level.Price, order.Id, thisOrder.Size, tradeSize)
		}
		order.Size -= tradeSize
		thisOrder.Size -= tradeSize
    level.OrderCount -= tradeSize
		if order.Size == 0 {
			if thisOrder.Size > 0 {
				level.Orders.Add(thisOrder)
			}
			return order
		}
	}
	return order
}

func NewLevel(order Order) *OrderLevel {
	newLevel := &OrderLevel{
		Price:        order.Price,
		Orders:       NewOrders(),
		GreaterLevel: nil,
		LesserLevel:  nil,
		OrderCount:   order.Size,
	}
	newLevel.Orders.Add(order)
	return newLevel
}

func (level *OrderLevel) Insert(order Order) *OrderLevel {
	// If orders are present then iterate
	// to the level least less than the required level
	if level == nil {
		return NewLevel(order)
	}
	if level.Price < order.Price {
		// Insert above
		for level.Price != order.Price && level.GreaterLevel != nil {
			if level.GreaterLevel.Price > order.Price {
				break
			}
			level = level.GreaterLevel
		}
		if level.Price == order.Price {
			level.Orders.Add(order)
			level.OrderCount += order.Size
			return level
		}
		newLevel := NewLevel(order)
		level.GreaterLevel = newLevel
		return newLevel
	} else if level.Price > order.Price {
		// Insert below
		for level.Price != order.Price && level.LesserLevel != nil {
			if level.LesserLevel.Price < order.Price {
				break
			}
			level = level.LesserLevel
		}
		if level.Price == order.Price {
			level.Orders.Add(order)
			level.OrderCount += order.Size
			return level
		}
		newLevel := NewLevel(order)
		level.LesserLevel = newLevel
		return newLevel
	}

	level.Orders.Add(order)
	level.OrderCount += order.Size
	return level
}
