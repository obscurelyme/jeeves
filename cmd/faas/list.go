package faas

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/obscurelyme/jeeves/config"
	"github.com/spf13/cobra"
)

var profile string
var listFaasCmd = &cobra.Command{
	Use:   "list",
	Short: "List available FaaS resources",
	Long:  "List the available FaaS resources (Lambdas) provisioned within AWS",
	RunE:  listFaasCmdHandler,
}

func init() {
	listFaasCmd.PersistentFlags().StringVar(&profile, "profile", "default", "AWS Profile to login to")
}

func listFaasCmdHandler(cmd *cobra.Command, args []string) error {
	return ListLambdas()
}

func ListLambdas() error {
	loginCfg := config.AWSConfigLoader{}
	cfg, err := loginCfg.LoadAWSConfig(profile)

	if err != nil {
		return err
	}

	lambdaClient := lambda.NewFromConfig(cfg)
	output, err := lambdaClient.ListFunctions(context.TODO(), &lambda.ListFunctionsInput{})

	if err != nil {
		return err
	}

	for _, lambdaFunction := range output.Functions {
		fmt.Printf(
			"---\nFunction: %s\nRuntime: %v\nVersion: %s\n---\n",
			*lambdaFunction.FunctionName,
			lambdaFunction.Runtime,
			*lambdaFunction.Version,
		)
	}

	return nil
}
