package java

import (
	_ "embed"
	"encoding/xml"
	"fmt"
	"os"
	"testing"

	"github.com/obscurelyme/jeeves/utils/java/pom"
)

//go:embed pom.xml
var testPomFile []byte

func setup(tmpDir string) error {
	return os.WriteFile(fmt.Sprintf("%s/pom.xml", tmpDir), testPomFile, 0644)
}

func TestJavaPom(t *testing.T) {
	tmpDir := t.TempDir()
	err := setup(tmpDir)
	if err != nil {
		t.Errorf("Error setting up test file: %s", err.Error())
		return
	}

	t.Run("Writes to a pom file", func(t *testing.T) {
		testPom, err := New(tmpDir)
		if err != nil {
			t.Errorf("error reading pom.xml: %s", err.Error())
			return
		}

		testPom.pom.Plugins = &pom.Plugins{Comment: "added by test", Plugin: []pom.Plugin{{
			XMLName: xml.Name{
				Space: "",
				Local: "plugin",
			},
			GroupId:    "org.apache.maven.plugins",
			ArtifactId: "maven-dependency-plugin",
			Version:    "3.1.2",
			Configuration: &pom.Any{
				XMLName: xml.Name{Space: "", Local: "configuration"},
				AnyElements: []pom.Any{
					{
						XMLName: xml.Name{Space: "", Local: "includeScope"},
						Value:   "runtime",
					},
				},
			},
			Executions: &pom.Executions{
				Execution: []pom.Execution{
					{
						XMLName: xml.Name{
							Space: "",
							Local: "execution",
						},
						Id:    "copy-dependencies-package",
						Phase: "package",
						Goals: &pom.Goals{
							Comment: "copy the dependencies",
							Goal:    []string{"copy-dependency"},
						},
					},
				},
			},
		}}}

		if testPom.pom.ArtifactId != "my-app" {
			t.Errorf("Artifact-ID mismatch. Expected: %s, Actual: %s", "my-app", testPom.pom.ArtifactId)
		}

		err = testPom.WriteFile()
		if err != nil {
			t.Errorf("error writing new pom: %s", err.Error())
			return
		}
	})
}
