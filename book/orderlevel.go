package book

import (
	"fmt"
	"io"
)

type OrderLevel struct {
	Price        float32
	Orders       []Order
	GreaterLevel *OrderLevel
	LesserLevel  *OrderLevel
	OrderCount   int
	Offset       int
}

// Return true if found false if not found
func (level *OrderLevel) Delete(id int, price float32, offset int) bool {
	for level.Price < price {
		newLevel := level.GreaterLevel
		if newLevel == nil {
			return false
		}
		if newLevel.Price > price {
			return false
		}
		level = newLevel
	}
	for level.Price > price {
		newLevel := level.LesserLevel
		if newLevel == nil {
			return false
		}
		if newLevel.Price < price {
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

func (level *OrderLevel) Match(order Order, outFile io.Writer) Order {
	for idx, thisOrder := range level.Orders {
		tradeSize := min(order.Size, thisOrder.Size)
		if tradeSize != 0 {
			outFile.Write([]byte(fmt.Sprintf("%d traded at %f\n", tradeSize, order.Price)))
		}
		order.Size -= tradeSize
		level.Orders[idx].Size -= tradeSize
		level.OrderCount -= tradeSize
		level.Offset++
		// Write the results to somewhere
		if order.Size == 0 {
			return order
		}
	}
	return order
}

func NewLevel(order Order) *OrderLevel {
	newLevel := &OrderLevel{
		Price:        order.Price,
		Orders:       make([]Order, 0, 10),
		GreaterLevel: nil,
		LesserLevel:  nil,
		OrderCount:   order.Size,
	}
	newLevel.Orders = append(newLevel.Orders, order)
	return newLevel
}

func (level *OrderLevel) Insert(order Order) *OrderLevel {
	// If orders are present then iterate
	// to the level least less than the required level
	if lvel == nil {
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
			level.Orders = append(level.Orders, order)
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
			level.Orders = append(level.Orders, order)
			level.OrderCount += order.Size
			return level
		}
		newLevel := NewLevel(order)
		level.LesserLevel = newLevel
		return newLevel
	}

	level.Orders = append(level.Orders, order)
	level.OrderCount += order.Size
	return level
}

/* TODO: convert OrderLevel into an interface to support multiple implementations
type OrderLevel interface {
  Delete(id int, price float32, offset int) bool
  Match(order Order, outFile io.Writer) Order
  Insert(order Order) *OrderLevel
}
*/
