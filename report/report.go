package report

import "github.com/keshavchand/orderbook/cti"

type TradeLogger struct {
	Orders []cti.TradedOrder
}

func (l *TradeLogger) Log(o []cti.TradedOrder) {
	l.Orders = append(l.Orders, o...)
}
