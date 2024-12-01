package ai

import (
	"github.com/spf13/cobra"
)

var modelCmd = &cobra.Command{
	Use:   "set-model",
	Short: "Sets the model-id for all \"ai\" commands",
	Long:  "Sets the model-id for all \"jeeves ai\" commands",
	RunE:  modelCmdHandler,
}

func modelCmdHandler(cmd *cobra.Command, args []string) error {
	// prompt.SelectPrompt()

	return nil
}
