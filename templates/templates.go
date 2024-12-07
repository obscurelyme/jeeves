package templates

import (
	_ "embed"
	"errors"
	"fmt"

	lambdaTypes "github.com/aws/aws-sdk-go-v2/service/lambda/types"
	pythonUtils "github.com/obscurelyme/jeeves/utils/python"
	"github.com/spf13/viper"
)

//go:embed files/dockerfile.nodejs.template
var dockerFileNodeJS string

//go:embed files/dockerfile.python.template
var dockerFilePython string

//go:embed files/dockerfile.java.template
var dockerFileJava string

//go:embed files/dockerfile.go.template
var dockerFileGo string

func GetDockerTemplate(runtime lambdaTypes.Runtime) (string, error) {
	switch runtime {
	case lambdaTypes.RuntimeNodejs20x:
		return dockerFileNodeJS, nil
	case lambdaTypes.RuntimeNodejs18x:
		return dockerFileNodeJS, nil
	case lambdaTypes.RuntimeJava21:
		return dockerFileJava, nil
	case lambdaTypes.RuntimeJava17:
		return dockerFileJava, nil
	case lambdaTypes.RuntimeProvidedal2023:
		return dockerFileGo, nil
	case lambdaTypes.RuntimePython39:
		return dockerFilePython, nil
	case lambdaTypes.RuntimePython310:
		return dockerFilePython, nil
	case lambdaTypes.RuntimePython311:
		return dockerFilePython, nil
	case lambdaTypes.RuntimePython312:
		return dockerFilePython, nil
	}

	return "", errors.New("no dockerfile template supports the provided runtime")
}

func GetDockerfile(runtime lambdaTypes.Runtime, image string, tag string, handler string) (string, error) {
	template, err := GetDockerTemplate(runtime)
	if err != nil {
		return "", err
	}

	if runtime == lambdaTypes.RuntimeProvidedal2023 || runtime == lambdaTypes.RuntimeProvidedal2 {
		// NOTE: Golang requires no handler to be specified, executable named "bootstrap" is mandatory
		return fmt.Sprintf(template, image, tag), nil
	}

	return fmt.Sprintf(template, image, tag, handler), nil
}

type ComposeTemplate struct {
	Cfg        *viper.Viper
	ConfigPath string
}

func (com *ComposeTemplate) ReadInConfig() error {
	return nil
}

func (com *ComposeTemplate) WriteConfig() error {
	return nil
}

func NewComposeTemplate(configPath string) *ComposeTemplate {
	template := new(ComposeTemplate)

	template.ConfigPath = configPath
	template.Cfg = viper.New()
	template.Cfg.AddConfigPath(configPath)
	template.Cfg.SetConfigName("compose")
	template.Cfg.SetConfigType("yaml")

	return template
}

type DockerFileWriter interface {
	WriteFile() error
}

type DockerFile struct {
	dockerFile string
}

func (df *DockerFile) WriteFile() error {
	return nil
}

// Creates a new JavaDockerFile writer ready to write a properly formatted Dockerfile for Java lambdas
func NewJavaDockerFile(runtime lambdaTypes.Runtime) (DockerFileWriter, error) {
	jdf := new(DockerFile)

	jdf.dockerFile = dockerFileJava

	return jdf, nil
}

// Creates a new PythonDockerFile writer ready to write a properly formatted Dockerfile for Python lambdas
func NewPythonDockerFile(runtime lambdaTypes.Runtime) (DockerFileWriter, error) {
	pdf := new(DockerFile)

	depsPath, err := pythonUtils.PythonDependenciesPath()
	if err != nil {
		return nil, err
	}

	pdf.dockerFile = fmt.Sprintf(dockerFilePython, "image", "tag", depsPath, "handler")

	return pdf, nil
}

// Creates a new GoDockerFile writer ready to write a properly formatted Dockerfile for Go lambdas
func NewGoDockerFile(runtime lambdaTypes.Runtime) (DockerFileWriter, error) {
	return nil, nil
}

// Creates a new NodeJSDockerFile writer ready to write a properly formatted Dockerfile for NodeJS lambdas
func NewNodeJSDockerFile(runtime lambdaTypes.Runtime) (DockerFileWriter, error) {
	return nil, nil
}
