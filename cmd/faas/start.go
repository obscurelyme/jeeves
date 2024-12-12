package faas

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	lambdaTypes "github.com/aws/aws-sdk-go-v2/service/lambda/types"
	jeevesEnv "github.com/obscurelyme/jeeves/env"
	"github.com/obscurelyme/jeeves/templates"
	"github.com/obscurelyme/jeeves/templates/scripts/python"
	"github.com/obscurelyme/jeeves/utils"
	"github.com/obscurelyme/jeeves/utils/java"
	pythonUtils "github.com/obscurelyme/jeeves/utils/python"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const FAAS_CONFIG_FILE string = "faas.yaml"
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
	faasConfig, err := ReadLambdaConfig()
	if err != nil {
		return err
	}

	faasRuntime := faasConfig.GetString("function.runtime")
	faasHandler := faasConfig.GetString("function.handler")
	isLoggedIn, _ := CheckAWSLogin()
	if !isLoggedIn {
		return errors.New("you need to login into AWS first, please run \"jeeves login\" then retry")
	}

	err = initializeDockerFiles(faasRuntime, faasHandler)
	if err != nil {
		return err
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

func initializeDockerFiles(faasRuntime string, faasHandler string) error {
	isConfigured := checkFaasDockerConfig()

	if !isConfigured {
		err := writeDockerfile(faasRuntime, faasHandler)
		if err != nil {
			return err
		}

		err = writeComposeFile(faasRuntime)
		if err != nil {
			return err
		}
	} else {
		fmt.Println("Dockerfile and docker-compose.yaml already written to, will not overwrite")
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

func writeDockerfile(faasRuntime string, faasHandler string) error {
	var dockerFile templates.DockerFileWriter
	var venv *pythonUtils.PythonVirtualEnv = nil
	var mvnPomDriver *java.MavenPomFileDriver = nil
	var mvnErr error = nil

	if strings.Contains(faasRuntime, "python") {
		venv = pythonUtils.NewPythonVirtualEnv()
		// Ensure we have an active venv
		// NOTE: This could most likely be removed
		if !pythonUtils.VirtualEnvActive() {
			return errors.New("no python venv is active")
		}
		// NOTE: write the bootstrap file
		script := python.New(ConfigPath, lambdaTypes.Runtime(faasRuntime))
		err := script.WriteFile()
		if err != nil {
			return err
		}
	}

	if strings.Contains(faasRuntime, "java") {
		fmt.Println("JAVA DETECTEED")
		mvnPomDriver, mvnErr = java.New(ConfigPath)
		if mvnErr != nil {
			return mvnErr
		}
	}

	dockerFile, err := templates.NewDockerFile(&templates.NewDockerFileInput{
		Runtime:        faasRuntime,
		Handler:        faasHandler,
		FilePath:       ConfigPath,
		VirtualEnv:     venv,
		MavenPomDriver: mvnPomDriver,
	})

	if err != nil {
		return err
	}

	return dockerFile.WriteFile()
}

func writeComposeFile(faasRuntime string) error {
	// TODO: need to work with the runtime param to determine what kind of compose file the user gets
	composeFilePath := fmt.Sprintf("%s/docker-compose.yaml", ConfigPath)
	return os.WriteFile(composeFilePath, []byte(COMPOSE_TEMPLATE), 0644)
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
