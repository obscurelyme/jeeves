package faas

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

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

	return CreateFaasResource(input)
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

	_, err = client.CreateFunction(context.TODO(), &lambda.CreateFunctionInput{
		Code:         &functionCode,
		FunctionName: &input.FunctionName,
		// Role: "", // TODO: each function will either use an existing role or make a new one
		Architectures: []types.Architecture{types.ArchitectureArm64},
		// Description: ""
		Runtime: input.Runtime.AWSRuntime,
		Timeout: &defaultTimeout,
		Handler: &input.Runtime.Handler,
	})

	if err != nil {
		return err
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
