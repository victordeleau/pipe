package pipe

import (
	"context"
	"github.com/google/uuid"
)

type Chan[T any] struct {
	Buffer chan T
	id     string
}

func NewChannel[T any]() ChannelInterface {
	return &Chan[T]{
		id: uuid.New().String(),
	}
}

func (c *Chan[T]) Make(size uint) ChannelInterface {
	if c.Buffer == nil {
		c.Buffer = make(chan T, size)
	}
	return c
}

func (c *Chan[T]) Send(payload ...any) {
	for _, p := range payload {
		c.Buffer <- p
	}
}

func (c *Chan[T]) Receive(ctx context.Context) any {
	select {
	case <-ctx.Done():
		return nil
	case v := <-c.Buffer:
		return v
	}
}

func (c *Chan[T]) ReceiveFrom(channel ReceiveChannel) error {
	// TODO
	return nil
}

func (c *Chan[T]) Id() string {
	return c.id
}

type ReceiveChannel interface {
	Make(size uint) ChannelInterface
	Id() string
	Receive(ctx context.Context) any
	ReceiveFrom(channel ReceiveChannel) error
	Send(...any)
}

type SendChannel interface {
	Make(size uint) ChannelInterface
	Id() string
	Send(...any)
}

type ChannelInterface interface {
	Make(size uint) ChannelInterface
	Id() string
	Send(...any)
	Receive(ctx context.Context) any
	ReceiveFrom(channel ReceiveChannel) error
}

type ChannelMap map[string]ChannelInterface
