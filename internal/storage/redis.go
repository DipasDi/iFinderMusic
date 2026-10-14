package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"time"

	"random/internal/config"
	"random/internal/model"

	"github.com/redis/go-redis/v9"
)

var (
	Ctx = context.Background()
	RDB *redis.Client
)

// Init scans candidate addresses via quick TCP probe and connects to the first reachable Redis instance.
func Init(cfg *config.Config) (*redis.Client, error) {
	candidates := []string{}
	if cfg.RedisAddr != "" {
		candidates = append(candidates, cfg.RedisAddr)
	} else if cfg.RedisHost != "" {
		candidates = append(candidates, cfg.RedisHost+":6379")
	}
	candidates = append(candidates, "127.0.0.1:6379", "127.0.1.20:6379", "redis:6379", "kotichi_redis:6379", "host.docker.internal:6379")

	passwords := []string{cfg.RedisPassword}
	if cfg.RedisPassword != "" {
		passwords = append(passwords, "") // Try with password first, then without
	}

	for _, addr := range candidates {
		// Quick TCP probe to check if the port is reachable before creating connection pool
		conn, netErr := net.DialTimeout("tcp", addr, 200*time.Millisecond)
		if netErr != nil {
			continue
		}
		_ = conn.Close()

		for _, pass := range passwords {
			client := redis.NewClient(&redis.Options{
				Addr:         addr,
				Password:     pass,
				DB:           0,
				MaxRetries:   2,
				DialTimeout:  1500 * time.Millisecond,
				ReadTimeout:  1500 * time.Millisecond,
				WriteTimeout: 1500 * time.Millisecond,
			})

			ctxPing, cancel := context.WithTimeout(Ctx, 1500*time.Millisecond)
			pong, err := client.Ping(ctxPing).Result()
			cancel()

			if err == nil {
				RDB = client
				log.Printf("✅ Successfully connected to Redis (%s): %s\n", addr, pong)
				return RDB, nil
			}
			_ = client.Close()
		}
	}

	return nil, fmt.Errorf("failed to connect to any available Redis address: %v", candidates)
}

// SaveUserTracks caches the user's searched tracks list in Redis for 3 hours.
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

// GetUserTracks retrieves the cached tracks list for a user from Redis.
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

// SaveUserQuery saves the user's latest search query in Redis for 3 hours.
func SaveUserQuery(chatID int64, query string) error {
	if RDB == nil {
		return nil
	}
	key := fmt.Sprintf("user_query:%d", chatID)
	return RDB.Set(Ctx, key, query, 3*time.Hour).Err()
}

// GetUserQuery retrieves the user's latest search query from Redis.
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

// SaveUserListMessageID saves the message ID containing the track list for future cleanup.
func SaveUserListMessageID(chatID int64, messageID int) error {
	if RDB == nil {
		return nil
	}
	key := fmt.Sprintf("user_list_music:%d", chatID)
	return RDB.Set(Ctx, key, messageID, 24*time.Hour).Err()
}

// GetUserListMessageID retrieves the saved track list message ID.
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
