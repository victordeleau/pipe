package main

import (
	"context"
	"fmt"
	pipe "github.com/victordeleau/pipe/pkg"
)

type Printer struct {
	*pipe.Stage
	Input pipe.ReceiveChannel
}

func newPrinter() *Printer {
	return &Printer{
		Stage: pipe.NewStage(),
		Input: pipe.NewChannel[struct{}](),
	}
}

func (p *Printer) Pipeline(ctx context.Context) {
	for {
		v := p.Input.Receive(ctx)
		if v == nil {
			return
		}
		fmt.Printf("print %v\n", v)
	}
}
