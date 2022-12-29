package pipe

import (
	"context"
	"github.com/google/uuid"
)

const bufferSize uint = 128

type StageInterface interface {
	To(via ReceiveChannel) *ToStage
	Pipeline(ctx context.Context)
	Id() string
}

type Stage struct {
	id string
}

func NewStage() *Stage {
	return &Stage{id: uuid.New().String()}
}

// Pipeline does nothing
func (s *Stage) Pipeline(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		}
	}
}

// Id returns the ID of the stage
func (s *Stage) Id() string {
	return s.id
}

func (s *Stage) To(via ReceiveChannel) *ToStage {
	return &ToStage{
		StageInterface: s,
		via:            via,
	}
}

type ToStage struct {
	StageInterface
	via ReceiveChannel
}
