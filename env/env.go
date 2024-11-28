package env

import (
	"fmt"

	"github.com/spf13/viper"
)

var ConfigPath = "."

// Reads in the .env file, if none exists then it will write a new one
func ReadEnv() (*viper.Viper, error) {
	envFilePath := fmt.Sprintf("%s/.env", ConfigPath)
	env := viper.New()

	env.SetConfigType("env")
	env.SetConfigFile(envFilePath)
	env.SafeWriteConfigAs(envFilePath)
	err := env.ReadInConfig()

	if err != nil {
		return nil, err
	}

	return env, nil
}
