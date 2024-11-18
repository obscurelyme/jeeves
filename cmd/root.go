package cmd

import (
	"log"
	"os"
	"path"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile     string
	userLicense string
	rootCmd     = &cobra.Command{
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

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.jeeves.yaml)")
	rootCmd.PersistentFlags().StringP("author", "a", "Mackenzie (obscurelyme) Greco", "author name for copyright attribution")
	rootCmd.PersistentFlags().StringVarP(&userLicense, "license", "l", "", "name of license for the project")
	rootCmd.PersistentFlags().Bool("viper", true, "use Viper for configuration")
	viper.BindPFlag("author", rootCmd.PersistentFlags().Lookup("author"))
	viper.BindPFlag("useViper", rootCmd.PersistentFlags().Lookup("viper"))
	viper.SetDefault("author", "Mackenzie (obscurelyme) Greco nico.frozencoffee@gmail.com")
	viper.SetDefault("license", "apache")
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
		// NOTE: do something I guess...
		log.Println("File changed...")
	})
	AWSConfig.WatchConfig()
}

func readConfigFile(cfg *viper.Viper) {
	err := cfg.ReadInConfig()

	if err != nil {
		cobra.CheckErr(err)
	}
}
