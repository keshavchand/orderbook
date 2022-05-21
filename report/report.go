package report

import (
	"fmt"
	"os"
	"time"

	"github.com/keshavchand/orderbook/cti"
)

type TradeLogger struct {
	Orders []cti.TradedOrder
	F      *os.File
	Offset int
}

func (l *TradeLogger) Log(o []cti.TradedOrder) {
	l.Orders = append(l.Orders, o...)
	if l.F == nil {
		const flag = os.O_APPEND | os.O_CREATE | os.O_WRONLY
		const perm = 0655
		file, err := os.OpenFile(fmt.Sprintf("TradeReport-%d.rep", time.Now()), flag, perm)
		if err != nil {
			panic(fmt.Sprintf("can't open file: %v", err))
		}
		l.F = file
	}
	for _, to := range o {
		w, err := l.F.Write([]byte(fmt.Sprintf("%x %v", l.Offset, to)))
		if err != nil {
			panic(fmt.Sprintf("can't write to file: %v", err))
		}
		l.Offset += w
	}
}
