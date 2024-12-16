package nova

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	nova "github.com/obscurelyme/jeeves/ai/nova/types"
	"github.com/obscurelyme/jeeves/ai/types"
)

type NovaAi struct {
	clientId   *bedrockruntime.Client
	modelId    string
	tokenCount int32
}

func (ai *NovaAi) Invoke(prompt string) (string, error) {
	return "", nil
}

func (ai *NovaAi) InvokeStream(prompt string, handler types.StreamingOutputHandler) error {
	return nil
}

func (ai *NovaAi) processStreamOutput(ctx context.Context,
	output *bedrockruntime.InvokeModelWithResponseStreamOutput,
	handler types.StreamingOutputHandler) error {
	return nil
}

func (ai *NovaAi) unmarshal(output *bedrockruntime.InvokeModelOutput) (*nova.Body, error) {
	var response nova.Body
	err := json.Unmarshal(output.Body, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}
