package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Env    string
	Server struct {
		Port int
	}
	Database struct {
		Host     string
		Port     int
		User     string
		Password string
		Name     string
		Sslmode  string
	}
	Redis struct {
		Host     string
		Port     int
		Password string
		Db       int
	}
	JWT struct {
		Secret               string `mapstructure:"secret"`
		ExpiresInMinutes     int    `mapstructure:"expires_in_minute"`
		RefreshExpiresInDays int    `mapstructure:"refresh_expires_in_days"`
	}
}

func LoadConfig() *Config {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %v", err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatalf("Unable to decode into struct, %v", err)
	}

	return &cfg
}
