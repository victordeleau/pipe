package main

import (
	"context"
	"fmt"
	"github.com/victordeleau/pipe/pkg"
	"time"
)

type Timer struct {
	*pipe.Stage
	Output pipe.SendChannel[struct{}]
}

func newTimer() *Timer {
	return &Timer{
		Stage:  pipe.NewStage(),
		Output: pipe.NewSendChannel[struct{}]()}
}

func (t *Timer) Pipeline(ctx context.Context) {
	fmt.Print("timer starting\n")
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(time.Second):
			fmt.Print("timer sending\n")
			t.Output.Send(struct{}{})
		}
	}
}
