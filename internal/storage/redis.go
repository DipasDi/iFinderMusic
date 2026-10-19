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

<<<<<<< HEAD
// Init initializes and verifies the connection to the Redis client using the provided configuration.
=======
>>>>>>> 0c59af3e9076fadbe77438343714ce617cca44d0
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
<<<<<<< HEAD
		Username:     cfg.RedisUser,
=======
>>>>>>> 0c59af3e9076fadbe77438343714ce617cca44d0
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

<<<<<<< HEAD
// SaveUserTracks serializes and stores the user's searched tracks list in Redis with a 3-hour TTL.
=======
>>>>>>> 0c59af3e9076fadbe77438343714ce617cca44d0
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

<<<<<<< HEAD
// GetUserTracks retrieves and deserializes the user's cached tracks list from Redis.
=======
>>>>>>> 0c59af3e9076fadbe77438343714ce617cca44d0
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

<<<<<<< HEAD
// SaveUserQuery stores the latest search query string for the given chat ID with a 3-hour TTL.
=======
>>>>>>> 0c59af3e9076fadbe77438343714ce617cca44d0
func SaveUserQuery(chatID int64, query string) error {
	if RDB == nil {
		return nil
	}
	key := fmt.Sprintf("user_query:%d", chatID)
	return RDB.Set(Ctx, key, query, 3*time.Hour).Err()
}

<<<<<<< HEAD
// GetUserQuery retrieves the latest search query string for the given chat ID.
=======
>>>>>>> 0c59af3e9076fadbe77438343714ce617cca44d0
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

<<<<<<< HEAD
// SaveUserListMessageID stores the Telegram message ID of the active track list for future UI updates.
=======
>>>>>>> 0c59af3e9076fadbe77438343714ce617cca44d0
func SaveUserListMessageID(chatID int64, messageID int) error {
	if RDB == nil {
		return nil
	}
	key := fmt.Sprintf("user_list_music:%d", chatID)
	return RDB.Set(Ctx, key, messageID, 24*time.Hour).Err()
}

<<<<<<< HEAD
// GetUserListMessageID returns the Telegram message ID of the active track list for the given chat ID.
=======
>>>>>>> 0c59af3e9076fadbe77438343714ce617cca44d0
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
