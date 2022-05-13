package report

import "log"

type Trade struct {
	To    int
	From  int
	Price float32
	Size  int
}

type TradeHistory struct {
	Trades []Trade
}

func (t *TradeHistory) Report(to int, from int, price float32, size int) {
	t.Trades = append(t.Trades, Trade{to, from, price, size})
	log.Printf("Traded at %0.2f from %d to %d of size %d",
		price, from, to, size)
}
