package config

import (
	"testing"
)

func TestLoad_Defaults(t *testing.T) {
	t.Setenv("BOT_TOKEN", "")
	t.Setenv("REDIS_PASSWORD", "")
	t.Setenv("REDIS_ADDR", "")
	t.Setenv("REDIS_HOST", "")

	cfg := Load()
	if cfg == nil {
		t.Fatal("expected non-nil config")
	}

	if cfg.BotToken != "" {
		t.Errorf("got %q, want empty", cfg.BotToken)
	}
	if cfg.RedisAddr != "" {
		t.Errorf("got %q, want empty", cfg.RedisAddr)
	}
	if cfg.RedisHost != "" {
		t.Errorf("got %q, want empty", cfg.RedisHost)
	}
	if cfg.RedisPassword != DefaultRedisPassword {
		t.Errorf("got %q, want %q", cfg.RedisPassword, DefaultRedisPassword)
	}
}

func TestLoad_CustomEnv(t *testing.T) {
	t.Setenv("BOT_TOKEN", "123456:TEST_TOKEN")
	t.Setenv("REDIS_ADDR", "custom-redis:6379")
	t.Setenv("REDIS_HOST", "custom-host")
	t.Setenv("REDIS_PASSWORD", "super_secret_password")

	cfg := Load()
	if cfg == nil {
		t.Fatal("expected non-nil config")
	}

	if cfg.BotToken != "123456:TEST_TOKEN" {
		t.Errorf("got %q, want '123456:TEST_TOKEN'", cfg.BotToken)
	}
	if cfg.RedisAddr != "custom-redis:6379" {
		t.Errorf("got %q, want 'custom-redis:6379'", cfg.RedisAddr)
	}
	if cfg.RedisHost != "custom-host" {
		t.Errorf("got %q, want 'custom-host'", cfg.RedisHost)
	}
	if cfg.RedisPassword != "super_secret_password" {
		t.Errorf("got %q, want 'super_secret_password'", cfg.RedisPassword)
	}
}
