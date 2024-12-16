package ai

import (
	"fmt"

	"github.com/manifoldco/promptui"
	"github.com/obscurelyme/jeeves/ai/types"
	"github.com/obscurelyme/jeeves/prompt"
	"github.com/obscurelyme/jeeves/utils"
	"github.com/spf13/cobra"
)

var modelCmd = &cobra.Command{
	Use:   "set-model",
	Short: "Sets the model-id for all \"ai\" commands",
	Long:  "Sets the model-id for all \"jeeves ai\" commands",
	RunE:  modelCmdHandler,
}

func modelCmdHandler(cmd *cobra.Command, args []string) error {
	options := []types.ModelSelect{}

	for _, modelId := range utils.Jeeves.ConfigSettings.AI.ApprovedModels {
		options = append(options, types.ModelSelect{
			Label:   modelId,
			ModelId: modelId,
		})
	}

	selectedModel, err := prompt.SelectPrompt("Select Model", options, &promptui.SelectTemplates{
		Label:    "{{ .Label }}",
		Active:   "{{ .Label | cyan}}",
		Inactive: "{{ .Label }}",
		Selected: "{{ .Label | cyan }}",
		Details: fmt.Sprintf(`
--------------------------------
Current: %s`, utils.Jeeves.ConfigSettings.AI.PreferredModel),
	})
	if err != nil {
		return err
	}

	utils.Jeeves.ConfigSettings.AI.PreferredModel = selectedModel.ModelId
	err = utils.Jeeves.WriteConfig()
	if err != nil {
		return err
	}

	return nil
}
