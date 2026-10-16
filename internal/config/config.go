package config

import "os"

const (
	// DefaultRedisPassword пароль к Redis по умолчанию
	DefaultRedisPassword = "2e33611d240a42828f6e93ee6da481b3"
)

// Config содержит параметры конфигурации бота и сервисов
type Config struct {
	BotToken      string
	RedisAddr     string
	RedisHost     string
	RedisPassword string
}

// Load считывает конфигурацию из окружения или использует дефолтные значения
func Load() *Config {
	token := os.Getenv("BOT_TOKEN")

	redisPassword := os.Getenv("REDIS_PASSWORD")
	if redisPassword == "" {
		redisPassword = DefaultRedisPassword
	}

	return &Config{
		BotToken:      token,
		RedisAddr:     os.Getenv("REDIS_ADDR"),
		RedisHost:     os.Getenv("REDIS_HOST"),
		RedisPassword: redisPassword,
	}
}
