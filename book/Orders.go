package book

import "errors"

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

func (o *Orders) Add(order Order) {
	parent := func(n int) int {
		if n <= 1 {
			return 1
		}
		return n / 2
	}

	lastWritten := 0
	if o.items == len(o.o)-1 {
		o.o = append(o.o, order)
		lastWritten = len(o.o) - 1
		o.items++
	} else {
		o.items++
		o.o[o.items] = order
		lastWritten = o.items
	}

	c := lastWritten
	p := parent(c)
	for o.o[p].Size < o.o[c].Size {
		o.o[p], o.o[c] = o.o[c], o.o[p]
		c = p
		p = parent(c)
	}
}

var (
	ErrNoOrder = errors.New("no orders at the level")
)

func (o *Orders) Pop() (Order, error) {
	var order Order
	if len(o.o) <= 1 || o.items == 0 {
		return order, ErrNoOrder
	}

	order = o.o[1]
	o.o[1] = o.o[o.items]

	child := func(parent int) (int, int) {
		return 2 * parent, 2*parent + 1
	}
	p := 1
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
