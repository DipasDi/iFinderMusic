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

// Init выполняет перебор адресов, быстрое TCP-сканирование и подключение к доступному Redis
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
		passwords = append(passwords, "") // Пробуем с паролем, затем без
	}

	for _, addr := range candidates {
		// Быстрый tcp-пробив порта перед созданием пула соединений
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
				log.Printf("✅ Успешное подключение к Redis (%s): %s\n", addr, pong)
				return RDB, nil
			}
			_ = client.Close()
		}
	}

	return nil, fmt.Errorf("не удалось подключиться ни к одному из доступных адресов Redis: %v", candidates)
}

// SaveUserTracks сохраняет список найденных треков пользователя в Redis на 3 часа
func SaveUserTracks(chatID int64, tracks []model.Music) error {
	if RDB == nil {
		return fmt.Errorf("клиент Redis не инициализирован")
	}

	data, err := json.Marshal(tracks)
	if err != nil {
		return fmt.Errorf("ошибка сериализации треков: %w", err)
	}

	key := fmt.Sprintf("user_music:%d", chatID)
	return RDB.Set(Ctx, key, data, 3*time.Hour).Err()
}

// GetUserTracks получает сохраненный список треков пользователя из Redis
func GetUserTracks(chatID int64) ([]model.Music, error) {
	if RDB == nil {
		return nil, fmt.Errorf("клиент Redis не инициализирован")
	}

	key := fmt.Sprintf("user_music:%d", chatID)
	data, err := RDB.Get(Ctx, key).Bytes()
	if err != nil {
		return nil, err
	}

	var tracks []model.Music
	if err := json.Unmarshal(data, &tracks); err != nil {
		return nil, fmt.Errorf("ошибка десериализации треков: %w", err)
	}

	return tracks, nil
}

// SaveUserQuery сохраняет последний поисковый запрос пользователя в Redis на 3 часа
func SaveUserQuery(chatID int64, query string) error {
	if RDB == nil {
		return nil
	}
	key := fmt.Sprintf("user_query:%d", chatID)
	return RDB.Set(Ctx, key, query, 3*time.Hour).Err()
}

// GetUserQuery получает последний поисковый запрос пользователя из Redis
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

// SaveUserListMessageID сохраняет ID сообщения со списком треков для очистки
func SaveUserListMessageID(chatID int64, messageID int) error {
	if RDB == nil {
		return nil
	}
	key := fmt.Sprintf("user_list_music:%d", chatID)
	return RDB.Set(Ctx, key, messageID, 24*time.Hour).Err()
}

// GetUserListMessageID получает сохраненный ID сообщения со списком треков
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
