package main

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Logger   LoggerConfig
	Database DatabaseConfig
	Telegram TelegramConfig
}

type LoggerConfig struct {
	Level string
}

type DatabaseConfig struct {
	Type string
	DSN  string
}

type TelegramConfig struct {
	Token          string
	AdminUsernames []string
	DefaultPhotoID string
}

func loadConfig() (*Config, error) {
	viper.SetDefault("Logger.Level", "info")

	viper.SetDefault("Database.Type", "sqlite")
	viper.SetDefault("Database.DSN", "data/db.sqlite")
	viper.SetDefault("Telegram.Token", "environment variable should be configured")
	viper.SetDefault("Telegram.AdminUsernames", []string{"mirzavaziri", "kmirzavaziri"})
	viper.SetDefault(
		"Telegram.DefaultPhotoID",
		"AgACAgQAAxkBAAIBuWiYkahnYgABEKsrpvLpL9wBUkjx_AACe8gxGzOlwFAgMZfZP1IV-gEAAwIAA3kAAzYE",
	)

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	var config Config

	err := viper.Unmarshal(&config)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &config, nil
}
