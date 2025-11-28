package config

import (
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

	_ = viper.ReadInConfig()

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
