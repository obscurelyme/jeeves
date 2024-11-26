package cmd

import (
	"os"
	"path"

	"github.com/fsnotify/fsnotify"
	"github.com/obscurelyme/jeeves/cmd/faas"
	"github.com/obscurelyme/jeeves/cmd/s3"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	rootCmd = &cobra.Command{
		Use:   "jeeves",
		Short: "A helpful CLI for your AWS infrastructure",
		Long:  "A helpful CLI for your AWS infrastructure",
	}
	JeevesConfig           *viper.Viper
	AWSConfig              *viper.Viper
	AWSCredentials         *viper.Viper
	AWSCredentialsFilePath string
)

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.AddCommand(faas.FaasRootCmd)
	rootCmd.AddCommand(s3.S3RootCmd)
}

func initConfig() {
	home, err := os.UserHomeDir()
	cobra.CheckErr(err)

	// Create Viper config structs
	JeevesConfig = viper.New()
	AWSConfig = viper.New()
	AWSCredentials = viper.New()

	JeevesConfig.AddConfigPath(home)
	JeevesConfig.SetConfigType("yaml")
	JeevesConfig.SetConfigName(".jeeves")

	configPath := path.Join(home, ".aws")

	AWSConfig.SetConfigType("ini")
	AWSConfig.AddConfigPath(configPath)
	AWSConfig.SetConfigName("config")

	AWSCredentials.SetConfigType("ini")
	AWSCredentials.AddConfigPath(configPath)
	AWSCredentials.SetConfigName("credentials")
	AWSCredentialsFilePath = path.Join(configPath, "credentials")

	viper.AutomaticEnv()

	readConfigFile(JeevesConfig)
	readConfigFile(AWSConfig)
	readConfigFile(AWSCredentials)

	// Watch the configs
	AWSConfig.OnConfigChange(func(e fsnotify.Event) {
		AWSConfig.ReadInConfig()
	})
	AWSConfig.WatchConfig()
}

func readConfigFile(cfg *viper.Viper) {
	err := cfg.ReadInConfig()

	if err != nil {
		cfg.WriteConfig()
	}
}
