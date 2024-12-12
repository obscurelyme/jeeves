package templates

import (
	"fmt"
	"os"
	"testing"
)

func readFile(tmpDir string, filename string) (string, error) {
	filepath := fmt.Sprintf("%s/%s", tmpDir, filename)
	data, err := os.ReadFile(filepath)

	return string(data), err
}

type MockPythonVirtualEnv struct {
	TmpDir            string
	MockPythonVersion string
}

func (p *MockPythonVirtualEnv) Path() string {
	return fmt.Sprintf("%s/venv", p.TmpDir)
}
func (p *MockPythonVirtualEnv) Name() (string, error) {
	return "venv", nil
}
func (p *MockPythonVirtualEnv) CwdContainsVenv() error {
	return nil
}
func (p *MockPythonVirtualEnv) PythonVersion() (string, error) {
	return p.MockPythonVersion, nil
}
func (p *MockPythonVirtualEnv) DependencyPath() (string, error) {
	return fmt.Sprintf("%s/venv/lib/%s/site-packages", p.TmpDir, p.MockPythonVersion), nil
}

const expectedFile string = `FROM amazon/aws-lambda-python:3.12

RUN pip3 install debugpy

# Copy the dependencies from the active venv
COPY %s/venv/lib/python3.12/site-packages ${LAMBDA_TASK_ROOT}
# Copy source code
COPY src/* ${LAMBDA_TASK_ROOT}

# Override the bootstrap script to allow for debugging
COPY bootstrap.sh /var/runtime/bootstrap

CMD [ "main.handler" ]`

func TestWriteDockerFiles(t *testing.T) {
	tmpDir := t.TempDir()

	t.Run("python", func(t *testing.T) {
		dockerFile, err := NewDockerFile(&NewDockerFileInput{
			Runtime:  "python3.12",
			Handler:  "main.handler",
			FilePath: tmpDir,
			VirtualEnv: &MockPythonVirtualEnv{
				TmpDir:            tmpDir,
				MockPythonVersion: "python3.12",
			},
		})

		if err != nil {
			t.Errorf("expected no errors but received, \"%s\"", err.Error())
			return
		}

		err = dockerFile.WriteFile()

		if err != nil {
			t.Errorf("expected no errors but received, \"%s\"", err.Error())
			return
		}

		file, err := readFile(tmpDir, "Dockerfile")
		if err != nil {
			t.Errorf("expected no errors but received, \"%s\"", err.Error())
			return
		}

		if file != fmt.Sprintf(expectedFile, tmpDir) {
			t.Errorf("mismatched file, \n%s\n%s", file, expectedFile)
		}
	})
}
