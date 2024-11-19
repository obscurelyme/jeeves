package faas

import "github.com/spf13/cobra"

var listFaasCmd = &cobra.Command{
	Use:   "list",
	Short: "List available FaaS resources",
	Long:  "List the available FaaS resources (Lambdas) provisioned within AWS",
	Run:   listFaasCmdHandler,
}

func listFaasCmdHandler(cmd *cobra.Command, args []string) {
	return
}
