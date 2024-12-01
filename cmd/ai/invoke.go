package ai

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime/types"
	"github.com/obscurelyme/jeeves/config"
	"github.com/obscurelyme/jeeves/prompt"
	"github.com/obscurelyme/jeeves/utils"
	"github.com/spf13/cobra"
)

type TitanTextResult struct {
	TokenCount       int32  `json:"tokenCount"`
	OutputText       string `json:"outputText"`
	CompletionReason string `json:"completionReason"`
}

type TitanTextResponse struct {
	InputTextTokenCount int32             `json:"inputTextTokenCount"`
	Results             []TitanTextResult `json:"results"`
}

type TitanTextGenerationConfig struct {
	// Use a lower value to decrease randomness in responses.
	Temperature float64 `json:"temperature"`
	// Use a lower value to ignore less probable options and decrease the diversity of responses.
	TopP float64 `json:"topP"`
	// Specify the maximum number of tokens to generate in the response. Maximum token limits are strictly enforced.
	MaxTokenCount int32 `json:"maxTokenCount"`
	// Specify a character sequence to indicate where the model should stop.
	StopSequences []string `json:"stopSequences,omitempty"`
}

type TitanTextRequest struct {
	// The prompt to provide the model for generating a response.
	// To generate responses in a conversational style, wrap the prompt by using the following format:
	//
	// "inputText": "User: <prompt>\nBot:
	InputText string `json:"inputText"`
	// Optional, used to configure inference parameters
	TextGenerationConfig *TitanTextGenerationConfig `json:"textGenerationConfig"`
}

type InvocationMetrics struct {
	InputTokenCount   int `json:"inputTokenCount"`
	OutputTokenCount  int `json:"outputTokenCount"`
	InvocationLatency int `json:"invocationLatency"`
	FirstByteLatency  int `json:"firstByteLatency"`
}

type TitanTextStreamedResponse struct {
	OutputText                string            `json:"outputText"`
	Index                     int               `json:"index"`
	TotalOutputTextTokenCount int               `json:"totalOutputTextTokenCount"`
	CompletionReason          string            `json:"completionReason"`
	InputTextTokenCount       int               `json:"inputTextTokenCount"`
	InvocationMetrics         InvocationMetrics `json:"amazon-bedrock-invocationMetrics"`
}

var invokeCmd = &cobra.Command{
	Use:   "invoke",
	Short: "Invoke the AWS Bedrock for a single use prompt.",
	Long:  "Invoke the AWS Bedrock for a single use prompt.",
	RunE:  invokeCmdHandler,
}

func tokenCount(modelId string) int32 {
	switch modelId {
	case "amazon.titan-text-premier-v1:0":
		return 521
	default:
		return 4096
	}
}

func invoke(cfg aws.Config, ctx context.Context) error {
	input, err := prompt.QuickPrompt("Input > ")
	if input == "exit" || input == "quit" {
		return nil
	}
	if input == "ty" {
		fmt.Println("You're very welcome!")
		return nil
	}
	if err != nil {
		return err
	}

	fmt.Print("thinking...\n\n")
	_, err = InvokeTitanText(cfg, ctx, input)
	if err != nil {
		return err
	}

	fmt.Print("\n---END---\n\n")

	return invoke(cfg, ctx)
}

func invokeCmdHandler(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	loader := config.AWSConfigLoader{}
	cfg, err := loader.LoadAWSConfig("default")
	if err != nil {
		return err
	}

	return invoke(cfg, ctx)
}

func InvokeTitanText(cfg aws.Config, ctx context.Context, prompt string) (string, error) {
	bedrockClient := bedrockruntime.NewFromConfig(cfg)
	modelId := utils.Jeeves.ConfigSettings.AI.PreferredModel

	body, err := json.Marshal(TitanTextRequest{
		InputText: prompt,
		TextGenerationConfig: &TitanTextGenerationConfig{
			Temperature:   0,
			TopP:          1,
			MaxTokenCount: tokenCount(modelId),
		},
	})
	if err != nil {
		return "", err
	}

	streamOutput, err := bedrockClient.InvokeModelWithResponseStream(ctx, &bedrockruntime.InvokeModelWithResponseStreamInput{
		ModelId:     aws.String(modelId),
		ContentType: aws.String("application/json"),
		Body:        body,
	})

	if err != nil {
		return "", err
	}

	resp, err := processStreamOutput(ctx, streamOutput, func(ctx context.Context, part TitanTextStreamedResponse) error {
		fmt.Print(part.OutputText)
		return nil
	})

	if err != nil {
		return "", err
	}
	return resp.OutputText, nil
}

type StreamingOutputHandler func(ctx context.Context, part TitanTextStreamedResponse) error

func processStreamOutput(ctx context.Context, output *bedrockruntime.InvokeModelWithResponseStreamOutput, handler StreamingOutputHandler) (TitanTextStreamedResponse, error) {
	resp := TitanTextStreamedResponse{}

	for event := range output.GetStream().Events() {
		switch v := event.(type) {
		case *types.ResponseStreamMemberChunk:
			var presp *TitanTextStreamedResponse
			err := json.Unmarshal(v.Value.Bytes, &presp)
			if err != nil {
				return *presp, err
			}
			err = handler(ctx, *presp)
			if err != nil {
				return *presp, err
			}
		case *types.UnknownUnionMember:
			fmt.Printf("unknown tag: %s", v.Tag)
		default:
			fmt.Print("union is nil or unknown type")
		}
	}

	return resp, nil
}
