package ai

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	"github.com/obscurelyme/jeeves/ai"
	"github.com/obscurelyme/jeeves/config"
	"github.com/obscurelyme/jeeves/prompt"
	"github.com/obscurelyme/jeeves/utils"
	"github.com/spf13/cobra"
)

var invokeCmd = &cobra.Command{
	Use:   "invoke",
	Short: "Invoke the AWS Bedrock for a single use prompt.",
	Long:  "Invoke the AWS Bedrock for a single use prompt.",
	RunE:  invokeCmdHandler,
}

func init() {
	invokeCmd.PersistentFlags().Bool("stream", false, "Toggles streaming responses")
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
	aiPrompt, _ := ai.New(&ai.NewInvokeDriverInput{
		Client:     bedrockruntime.NewFromConfig(cfg),
		ModelId:    utils.Jeeves.ConfigSettings.AI.PreferredModel,
		TokenCount: 1024,
	})

	output, err := aiPrompt.Invoke(input)
	fmt.Println(output)

	// err = aiPrompt.InvokeStream(input, func(ctx context.Context, part string) error {
	// 	fmt.Print(part)
	// 	return nil
	// })
	if err != nil {
		return err
	}

	fmt.Print("\n---END---\n\n")

	return invoke(cfg, ctx)
}

func invokeCmdHandler(cmd *cobra.Command, args []string) error {
	isLoggedIn, _ := utils.CheckAWSLogin()
	if !isLoggedIn {
		return utils.ErrNotLoggedIn
	}

	ctx := context.Background()
	loader := config.AWSConfigLoader{}
	cfg, err := loader.LoadAWSConfig("default")
	if err != nil {
		return err
	}

	return invoke(cfg, ctx)
}
