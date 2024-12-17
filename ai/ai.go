package ai

import (
	"errors"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	"github.com/obscurelyme/jeeves/ai/nova"
	"github.com/obscurelyme/jeeves/ai/titan"
	"github.com/obscurelyme/jeeves/ai/types"
)

type NewInvokeDriverInput struct {
	ModelId    string
	Client     *bedrockruntime.Client
	TokenCount int32
}

func New(input *NewInvokeDriverInput) (types.InvokeDriver, error) {
	if input == nil {
		return nil, errors.New("input interface is required to create a new invoke driver")
	}

	if strings.Contains(input.ModelId, string(types.Titan)) {
		p := titan.New(&titan.TitanAiInput{
			Client:     input.Client,
			ModelId:    input.ModelId,
			TokenCount: input.TokenCount,
		})
		return p, nil
	} else if strings.Contains(input.ModelId, string(types.Nova)) {
		p := nova.New(&nova.NovaAiInput{
			Client:     input.Client,
			ModelId:    input.ModelId,
			TokenCount: input.TokenCount,
		})
		return p, nil
	}

	return nil, errors.New("not implemented")
}
