package faas

import (
	"fmt"

	lambdaTypes "github.com/aws/aws-sdk-go-v2/service/lambda/types"
)

type LambdaLanguage string

const NodeJs LambdaLanguage = "nodejs"
const Golang LambdaLanguage = "golang"
const Java LambdaLanguage = "java"
const Python LambdaLanguage = "python"

type LambdaRuntime struct {
	// Internal AWS runtime for the lambda function
	AWSRuntime lambdaTypes.Runtime
	Language   LambdaLanguage
	// Handler for the lambda function
	Handler string
	// Name of the S3 key (zip file) stored in in the examples S3 bucket
	Example string
	// Template repository from which a new repo will be provisioned from
	TemplateRepo string
}

var runtimeSelection []LambdaRuntime = []LambdaRuntime{
	{
		AWSRuntime:   lambdaTypes.RuntimeNodejs20x,
		Language:     NodeJs,
		Handler:      "dist/index.handler",
		Example:      fmt.Sprintf("%s-function.zip", NodeJs),
		TemplateRepo: fmt.Sprintf("%s-lambda", NodeJs),
	},
	{
		AWSRuntime:   lambdaTypes.RuntimeProvidedal2023,
		Language:     Golang,
		Handler:      "main",
		Example:      fmt.Sprintf("%s-function.zip", Golang),
		TemplateRepo: fmt.Sprintf("%s-lambda", Golang),
	},
	{
		AWSRuntime:   lambdaTypes.RuntimeJava21,
		Language:     Java,
		Handler:      "com.example.app.Handler::handleRequest",
		Example:      fmt.Sprintf("%s-function.jar", Java),
		TemplateRepo: fmt.Sprintf("%s-lambda", Java),
	},
	{
		AWSRuntime:   lambdaTypes.RuntimePython310,
		Language:     Python,
		Handler:      "handler",
		Example:      fmt.Sprintf("%s-function.zip", Python),
		TemplateRepo: fmt.Sprintf("%s-lambda", Python),
	},
}

type CreateFaaSResourceInput struct {
	FunctionName string
	Runtime      *LambdaRuntime
}

// Payload to send when provisioning a new template repository for a new FaaS resource
type Payload struct {
	TemplateRepo          string `json:"templateRepo"`
	TemplateOwner         string `json:"templateOwner"`
	Owner                 string `json:"owner"`
	RepositoryName        string `json:"repositoryName"`
	RepositoryDescription string `json:"repositoryDescription"`
	Visibility            string `json:"visibility"`
}

// Payload to send when deleting a repository after deleting an FaaS resource
type DeleteRepositoryPayload struct {
	RepositoryOwner string `json:"repositoryOwner"`
	RepositoryName  string `json:"repositoryName"`
}

type DockerImage string

const (
	DockerImageNodeJS DockerImage = "amazon/aws-lambda-nodejs"
	DockerImageJava   DockerImage = "amazon/aws-lambda-java"
	DockerImagePython DockerImage = "amazon/aws-lambda-python"
	DockerImageGo     DockerImage = "amazon/aws-lambda-go"
)
