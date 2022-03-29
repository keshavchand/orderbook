package main

type OrderBook struct {
	BuyOrders  *OrderLevel
	BestBuy    float32 // Highest
	SellOrders *OrderLevel
	BestSell   float32 // Lowest
}

type OrderLevel struct {
	Price        float32
	CurrentOrder []Order
	GreaterLevel *OrderLevel
	LesserLevel  *OrderLevel
}

func (currentLevel *OrderLevel) Insert(order Order) {
	// If orders are present then iterate
	// to the level least less than the required level
	level := currentLevel
	if level.Price < order.Price {
		for level.GreaterLevel != nil {
			if level.Price == order.Price {
				level.CurrentOrder = append(level.CurrentOrder, order)
				return
			}
			if level.GreaterLevel.Price > order.Price {
				break
			}
			level = level.GreaterLevel
		}
		newLevel := OrderLevel{
			Price:        order.Price,
			CurrentOrder: make([]Order, 0, 100),
			GreaterLevel: level.GreaterLevel,
			LesserLevel:  level,
		}
		level.GreaterLevel = &newLevel
		return
	}
	if level.Price > order.Price {
		for level.LesserLevel != nil {
			if level.Price == order.Price {
				level.CurrentOrder = append(level.CurrentOrder, order)
				return
			}
			if level.LesserLevel.Price < order.Price {
				break
			}
			level = level.GreaterLevel
		}
		newLevel := OrderLevel{
			Price:        order.Price,
			CurrentOrder: make([]Order, 0, 100),
			GreaterLevel: level,
			LesserLevel:  level.LesserLevel,
		}
		level.LesserLevel = &newLevel
		return
	}

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
					CurrentOrder: make([]Order, 0, 100),
				}
				level.CurrentOrder = append(level.CurrentOrder, order)
				book.BuyOrders = &level
				book.BestBuy = order.Price
				return
			}
		}
    if order.Price < book.BestSell {
      book.BuyOrders.Insert(order)
      if order.Price > book.BestBuy {
        book.BestBuy = order.Price
      }
    }
	case SELL:
		if book.BuyOrders == nil {
			if order.Type == MARKET {
				return /* TODO: Some sort of error ?? */
			}
			if book.SellOrders == nil {
				level := OrderLevel{
					Price:        order.Price,
					CurrentOrder: make([]Order, 0, 100),
				}
				level.CurrentOrder = append(level.CurrentOrder, order)
				book.SellOrders = &level
				book.BestBuy = order.Price
				return
			}
		}
    if order.Price > book.BestBuy {
      book.SellOrders.Insert(order)
      if order.Price < book.BestBuy {
        book.BestBuy = order.Price
      }
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
