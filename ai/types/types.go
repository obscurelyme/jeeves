package types

import (
	"context"
)

type StreamingOutputHandler func(ctx context.Context, part interface{}) error

type InvokeDriver interface {
	Invoke(prompt string) (string, error)
	InvokeStream(prompt string, handler StreamingOutputHandler) error
	// ProcessStreamOutput(ctx context.Context, output *bedrockruntime.InvokeModelWithResponseStreamOutput, handler StreamingOutputHandler) error
}

type AiType string

const (
	Nova  AiType = "nova"
	Titan AiType = "titan"
)

type ModelSelect struct {
	Label   string
	ModelId string
}
