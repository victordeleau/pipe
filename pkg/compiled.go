package pipe

import (
	"context"
	"fmt"
	"sync"
)

type CompiledPipeline struct {
	*Pipeline
	topologicalOrder []string
	wg               sync.WaitGroup
}

func (c *CompiledPipeline) Start(ctx context.Context) {
	for _, stageId := range c.topologicalOrder {
		c.wg.Add(1)
		sid := stageId
		go func() {
			defer c.wg.Done()
			c.stages[sid].Pipeline(ctx)
		}()
	}
}

// Stop returns a channel that is closed after every stage goroutine has returned
// (typically after ctx is cancelled and each Pipeline method exits).
func (c *CompiledPipeline) Stop() <-chan struct{} {
	done := make(chan struct{})
	go func() {
		c.wg.Wait()
		close(done)
	}()
	return done
}

// Wait blocks until every stage goroutine has returned.
func (c *CompiledPipeline) Wait() {
	c.wg.Wait()
}

func (c *CompiledPipeline) Log() {
	fmt.Printf("\n=== Pipeline Id %s\n", c.id)
	for _, stageId := range c.topologicalOrder {
		fmt.Printf("  === Stage ID %s\n", stageId)
		if len(c.stages[stageId].InputChannels) != 0 {
			fmt.Printf("    === Input Channels\n")
			for inputChannelId, _ := range c.stages[stageId].InputChannels {
				fmt.Printf("      === ID %s\n", inputChannelId)
			}
		}
		if len(c.stages[stageId].OutputChannels) != 0 {
			fmt.Printf("    === Output Channels\n")
			for outputChannelId, outputChannel := range c.stages[stageId].OutputChannels {
				fmt.Printf("      === ID %s is feeding channels %+v\n", outputChannelId, outputChannel.isFeedingChannels.Slice())
			}
		}
	}
	fmt.Printf("\n")
}
