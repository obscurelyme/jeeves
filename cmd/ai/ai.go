package ai

import "github.com/spf13/cobra"

var AIRootCmd = &cobra.Command{
	Use:   "ai",
	Short: "Execute AI commands that work with AWS Bedrock",
	Long:  "Work with AWS AI tools such as Bedrock",
}

func init() {
	// NOTE: attach AI subcommands
}
