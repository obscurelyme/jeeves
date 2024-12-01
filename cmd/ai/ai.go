package ai

import (
	"fmt"
	"slices"

	"github.com/obscurelyme/jeeves/utils"
	"github.com/spf13/cobra"
)

var AIRootCmd = &cobra.Command{
	Use:   "ai",
	Short: "Execute AI commands that work with AWS Bedrock",
	Long:  "Work with AWS AI tools such as Bedrock",
	RunE:  aiRootCmdHandler,
}

func init() {
	AIRootCmd.AddCommand(invokeCmd)
	AIRootCmd.AddCommand(converseCmd)
	AIRootCmd.AddCommand(modelCmd)

	AIRootCmd.PersistentFlags().String("model-id", "", "Sets the model-id for all \"jeeves ai\" commands")
}

func aiRootCmdHandler(cmd *cobra.Command, args []string) error {
	modelId, err := cmd.PersistentFlags().GetString("model-id")
	if err != nil {
		return err
	}

	if modelId == "" {
		return nil
	}

	if slices.Contains(utils.Jeeves.ConfigSettings.AI.ApprovedModels, modelId) {
		return utils.Jeeves.WriteConfig()
	} else {
		return fmt.Errorf("%s is not an approved valid model", modelId)
	}
}
