package book

import "github.com/keshavchand/orderbook/cti"

type OrderLevel struct {
	Price        float32
	Orders       Orders
	GreaterLevel *OrderLevel
	LesserLevel  *OrderLevel
	OrderCount   int
}

func (level *OrderLevel) match(order Order) ([]cti.TradedOrder, Order) {
  var t []cti.TradedOrder
	for order.Size > 0 {
		thisOrder, err := level.Orders.Pop()
		if err != nil {
			switch err {
			case ErrNoOrder:
				return t, order
			default:
				return t, order // XXX: PANIC OR SOMETHING
			}
		}

		tradeSize := min(order.Size, thisOrder.Size)
		if tradeSize == 0 {
			continue
		}

		to := thisOrder.Id
		from := order.Id
		if order.Side == BUY {
			to, from = from, to
		}
    t = append(t, cti.TradedOrder{
      To: to,
      From: from,
      Size: tradeSize,
      Price: thisOrder.Price,
    })


		order.Size -= tradeSize
		thisOrder.Size -= tradeSize
		level.OrderCount -= tradeSize

		if thisOrder.Size > 0 {
			level.Orders.Add(thisOrder)
		}
	}
	return t, order
}

func newLevel(order Order) *OrderLevel {
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

func (level *OrderLevel) insert(order Order) *OrderLevel {
	// If orders are present then iterate
	// to the level least less than the required level
	if level == nil {
		return newLevel(order)
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
		newLevel := newLevel(order)
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
		newLevel := newLevel(order)
		level.LesserLevel = newLevel
		return newLevel
	}

	level.Orders.Add(order)
	level.OrderCount += order.Size
	return level
}
