package utils

import (
	"os"
	"path"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var AWSConfig *viper.Viper
var AWSConfigPath string

func init() {
	home, err := os.UserHomeDir()
	cobra.CheckErr(err)
	configPath := path.Join(home, ".aws")

	AWSConfig = viper.New()
	AWSConfig.AddConfigPath(home)
	AWSConfig.SetConfigType("ini")
	AWSConfig.AddConfigPath(configPath)
	AWSConfig.SetConfigName("config")

	err = AWSConfig.ReadInConfig()
	cobra.CheckErr(err)
}
