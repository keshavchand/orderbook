package book

import (
	"testing"
)

func Test_OrdersInsert(t *testing.T) {
	orders := NewOrders()
	for i := 0; i <= 10; i++ {
		order := Order{10.0, BUY, LIMIT, i, i}
		orders.Add(order)
	}
	for i := 10; i >= 0; i-- {
		order, err := orders.Pop()
		if err != nil {
			t.Fatalf("%v\n", err)
		}
		if order.Size != i {
			t.Fatalf("Size diff than expection %d vs %d\n", order.Size, i)
		}
	}
	_, err := orders.Pop()
	if err == nil {
		t.Fatalf("Reading after all items had been poped creates no nil")
	}
}

func Test_OrdersInsertAfterPop(t *testing.T) {
	orders := NewOrders()
	for i := 0; i <= 10; i++ {
		order := Order{10.0, BUY, LIMIT, i, 0}
		orders.Add(order)
	}
	for i := 10; i >= 5; i-- {
		order, err := orders.Pop()
		if err != nil {
			t.Fatalf("%v\n", err)
		}
		if order.Size != i {
			t.Fatalf("Size diff than expection %d vs %d\n", order.Size, i)
		}
	}

	for i := 5; i <= 10; i++ {
		order := Order{10.0, BUY, LIMIT, i, 1}
		orders.Add(order)
	}
	for i := 10; i >= 0; i-- {
		order, err := orders.Pop()
		if err != nil {
			t.Fatalf("%v\n", err)
		}
		if order.Size != i {
			t.Fatalf("Size diff than expection %d vs %d\n", order.Size, i)
		}
	}
}
