package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"random/internal/config"
	"random/internal/model"

	"github.com/redis/go-redis/v9"
)

var (
	Ctx = context.Background()
	RDB *redis.Client
)

func Init(cfg *config.Config) (*redis.Client, error) {
	addr := cfg.RedisAddr
	if addr == "" && cfg.RedisHost != "" {
		addr = cfg.RedisHost + ":6379"
	}
	if addr == "" {
		addr = "127.0.0.1:6379"
	}

	client := redis.NewClient(&redis.Options{
		Addr:         addr,
		Password:     cfg.RedisPassword,
		DB:           0,
		DialTimeout:  2 * time.Second,
		ReadTimeout:  2 * time.Second,
		WriteTimeout: 2 * time.Second,
	})

	ctx, cancel := context.WithTimeout(Ctx, 2*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		_ = client.Close()
		return nil, fmt.Errorf("failed to connect to redis at %s: %w", addr, err)
	}

	RDB = client
	return RDB, nil
}

func SaveUserTracks(chatID int64, tracks []model.Music) error {
	if RDB == nil {
		return fmt.Errorf("redis client is not initialized")
	}

	data, err := json.Marshal(tracks)
	if err != nil {
		return fmt.Errorf("failed to serialize tracks: %w", err)
	}

	key := fmt.Sprintf("user_music:%d", chatID)
	return RDB.Set(Ctx, key, data, 3*time.Hour).Err()
}

func GetUserTracks(chatID int64) ([]model.Music, error) {
	if RDB == nil {
		return nil, fmt.Errorf("redis client is not initialized")
	}

	key := fmt.Sprintf("user_music:%d", chatID)
	data, err := RDB.Get(Ctx, key).Bytes()
	if err != nil {
		return nil, err
	}

	var tracks []model.Music
	if err := json.Unmarshal(data, &tracks); err != nil {
		return nil, fmt.Errorf("failed to deserialize tracks: %w", err)
	}

	return tracks, nil
}

func SaveUserQuery(chatID int64, query string) error {
	if RDB == nil {
		return nil
	}
	key := fmt.Sprintf("user_query:%d", chatID)
	return RDB.Set(Ctx, key, query, 3*time.Hour).Err()
}

func GetUserQuery(chatID int64) string {
	if RDB == nil {
		return ""
	}
	key := fmt.Sprintf("user_query:%d", chatID)
	val, err := RDB.Get(Ctx, key).Result()
	if err != nil {
		return ""
	}
	return val
}

func SaveUserListMessageID(chatID int64, messageID int) error {
	if RDB == nil {
		return nil
	}
	key := fmt.Sprintf("user_list_music:%d", chatID)
	return RDB.Set(Ctx, key, messageID, 24*time.Hour).Err()
}

func GetUserListMessageID(chatID int64) int {
	if RDB == nil {
		return 0
	}
	key := fmt.Sprintf("user_list_music:%d", chatID)
	data, err := RDB.Get(Ctx, key).Bytes()
	if err != nil {
		return 0
	}
	var msgID int
	if err := json.Unmarshal(data, &msgID); err != nil {
		return 0
	}
	return msgID
}
