package config

type Database struct {
	Postgres Postgres `mapstructure:"postgres"`
}

type Postgres struct {
	Host     string `mapstructure:"POSTGRES_HOST"`
	Port     string `mapstructure:"POSTGRES_PORT"`
	User     string `mapstructure:"POSTGRES_USER"`
	Password string `mapstructure:"POSTGRES_PASSWORD"`
	DBName   string `mapstructure:"POSTGRES_DBNAME"`
	Schema   string `mapstructure:"POSTGRES_SCHEMA"`
}
