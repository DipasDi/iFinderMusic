package ui

import (
	"fmt"
	"strings"
	"testing"
	"unicode/utf8"

	"random/internal/model"
)

func TestButtonTrackText(t *testing.T) {
	tests := []struct {
		name     string
		track    model.Music
		expected string
	}{
		{
			name: "short track title and author",
			track: model.Music{
				Author:   "Queen",
				Title:    "Bohemian Rhapsody",
				Duration: "05:55",
			},
			expected: "Queen — Bohemian Rhapsody [05:55]",
		},
		{
			name: "cyrillic track text",
			track: model.Music{
				Author:   "Кино",
				Title:    "Группа крови",
				Duration: "04:45",
			},
			expected: "Кино — Группа крови [04:45]",
		},
		{
			name: "empty music fields",
			track: model.Music{
				Author:   "",
				Title:    "",
				Duration: "",
			},
			expected: " —  []",
		},
		{
			name: "exactly 60 runes does not truncate",
			track: model.Music{
				Author:   "12345678901234567890",
				Title:    "12345678901234567890123456789",
				Duration: "05:00",
			},
			expected: "12345678901234567890 — 12345678901234567890123456789 [05:00]",
		},
		{
			name: "exceeds 60 runes gets truncated to 57 runes plus ellipsis",
			track: model.Music{
				Author:   "Very Long Artist Name With Multiple Words And Features",
				Title:    "Extremely Long Track Title Remastered Deluxe Edition 2026",
				Duration: "03:45",
			},
			expected: "Very Long Artist Name With Multiple Words And Features — ...",
		},
		{
			name: "long cyrillic track gets truncated by runes correctly",
			track: model.Music{
				Author:   "Очень Длинный Российский Исполнитель С Длинным Псевдонимом",
				Title:    "Невероятно Длинное Название Песни С Описанием Ремикса",
				Duration: "04:12",
			},
			expected: "Очень Длинный Российский Исполнитель С Длинным Псевдонимо...",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := ButtonTrackText(tt.track)

			if actual != tt.expected {
				t.Errorf("ButtonTrackText() = %q, want %q", actual, tt.expected)
			}

			runeCount := utf8.RuneCountInString(actual)
			if runeCount > 60 {
				t.Errorf("ButtonTrackText() length in runes = %d, exceeds max 60", runeCount)
			}

			fullText := fmt.Sprintf("%s — %s [%s]", tt.track.Author, tt.track.Title, tt.track.Duration)
			if utf8.RuneCountInString(fullText) > 60 {
				if !strings.HasSuffix(actual, "...") {
					t.Errorf("expected ellipsis '...' at the end of truncated text, got %q", actual)
				}
				if runeCount != 60 {
					t.Errorf("expected truncated string length to be exactly 60 runes, got %d", runeCount)
				}
			}
		})
	}
}
