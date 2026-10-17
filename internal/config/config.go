package config

import "os"

const DefaultRedisPassword = ""

type Config struct {
	BotToken      string
	RedisAddr     string
	RedisHost     string
	RedisPassword string
}

func Load() *Config {
	pass := os.Getenv("REDIS_PASSWORD")
	if pass == "" {
		pass = DefaultRedisPassword
	}

	return &Config{
		BotToken:      os.Getenv("BOT_TOKEN"),
		RedisAddr:     os.Getenv("REDIS_ADDR"),
		RedisHost:     os.Getenv("REDIS_HOST"),
		RedisPassword: pass,
	}
}
