package java

import (
	"encoding/xml"
	"fmt"
	"os"

	"github.com/obscurelyme/jeeves/utils/java/pom"
)

// Required plugin in order to copy dependencies to the Docker container
const REQUIRED_MAVEN_PLUGIN string = "maven-dependency-plugin"

var plugins []pom.Plugin = []pom.Plugin{{
	Comment:    "Required plugin for Jeeves - maven-dependency-plugin",
	GroupId:    "org.apache.maven.plugins",
	ArtifactId: "maven-dependency-plugin",
	Version:    "3.1.2",
	Configuration: &pom.Any{
		XMLName: xml.Name{Space: "", Local: "configuration"},
		Children: []pom.Any{
			{
				XMLName: xml.Name{Space: "", Local: "includeScope"},
				Value:   "runtime",
			},
		},
	},
	Executions: &pom.Executions{
		Execution: []pom.Execution{
			{
				Id:    "copy-dependencies-package",
				Phase: "package",
				Goals: &pom.Goals{
					Comment: "copy the dependencies",
					Goal:    []string{"copy-dependency"},
				},
			},
			{
				Id:    "copy-dependencies-compile",
				Phase: "compile",
				Goals: &pom.Goals{
					Comment: "copy the dependencies",
					Goal:    []string{"copy-dependency"},
				},
			},
		},
	},
}}

type MavenPomDriver interface {
	HasRequiredPlugins() bool
	AddRequiredPlugins()
	WriteFile() error
}

type MavenPomFileDriver struct {
	pom        pom.Project
	configFile string
}

func (jpd *MavenPomFileDriver) WriteFile() error {
	data, err := xml.MarshalIndent(jpd.pom, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(fmt.Sprintf("%s/pom.xml", jpd.configFile), data, 0644)
}

// Checks if the pom file has the required maven plugin, REQUIRED_MAVEN_PLUGIN
func (jpd *MavenPomFileDriver) HasRequiredPlugins() bool {
	for _, plugin := range jpd.pom.Build.Plugins.Plugin {
		if plugin.ArtifactId == REQUIRED_MAVEN_PLUGIN {
			return true
		}
	}

	return false
}

func (jpd *MavenPomFileDriver) AddRequiredPlugins() {
	if jpd.pom.Build.Plugins == nil {
		jpd.pom.Build.Plugins = &pom.Plugins{
			Comment: "Plugin[s] added by Jeeves",
			Plugin:  plugins,
		}
		return
	}

	if jpd.pom.Build.Plugins.Plugin == nil {
		jpd.pom.Build.Plugins.Plugin = plugins
		return
	}

	jpd.pom.Build.Plugins.Plugin = append(jpd.pom.Build.Plugins.Plugin, plugins...)
}

func New(configPath string) (*MavenPomFileDriver, error) {
	j := new(MavenPomFileDriver)
	j.configFile = configPath

	data, err := os.ReadFile(fmt.Sprintf("%s/pom.xml", configPath))
	if err != nil {
		return nil, err
	}

	err = xml.Unmarshal(data, &j.pom)

	if err != nil {
		return nil, err
	}

	fmt.Println("returning pointer")
	return j, nil
}
