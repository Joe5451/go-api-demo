package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Debug    bool
	Database Database
}

func LoadConfig() (*Config, error) {
	viper.AutomaticEnv()

	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load config file, %v", err)
	}

	return &Config{
		Debug: viper.GetBool("DEBUG"),
		Database: Database{
			Postgres: Postgres{
				Host:     viper.GetString("POSTGRES_HOST"),
				Port:     viper.GetString("POSTGRES_PORT"),
				User:     viper.GetString("POSTGRES_USER"),
				Password: viper.GetString("POSTGRES_PASSWORD"),
				DBName:   viper.GetString("POSTGRES_DBNAME"),
				Schema:   viper.GetString("POSTGRES_SCHEMA"),
			},
		},
	}, nil
}
