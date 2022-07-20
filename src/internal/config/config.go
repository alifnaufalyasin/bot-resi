package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Database         DatabaseConfig
	Secret           string
	Port             string `env:"PORT,default=80"`
	TelegramApiToken string `env:"TELEGRAM_APITOKEN"`
	UriSicepat       string `env:"URI_SICEPAT"`
}

type DatabaseConfig struct {
	URL string `env:"DATABASE_URL,default=localhost:5432"`
}

func GetConfig() Config {
	err := godotenv.Load()
	if err != nil {

	}

	return Config{
		Database: DatabaseConfig{
			URL: os.Getenv("DATABASE_URL"),
		},
		Secret:           os.Getenv("SECRET"),
		Port:             os.Getenv("PORT"),
		TelegramApiToken: os.Getenv("TELEGRAM_APITOKEN"),
		UriSicepat:       os.Getenv("URI_SICEPAT"),
	}
}
