package pipe

import (
	"context"
	"fmt"
)

type PipelineInterface interface {
	Add(stages ...StageInterface) *Pipeline
	Link(from ChannelInterface, to ...ChannelInterface) error
	Compile() error
	Start(ctx context.Context)
	Stop() <-chan struct{}
}

type Pipeline struct {
	links      map[string][]ReceiveChannel
	stages     map[string][]StageInterface
	isConsumed map[string]bool

	isCompiled bool

	visited map[string]bool
	stack   map[string]bool

	option Option
}

func NewPipe(option ...Option) *Pipeline {
	pipeline := &Pipeline{
		links:      make(map[string][]ReceiveChannel),
		stages:     make(map[string][]StageInterface),
		isConsumed: make(map[string]bool),
		visited:    make(map[string]bool),
		stack:      make(map[string]bool),
	}
	if len(option) != 0 {
		pipeline.option = option[0]
	}
	return pipeline
}

func (p *Pipeline) Add(stages ...StageInterface) *Pipeline {

	// add stages
	for _, s := range stages {
		if _, ok := p.stages[s.Id()]; !ok {
			p.stages[s.Id()] = append(p.stages[s.Id()], s)
		}
	}

	// TODO set that stage outputs must be consumed

	return p
}

// Link configure stage 'to' to receive from the provided channels.
// The provided channels must appear as fields with type 'ChannelInterface' in the 'to' main.
func (p *Pipeline) Link(from ReceiveChannel, to ...ReceiveChannel) error {

	if p == nil {
		return nil
	}

	for _, t := range to {

		// add link to main
		if _, ok := (p.links)[from.Id()]; !ok {
			(p.links)[from.Id()] = []ReceiveChannel{t}
		} else {
			(p.links)[from.Id()] = append((p.links)[t.Id()], t)
		}

		t.ReceiveFrom(from) // TODO set inputs in 'to' stage

		// set that 'from' channel is consumed to later detect un-consumed channels
		if isConsumed, ok := p.isConsumed[t.Id()]; !ok {
			return fmt.Errorf("'from' channel does not belong to any known stage")
		} else if !isConsumed {
			p.isConsumed[t.Id()] = true
		}
	}

	return nil
}

// Compile the current graph and verify its validity.
func (p *Pipeline) Compile() error {

	// check that there is no cycle
	p.visited = make(map[string]bool, len(p.stages)+1)
	p.stack = make(map[string]bool, len(p.stages)+1)
	for nodeId, _ := range p.stages {
		if _, ok := p.visited[nodeId]; !ok {
			if p.isCycle(nodeId) {
				return fmt.Errorf("cycle detected in graph")
			}
		}
	}

	// check no open-ended
	for id, isConsumed := range p.isConsumed {
		if !isConsumed {
			return fmt.Errorf("channel %s is not consumed by any main", id)
		}
	}

	// get topological order

	// from beginning to end, initialize

	p.isCompiled = true

	return nil
}

func (p *Pipeline) isCycle(nodeId string) bool {

	p.visited[nodeId], p.stack[nodeId] = true, true

	for _, neighbor := range p.stages[nodeId] {
		if !p.visited[neighbor.Id()] {
			if p.isCycle(neighbor.Id()) {
				return true
			} else if p.stack[neighbor.Id()] == true {
				return true
			}
		}
	}

	p.stack[nodeId] = false
	return false
}

func (p *Pipeline) Start(ctx context.Context) {

	if !p.isCompiled {
		return
	}

	return
}

func (p *Pipeline) Stop() <-chan struct{} {
	done := make(chan struct{}, 1)
	return done
}
