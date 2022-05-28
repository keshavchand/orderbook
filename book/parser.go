package book

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type RequestType int

const (
	ReqNewOrder RequestType = iota
	ReqRemoveOrder
	ReqUpdateOrder

	ReqGetUnits

	ReqTypeCount
)

func (t RequestType) Valid() bool {
	if t < 0 && t >= ReqTypeCount {
		return false
	}

	return true
}

type ParsedInfo interface {
	// This interface is inteded to be returned by the parse function
	// this function doesn't do anything its only purpose is to resitrict
	// the struct types returned by the parse function
	ParsedDoNothing()
}

// Parser Reqest format
// ReqType:RequestType(int);<remaining info>
func Parse(s string) (ParsedInfo, error) {
	parts := strings.Split(s, ";")
	if len(parts) <= 0 {
		return nil, errors.New("not enough parameters")
	}

	rt, err := strconv.Atoi(parts[0])
	if err != nil {
		return nil, fmt.Errorf("cant parse request type %w", err)
	}

	reqType := RequestType(rt)
	if !reqType.Valid() {
		return nil, errors.New("invalid req type")
	}

	switch reqType {
	case ReqNewOrder:
		return ParseNewOrder(parts[1:])
	case ReqRemoveOrder:
		return ParseRemoveOrder(parts[1:])
	case ReqUpdateOrder:
		return ParseUpdateOrder(parts[1:])
	case ReqGetUnits:
		return ParseGetUnits(parts[1:])
	}

	return nil, errors.New("how tf it reached here?!?!?!")
}

type RequestParserFunc func(*Order, string) error

func parsePrice(o *Order, s string) error {
	price, err := strconv.ParseFloat(s, 32)
	if err != nil {
		return fmt.Errorf("parsing price: error %w", err)
	}
	o.Price = float32(price)
	return nil
}

func parseSide(o *Order, s string) error {
	pside, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
		return fmt.Errorf("parsing side: error %w", err)
	}
	side := OrderSide(int(pside))
	if !side.Valid() {
		return fmt.Errorf("invalid order side")
	}
	o.Side = side
	return nil
}

func parseType(o *Order, s string) error {
	potype, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
		return fmt.Errorf("parsing type: error %w", err)
	}
	otype := OrderType(int(potype))
	if !otype.Valid() {
		return fmt.Errorf("side invalid")
	}
	o.Type = otype
	return nil
}

func parseSize(o *Order, s string) error {
	psize, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
		return fmt.Errorf("parsing size: error %w", err)
	}
	size := int(psize)
	o.Size = size
	return nil
}

func parseId(o *Order, s string) error {
	pid, err := strconv.ParseUint(s, 10, 32)
	if err != nil {
		return fmt.Errorf("parsing size: error %w", err)
	}
	id := pid
	o.Id = id
	return nil
}

type NewOrder struct {
	O Order
}

// Price:float;Side:OrderSide(int);Type:OrderType(int);Size:int
// NOTE: Caller must assign the id to the order
func (_ NewOrder) ParsedDoNothing() {}
func ParseNewOrder(s []string) (NewOrder, error) {
	var o Order
	if len(s) < 4 {
		return NewOrder{}, errors.New("parsing order: not enough args")
	}

	newOrderConds := []RequestParserFunc{parsePrice, parseSide, parseType, parseSize}
	for idx, c := range newOrderConds {
		if err := c(&o, s[idx]); err != nil {
			return NewOrder{}, err
		}
	}
	return NewOrder{o}, nil
}

type RemoveOrder struct {
	Id uint64
}

func (_ RemoveOrder) ParsedDoNothing() {}

// id
func ParseRemoveOrder(s []string) (RemoveOrder, error) {
	var o Order
	if len(s) < 1 {
		return RemoveOrder{0}, errors.New("parsing order: not enough args")
	}
	if err := parseId(&o, s[0]); err != nil {
		return RemoveOrder{0}, err
	}
	return RemoveOrder{o.Id}, nil
}

type UpdateOrder struct {
	O Order
}

func (_ UpdateOrder) ParsedDoNothing() {}

// id:int;Price:float;Side:OrderSide(int);Type:OrderType(int);Size:int
func ParseUpdateOrder(s []string) (UpdateOrder, error) {
	var o Order
	if len(s) < 5 {
		return UpdateOrder{}, errors.New("parsing order: not enough args")
	}

	newOrderConds := []RequestParserFunc{parseId, parsePrice, parseSide, parseType, parseSize}
	for idx, c := range newOrderConds {
		if err := c(&o, s[idx]); err != nil {
			return UpdateOrder{}, err
		}
	}

	return UpdateOrder{o}, nil
}

type GetUnits struct{}

func (_ GetUnits) ParsedDoNothing()              {}
func ParseGetUnits(s []string) (GetUnits, error) { return GetUnits{}, nil }
