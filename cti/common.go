package cti

import "errors"

type TradedOrder struct {
	To    uint64
	From  uint64
	Size  int
	Price float32
}

const (
	MaxSenderId   = (1 << 16) - 1
	MaxOrderCount = (1 << (64 - 16)) - 1
)

var (
	ErrSenderIdOutOfRange   = errors.New("sender id out of range")
	ErrOrderCountOutOfRange = errors.New("order count out of range")
)

func CreateOrderId(sender_id, count uint64) (uint64, error) {
	if sender_id > MaxSenderId {
		return 0, ErrSenderIdOutOfRange
	}
	if sender_id > MaxOrderCount {
		return 0, ErrOrderCountOutOfRange
	}

	return (sender_id << (64 - 16)) | (count), nil
}

func ParseOrderId(id uint64) (sender_id, count uint64) {
	sender_id = (id >> (64 - 16)) & MaxSenderId
	count = ((id) & MaxOrderCount)

	return sender_id, count
}
