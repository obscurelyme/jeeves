package faas

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"

	lambdaTypes "github.com/aws/aws-sdk-go-v2/service/lambda/types"
	jeevesEnv "github.com/obscurelyme/jeeves/env"
	"github.com/obscurelyme/jeeves/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const FAAS_CONFIG_FILE string = "faas.yaml"
const DOCKERFILE_TEMPLATE string = `FROM %s:%s

%s

CMD [ "%s" ]`

const NODEJS_DOCKERFILE_COPY string = `COPY node_modules ${LAMBDA_TASK_ROOT}/node_modules
COPY dist ${LAMBDA_TASK_ROOT}/dist
COPY package.json ${LAMBDA_TASK_ROOT}`
const GOLANG_DOCKERFILE_COPY string = `COPY bootstrap ${LAMBDA_TASK_ROOT}`
const JAVA_DOCKERFILE_COPY string = `COPY target/classes ${LAMBDA_TASK_ROOT}
COPY target/dependency/* ${LAMBDA_TASK_ROOT}/lib/`
const PYTHON_DOCKERFILE_COPY string = `COPY app.py ${LAMBDA_TASK_ROOT}`

const COMPOSE_TEMPLATE string = `services:
  lambda:
    build: .
    ports:
      - 9000:8080
    env_file:
      - .env`

var ConfigPath = "."
var startFaasCmd = &cobra.Command{
	Use:   "start",
	Short: "Starts a local FaaS resource",
	Long:  "Starts a FaaS resource locally, using docker",
	RunE:  startFaasCmdHandler,
}

var CheckAWSLogin func() (bool, error)

func init() {
	CheckAWSLogin = utils.CheckAWSLogin
}

func startFaasCmdHandler(cmd *cobra.Command, args []string) error {
	err := initializeDockerFiles()
	if err != nil {
		return err
	}

	isLoggedIn, _ := CheckAWSLogin()
	if !isLoggedIn {
		return errors.New("you need to login into AWS first, please run \"jeeves login\" then retry")
	}

	err = initializeEnvFile()
	if err != nil {
		return err
	}

	return dockerCompose()
}

// Set up the .env file to contain AWS ENV vars
//
// AWS_ACCESS_KEY_ID
//
// AWS_SECRET_ACCESS_KEY
//
// AWS_SESSION_TOKEN
func initializeEnvFile() error {
	envFile, err := jeevesEnv.ReadEnv()
	if err != nil {
		return err
	}

	ssoCreds, err := utils.GetSSOSessionCredentials("default")
	if err != nil {
		return err
	}

	envFile.Set("AWS_ACCESS_KEY_ID", ssoCreds.AccessKeyID)
	envFile.Set("AWS_SECRET_ACCESS_KEY", ssoCreds.SecretAccessKey)
	envFile.Set("AWS_SESSION_TOKEN", ssoCreds.SessionToken)

	return envFile.WriteConfig()
}

func initializeDockerFiles() error {
	isConfigured := checkFaasDockerConfig()

	if !isConfigured {
		err := writeDockerfile()
		if err != nil {
			return err
		}

		err = writeComposeFile()
		if err != nil {
			return err
		}
	}

	return nil
}

// Very naive, just stating the Dockerfile and the compose.yaml file[s] to see if they exist.
func checkFaasDockerConfig() bool {
	dockerFile := fmt.Sprintf("%s/Dockerfile", ConfigPath)
	composeFiles := []string{
		fmt.Sprintf("%s/docker-compose.yaml", ConfigPath),
		fmt.Sprintf("%s/docker-compose.yml", ConfigPath),
		fmt.Sprintf("%s/compose.yaml", ConfigPath),
		fmt.Sprintf("%s/compose.yml", ConfigPath),
	}

	_, err := os.Stat(dockerFile)
	if err != nil {
		return false
	}

	for _, composeFile := range composeFiles {
		_, err := os.Stat(composeFile)
		if err == nil {
			return true
		}
	}

	return false
}

func writeDockerfile() error {
	faasConfig, err := ReadLambdaConfig()
	if err != nil {
		return err
	}

	faasRuntime := faasConfig.GetString("function.runtime")
	faasHandler := faasConfig.GetString("function.handler")

	dockerImage, err := getDockerImage(lambdaTypes.Runtime(faasRuntime))
	if err != nil {
		return err
	}
	dockerImageTag, err := getDockerImageTag(lambdaTypes.Runtime(faasRuntime))
	if err != nil {
		return err
	}
	copyContent, err := getCopyContent(faasRuntime)
	if err != nil {
		return err
	}

	dockerFilePath := fmt.Sprintf("%s/Dockerfile", ConfigPath)
	dockerFile := fmt.Sprintf(DOCKERFILE_TEMPLATE, dockerImage, dockerImageTag, copyContent, faasHandler)
	return os.WriteFile(dockerFilePath, []byte(dockerFile), 0644)
}

func getCopyContent(runtime string) (string, error) {
	if strings.Contains(runtime, "nodejs") {
		return NODEJS_DOCKERFILE_COPY, nil
	} else if strings.Contains(runtime, "java") {
		return JAVA_DOCKERFILE_COPY, nil
	} else if strings.Contains(runtime, "provided.al2") {
		return GOLANG_DOCKERFILE_COPY, nil
	} else if strings.Contains(runtime, "python") {
		return PYTHON_DOCKERFILE_COPY, nil
	}

	return "", errors.New("no dockerfile copy content supports given runtime")
}

func writeComposeFile() error {
	composeFilePath := fmt.Sprintf("%s/docker-compose.yaml", ConfigPath)
	return os.WriteFile(composeFilePath, []byte(COMPOSE_TEMPLATE), 0644)
}

func getDockerImageTag(runtime lambdaTypes.Runtime) (string, error) {
	const latestTag string = "latest"

	if strings.Contains(string(runtime), "nodejs") {
		regx := regexp.MustCompile("[0-9]+")
		tag := regx.FindString(string(runtime))
		return tag, nil
	} else if strings.Contains(string(runtime), "python") {
		return strings.Split(string(runtime), "python")[1], nil
	} else if strings.Contains(string(runtime), "provided.al2") {
		return latestTag, nil
	} else if strings.Contains(string(runtime), "java") {
		return strings.Split(string(runtime), "java")[1], nil
	}

	return "", errors.New("no docker image supports given runtime")
}

func getDockerImage(runtime lambdaTypes.Runtime) (string, error) {
	if strings.Contains(string(runtime), "nodejs") {
		return string(DockerImageNodeJS), nil
	} else if strings.Contains(string(runtime), "python") {
		return string(DockerImagePython), nil
	} else if strings.Contains(string(runtime), "provided.al2") {
		return string(DockerImageGo), nil
	} else if strings.Contains(string(runtime), "java") {
		return string(DockerImageJava), nil
	}

	return "", errors.New("no docker image tag supports given runtime")
}

func ReadLambdaConfig() (*viper.Viper, error) {
	config := viper.New()
	config.AddConfigPath(ConfigPath)
	config.SetConfigName("faas")
	config.SetConfigType("yaml")
	err := config.ReadInConfig()
	return config, err
}

func dockerCompose() error {
	dockerCmd := exec.Command("docker", "compose", "up", "--build")

	dockerCmd.Stdin = os.Stdin
	dockerCmd.Stdout = os.Stdout
	dockerCmd.Stderr = os.Stderr

	return dockerCmd.Run()
}
