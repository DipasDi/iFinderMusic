package config

import (
	"testing"
)

// TestLoad_Defaults verifies that default configuration values are populated when environment variables are omitted.
func TestLoad_Defaults(t *testing.T) {
	t.Setenv("BOT_TOKEN", "")
	t.Setenv("REDIS_PASSWORD", "")
	t.Setenv("REDIS_ADDR", "")
	t.Setenv("REDIS_HOST", "")
	t.Setenv("REDIS_USER", "")
	t.Setenv("REDIS_USERNAME", "")
	t.Setenv("USERS_FILE_PATH", "")

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
	if cfg.RedisUser != DefaultRedisUser {
		t.Errorf("got %q, want %q", cfg.RedisUser, DefaultRedisUser)
	}
	if cfg.UsersFilePath != DefaultUsersFilePath {
		t.Errorf("got %q, want %q", cfg.UsersFilePath, DefaultUsersFilePath)
	}
	if !cfg.TrackUsers {
		t.Errorf("expected TrackUsers to default to true, got false")
	}
}

// TestLoad_CustomEnv tests parsing custom configuration values set in the environment.
func TestLoad_CustomEnv(t *testing.T) {
	t.Setenv("BOT_TOKEN", "123456:TEST_TOKEN")
	t.Setenv("REDIS_ADDR", "custom-redis:6379")
	t.Setenv("REDIS_HOST", "custom-host")
	t.Setenv("REDIS_PASSWORD", "super_secret_password")
	t.Setenv("REDIS_USER", "custom_user")
	t.Setenv("USERS_FILE_PATH", "custom/path/users.json")
	t.Setenv("TRACK_USERS", "true")

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
	if cfg.RedisUser != "custom_user" {
		t.Errorf("got %q, want 'custom_user'", cfg.RedisUser)
	}
	if cfg.UsersFilePath != "custom/path/users.json" {
		t.Errorf("got %q, want 'custom/path/users.json'", cfg.UsersFilePath)
	}
	if !cfg.TrackUsers {
		t.Errorf("expected TrackUsers to be true")
	}
}

// TestLoad_TrackUsersDisabled tests disabling user tracking via TRACK_USERS or ENABLE_USER_TRACKING environment variables.
func TestLoad_TrackUsersDisabled(t *testing.T) {
	t.Run("TRACK_USERS=false disables tracking", func(t *testing.T) {
		t.Setenv("TRACK_USERS", "false")
		t.Setenv("ENABLE_USER_TRACKING", "")
		cfg := Load()
		if cfg.TrackUsers {
			t.Errorf("expected TrackUsers=false, got true")
		}
	})

	t.Run("ENABLE_USER_TRACKING=0 disables tracking", func(t *testing.T) {
		t.Setenv("TRACK_USERS", "")
		t.Setenv("ENABLE_USER_TRACKING", "0")
		cfg := Load()
		if cfg.TrackUsers {
			t.Errorf("expected TrackUsers=false, got true")
		}
	})
}

// TestLoad_RedisUsernameFallback tests fallback from REDIS_USER to REDIS_USERNAME.
func TestLoad_RedisUsernameFallback(t *testing.T) {
	t.Setenv("REDIS_USER", "")
	t.Setenv("REDIS_USERNAME", "acl_username")

	cfg := Load()
	if cfg == nil {
		t.Fatal("expected non-nil config")
	}
	if cfg.RedisUser != "acl_username" {
		t.Errorf("got %q, want 'acl_username'", cfg.RedisUser)
	}
}
