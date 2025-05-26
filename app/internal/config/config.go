package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Auth struct {
		Email     string
		Password  string
		ServerUrl string
	}
}

func Load() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading congig")
	}

	cfg := &Config{}
	cfg.Auth.Email = getEnv("EMAIL", "your@email.addr")
	cfg.Auth.Password = getEnv("PASSWORD", "password")
	cfg.Auth.ServerUrl = getEnv("SERVER_URL", "http://yourserver.com")

	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
