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
	Receive(ctx context.Context) *T
}

type SendChannel[T any] interface {
	ChannelInterface
	Send(payload ...T)
}

// Channel is a basic Channel
type Channel[T any] struct {
	Buffer      *chan T
	id          string
	channelType ChannelType
}

func NewReceiveChannel[T any]() ReceiveChannel[T] {
	return &Channel[T]{
		Buffer:      new(chan T),
		id:          uuid.New().String(),
		channelType: ReceiveChannelType,
	}
}

func NewSendChannel[T any]() SendChannel[T] {
	return &Channel[T]{
		Buffer:      new(chan T),
		id:          uuid.New().String(),
		channelType: SendChannelType,
	}
}

func (c *Channel[T]) Len() int {
	return len(*c.Buffer)
}

func (c *Channel[T]) Make(size uint) {
	*c.Buffer = make(chan T, size)
}

func (c *Channel[T]) Id() string {
	return c.id
}

func (c *Channel[T]) Send(payload ...T) {
	for _, p := range payload {
		*c.Buffer <- p
	}
}

func (c *Channel[T]) Receive(ctx context.Context) *T {
	select {
	case <-ctx.Done():
		return nil
	case v := <-*c.Buffer:
		return &v
	}
}

func (c *Channel[T]) ReceiveFrom(channel ChannelInterface) error {
	channelInferred, ok := channel.(*Channel[T])
	if !ok {
		return fmt.Errorf("'from' channel doesn't have the same type as the 'to' channel")
	}
	(*c).Buffer = channelInferred.Buffer
	return nil
}

func (c *Channel[T]) Type() ChannelType {
	return c.channelType
}
