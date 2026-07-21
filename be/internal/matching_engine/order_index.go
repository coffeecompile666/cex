package matching_engine

import "icon_exchange/internal/shared"

type IOrderIndex interface {
	Push(cmd Command) error
	Pop() Command
	GetLength() uint
}

type Command struct {
	ID uint
}

const MaxBufferSize = 100

type OrderIndex struct {
	oi chan Command
}

func NewOrderIndex() *OrderIndex {
	return &OrderIndex{
		oi: make(chan Command, MaxBufferSize),
	}
}

func (o *OrderIndex) Push(cmd Command) error {
	select {
	case o.oi <- cmd:
		return nil
	default:
		return shared.ErrQueueFull
	}
}

func (o *OrderIndex) Pop() Command {
	return <-o.oi
}

func (o *OrderIndex) GetLength() uint {
	return uint(len(o.oi))
}
