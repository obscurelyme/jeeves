package python

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func VirtualEnvActive() bool {
	return os.Getenv("VIRTUAL_ENV") != ""
}

func VirtualEnv() string {
	return os.Getenv("VIRTUAL_ENV")
}

func VirtualEnvPath() (string, error) {
	venvFullPath := VirtualEnv()
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	if venvFullPath == "" {
		return "", errors.New("no python venv is active")
	}
	strs := strings.Split(venvFullPath, dir+"/")
	if len(strs) < 2 {
		return "", errors.New("venv is not located within the project directory")
	}
	return strs[1], nil
}

func PythonVersion() (string, error) {
	cmd := exec.Command("python", "-V")

	data, err := cmd.Output()

	if err != nil {
		return "", err
	}

	return formatPythonVersion(string(data)), nil
}

func formatPythonVersion(version string) string {
	split := strings.Split(version, " ")
	versionNumber := split[1]
	nums := strings.Split(versionNumber, ".")

	return fmt.Sprintf("python%s.%s", nums[0], nums[1])
}

// Returns the path of the python dependencies based on the currently active python venv
func PythonDependenciesPath() (string, error) {
	venvPath, err := VirtualEnvPath()
	if err != nil {
		return "", err
	}

	version, err := PythonVersion()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s/lib/%s/site-packages", venvPath, version), nil
}
