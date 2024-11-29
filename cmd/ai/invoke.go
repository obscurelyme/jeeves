package ai

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	"github.com/obscurelyme/jeeves/config"
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

var invokeCmd = &cobra.Command{
	Use:   "invoke",
	Short: "Invoke the AWS Bedrock for a single use prompt.",
	Long:  "Invoke the AWS Bedrock for a single use prompt.",
	RunE:  invokeCmdHandler,
}

func invokeCmdHandler(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	loader := config.AWSConfigLoader{}
	cfg, err := loader.LoadAWSConfig("default")
	if err != nil {
		return err
	}

	res, err := InvokeTitanText(cfg, ctx, "Can you write up a \"Hello, World\" program in Golang?")
	if err != nil {
		return err
	}
	fmt.Println(res)
	return nil
}

func InvokeTitanText(cfg aws.Config, ctx context.Context, prompt string) (string, error) {
	bedrockClient := bedrockruntime.NewFromConfig(cfg)
	modelId := "amazon.titan-text-lite-v1"

	body, err := json.Marshal(TitanTextRequest{
		InputText: prompt,
		TextGenerationConfig: &TitanTextGenerationConfig{
			Temperature:   0,
			TopP:          1,
			MaxTokenCount: 4096,
		},
	})
	if err != nil {
		return "", err
	}

	output, err := bedrockClient.InvokeModel(ctx, &bedrockruntime.InvokeModelInput{
		ModelId:     aws.String(modelId),
		ContentType: aws.String("application/json"),
		Body:        body,
	})

	if err != nil {
		return "", err
	}

	var response TitanTextResponse
	err = json.Unmarshal(output.Body, &response)
	if err != nil {
		return "", err
	}

	return response.Results[0].OutputText, nil
}
