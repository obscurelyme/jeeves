package titan

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	bedrockruntimeTypes "github.com/aws/aws-sdk-go-v2/service/bedrockruntime/types"
	titan "github.com/obscurelyme/jeeves/ai/titan/types"
	"github.com/obscurelyme/jeeves/ai/types"
)

type TitanAI struct {
	client     *bedrockruntime.Client
	modelId    string
	tokenCount int32
}

func (ai *TitanAI) Invoke(prompt string) (string, error) {
	body, err := json.Marshal(titan.TextRequest{
		InputText: prompt,
		TextGenerationConfig: &titan.TextGenerationConfig{
			Temperature:   0,
			TopP:          1,
			MaxTokenCount: ai.tokenCount,
		},
	})
	if err != nil {
		return "", err
	}

	output, err := ai.client.InvokeModel(context.Background(), &bedrockruntime.InvokeModelInput{
		ModelId:     aws.String(ai.modelId),
		ContentType: aws.String("application/json"),
		Body:        body,
	})
	if err != nil {
		return "", err
	}

	res, err := ai.unmarshal(output)
	if err != nil {
		return "", err
	}
	return res.Results[0].OutputText, nil
}

func (ai *TitanAI) InvokeStream(prompt string, handler types.StreamingOutputHandler) error {
	body, err := json.Marshal(titan.TextRequest{
		InputText: prompt,
		TextGenerationConfig: &titan.TextGenerationConfig{
			Temperature:   0,
			TopP:          1,
			MaxTokenCount: ai.tokenCount,
		},
	})
	if err != nil {
		return err
	}

	output, err := ai.client.InvokeModelWithResponseStream(context.Background(), &bedrockruntime.InvokeModelWithResponseStreamInput{
		ModelId:     aws.String(ai.modelId),
		ContentType: aws.String("application/json"),
		Body:        body,
	})
	if err != nil {
		return err
	}

	return ai.processStreamOutput(context.Background(), output, handler)
}

func (ai *TitanAI) unmarshal(output *bedrockruntime.InvokeModelOutput) (*titan.TextResponse, error) {
	var response titan.TextResponse
	err := json.Unmarshal(output.Body, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (ai *TitanAI) processStreamOutput(
	ctx context.Context,
	output *bedrockruntime.InvokeModelWithResponseStreamOutput,
	handler types.StreamingOutputHandler) error {
	for event := range output.GetStream().Events() {
		switch v := event.(type) {
		case *bedrockruntimeTypes.ResponseStreamMemberChunk:
			var presp *titan.TextStreamedResponse
			err := json.Unmarshal(v.Value.Bytes, &presp)
			if err != nil {
				return err
			}
			err = handler(ctx, presp.OutputText)
			if err != nil {
				return err
			}
		case *bedrockruntimeTypes.UnknownUnionMember:
			fmt.Printf("unknown tag: %s", v.Tag)
		default:
			fmt.Print("union is nil or unknown type")
		}
	}

	return nil
}

type TitanAiInput struct {
	Client     *bedrockruntime.Client
	ModelId    string
	TokenCount int32
}

func New(input *TitanAiInput) *TitanAI {
	titanAi := new(TitanAI)
	titanAi.client = input.Client
	titanAi.modelId = input.ModelId
	titanAi.tokenCount = input.TokenCount
	return titanAi
}
