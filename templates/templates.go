package templates

import (
	_ "embed"
	"errors"
	"fmt"

	lambdaTypes "github.com/aws/aws-sdk-go-v2/service/lambda/types"
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

//go:embed scripts/bootstrap.python.sh
var bootstrapPython string

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
