package main

import (
	"github.com/keshavchand/orderbook/book"
	"github.com/keshavchand/orderbook/report"
)

func main() {
	logger := report.TradeLogger{}
	comp := book.OrderBook{}
	var sample []book.Order
	sample = append(sample, book.Order{20.0, book.BUY, book.LIMIT, 10, 1})
	sample = append(sample, book.Order{10.0, book.SELL, book.LIMIT, 10, 0})
	for _, s := range sample {
		logger.Log(comp.Insert(s))
	}
}
