package cmd

import (
	"github.com/obscurelyme/jeeves/cmd/ai"
	"github.com/obscurelyme/jeeves/cmd/faas"
	"github.com/obscurelyme/jeeves/cmd/s3"
	"github.com/obscurelyme/jeeves/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	rootCmd = &cobra.Command{
		Use:   "jeeves",
		Short: "A helpful CLI for your AWS infrastructure",
		Long:  "A helpful CLI for your AWS infrastructure",
	}
)

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.AddCommand(faas.FaasRootCmd)
	rootCmd.AddCommand(s3.S3RootCmd)
	rootCmd.AddCommand(ai.AIRootCmd)
}

func initConfig() {
	utils.LoadAWSConfig()
	// TODO: need to check for an invalid config, else quit app and print err here
	utils.LoadJeevesConfig()

	viper.AutomaticEnv()
}
