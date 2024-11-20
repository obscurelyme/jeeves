package faas

import "github.com/spf13/cobra"

var FaasRootCmd = &cobra.Command{
	Use:   "faas",
	Short: "Execute FaaS commands",
	Long:  "Configure, Create, and Execute Functions as a Service (AWS Lambdas)",
}

func init() {
	FaasRootCmd.AddCommand(listFaasCmd)
	FaasRootCmd.AddCommand(createFaasCmd)
}
