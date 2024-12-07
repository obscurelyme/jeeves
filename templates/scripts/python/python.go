package python

import (
	_ "embed"
	"fmt"
	"os"

	lambdaTypes "github.com/aws/aws-sdk-go-v2/service/lambda/types"
)

//go:embed bootstrap.sh
var boostrapScript string

type PythonBootstrapScript struct {
	boostrapScript string
	configPath     string
}

func New(configPath string, runtime lambdaTypes.Runtime) *PythonBootstrapScript {
	script := new(PythonBootstrapScript)

	script.configPath = configPath
	script.boostrapScript = fmt.Sprintf(boostrapScript, runtime, runtime, runtime)

	return script
}

func (pbs *PythonBootstrapScript) WriteFile() error {
	return os.WriteFile(fmt.Sprintf("%s/bootstrap.sh", pbs.configPath), []byte(pbs.boostrapScript), 0755)
}
