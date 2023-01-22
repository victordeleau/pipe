package pipe

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"reflect"
)

type PipelineInterface interface {
	Add(stages ...StageInterface) *Pipeline
	Link(from ChannelInterface, to ...ChannelInterface) error
	Compile() error
	Start(ctx context.Context)
	Stop() <-chan struct{}
}

// StageWrapper wraps the StageInterface with information regarding the stages is feeds.
type StageWrapper struct {
	StageInterface
	isFeedingStages set
	OutputChannels  map[string]*ChannelWrapper
	InputChannels   map[string]*ChannelWrapper
}

// ChannelWrapper wraps the ChannelInterface with information regarding the channels it feeds,
// and the stage it belongs to.
type ChannelWrapper struct {
	ChannelInterface
	isFeedingChannels set
	parentStage       string
}

func (c *ChannelWrapper) isConsumed() bool {
	return len(c.isFeedingChannels) != 0
}

type Pipeline struct {
	id       string
	stages   map[string]*StageWrapper   // stage ID indexing stages
	channels map[string]*ChannelWrapper // channel ID indexing channels
	option   Option
}

func NewPipe(option ...Option) *Pipeline {
	pipeline := &Pipeline{
		id:       uuid.New().String(),
		stages:   make(map[string]*StageWrapper),
		channels: make(map[string]*ChannelWrapper),
	}
	if len(option) != 0 {
		pipeline.option = option[0]
	}
	return pipeline
}

func (p *Pipeline) Add(stages ...StageInterface) *Pipeline {

	for _, s := range stages {

		// add stage
		if _, ok := p.stages[s.Id()]; ok {
			continue // already added
		}
		p.stages[s.Id()] = &StageWrapper{
			StageInterface:  s,
			isFeedingStages: make(set),
			InputChannels:   make(map[string]*ChannelWrapper),
			OutputChannels:  make(map[string]*ChannelWrapper),
		}

		// add channels
		stageType, stageValue := reflect.TypeOf(s).Elem(), reflect.ValueOf(s).Elem()
		for i := 0; i < stageType.NumField(); i++ {

			if stageType.Field(i).Type.Implements(reflect.TypeOf((*ChannelInterface)(nil)).Elem()) {

				channel := stageValue.Field(i).Interface().(ChannelInterface)

				channelWrapper := &ChannelWrapper{
					ChannelInterface:  channel,
					isFeedingChannels: make(set),
					parentStage:       s.Id(),
				}

				p.channels[channelWrapper.Id()] = channelWrapper

				if channel.Type() == ReceiveChannelType {
					p.stages[s.Id()].InputChannels[channel.Id()] = channelWrapper
				} else { // then SendChannelType
					p.stages[s.Id()].OutputChannels[channel.Id()] = channelWrapper
				}
			}
		}
	}

	return p
}

// Link configure stage 'to' to receive from the provided channels.
// The provided channels must appear as fields with type 'ChannelInterface' in the 'to' pipeline.
func (p *Pipeline) Link(from ChannelInterface, to ...ChannelInterface) error {

	if p == nil {
		return nil
	}

	var found bool

	// is the channel is found, then it means the stage it belongs to was added
	fromChannelWrapper, fromStageAdded := p.channels[from.Id()]

	for _, t := range to {

		if from.Type() != SendChannelType {
			return fmt.Errorf("from channel is not a receive channel")
		}

		if t.Type() != ReceiveChannelType {
			return fmt.Errorf("to channel %s is not a receive channel", t.Id())
		}

		if _, found = p.channels[t.Id()]; !found {
			return fmt.Errorf("'to' channel %s belongs to a stage was not previously added", t.Id())
		}

		if fromStageAdded {
			fromStage := p.stages[fromChannelWrapper.parentStage]
			fromStage.isFeedingStages[p.channels[t.Id()].parentStage] = struct{}{}
			p.stages[fromChannelWrapper.parentStage] = fromStage

			fromChannelWrapper.isFeedingChannels[t.Id()] = struct{}{}
		}

		if err := t.ReceiveFrom(from); err != nil {
			return err
		}
	}

	if fromStageAdded {
		p.channels[from.Id()] = fromChannelWrapper
	}

	return nil
}

// Compile the current graph and verify its validity.
func (p *Pipeline) Compile(bufferSize uint) (*CompiledPipeline, error) {

	// check that there is no cycle
	if err := p.isCycle(); err != nil {
		return nil, err
	}

	// check no open-ended
	for channelId, channelWrapper := range p.channels {
		if !channelWrapper.isConsumed() && channelWrapper.Type() == SendChannelType {
			return nil, fmt.Errorf("channel %s is not consumed by any pipeline", channelId)
		}
	}

	// get topological order
	compiled := &CompiledPipeline{
		Pipeline:         p,
		topologicalOrder: p.topologicalSort(),
	}

	// from beginning to end, initialize
	for _, stageId := range compiled.topologicalOrder {
		for _, outputChannel := range compiled.stages[stageId].OutputChannels {
			outputChannel.Make(bufferSize)
		}
	}

	return compiled, nil
}

func (p *Pipeline) isCycleUtility(nodeId string, visited set, stack set) bool {

	visited[nodeId], stack[nodeId] = struct{}{}, struct{}{}

	for stageId, _ := range p.stages[nodeId].isFeedingStages {
		if _, found := visited[stageId]; !found {
			if p.isCycleUtility(stageId, visited, stack) {
				return true
			} else if _, found = stack[stageId]; found {
				return true
			}
		}
	}

	delete(stack, nodeId)
	return false
}

func (p *Pipeline) isCycle() error {

	visited, stack := make(set, len(p.stages)+1), make(set, len(p.stages)+1)

	for nodeId, _ := range p.stages {
		if _, ok := visited[nodeId]; !ok {
			if p.isCycleUtility(nodeId, visited, stack) {
				return fmt.Errorf("cycle detected in graph")
			}
		}
	}

	return nil
}

func (p *Pipeline) topologicalSortUtility(stageId string, visited *set, order *[]string) {

	(*visited)[stageId] = struct{}{}

	for neighborId, _ := range p.stages[stageId].isFeedingStages {
		if _, ok := (*visited)[neighborId]; !ok {
			p.topologicalSortUtility(neighborId, visited, order)
		}
	}

	*order = append(*order, stageId)
}

func (p *Pipeline) topologicalSort() []string {

	visited := make(set, len(p.stages))
	order := make([]string, 0, len(p.stages))

	for stageId, _ := range p.stages {
		if _, ok := visited[stageId]; !ok {
			p.topologicalSortUtility(stageId, &visited, &order)
		}
	}

	reversedStack := make([]string, 0, len(p.stages))
	for i := len(order) - 1; i >= 0; i-- {
		reversedStack = append(reversedStack, order[i])
	}

	return reversedStack
}
