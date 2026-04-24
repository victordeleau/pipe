package pipe

import (
	"context"
	"fmt"
	"github.com/google/uuid"
)

type ChannelType uint

const (
	ReceiveChannelType ChannelType = 0
	SendChannelType    ChannelType = 1
)

func (c ChannelType) IsValid() error {
	if c != ReceiveChannelType && c != SendChannelType {
		return fmt.Errorf("invalid channel type")
	}
	return nil
}

type ChannelInterface interface {
	Make(size uint)
	Id() string
	ReceiveFrom(channel ChannelInterface) error
	Type() ChannelType
	Len() int
}

type ReceiveChannel[T any] interface {
	ChannelInterface
	// Receive returns (v, true) when a value was received from the channel.
	// It returns (zero, false) when ctx is done or the channel is closed.
	Receive(ctx context.Context) (T, bool)
	// ReadChan exposes the stream for range/select without leaking channel internals to helpers like Fanin.
	ReadChan() <-chan T
}

type SendChannel[T any] interface {
	ChannelInterface
	Send(payload ...T)
}

// Channel is a basic Channel
type Channel[T any] struct {
	buffer      *chan T
	id          string
	channelType ChannelType
}

func NewReceiveChannel[T any]() ReceiveChannel[T] {
	return &Channel[T]{
		buffer:      new(chan T),
		id:          uuid.New().String(),
		channelType: ReceiveChannelType,
	}
}

func NewSendChannel[T any]() SendChannel[T] {
	return &Channel[T]{
		buffer:      new(chan T),
		id:          uuid.New().String(),
		channelType: SendChannelType,
	}
}

func (c *Channel[T]) Len() int {
	return len(*c.buffer)
}

func (c *Channel[T]) Make(size uint) {
	*c.buffer = make(chan T, size)
}

func (c *Channel[T]) Id() string {
	return c.id
}

func (c *Channel[T]) Send(payload ...T) {
	for _, p := range payload {
		*c.buffer <- p
	}
}

func (c *Channel[T]) Receive(ctx context.Context) (T, bool) {
	var zero T
	select {
	case <-ctx.Done():
		return zero, false
	case v, ok := <-c.ReadChan():
		if !ok {
			return zero, false
		}
		return v, true
	}
}

func (c *Channel[T]) ReadChan() <-chan T {
	return *c.buffer
}

func (c *Channel[T]) ReceiveFrom(channel ChannelInterface) error {
	channelInferred, ok := channel.(*Channel[T])
	if !ok {
		return fmt.Errorf("'from' channel doesn't have the same type as the 'to' channel")
	}
	c.buffer = channelInferred.buffer
	return nil
}

func (c *Channel[T]) Type() ChannelType {
	return c.channelType
}
