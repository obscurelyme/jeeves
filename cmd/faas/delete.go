package faas

import "github.com/spf13/cobra"

var deleteFaasCmd = &cobra.Command{
	Use:   "delete",
	Short: "Deletes an existing FaaS resource",
	Long:  "Opens a prompt to delete an FaaS resource and its corresponding IAM roles",
	RunE:  deleteFassCmdHandler,
}

func deleteFassCmdHandler(cmd *cobra.Command, args []string) error {
	return nil
}
