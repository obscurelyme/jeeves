package faas

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/spf13/cobra"
)

var listFaasCmd = &cobra.Command{
	Use:   "list",
	Short: "List available FaaS resources",
	Long:  "List the available FaaS resources (Lambdas) provisioned within AWS",
	RunE:  listFaasCmdHandler,
}

func listFaasCmdHandler(cmd *cobra.Command, args []string) error {
	return ListLambdas()
}

func ListLambdas() error {
	cfg := aws.Config{}
	lambda.NewFromConfig(cfg)
	return nil
}
