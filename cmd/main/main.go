package main

import (
	"context"
	"github.com/victordeleau/pipe/pkg"
	"log"
	"time"
)

func main() {

	pipeline := pipe.NewPipe()

	timerStage, printerStage := newTimer(), newPrinter()
	pipeline.Add(timerStage, printerStage)

	err := pipeline.Add(timerStage, printerStage).Link(timerStage.Output, printerStage.Input)
	if err != nil {
		log.Fatalf(err.Error())
	}

	if err = pipeline.Compile(); err != nil {
		log.Fatalf(err.Error())
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	pipeline.Start(ctx)
}
