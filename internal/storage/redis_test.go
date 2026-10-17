package storage

import (
	"encoding/json"
	"testing"

	"random/internal/config"
	"random/internal/model"
)

func TestInit_ConnectionFailure(t *testing.T) {
	oldRDB := RDB
	defer func() {
		RDB = oldRDB
	}()

	cfg := &config.Config{
		RedisAddr: "127.0.0.1:54321", // unreachable port
	}

	client, err := Init(cfg)
	if err == nil {
		if client != nil {
			_ = client.Close()
		}
		t.Fatal("expected error connecting to dead port, got nil")
	}
}

func TestStorage_UninitializedClient(t *testing.T) {
	oldRDB := RDB
	RDB = nil
	defer func() {
		RDB = oldRDB
	}()

	chatID := int64(123456789)
	dummyTracks := []model.Music{
		{Title: "Test", Author: "Artist", Duration: "03:00", DownloadLink: "http://example.com/audio.mp3"},
	}

	t.Run("SaveUserTracks handles nil client", func(t *testing.T) {
		err := SaveUserTracks(chatID, dummyTracks)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		expectedErr := "redis client is not initialized"
		if err.Error() != expectedErr {
			t.Errorf("got %q, want %q", err.Error(), expectedErr)
		}
	})

	t.Run("GetUserTracks handles nil client", func(t *testing.T) {
		tracks, err := GetUserTracks(chatID)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if tracks != nil {
			t.Errorf("expected nil tracks, got %+v", tracks)
		}
	})

	t.Run("SaveUserQuery handles nil client", func(t *testing.T) {
		err := SaveUserQuery(chatID, "Queen")
		if err != nil {
			t.Errorf("expected nil error, got %v", err)
		}
	})

	t.Run("GetUserQuery handles nil client", func(t *testing.T) {
		query := GetUserQuery(chatID)
		if query != "" {
			t.Errorf("expected empty string, got %q", query)
		}
	})

	t.Run("SaveUserListMessageID handles nil client", func(t *testing.T) {
		err := SaveUserListMessageID(chatID, 999)
		if err != nil {
			t.Errorf("expected nil error, got %v", err)
		}
	})

	t.Run("GetUserListMessageID handles nil client", func(t *testing.T) {
		msgID := GetUserListMessageID(chatID)
		if msgID != 0 {
			t.Errorf("expected 0, got %d", msgID)
		}
	})
}

func TestTracksSerializationRoundTrip(t *testing.T) {
	original := []model.Music{
		{
			Title:        "Bohemian Rhapsody",
			Author:       "Queen",
			Duration:     "05:55",
			DownloadLink: "https://example.com/bohemian.mp3",
		},
		{
			Title:        "Группа крови",
			Author:       "Кино",
			Duration:     "04:45",
			DownloadLink: "/track/download/123",
		},
	}

	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("json.Marshal failed: %v", err)
	}

	var restored []model.Music
	if err := json.Unmarshal(data, &restored); err != nil {
		t.Fatalf("json.Unmarshal failed: %v", err)
	}

	if len(restored) != len(original) {
		t.Fatalf("got %d tracks, want %d", len(restored), len(original))
	}

	for i := range original {
		if restored[i] != original[i] {
			t.Errorf("track %d mismatch:\n got:  %+v\n want: %+v", i, restored[i], original[i])
		}
	}
}
