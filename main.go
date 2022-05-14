package main

import (
	"log"

	"github.com/keshavchand/orderbook/book"
)

func main() {
  logger := TradeLogger{}
	comp := book.OrderBook{}
	var sample []book.Order
	sample = append(sample, book.Order{20.0, book.BUY, book.LIMIT, 10, 1})
	sample = append(sample, book.Order{10.0, book.SELL, book.LIMIT, 10, 0})
	for _, s := range sample {
		logger.Log(comp.Insert(s))
	}
}
