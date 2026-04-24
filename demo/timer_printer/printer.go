package main

import (
	"context"
	"fmt"
	"github.com/victordeleau/pipe/pkg"
)

type Printer struct {
	*pipe.Stage
	Input pipe.ReceiveChannel[struct{}]
}

func newPrinter() *Printer {
	return &Printer{
		Stage: pipe.NewStage(),
		Input: pipe.NewReceiveChannel[struct{}](),
	}
}

func (p *Printer) Pipeline(ctx context.Context) {
	fmt.Print("printer starting\n")
	for {
		v, ok := p.Input.Receive(ctx)
		if !ok {
			return
		}
		fmt.Printf("printer received %v\n", v)
	}
}
