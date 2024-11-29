package utils

import (
	"fmt"
	"os"
	"path"

	"github.com/spf13/viper"
)

var (
	AWSConfig     *viper.Viper
	AWSConfigPath string
	JeevesConfig  *viper.Viper
)

// Loads the .jeeves.yaml config file, creates it if it does not exist
func LoadJeevesConfig() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	JeevesConfig = viper.New()
	JeevesConfig.AddConfigPath(home)
	JeevesConfig.SetConfigType("yaml")
	JeevesConfig.SetConfigName(".jeeves")
	JeevesConfig.SafeWriteConfig()

	return JeevesConfig.ReadInConfig()
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
