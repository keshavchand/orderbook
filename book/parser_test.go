package book

import "testing"

func TestNewOrder(t *testing.T) {
	s := "0;10.0;0;0;10"
	po, err := Parse(s)
	if err != nil {
		t.Fatal(err)
	}
	o_ret, success := po.(Order)
	if !success {
		t.Fatalf("returned type not Order")
	}

	o_req := Order{
		10, BUY, LIMIT, 10, 0,
	}

	if o_ret.Price != o_req.Price {
		t.Fatalf("Prices do not match %f vs %f", o_ret.Price, o_req.Price)
	}
	if o_ret.Side != o_req.Side {
		t.Fatalf("Sides do not match %d vs %d", o_ret.Side, o_req.Side)
	}
	if o_ret.Type != o_req.Type {
		t.Fatalf("Types do not match %d vs %d", o_ret.Type, o_req.Type)
	}
	if o_ret.Size != o_req.Size {
		t.Fatalf("Sizes do not match %d vs %d", o_ret.Size, o_req.Size)
	}
}

func TestRemoveOrder(t *testing.T) {
	s := "1;10"
	po, err := Parse(s)
	if err != nil {
		t.Fatal(err)
	}
	o_ret, success := po.(Order)
	if !success {
		t.Fatalf("returned type not Remove Order")
	}

	if o_ret.Id != 10 {
		t.Fatalf("Ids don't match")
	}
}

func TestUpdateOrder(t *testing.T) {
	s := "2;10;10.0;0;0;10"
	po, err := Parse(s)
	if err != nil {
		t.Fatal(err)
	}
	o_ret, success := po.(Order)
	if !success {
		t.Fatalf("returned type not Update Order")
	}

	o_req := Order{
		10, BUY, LIMIT, 10, 10,
	}

	if o_ret.Id != o_req.Id {
		t.Fatalf("Ids do not match %d vs %d", o_ret.Id, o_req.Id)
	}
	if o_ret.Price != o_req.Price {
		t.Fatalf("Prices do not match %f vs %f", o_ret.Price, o_req.Price)
	}
	if o_ret.Side != o_req.Side {
		t.Fatalf("Sides do not match %d vs %d", o_ret.Side, o_req.Side)
	}
	if o_ret.Type != o_req.Type {
		t.Fatalf("Types do not match %d vs %d", o_ret.Type, o_req.Type)
	}
	if o_ret.Size != o_req.Size {
		t.Fatalf("Sizes do not match %d vs %d", o_ret.Size, o_req.Size)
	}
}

func BenchmarkParser(b *testing.B) {
	// run the Fib function b.N times
	for n := 0; n < b.N; n++ {
		s := "0;10.0;0;0;10"
		Parse(s)
	}
}
