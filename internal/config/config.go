package config

import (
	"os"
	"strconv"
)

const (
	DefaultRedisPassword = ""
	DefaultRedisUser     = "default"
	DefaultUsersFilePath = "users.json"
	DefaultTrackUsers    = true
)

// Config holds the application configuration loaded from environment variables.
type Config struct {
	BotToken      string
	RedisAddr     string
	RedisHost     string
	RedisPassword string
	RedisUser     string
	UsersFilePath string
	TrackUsers    bool
}

// getEnv returns the first non-empty value among the given environment keys, or fallback.
func getEnv(fallback string, keys ...string) string {
	for _, key := range keys {
		if val := os.Getenv(key); val != "" {
			return val
		}
	}
	return fallback
}

// getEnvBool parses the first non-empty boolean value among keys, or returns fallback.
func getEnvBool(fallback bool, keys ...string) bool {
	for _, key := range keys {
		if val := os.Getenv(key); val != "" {
			if b, err := strconv.ParseBool(val); err == nil {
				return b
			}
		}
	}
	return fallback
}

// Load reads application configuration from environment variables with sensible defaults.
func Load() *Config {
	return &Config{
		BotToken:      os.Getenv("BOT_TOKEN"),
		RedisAddr:     os.Getenv("REDIS_ADDR"),
		RedisHost:     os.Getenv("REDIS_HOST"),
		RedisPassword: getEnv(DefaultRedisPassword, "REDIS_PASSWORD"),
		RedisUser:     getEnv(DefaultRedisUser, "REDIS_USER", "REDIS_USERNAME"),
		UsersFilePath: getEnv(DefaultUsersFilePath, "USERS_FILE_PATH"),
		TrackUsers:    getEnvBool(DefaultTrackUsers, "TRACK_USERS", "ENABLE_USER_TRACKING"),
	}
}
