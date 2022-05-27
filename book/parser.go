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

	ReqTypeCount
)

func (t RequestType) Valid() bool {
	if t < 0 && t >= ReqTypeCount {
		return false
	}

	return true
}

// Parser Reqest format
// ReqType:RequestType(int);<remaining info>
func Parse(s string) (interface{}, error) {
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

// Price:float;Side:OrderSide(int);Type:OrderType(int);Size:int
// NOTE: Caller must assign the id to the order
func ParseNewOrder(s []string) (Order, error) {
	var o Order
	if len(s) < 4 {
		return Order{}, errors.New("parsing order: not enough args")
	}

	newOrderConds := []RequestParserFunc{parsePrice, parseSide, parseType, parseSize}
	for idx, c := range newOrderConds {
		if err := c(&o, s[idx]); err != nil {
			return Order{}, err
		}
	}
	return o, nil
}

// id
func ParseRemoveOrder(s []string) (Order, error) {
	var o Order
	if len(s) < 1 {
		return Order{}, errors.New("parsing order: not enough args")
	}
	if err := parseId(&o, s[0]); err != nil {
		return Order{}, err
	}
	return o, nil
}

// id:int;Price:float;Side:OrderSide(int);Type:OrderType(int);Size:int
func ParseUpdateOrder(s []string) (Order, error) {
	var o Order
	if len(s) < 5 {
		return Order{}, errors.New("parsing order: not enough args")
	}
	if err := parseId(&o, s[0]); err != nil {
		return Order{}, err
	}

	newOrderConds := []RequestParserFunc{parsePrice, parseSide, parseType, parseSize}
	for idx, c := range newOrderConds {
		if err := c(&o, s[idx+1]); err != nil {
			return Order{}, err
		}
	}

	return o, nil
}
