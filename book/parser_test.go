package book

import "testing"

func TestNewOrder(t *testing.T) {
	s := "0;10.0;0;0;10"
	po, err := Parse(s)
	if err != nil {
		t.Fatal(err)
	}
	o_ret, success := po.(NewOrder)
	if !success {
		t.Fatalf("returned type not Order")
	}

	o_req := Order{
		10, BUY, LIMIT, 10, 0,
	}

	if o_ret.O.Price != o_req.Price {
		t.Fatalf("Prices do not match %f vs %f", o_ret.O.Price, o_req.Price)
	}
	if o_ret.O.Side != o_req.Side {
		t.Fatalf("Sides do not match %d vs %d", o_ret.O.Side, o_req.Side)
	}
	if o_ret.O.Type != o_req.Type {
		t.Fatalf("Types do not match %d vs %d", o_ret.O.Type, o_req.Type)
	}
	if o_ret.O.Size != o_req.Size {
		t.Fatalf("Sizes do not match %d vs %d", o_ret.O.Size, o_req.Size)
	}
}

func TestRemoveOrder(t *testing.T) {
	s := "1;10"
	po, err := Parse(s)
	if err != nil {
		t.Fatal(err)
	}
	p, ok := po.(RemoveOrder)
	if !ok {
		t.Fatal("return type not uint64")
	}
	if p.Id != 10 {
		t.Fatal("id not parsed correctly")
	}
}

func TestUpdateOrder(t *testing.T) {
	s := "2;10;10.0;0;0;10"
	po, err := Parse(s)
	if err != nil {
		t.Fatal(err)
	}
	o_ret, success := po.(UpdateOrder)
	if !success {
		t.Fatalf("returned type not Update Order")
	}

	o_req := Order{
		10, BUY, LIMIT, 10, 10,
	}

	if o_ret.O.Id != o_req.Id {
		t.Fatalf("Ids do not match %d vs %d", o_ret.O.Id, o_req.Id)
	}
	if o_ret.O.Price != o_req.Price {
		t.Fatalf("Prices do not match %f vs %f", o_ret.O.Price, o_req.Price)
	}
	if o_ret.O.Side != o_req.Side {
		t.Fatalf("Sides do not match %d vs %d", o_ret.O.Side, o_req.Side)
	}
	if o_ret.O.Type != o_req.Type {
		t.Fatalf("Types do not match %d vs %d", o_ret.O.Type, o_req.Type)
	}
	if o_ret.O.Size != o_req.Size {
		t.Fatalf("Sizes do not match %d vs %d", o_ret.O.Size, o_req.Size)
	}
}

func BenchmarkParser(b *testing.B) {
	// run the Fib function b.N times
	for n := 0; n < b.N; n++ {
		s := "0;10.0;0;0;10"
		Parse(s)
	}
}
