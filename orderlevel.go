package main
type OrderLevel struct {
	Price        float32
	Orders       []Order
	GreaterLevel *OrderLevel
	LesserLevel  *OrderLevel
}

func (level *OrderLevel) OrderCount() int {
	levelOrderSize := 0
	for _, buyOrder := range level.Orders {
		levelOrderSize += buyOrder.Size
	}
	return levelOrderSize
}

func (level *OrderLevel) MatchOrder(order Order) Order {
	for idx, thisOrder := range level.Orders {
		tradeSize := min(order.Size, thisOrder.Size)
		order.Size -= tradeSize
		level.Orders[idx].Size -= tradeSize
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
			return nil
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
		// Insert below
		for level.Price != order.Price && level.LesserLevel != nil {
			if level.LesserLevel.Price < order.Price {
				break
			}
			level = level.LesserLevel
		}
		if level.Price == order.Price {
			level.Orders = append(level.Orders, order)
			return nil
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

	level.Orders = append(level.Orders, order)
	return nil
}
