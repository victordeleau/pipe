package pipe

import (
	"context"
	"fmt"
)

type CompiledPipeline struct {
	*Pipeline
	topologicalOrder []string
}

func (c *CompiledPipeline) Start(ctx context.Context) {
	for _, stageId := range c.topologicalOrder {
		go c.stages[stageId].Pipeline(ctx)
	}
}

func (c *CompiledPipeline) Stop() <-chan struct{} {
	done := make(chan struct{}, 1)
	return done
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
