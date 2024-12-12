package templates

import (
	_ "embed"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"

	lambdaTypes "github.com/aws/aws-sdk-go-v2/service/lambda/types"
	"github.com/obscurelyme/jeeves/types"
	"github.com/obscurelyme/jeeves/utils/java"
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
	filePath   string
}

func (df *DockerFile) WriteFile() error {
	return os.WriteFile(df.filePath, []byte(df.dockerFile), 0644)
}

type NewDockerFileInput struct {
	// Lambda runtime
	Runtime string
	// Handler of the lambda function
	Handler string
	// Directory location the dockerfile will be written to
	FilePath string
	// Optional: Virtual Environment, used for Python
	VirtualEnv pythonUtils.PythonVirtualEnvDriver
	// Optional: Driver to modify the project's Maven POM file, used for Java
	MavenPomDriver java.MavenPomDriver
}

func NewDockerFile(input *NewDockerFileInput) (DockerFileWriter, error) {
	if strings.Contains(string(input.Runtime), "nodejs") {
		return NewNodeJSDockerFile(input)
	} else if strings.Contains(string(input.Runtime), "python") {
		return NewPythonDockerFile(input)
	} else if strings.Contains(string(input.Runtime), "provided.al2") {
		return NewGoDockerFile(input)
	} else if strings.Contains(string(input.Runtime), "java") {
		return NewJavaDockerFile(input)
	}

	return nil, errors.New("no docker image supports given runtime")
}

// Creates a new JavaDockerFile writer ready to write a properly formatted Dockerfile for Java lambdas
func NewJavaDockerFile(input *NewDockerFileInput) (DockerFileWriter, error) {
	jdf := new(DockerFile)

	if !input.MavenPomDriver.HasRequiredPlugins() {
		input.MavenPomDriver.AddRequiredPlugins()
	}
	err := input.MavenPomDriver.WriteFile()
	if err != nil {
		return nil, err
	}

	tag, err := getDockerImageTag(lambdaTypes.Runtime(input.Runtime))
	if err != nil {
		return nil, err
	}

	jdf.dockerFile = fmt.Sprintf(dockerFileJava, string(types.DockerImageJava), tag, input.Handler)
	jdf.filePath = fmt.Sprintf("%s/Dockerfile", input.FilePath)

	return jdf, nil
}

// Creates a new PythonDockerFile writer ready to write a properly formatted Dockerfile for Python lambdas
func NewPythonDockerFile(input *NewDockerFileInput) (DockerFileWriter, error) {
	pdf := new(DockerFile)

	if input.VirtualEnv == nil {
		return nil, errors.New("creation of python dockerfile requires a venv driver")
	}

	// NOTE: you need to be in the root workspace of the venv for this to work
	err := input.VirtualEnv.CwdContainsVenv()
	if err != nil {
		return nil, err
	}

	depsPath, err := input.VirtualEnv.DependencyPath()
	if err != nil {
		return nil, err
	}

	tag, err := getDockerImageTag(lambdaTypes.Runtime(input.Runtime))
	if err != nil {
		return nil, err
	}

	pdf.filePath = fmt.Sprintf("%s/Dockerfile", input.FilePath)
	pdf.dockerFile = fmt.Sprintf(
		dockerFilePython,
		string(types.DockerImagePython),
		tag,
		depsPath,
		input.Handler,
	)

	return pdf, nil
}

// Creates a new GoDockerFile writer ready to write a properly formatted Dockerfile for Go lambdas
func NewGoDockerFile(input *NewDockerFileInput) (DockerFileWriter, error) {
	gdf := new(DockerFile)

	tag, err := getDockerImageTag(lambdaTypes.Runtime(input.Runtime))
	if err != nil {
		return nil, err
	}

	gdf.dockerFile = fmt.Sprintf(dockerFileGo, string(types.DockerImageGo), tag)
	gdf.filePath = fmt.Sprintf("%s/Dockerfile", input.FilePath)

	return gdf, nil
}

// Creates a new NodeJSDockerFile writer ready to write a properly formatted Dockerfile for NodeJS lambdas
func NewNodeJSDockerFile(input *NewDockerFileInput) (DockerFileWriter, error) {
	ndf := new(DockerFile)

	tag, err := getDockerImageTag(lambdaTypes.Runtime(input.Runtime))
	if err != nil {
		return nil, err
	}

	ndf.filePath = fmt.Sprintf("%s/Dockerfile", input.FilePath)
	ndf.dockerFile = fmt.Sprintf(dockerFileNodeJS, string(types.DockerImageNodeJS), tag, input.Handler)

	return ndf, nil
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
