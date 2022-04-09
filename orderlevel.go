package main

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
}


func (level *OrderLevel) Delete(Id int, price float32) bool {
  for level.Price < price {
    newLevel := level.GreaterLevel
    if newLevel.Price > price {
      return false
    }
    level = newLevel
  }
  for level.Price > price {
    newLevel := level.LesserLevel
    if newLevel.Price < price {
      return false
    }
    level = newLevel
  }

  for idx, order := range level.Orders {
    if order.Id == Id {
      level.OrderCount -= level.Orders[idx].Size
      level.Orders[idx].Size = 0;
      return true
    }
  }
  return false
}

func (level *OrderLevel) MatchOrder(order Order, outFile io.Writer) Order {
	for idx, thisOrder := range level.Orders {
		tradeSize := min(order.Size, thisOrder.Size)
		if tradeSize != 0 {
      outFile.Write([]byte(fmt.Sprintf("%d traded at %f\n", tradeSize, order.Price)))
		}
		order.Size -= tradeSize
		level.Orders[idx].Size -= tradeSize
		level.OrderCount -= tradeSize
		// Write the results to somewhere
		if order.Size == 0 {
			return order
		}
	}
	return order
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
			OrderCount:   order.Size,
		}
		newLevel.Orders = append(newLevel.Orders, order)
		return newLevel
	}
	level := currentLevel
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
		newLevel := &OrderLevel{
			Price:        order.Price,
			Orders:       make([]Order, 0, 10),
			GreaterLevel: level.GreaterLevel,
			LesserLevel:  level,
			OrderCount:   order.Size,
		}
		newLevel.Orders = append(newLevel.Orders, order)
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
		newLevel := &OrderLevel{
			Price:        order.Price,
			Orders:       make([]Order, 0, 100),
			GreaterLevel: level,
			LesserLevel:  level.LesserLevel,
			OrderCount:   order.Size,
		}
		newLevel.Orders = append(newLevel.Orders, order)
		level.LesserLevel = newLevel
		return newLevel
	}

	level.Orders = append(level.Orders, order)
	level.OrderCount += order.Size
	return nil
}
