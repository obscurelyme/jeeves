package java

import (
	"encoding/xml"
	"fmt"
	"os"

	"github.com/obscurelyme/jeeves/utils/java/pom"
)

type JavaPomFileDriver struct {
	pom        pom.Project
	configFile string
}

func (jpd *JavaPomFileDriver) WriteFile() error {
	data, err := xml.MarshalIndent(jpd.pom, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(fmt.Sprintf("%s/pom.xml", jpd.configFile), data, 0644)
}

func New(configPath string) (*JavaPomFileDriver, error) {
	j := new(JavaPomFileDriver)
	j.configFile = configPath

	data, err := os.ReadFile(fmt.Sprintf("%s/pom.xml", configPath))
	if err != nil {
		return nil, err
	}

	err = xml.Unmarshal(data, &j.pom)

	if err != nil {
		return nil, err
	}

	return j, nil
}
