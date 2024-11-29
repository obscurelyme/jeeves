package ai

import "github.com/spf13/cobra"

var converseCmd = &cobra.Command{
	Use:   "converse",
	Short: "Converse with AWS Bedrock.",
	Long:  "Converse with AWS Bedrock for multi chat prompts.",
	RunE:  invokeCmdHandler,
}

func converseCmdHandler(cmd *cobra.Command, args []string) error {
	return nil
}
