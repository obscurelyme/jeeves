package env

import (
	"errors"
	"io/fs"

	"github.com/spf13/viper"
)

var ConfigPath = "."

// Reads in the .env file, if none exists then it will write a new one
func ReadEnv() (*viper.Viper, error) {
	env := viper.New()

	env.AddConfigPath(ConfigPath)
	env.SetConfigFile(".env")
	err := env.ReadInConfig()

	var pathErr *fs.PathError
	if errors.As(err, &pathErr) {
		env.WriteConfig()
	}

	return env, nil
}
