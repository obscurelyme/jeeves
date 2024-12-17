package nova

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-sdk-go-v2/aws"

	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	nova "github.com/obscurelyme/jeeves/ai/nova/types"
	"github.com/obscurelyme/jeeves/ai/types"
)

type NovaAi struct {
	client     *bedrockruntime.Client
	modelId    string
	tokenCount int32
}

func (ai *NovaAi) Invoke(prompt string) (string, error) {
	body, err := json.Marshal(nova.Body{
		System:   []nova.System{},
		Messages: []nova.Message{},
		InferenceConfig: &nova.InferenceConfig{
			MaxNewTokens: int(ai.tokenCount),
			Temperature:  0.7,
			TopP:         0.9,
			TopK:         20,
		},
	})
	if err != nil {
		return "", err
	}

	output, err := ai.client.InvokeModel(context.Background(), &bedrockruntime.InvokeModelInput{
		ModelId:     aws.String(ai.modelId),
		Body:        body,
		ContentType: aws.String("application/json"),
	})
	if err != nil {
		return "", nil
	}

	res, err := ai.unmarshal(output)
	if err != nil {
		return "", err
	}

	return res.Messages[len(res.Messages)-1].Content[len(res.Messages[len(res.Messages)-1].Content)-1].Text, nil
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

type NovaAiInput struct {
	Client     *bedrockruntime.Client
	ModelId    string
	TokenCount int32
}

func New(input *NovaAiInput) *NovaAi {
	novaAi := new(NovaAi)
	novaAi.client = input.Client
	novaAi.modelId = input.ModelId
	novaAi.tokenCount = input.TokenCount
	return novaAi
}
