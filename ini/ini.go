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
				insertKey(cfg, key, innerKey, innerValue)
			}
		}
	}

	return cfg.SaveTo(filename)
}

func insertKey(cfg *ini.File, sectionName string, key string, value any) {
	if strValue, ok := value.(string); ok {
		cfg.Section(sectionName).NewKey(key, strValue)
	}

	if _, ok := value.(bool); ok {
		cfg.Section(sectionName).NewBooleanKey(key)
	}
}
