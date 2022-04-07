package main

import "testing"

func Benchmark_OrderWriting(b *testing.B) {
  someComp := OrderBook{}
  someComp.OutFile = WriterStub{}
	for n := 0; n < b.N; n++ {
		someComp.Insert(Order{10.0, BUY, LIMIT, 10})
		someComp.Insert(Order{10.0, BUY, LIMIT, 10})
		someComp.Insert(Order{10.0, BUY, LIMIT, 10})
		someComp.Insert(Order{11.0, BUY, LIMIT, 10})
		someComp.Insert(Order{12.0, BUY, LIMIT, 10})
		someComp.Insert(Order{9.0, BUY, LIMIT, 10})
		someComp.Insert(Order{8.0, BUY, LIMIT, 10})
		someComp.Insert(Order{10.0, SELL, LIMIT, 10})
		someComp.Insert(Order{10.0, SELL, LIMIT, 15})
		someComp.Insert(Order{11.0, SELL, LIMIT, 10})
		someComp.Insert(Order{13.0, SELL, LIMIT, 10})
		someComp.Insert(Order{9.0, SELL, LIMIT, 10})
		someComp.Insert(Order{8.0, SELL, LIMIT, 10})
	}
}
