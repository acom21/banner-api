package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	Host string
	Port int
	DB   string
}

func Load() *Config {
	return &Config{
		Host: getEnv("HOST", "0.0.0.0"),
		Port: getEnvInt("PORT", 8080),
		DB:   getEnv("DB_URL", "postgres://user:pass@localhost:5432/banner_db?sslmode=disable"),
	}
}

func (c *Config) Address() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
