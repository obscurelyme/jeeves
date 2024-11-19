package faas

import "github.com/spf13/cobra"

var FaasRootCmd = &cobra.Command{
	Use:   "faas",
	Short: "Execute FaaS commands",
	Long:  "Configure, Create, and Execute Functions as a Service (AWS Lambdas)",
	Run:   faasHandler,
}

func init() {
	FaasRootCmd.AddCommand(listFaasCmd)
}

func faasHandler(cmd *cobra.Command, args []string) {
	return
}
