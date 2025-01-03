package faas

import (
	"fmt"
	"os"
	"testing"
)

const nodejsYaml = `function:
  runtime: nodejs20.x
  handler: dist/index.js`
const pythonYaml = `function:
  runtime: python3.12
  handler: pkg.main.handler`
const javaYaml = `function:
  runtime: java21
  handler: com.example.app.Handler::handleRequest`
const goYaml = `function:
  runtime: provided.al2023
  handler: main`

// Writes the faas.yaml file to the tmp directory for unit tests
func setup(tmpDir string, yamlFile string) {
	filepath := fmt.Sprintf("%s/%s", tmpDir, "faas.yaml")
	os.WriteFile(filepath, []byte(yamlFile), 0644)
}

func readFile(tmpDir string, filename string) (string, error) {
	filepath := fmt.Sprintf("%s/%s", tmpDir, filename)
	data, err := os.ReadFile(filepath)

	return string(data), err
}

func TestStartFaaS(t *testing.T) {
	tmpDir := t.TempDir()
	ConfigPath = tmpDir

	t.Run("should write up a Dockerfile and docker-compose.yaml file for nodejs", func(t *testing.T) {
		const expectedDockerFile = `FROM amazon/aws-lambda-nodejs:20

COPY node_modules ${LAMBDA_TASK_ROOT}/node_modules
COPY dist ${LAMBDA_TASK_ROOT}/dist
COPY package.json ${LAMBDA_TASK_ROOT}

CMD [ "dist/index.js" ]`

		setup(tmpDir, nodejsYaml)

		err := initializeDockerFiles("nodejs20.x", "dist/index.js")
		if err != nil {
			t.Errorf("expected no errors, but received \"%s\"", err.Error())
			return
		}

		dockerFile, err := readFile(tmpDir, "Dockerfile")
		if err != nil {
			t.Errorf("expected no errors reading Dockerfile, but receieved \"%s\"", err.Error())
			return
		}
		composeFile, err := readFile(tmpDir, "docker-compose.yaml")
		if err != nil {
			t.Errorf("expected no errors reading docker-compose.yaml, but receieved \"%s\"", err.Error())
			return
		}

		if composeFile != COMPOSE_TEMPLATE {
			t.Errorf("compose file written did not match the template")
			return
		}

		if dockerFile != expectedDockerFile {
			t.Errorf("dockerfile written did not match the expected value")
		}
	})
}
