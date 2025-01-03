package faas

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/obscurelyme/jeeves/config"
	"github.com/obscurelyme/jeeves/types"
	"github.com/spf13/cobra"
)

var resourceName string
var deleteFaasCmd = &cobra.Command{
	Use:   "delete",
	Short: "Deletes an existing FaaS resource",
	Long:  "Opens a prompt to delete an FaaS resource and its corresponding IAM roles",
	RunE:  deleteFassCmdHandler,
}

func init() {
	deleteFaasCmd.PersistentFlags().StringVar(&profile, "profile", "default", "AWS Profile to work with")
	deleteFaasCmd.PersistentFlags().StringVar(&resourceName, "resource-name", "", "Name of the FaaS resource to delete. (required)")
}

func deleteFassCmdHandler(cmd *cobra.Command, args []string) error {
	if resourceName == "" {
		return errors.New("resource name is a required flag --resource-name [RESOURCE_TO_DELETE]")
	}

	// NOTE: get the config
	loader := config.AWSConfigLoader{}

	cfg, err := loader.LoadAWSConfig(profile)
	if err != nil {
		return err
	}

	// Detach policies
	err = DetachFaaSPolicies(cfg)
	if err != nil {
		return err
	}

	// Run the delete FaaS resource
	err = DeleteFaaSResource(cfg)
	if err != nil {
		return err
	}

	// Delete the IAM role associated with the FaaS resource
	err = DeleteFaaSResourceRole(cfg)

	if err == nil {
		fmt.Printf("FaaS resource, %s, was successfully deleted!\n", resourceName)
	}

	err = DeleteFaaSRepo(cfg, resourceName)

	return err
}

func DeleteFaaSRepo(cfg aws.Config, name string) error {
	lambdaClient := lambda.NewFromConfig(cfg)
	functionName := "delete-lambda-repository"
	payload, err := json.Marshal(&types.DeleteRepositoryPayload{
		RepositoryOwner: "obscurelyme",
		RepositoryName:  fmt.Sprintf("%s.lambda", name),
	})

	if err != nil {
		return nil
	}

	output, err := lambdaClient.Invoke(context.TODO(), &lambda.InvokeInput{
		FunctionName: &functionName,
		Payload:      payload,
	})

	if err != nil {
		return err
	}

	if output.StatusCode == http.StatusOK {
		fmt.Printf("Github repository for %s deleted!\n", resourceName)
		return nil
	}

	return fmt.Errorf("lambda failed with status code: %d", output.StatusCode)
}

func DeleteFaaSResource(cfg aws.Config) error {
	lambdaClient := lambda.NewFromConfig(cfg)

	_, err := lambdaClient.DeleteFunction(context.TODO(), &lambda.DeleteFunctionInput{
		FunctionName: &resourceName,
	})

	return err
}

func DeleteFaaSResourceRole(cfg aws.Config) error {
	iamClient := iam.NewFromConfig(cfg)

	roleName := fmt.Sprintf("%s-IamRole", resourceName)

	_, err := iamClient.DeleteRole(context.TODO(), &iam.DeleteRoleInput{
		RoleName: &roleName,
	})

	return err
}

func DetachFaaSPolicies(cfg aws.Config) error {
	iamClient := iam.NewFromConfig(cfg)

	roleName := fmt.Sprintf("%s-IamRole", resourceName)

	_, err := iamClient.DetachRolePolicy(context.TODO(), &iam.DetachRolePolicyInput{
		PolicyArn: &BASIC_LAMBDA_POLICY_ARN,
		RoleName:  &roleName,
	})

	return err
}
