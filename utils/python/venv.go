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
	if !strings.Contains(venvFullPath, dir+"/") {
		return "", errors.New("cwd is not within your venv")
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
	venvName, err := VirtualEnvName()
	if err != nil {
		return "", err
	}

	version, err := PythonVersion()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s/lib/%s/site-packages", venvName, version), nil
}

// Checks if the current working directory contains the active python venv
func CwdIsVenv() error {
	venvPath := VirtualEnv()
	venvName := os.Getenv("VIRTUAL_ENV_PROMPT")

	if venvPath == "" || venvName == "" {
		return errors.New("no python venv is active")
	}

	// NOTE: remove ".tox" if it exists because tox will place the venv in its own directory.
	// which the user does not need to cd into
	venvPath = strings.Replace(venvPath, ".tox/", "", 1)

	venvName = strings.ReplaceAll(venvName, "(", "")
	venvName = strings.ReplaceAll(venvName, ")", "")
	venvName = strings.TrimSpace(venvName)

	venvCwd, _ := strings.CutSuffix(venvPath, venvName)

	dir, err := os.Getwd()
	if err != nil {
		return err
	}

	if venvCwd != dir+"/" {
		return errors.New("python venv is not within your cwd")
	}

	return nil
}

func VirtualEnvName() (string, error) {
	venvName := os.Getenv("VIRTUAL_ENV_PROMPT")

	if venvName == "" {
		return "", errors.New("no python venv is active")
	}

	venvName = strings.ReplaceAll(venvName, "(", "")
	venvName = strings.ReplaceAll(venvName, ")", "")
	venvName = strings.TrimSpace(venvName)

	return venvName, nil
}
