package faas

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/lambda/types"
	"github.com/manifoldco/promptui"
	"github.com/obscurelyme/jeeves/config"
	"github.com/spf13/cobra"
)

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
	Label:    "{{ . }}",
	Selected: "{{ . | cyan }}",
}

func createFassCmdHandler(cmd *cobra.Command, args []string) error {
	functionName, err := promptInput()
	if err != nil {
		return err
	}

	_, runtimeName, err := promptLambdaRuntimeSelect()
	if err != nil {
		return err
	}

	confirmed, err := promptConfirm(functionName)
	if err != nil {
		return err
	}

	if confirmed == "y" {
		fmt.Printf("Creating FaaS resource (%s) with %s runtime...\n", functionName, runtimeName)
	} else {
		fmt.Printf("Cancelling creation of FaaS resource (%s) with %s runtime...\n", functionName, runtimeName)
	}

	return nil
	// return CreateFaasResource()
}

func promptInput() (string, error) {
	prompt := promptui.Prompt{
		Label:     "Function Name: ",
		Templates: promptTemplate,
		Validate:  ValidateFunctionName,
	}

	return prompt.Run()
}

func promptLambdaRuntimeSelect() (int, string, error) {
	items := types.Runtime.Values(types.RuntimeGo1x)

	prompt := promptui.Select{
		Label:     "Function Runtime: ",
		Templates: selectTemplate,
		Items:     items,
	}

	return prompt.Run()
}

func promptConfirm(name string) (string, error) {
	fmt.Printf("You are about create this new resource.\nFunction Name: %s\n", name)

	prompt := promptui.Prompt{
		Label:     "Are you sure?",
		IsConfirm: true,
	}

	return prompt.Run()
}

func CreateFaasResource() error {
	loader := &config.AWSConfigLoader{}
	cfg, err := loader.LoadAWSConfig(profile)

	if err != nil {
		return err
	}
	client := lambda.NewFromConfig(cfg)

	var defaultTimeout int32 = 30
	s3BucketName := "example-lambda-apps"
	s3Key := "nodejs-lambda.zip"
	var functionCode = types.FunctionCode{
		S3Bucket: &s3BucketName,
		S3Key:    &s3Key,
	}

	_, err = client.CreateFunction(context.TODO(), &lambda.CreateFunctionInput{
		Code: &functionCode,
		// FunctionName: "",
		// Role: "", // TODO: each function will either use an existing role or make a new one
		Architectures: []types.Architecture{types.ArchitectureArm64},
		// Description: ""
		// Runtime: "", // TODO: based on the user defined template to make
		Timeout: &defaultTimeout,
		// Handler: "", // TODO: based on the user defined template to make
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
