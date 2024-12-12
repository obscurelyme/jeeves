package java

import (
	_ "embed"
	"fmt"
	"os"
	"testing"
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

	t.Run("Adds required plugins to a pom file", func(t *testing.T) {
		testPom, err := New(tmpDir)
		if err != nil {
			t.Errorf("error reading pom.xml: %s", err.Error())
			return
		}

		testPom.AddRequiredPlugins()
		testPom.WriteFile()

		if testPom.pom.Build.Plugins.Plugin[0].ArtifactId != "maven-dependency-plugin" {
			t.Errorf("Artifact-ID mismatch. Expected: %s, Actual: %s", "maven-dependency-plugin", testPom.pom.Build.Plugins.Plugin[0].ArtifactId)
		}
	})
}
