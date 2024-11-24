package faas

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/lambda/types"
	"github.com/manifoldco/promptui"
	"github.com/obscurelyme/jeeves/config"
	"github.com/spf13/cobra"
)

// Name of the S3 bucket which holds all example lambda zips
var S3_BUCKET_NAME string = "example-lambda-apps"

// Owner of the template repos
var TEMPLATE_REPO_OWNER string = "obscurelyme"

// Basic execution policy for Lambda Functions.
var LAMBDA_BASIC_EXECUTION_ROLE string = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"

var CREATE_LAMBDA_REPOSITORY string = "create-lambda-repository"

var TRUST_POLICY_DOC string = `{
	"Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "Service": "lambda.amazonaws.com"
      },
      "Action": "sts:AssumeRole"
    }
  ]
}`

var BASIC_LAMBDA_POLICY_ARN = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"

var createFaasCmd = &cobra.Command{
	Use:   "create",
	Short: "Create and provison new FaaS resources",
	Long:  "Opens a prompt to create and provision brand new FaaS functions",
	RunE:  createFassCmdHandler,
}

var promptTemplate = &promptui.PromptTemplates{
	Prompt:  "{{ . }}",
	Valid:   "{{ . | green }}",
	Invalid: "{{ . | red }}",
}

var selectTemplate = &promptui.SelectTemplates{
	Label:    "{{ .Language }}",
	Active:   "{{ .Language | cyan}}",
	Inactive: "{{ .Language }}",
	Selected: "{{ .Language | cyan }}",
	Details: `
---------- Function Runtime -----------
Language: {{ .Language }}
Runtime: {{ .AWSRuntime }}`,
}

func createFassCmdHandler(cmd *cobra.Command, args []string) error {
	functionName, err := promptInput()
	if err != nil {
		return err
	}

	runtimeSelection, err := promptLambdaRuntimeSelect()
	if err != nil {
		return err
	}

	confirmed, err := promptConfirm(functionName)
	if err != nil {
		return err
	}

	if confirmed == "y" {
		fmt.Printf("Creating FaaS resource (%s)...\n", functionName)
	} else {
		fmt.Printf("Cancelling creation of FaaS resource (%s)...\n", functionName)
	}

	input := CreateFaaSResourceInput{
		FunctionName: functionName,
		Runtime:      &runtimeSelection,
	}

	err = ProvisionFaasRepo(input)
	if err != nil {
		return err
	}

	fmt.Printf("Provisioning %s FaaS!\n", input.FunctionName)
	err = CreateFaasResource(input)
	if err != nil {
		return err
	}

	return nil
}

func promptInput() (string, error) {
	prompt := promptui.Prompt{
		Label:     "Function Name: ",
		Templates: promptTemplate,
		Validate:  ValidateFunctionName,
	}

	return prompt.Run()
}

func promptLambdaRuntimeSelect() (LambdaRuntime, error) {
	prompt := promptui.Select{
		Label:     "Function Runtime: ",
		Templates: selectTemplate,
		Items:     runtimeSelection,
	}

	index, _, err := prompt.Run()

	return runtimeSelection[index], err
}

func promptConfirm(name string) (string, error) {
	fmt.Printf("You are about create this new resource.\nFunction Name: %s\n", name)

	prompt := promptui.Prompt{
		Label:     "Are you sure?",
		IsConfirm: true,
	}

	return prompt.Run()
}

func ProvisionFaasRepo(input CreateFaaSResourceInput) error {
	loader := &config.AWSConfigLoader{}
	cfg, err := loader.LoadAWSConfig(profile)

	if err != nil {
		return err
	}

	client := lambda.NewFromConfig(cfg)
	creds, _ := cfg.Credentials.Retrieve(context.TODO())
	functionName := fmt.Sprintf("arn:aws:lambda:%s:%s:function:%s", cfg.Region, creds.AccountID, CREATE_LAMBDA_REPOSITORY)
	data, _ := json.Marshal(&Payload{
		TemplateRepo:          input.Runtime.TemplateRepo,
		TemplateOwner:         TEMPLATE_REPO_OWNER,
		Owner:                 TEMPLATE_REPO_OWNER,
		RepositoryName:        fmt.Sprintf("%v.lambda", input.FunctionName),
		RepositoryDescription: "",
		Visibility:            "public",
	})

	output, err := client.Invoke(context.TODO(), &lambda.InvokeInput{
		FunctionName: &functionName,
		Payload:      data,
	})

	if err != nil {
		return err
	}

	if output.StatusCode == http.StatusOK {
		fmt.Printf("Github repository for %s created!\n", input.FunctionName)
		return nil
	}

	return fmt.Errorf("lambda failed with status code: %d", output.StatusCode)
}

// Creates the Lambda Function, this function can take some time due to having
// to wait for the AWS Policies and Role to take effect so that the lambda
// may assume the role.
func CreateFaasResource(input CreateFaaSResourceInput) error {
	loader := &config.AWSConfigLoader{}
	cfg, err := loader.LoadAWSConfig(profile)

	if err != nil {
		return err
	}
	client := lambda.NewFromConfig(cfg)

	var defaultTimeout int32 = 30
	var functionCode = types.FunctionCode{
		S3Bucket: &S3_BUCKET_NAME,
		S3Key:    &input.Runtime.Example,
	}

	roleArn, _, roleErr := CreateLambdaRole(&input)
	if roleErr != nil {
		return roleErr
	}

	retry := 3
	done := false
	fmt.Println("This may take some time please be patient\nIAM roles and Policies can take up to 30 seconds\nbefore taking effect...")
	for retry > -1 && !done {
		done, _ = TryMakeFaaSResource(client, functionCode, defaultTimeout, roleArn, &input)
		if !done {
			fmt.Println("...")
			time.Sleep(5 * time.Second)
		} else {
			fmt.Println("...Complete!")
		}
	}

	return nil
}

func ValidateFunctionName(input string) error {
	if len(input) <= 0 {
		return errors.New("function name is required")
	}
	if strings.ContainsAny(input, " ,.\"'[]{}/:;_=+`~&^%$#@!*()\\?<>") {
		return errors.New("spaces and special characters other than \"-\" are not allowed")
	}
	return nil
}

func CreateLambdaRole(input *CreateFaaSResourceInput) (string, string, error) {
	loader := config.AWSConfigLoader{}
	cfg, err := loader.LoadAWSConfig(profile)
	if err != nil {
		return "", "", err
	}
	iamClient := iam.NewFromConfig(cfg)

	roleName := fmt.Sprintf("%s-IamRole", input.FunctionName)
	roleOutput, roleErr := iamClient.CreateRole(context.TODO(), &iam.CreateRoleInput{
		AssumeRolePolicyDocument: &TRUST_POLICY_DOC,
		RoleName:                 &roleName,
	})

	if roleErr != nil {
		return "", "", err
	}

	_, policyErr := iamClient.AttachRolePolicy(context.TODO(), &iam.AttachRolePolicyInput{
		PolicyArn: &BASIC_LAMBDA_POLICY_ARN,
		RoleName:  roleOutput.Role.RoleName,
	})

	if policyErr != nil {
		return "", "", policyErr
	}

	return *roleOutput.Role.Arn, *roleOutput.Role.RoleName, nil
}

func TryMakeFaaSResource(client *lambda.Client, functionCode types.FunctionCode, defaultTimeout int32, roleArn string, input *CreateFaaSResourceInput) (bool, error) {
	_, err := client.CreateFunction(context.TODO(), &lambda.CreateFunctionInput{
		Code:          &functionCode,
		FunctionName:  &input.FunctionName,
		Role:          &roleArn,
		Architectures: []types.Architecture{types.ArchitectureArm64},
		// Description: ""
		Runtime: input.Runtime.AWSRuntime,
		Timeout: &defaultTimeout,
		Handler: &input.Runtime.Handler,
	})

	if err != nil {
		var conflict *types.InvalidParameterValueException
		if errors.As(err, &conflict) {
			return false, err
		}
		return false, err
	}

	return true, nil
}
