package ai

import (
	"fmt"

	"github.com/spf13/cobra"
)

var converseCmd = &cobra.Command{
	Use:   "converse",
	Short: "Converse with AWS Bedrock.",
	Long:  "Converse with AWS Bedrock for multi chat prompts.",
	RunE:  converseCmdHandler,
}

func converseCmdHandler(cmd *cobra.Command, args []string) error {
	modelId, err := cmd.PersistentFlags().GetString("model-identifier")
	if err != nil {
		return err
	}

	fmt.Println(modelId)
	return nil
}
