package cti

import "testing"

func TestOrderId(t *testing.T) {
	samples := []struct {
		s, e      uint64
		shouldErr bool
	}{
		{10, 10, false},
		{MaxSenderId, MaxOrderCount, false},
	}
	for _, sample := range samples {
		o, err := CreateOrderId(sample.s, sample.e)
		if sample.shouldErr && err == nil {
			t.Errorf("Should've errored s:%d e:%d", sample.s, sample.e)
			continue
		}
		s, e := ParseOrderId(o)
		if s != sample.s {
			t.Errorf("s not same %d vs %d", s, sample.s)
		}
		if e != sample.e {
			t.Errorf("e not same %d vs %d", e, sample.e)
		}
	}
}
