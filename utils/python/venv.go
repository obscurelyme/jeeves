package python

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

const VIRTUAL_ENV string = "VIRTUAL_ENV"
const VIRTUAL_ENV_PROMPT string = "VIRTUAL_ENV_PROMPT"

type PythonVirtualEnvDriver interface {
	// The full filepath to the virtual environment
	Path() string
	// The name of the virtual environment
	Name() (string, error)
	// Checks if the current working directory contains the active virtual environment
	CwdContainsVenv() error
	// Returns a formatted version of the python bin used for this virtual environment
	//
	// format style examples: python3.9, python3.10, python3.11, python3.12
	PythonVersion() (string, error)
	// The path to the python dependencies for the active virtual environment
	DependencyPath() (string, error)
}

type PythonVirtualEnv struct {
}

func (venv *PythonVirtualEnv) Path() string {
	return os.Getenv(VIRTUAL_ENV)
}

func (venv *PythonVirtualEnv) Name() (string, error) {
	venvName := os.Getenv(VIRTUAL_ENV_PROMPT)

	if venvName == "" {
		return "", errors.New("no python venv is active")
	}

	venvName = strings.ReplaceAll(venvName, "(", "")
	venvName = strings.ReplaceAll(venvName, ")", "")
	venvName = strings.TrimSpace(venvName)

	return venvName, nil
}

func (venv *PythonVirtualEnv) CwdContainsVenv() error {
	path := venv.Path()
	name, err := venv.Name()
	if err != nil {
		return err
	}

	venvCwd, _ := strings.CutSuffix(path, name)

	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	if venvCwd != cwd+"/" {
		return errors.New("python venv is not within your cwd")
	}

	return nil
}

func (venv *PythonVirtualEnv) PythonVersion() (string, error) {
	cmd := exec.Command("python", "-V")

	data, err := cmd.Output()

	if err != nil {
		return "", err
	}

	return formatPythonVersion(string(data)), nil
}

func (venv *PythonVirtualEnv) DependencyPath() (string, error) {
	venvName, err := venv.Name()
	if err != nil {
		return "", err
	}

	version, err := venv.PythonVersion()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s/lib/%s/site-packages", venvName, version), nil
}

func NewPythonVirtualEnv() *PythonVirtualEnv {
	venv := new(PythonVirtualEnv)
	return venv
}

func VirtualEnvActive() bool {
	return os.Getenv("VIRTUAL_ENV") != ""
}

func formatPythonVersion(version string) string {
	split := strings.Split(version, " ")
	versionNumber := split[1]
	nums := strings.Split(versionNumber, ".")

	return fmt.Sprintf("python%s.%s", nums[0], nums[1])
}
