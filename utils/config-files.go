package utils

import (
	"fmt"
	"os"
	"path"

	"github.com/goccy/go-yaml"
	"github.com/spf13/viper"
)

var (
	AWSConfig     *viper.Viper
	AWSConfigPath string
	Jeeves        *YamlConfigFile
)

type JeevesAI struct {
	// List of approved models for use with Jeeves
	ApprovedModels []string `yaml:"ApprovedModels"`
	// The currently set model to use with all Jeeves AI commands
	PreferredModel string `yaml:"PreferredModel"`
}

type JeevesSSO struct {
	// The start url for SSO
	Start string `yaml:"Start"`
}

type JeevesConfig struct {
	// Jeeves AI configuration
	AI JeevesAI `yaml:"AI"`
	// Jeeves SSO configuration
	SSO JeevesSSO `yaml:"SSO"`
}

type YamlConfigFile struct {
	configFile     string
	rawData        []byte
	ConfigSettings JeevesConfig
}

func (cfg *YamlConfigFile) SetConfigFile(filepath string) {
	cfg.configFile = filepath
}

func (cfg *YamlConfigFile) ReadInConfig() error {
	var err error
	cfg.rawData, err = os.ReadFile(cfg.configFile)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(cfg.rawData, &cfg.ConfigSettings)
}

func (cfg *YamlConfigFile) WriteConfig() error {
	data, err := yaml.Marshal(cfg.ConfigSettings)
	if err != nil {
		return err
	}

	return os.WriteFile(cfg.configFile, data, 0644)
}

// Loads the .jeeves.yaml config file, creates it if it does not exist
func LoadJeevesConfig() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	Jeeves = new(YamlConfigFile)
	Jeeves.SetConfigFile(fmt.Sprintf("%s/.jeeves.yaml", home))
	return Jeeves.ReadInConfig()
}

// Loads the .aws/config file into a Viper struct
func LoadAWSConfig() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	configPath := path.Join(home, ".aws")

	AWSConfig = viper.New()
	AWSConfig.AddConfigPath(home)
	AWSConfig.SetConfigType("ini")
	AWSConfig.AddConfigPath(configPath)
	AWSConfig.SetConfigName("config")

	err = AWSConfig.ReadInConfig()
	if err != nil {
		fmt.Println("Could not read .aws/config file, please run \"jeeves login\" first")
		return err
	}

	return nil
}
