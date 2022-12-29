package main

import (
	"context"
	pipe "github.com/victordeleau/pipe/pkg"
	"time"
)

type Timer struct {
	*pipe.Stage
	Output pipe.ReceiveChannel
}

func newTimer() *Timer {
	return &Timer{
		Stage:  pipe.NewStage(),
		Output: pipe.NewChannel[struct{}]()}
}

func (t *Timer) Pipeline(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(time.Second):
			t.Output.Send(struct{}{})
		}
	}
}
