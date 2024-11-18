package ini

import (
	"gopkg.in/ini.v1"
)

// Patch to cover the fact that viper cannot handle AWS INI files.
//
// VERY NAIVE CODE, do not try to use this for anything other than AWS Credential files.
func WriteIniFile(filename string, settings map[string]any) error {
	cfg := ini.Empty()

	for key, value := range settings {
		cfg.NewSection(key)
		if v, ok := value.(map[string]any); ok {
			for innerKey, innerValue := range v {
				cfg.Section(key).NewKey(innerKey, innerValue.(string))
			}
		}
	}

	return cfg.SaveTo(filename)
}
