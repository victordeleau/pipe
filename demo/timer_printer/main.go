package main

import (
	"context"
	"fmt"
	"github.com/victordeleau/pipe/pkg"
	"log"
	"time"
)

func main() {

	timerStage, printerStage := newTimer(), newPrinter()
	fmt.Printf("timer stage %s - printer stage %s\n", timerStage.Id(), printerStage.Id())

	pipeline := pipe.NewPipe()
	err := pipeline.Add(timerStage, printerStage).Link(timerStage.Output, printerStage.Input)
	if err != nil {
		log.Fatalf(err.Error())
	}

	compiled, err := pipeline.Compile(16)
	if err != nil {
		log.Fatalf(err.Error())
	}

	compiled.Log()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	compiled.Start(ctx)

	select {
	case <-ctx.Done():
		break
	}
}
