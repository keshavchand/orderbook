package book

import (
	"errors"
)

//TODO: REMOVE CAPITAL LETTERS FROM FUNCTION NAMES

// Implements priority heap data structure
type Orders struct {
	o     []Order
	items int
}

// The default value of 10
func NewOrders() Orders {
	var o Orders
	o.o = make([]Order, 0, 10)
	o.o = append(o.o, Order{})
	return o
}

func (o *Orders) propagateUp(l int) {
	parent := func(n int) int {
		if n <= 1 {
			return 1
		}
		return n / 2
	}

	c := l
	p := parent(c)
	for o.o[p].Size < o.o[c].Size {
		o.o[p], o.o[c] = o.o[c], o.o[p]
		c = p
		p = parent(c)
	}
}

func (o *Orders) Add(order Order) {
	lastWritten := 0
	if o.items == len(o.o)-1 {
		o.items++
		o.o = append(o.o, order)
		lastWritten = len(o.o) - 1
	} else {
		o.items++
		o.o[o.items] = order
		lastWritten = o.items
	}
	o.propagateUp(lastWritten)
}

var (
	ErrNoOrder = errors.New("no orders at the level")
)

func (o *Orders) propagateDown(l int) {
	child := func(parent int) (int, int) {
		return 2 * parent, 2*parent + 1
	}
	p := l
	c1, c2 := child(p)
	for {
		if c1 < o.items && c2 < o.items {
			if o.o[c1].Size > o.o[c2].Size && o.o[c1].Size > o.o[p].Size {
				o.o[c1], o.o[p] = o.o[p], o.o[c1]
				p = c1
				c1, c2 = child(p)
			} else if o.o[c2].Size > o.o[c1].Size && o.o[c2].Size > o.o[p].Size {
				o.o[c2], o.o[p] = o.o[p], o.o[c2]
				p = c2
				c1, c2 = child(p)
			} else {
				break
			}
		} else if c1 < o.items {
			if o.o[c1].Size > o.o[p].Size {
				o.o[c1], o.o[p] = o.o[p], o.o[c1]
				p = c1
				c1, c2 = child(p)
			} else {
				break
			}
		} else if c2 < o.items {
			if o.o[c2].Size > o.o[p].Size {
				o.o[c2], o.o[p] = o.o[p], o.o[c2]
				p = c2
				c1, c2 = child(p)
			} else {
				break
			}
		} else {
			break
		}
	}
}

func (o *Orders) Pop() (Order, error) {
	var order Order
	if len(o.o) <= 1 || o.items == 0 {
		return order, ErrNoOrder
	}

	order = o.o[1]
	o.o[1] = o.o[o.items]
	o.propagateDown(1)
	o.items--
	return order, nil
}

func (o *Orders) Peek() (Order, error) {
	var order Order
	if len(o.o) <= 1 {
		return order, ErrNoOrder
	}
	order = o.o[1]
	return order, nil
}

func (o *Orders) Remove(id uint64) (Order, bool) {
	var tOrder Order
	for i, order := range o.o[1:] {
		i := i + 1
		if order.Id == id {
			tOrder = order
			o.o[i].Size = 0
			o.o[i], o.o[o.items] = o.o[o.items], o.o[i]
			o.propagateDown(i)
			o.items--
			return tOrder, true
		}
	}

	return tOrder, false
}
