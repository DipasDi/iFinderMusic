package config

import (
	"bufio"
	"os"
	"strings"
)

// Config holds all application configuration parameters.
type Config struct {
	BotToken      string
	RedisAddr     string
	RedisHost     string
	RedisPassword string
}

// Load reads configuration from environment variables (loading from .env if present).
func Load() *Config {
	loadEnvFile(".env")

	return &Config{
		BotToken:      os.Getenv("BOT_TOKEN"),
		RedisAddr:     os.Getenv("REDIS_ADDR"),
		RedisHost:     os.Getenv("REDIS_HOST"),
		RedisPassword: os.Getenv("REDIS_PASSWORD"),
	}
}

// loadEnvFile reads a key=value .env file and sets environment variables if not already defined.
func loadEnvFile(path string) {
	file, err := os.Open(path)
	if err != nil {
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			val := strings.TrimSpace(parts[1])
			val = strings.Trim(val, `"'`)
			if os.Getenv(key) == "" {
				_ = os.Setenv(key, val)
			}
		}
	}
}
