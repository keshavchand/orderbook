package book

import (
	"fmt"
	"sync/atomic"
	"testing"

)

type WriterStub struct{}

func (w WriterStub) Write(b []byte) (int, error) {
	return len(b), nil
}

func CreateOrderBook(book OrderBook) {
	fmt.Println(`
	book := OrderBook{}
  `)

	buySide := book.BuyOrders
	for buySide != nil {
		for _, order := range buySide.Orders.o[1:] {
			t := "LIMIT"
			if order.Type == MARKET {
				t = "MARKET"
			}

			s := "BUY"
			if order.Side == SELL {
				s = "SELL"
			}
			fmt.Printf("book.Insert(OrderBook{%f, %s, %s, %d, %d})\n",
				order.Price, s, t, order.Size, order.Id)
		}
		buySide = buySide.LesserLevel
	}

	sellSide := book.SellOrders
	for sellSide != nil {
		for _, order := range sellSide.Orders.o[1:] {
			t := "LIMIT"
			if order.Type == MARKET {
				t = "MARKET"
			}

			s := "BUY"
			if order.Side == SELL {
				s = "SELL"
			}
			fmt.Printf("book.Insert(OrderBook{%f, %s, %s, %d, %d})\n",
				order.Price, s, t, order.Size, order.Id)
		}
		sellSide = sellSide.GreaterLevel
	}
}

func SameOrderBook(t *testing.T, book1, book2 OrderBook) {
	if (book1.BuyOrders == nil && book2.BuyOrders != nil) ||
		(book1.BuyOrders != nil && book2.BuyOrders == nil) {
		t.Errorf("Books BuyOrders Mismatch %v vs %v", book1.BuyOrders, book2.BuyOrders)
		return
	}
	if (book1.SellOrders == nil && book2.SellOrders != nil) ||
		(book1.SellOrders != nil && book2.SellOrders == nil) {
		t.Errorf("Books SellOrders Mismatch %v vs %v", book1.SellOrders, book2.SellOrders)
		return
	}

	book1_lvl := book1.BuyOrders
	book2_lvl := book2.BuyOrders
	if book1_lvl != nil && book2_lvl != nil {
		for !(book1_lvl == nil && book2_lvl == nil) {
			if (book1_lvl == nil && book2_lvl != nil) ||
				(book1_lvl != nil && book2_lvl == nil) {
				t.Errorf("Book level Mismatch %v %v", book1_lvl, book2_lvl)
			}
			for book1_lvl.OrderCount == 0 {
				book1_lvl = book1_lvl.LesserLevel
			}
			for book2_lvl.OrderCount == 0 {
				book2_lvl = book2_lvl.LesserLevel
			}
			if book1_lvl.Price != book2_lvl.Price {
				t.Errorf(`Books probably misses levels book1 %f vs book2 %f`, book1_lvl.Price, book2_lvl.Price)
			}
			lvl1_orders := make([]Order, 0, book1_lvl.OrderCount)
			lvl2_orders := make([]Order, 0, book2_lvl.OrderCount)

			for _, order := range book1_lvl.Orders.o {
				if order.Size == 0 {
					continue
				}
				lvl1_orders = append(lvl1_orders, order)
			}
			for _, order := range book2_lvl.Orders.o {
				if order.Size == 0 {
					continue
				}
				lvl2_orders = append(lvl2_orders, order)
			}

			if len(lvl1_orders) != len(lvl2_orders) {
				t.Errorf("Orders level Size different %d vs %d", len(lvl1_orders), len(lvl2_orders))
				return
			}

			for i := 0; i < len(lvl1_orders); i++ {
				o1 := lvl1_orders[i]
				o2 := lvl2_orders[i]
				if !SameOrders(t, o1, o2) {
					return
				}
			}
			book1_lvl = book1_lvl.LesserLevel
			book2_lvl = book2_lvl.LesserLevel
		}
	}
	book1_lvl = book1.SellOrders
	book2_lvl = book2.SellOrders
	if book1_lvl != nil && book2_lvl != nil {
		for !(book1_lvl == nil && book2_lvl == nil) {
			if (book1_lvl == nil && book2_lvl != nil) ||
				(book1_lvl != nil && book2_lvl == nil) {
				t.Errorf("Book level Mismatch %v %v", book1_lvl, book2_lvl)
			}
			for book1_lvl.OrderCount == 0 {
				book1_lvl = book1_lvl.GreaterLevel
			}
			for book2_lvl.OrderCount == 0 {
				book2_lvl = book2_lvl.GreaterLevel
			}
			if book1_lvl.Price != book2_lvl.Price {
				t.Errorf(`Books probably misses levels book1 %f vs book2 %f`, book1_lvl.Price, book2_lvl.Price)
				return
			}
			lvl1_orders := make([]Order, 0, book1_lvl.OrderCount)
			lvl2_orders := make([]Order, 0, book2_lvl.OrderCount)

			for _, order := range book1_lvl.Orders.o[1:] {
				if order.Size == 0 {
					continue
				}
				lvl1_orders = append(lvl1_orders, order)
			}
			for _, order := range book2_lvl.Orders.o[1:] {
				if order.Size == 0 {
					continue
				}
				lvl2_orders = append(lvl2_orders, order)
			}

			if len(lvl1_orders) != len(lvl2_orders) {
				t.Errorf("Orders level Size different %d vs %d", len(lvl1_orders), len(lvl2_orders))
				return
			}

			for i := 0; i < len(lvl1_orders); i++ {
				o1 := lvl1_orders[i]
				o2 := lvl2_orders[i]
				if !SameOrders(t, o1, o2) {
					return
				}
			}
			book1_lvl = book1_lvl.GreaterLevel
			book2_lvl = book2_lvl.GreaterLevel
		}
	}
}

func SameOrders(t *testing.T, o1, o2 Order) bool {
	if o1.Price != o2.Price {
		t.Errorf("Order Prices different")
		return false
	}
	if o1.Side != o2.Side {
		t.Errorf("Order Sides different")
		return false
	}
	if o1.Type != o2.Type {
		t.Errorf("Order Types different")
		return false
	}
	if o1.Size != o2.Size {
		t.Errorf("Order Sizes different")
		return false
	}
	if o1.Id != o2.Id {
		t.Errorf("Order Ids different")
		return false
	}

	return true
}

func Test_OrderBook_Add_Buy(t *testing.T) {
	book := OrderBook{}
	//book.IdToPrice = make(map[int]PriceSide)
	book.Insert(Order{10.0, BUY, LIMIT, 10, 0})
	if book.BuyOrders.Price == 10.0 &&
		book.BuyOrders.OrderCount == 10 &&
		book.BuyOrders.Orders.o[1].Side == BUY &&
		book.BuyOrders.Orders.o[1].Type == LIMIT {
		return
	}
	t.Errorf("Order isn't Inserted correctly %v", book.BuyOrders)
}

func Test_OrderBook_Add_Sell(t *testing.T) {
	book := OrderBook{}
	//book.IdToPrice = make(map[int]PriceSide)
	book.Insert(Order{10.0, SELL, LIMIT, 10, 0})
	if book.SellOrders.Price == 10.0 &&
		book.SellOrders.OrderCount == 10 &&
		book.SellOrders.Orders.o[1].Side == SELL &&
		book.SellOrders.Orders.o[1].Type == LIMIT {
		return
	}
	t.Errorf("Order isn't Inserted correctly %v", book.SellOrders)
}

func Test_OrderBook_Add_Buy_Market(t *testing.T) {
	book := OrderBook{}
	//book.IdToPrice = make(map[int]PriceSide)
	book.Insert(Order{10.0, SELL, MARKET, 10, 0})
	if book.SellOrders == nil {
		return
	}
	t.Errorf("Market Order isn't Inserted correctly %v", book.SellOrders)
}

func Test_OrderBook_MatchOrder_Sell(t *testing.T) {
	book1 := OrderBook{}
	//book1.IdToPrice = make(map[int]PriceSide)
	book1.Insert(Order{10.0, SELL, LIMIT, 10, 0})
	book1.Insert(Order{15.0, SELL, LIMIT, 10, 1})
	book1.Insert(Order{20.0, BUY, LIMIT, 15, 2})

	book2 := OrderBook{}
	//book2.IdToPrice = make(map[int]PriceSide)
	book2.Insert(Order{15.0, SELL, LIMIT, 5, 1})
	SameOrderBook(t, book1, book2)
}

// TODO: convert every test into this format
func Test_OrderBook_MatchOrder_Buy(t *testing.T) {
	cases := []struct {
		name    string
		tests   []Order
		results []Order
	}{
		{
			name: "Matching Buy Order",
			tests: []Order{
				Order{20.0, BUY, LIMIT, 10, 0},
				Order{10.0, BUY, LIMIT, 10, 1},
				Order{10.0, SELL, LIMIT, 100, 2},
			},
			results: []Order{
				Order{10.0, SELL, LIMIT, 80, 2},
			},
		},
	}
	for _, c := range cases {
		book1 := OrderBook{}
		book2 := OrderBook{}

		for _, o := range c.tests {
			book1.Insert(o)
		}
		for _, o := range c.results {
			book2.Insert(o)
		}

		SameOrderBook(t, book1, book2)
	}
}

/* Delete is UNIMPLEMENTED
func Test_OrderDelete(t *testing.T) {
	book1 := OrderBook{}
	//book1.IdToPrice = make(map[int]PriceSide)
	book1.Insert(Order{20.0, BUY, LIMIT, 10, 0})
	book1.Insert(Order{10.0, BUY, LIMIT, 10, 1})
	book1.Insert(Order{10.0, BUY, LIMIT, 10, 2})
	book1.Delete(1)

	book2 := OrderBook{}
	//book2.IdToPrice = make(map[int]PriceSide)
	book2.Insert(Order{20.0, BUY, LIMIT, 10, 0})
	book2.Insert(Order{10.0, BUY, LIMIT, 10, 2})

	SameOrderBook(t, book1, book2)
}

func Test_OrderDelete_AfterMatching(t *testing.T) {
	book1 := OrderBook{}
	//book1.IdToPrice = make(map[int]PriceSide)
	book1.Insert(Order{20.0, BUY, LIMIT, 10, 0})
	book1.Insert(Order{10.0, BUY, LIMIT, 10, 1})
	book1.Insert(Order{5.0, SELL, LIMIT, 10, 2})
	book1.Delete(0)

	book2 := OrderBook{}
	//book2.IdToPrice = make(map[int]PriceSide)
	book2.Insert(Order{10.0, BUY, LIMIT, 10, 1})

	SameOrderBook(t, book1, book2)
}
*/

func FuzzOrderBook(f *testing.F) {
	someComp := OrderBook{}
	//someComp.IdToPrice = make(map[int]PriceSide)

	f.Add(float32(10.0), uint(10))
	var i int32 = 0
	f.Fuzz(func(t *testing.T, price float32, size uint) {
		someComp.Insert(Order{price, BUY, LIMIT, int(size), int(i)})
		atomic.AddInt32(&i, 1)
	})
}

func Test_BenchMark(t *testing.T) {
	someComp := OrderBook{}
	//someComp.IdToPrice = make(map[int]PriceSide)

	i := 0
	for n := 0; n < 10; n++ {
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
	CreateOrderBook(someComp)
}
func Benchmark_OrderWriting(b *testing.B) {
	someComp := OrderBook{}
	//someComp.IdToPrice = make(map[int]PriceSide)

	i := 0
	for n := 0; n < b.N; n++ {
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
}
