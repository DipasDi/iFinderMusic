package ui

import (
	"fmt"
	"strings"
	"testing"
	"unicode/utf8"

	"random/internal/model"
)

<<<<<<< HEAD
// TestButtonTrackText tests truncation and formatting of music button text.
=======
>>>>>>> 0c59af3e9076fadbe77438343714ce617cca44d0
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

<<<<<<< HEAD
// TestMainMenuKeyboard tests the structure and layout of the main reply keyboard.
=======
>>>>>>> 0c59af3e9076fadbe77438343714ce617cca44d0
func TestMainMenuKeyboard(t *testing.T) {
	kb := MainMenuKeyboard()
	if kb == nil {
		t.Fatal("MainMenuKeyboard() returned nil")
	}

	if !kb.ResizeKeyboard {
		t.Errorf("expected ResizeKeyboard to be true, got false")
	}

	if len(kb.Keyboard) != 2 {
		t.Fatalf("expected 2 rows in keyboard, got %d", len(kb.Keyboard))
	}

<<<<<<< HEAD
=======
	// Row 0: "🔥 Популярные", "✨ Новинки"
>>>>>>> 0c59af3e9076fadbe77438343714ce617cca44d0
	if len(kb.Keyboard[0]) != 2 {
		t.Fatalf("expected 2 buttons in row 0, got %d", len(kb.Keyboard[0]))
	}
	if kb.Keyboard[0][0].Text != "🔥 Популярные" {
		t.Errorf("expected button 0,0 text '🔥 Популярные', got %q", kb.Keyboard[0][0].Text)
	}
	if kb.Keyboard[0][1].Text != "✨ Новинки" {
		t.Errorf("expected button 0,1 text '✨ Новинки', got %q", kb.Keyboard[0][1].Text)
	}

<<<<<<< HEAD
=======
	// Row 1: "❓ Помощь"
>>>>>>> 0c59af3e9076fadbe77438343714ce617cca44d0
	if len(kb.Keyboard[1]) != 1 {
		t.Fatalf("expected 1 button in row 1, got %d", len(kb.Keyboard[1]))
	}
	if kb.Keyboard[1][0].Text != "❓ Помощь" {
		t.Errorf("expected button 1,0 text '❓ Помощь', got %q", kb.Keyboard[1][0].Text)
	}
}

<<<<<<< HEAD
// TestBuildTracksKeyboard tests inline keyboard pagination and callback formatting.
=======
>>>>>>> 0c59af3e9076fadbe77438343714ce617cca44d0
func TestBuildTracksKeyboard(t *testing.T) {
	generateTracks := func(count int) []model.Music {
		tracks := make([]model.Music, count)
		for i := 0; i < count; i++ {
			tracks[i] = model.Music{
				Title:    fmt.Sprintf("Track %d", i+1),
				Author:   fmt.Sprintf("Artist %d", i+1),
				Duration: "03:30",
			}
		}
		return tracks
	}

	t.Run("empty tracks slice returns empty keyboard", func(t *testing.T) {
		kb := BuildTracksKeyboard(nil, 0)
		if kb == nil {
			t.Fatal("expected non-nil keyboard")
		}
		if len(kb.InlineKeyboard) != 0 {
			t.Errorf("expected 0 rows for empty tracks, got %d", len(kb.InlineKeyboard))
		}
	})

	t.Run("single page of tracks has no pagination row", func(t *testing.T) {
		tracks := generateTracks(5)
		kb := BuildTracksKeyboard(tracks, 0)

<<<<<<< HEAD
=======
		// 5 tracks -> 5 rows, no pagination row because totalPages == 1
>>>>>>> 0c59af3e9076fadbe77438343714ce617cca44d0
		if len(kb.InlineKeyboard) != 5 {
			t.Fatalf("expected 5 rows, got %d", len(kb.InlineKeyboard))
		}

		for i := 0; i < 5; i++ {
			row := kb.InlineKeyboard[i]
			if len(row) != 1 {
				t.Fatalf("expected row %d to have 1 button, got %d", i, len(row))
			}
			expectedCallback := fmt.Sprintf("dl_%d", i)
			if row[0].CallbackData != expectedCallback {
				t.Errorf("expected callback %q, got %q", expectedCallback, row[0].CallbackData)
			}
		}
	})

	t.Run("first page of multi-page tracks has next button but no prev button", func(t *testing.T) {
<<<<<<< HEAD
		tracks := generateTracks(25)
		kb := BuildTracksKeyboard(tracks, 0)

=======
		tracks := generateTracks(25) // 25 tracks -> 3 pages (10 + 10 + 5)
		kb := BuildTracksKeyboard(tracks, 0)

		// 10 track rows + 1 pagination row = 11 rows
>>>>>>> 0c59af3e9076fadbe77438343714ce617cca44d0
		if len(kb.InlineKeyboard) != 11 {
			t.Fatalf("expected 11 rows, got %d", len(kb.InlineKeyboard))
		}

		navRow := kb.InlineKeyboard[10]
<<<<<<< HEAD
=======
		// On page 0: only "📄 1/3" and "Вперёд ➡️" (2 buttons)
>>>>>>> 0c59af3e9076fadbe77438343714ce617cca44d0
		if len(navRow) != 2 {
			t.Fatalf("expected 2 nav buttons on page 0, got %d", len(navRow))
		}

		if navRow[0].Text != "📄 1/3" || navRow[0].CallbackData != "noop_1" {
			t.Errorf("nav button 0 mismatch: text=%q, callback=%q", navRow[0].Text, navRow[0].CallbackData)
		}
		if navRow[1].Text != "Вперёд ➡️" || navRow[1].CallbackData != "page_1" {
			t.Errorf("nav button 1 mismatch: text=%q, callback=%q", navRow[1].Text, navRow[1].CallbackData)
		}
	})

	t.Run("middle page has both prev and next buttons", func(t *testing.T) {
		tracks := generateTracks(25)
		kb := BuildTracksKeyboard(tracks, 1)

		if len(kb.InlineKeyboard) != 11 {
			t.Fatalf("expected 11 rows, got %d", len(kb.InlineKeyboard))
		}

		navRow := kb.InlineKeyboard[10]
<<<<<<< HEAD
=======
		// On page 1: "⬅️ Назад", "📄 2/3", "Вперёд ➡️" (3 buttons)
>>>>>>> 0c59af3e9076fadbe77438343714ce617cca44d0
		if len(navRow) != 3 {
			t.Fatalf("expected 3 nav buttons on middle page, got %d", len(navRow))
		}

		if navRow[0].Text != "⬅️ Назад" || navRow[0].CallbackData != "page_0" {
			t.Errorf("prev button mismatch: text=%q, callback=%q", navRow[0].Text, navRow[0].CallbackData)
		}
		if navRow[1].Text != "📄 2/3" || navRow[1].CallbackData != "noop_2" {
			t.Errorf("counter button mismatch: text=%q, callback=%q", navRow[1].Text, navRow[1].CallbackData)
		}
		if navRow[2].Text != "Вперёд ➡️" || navRow[2].CallbackData != "page_2" {
			t.Errorf("next button mismatch: text=%q, callback=%q", navRow[2].Text, navRow[2].CallbackData)
		}

<<<<<<< HEAD
=======
		// First track on page 1 should have callback dl_10
>>>>>>> 0c59af3e9076fadbe77438343714ce617cca44d0
		if kb.InlineKeyboard[0][0].CallbackData != "dl_10" {
			t.Errorf("expected first track on page 1 to have callback dl_10, got %q", kb.InlineKeyboard[0][0].CallbackData)
		}
	})

	t.Run("last page has prev button but no next button", func(t *testing.T) {
		tracks := generateTracks(25)
		kb := BuildTracksKeyboard(tracks, 2)

<<<<<<< HEAD
=======
		// 5 tracks on last page + 1 pagination row = 6 rows
>>>>>>> 0c59af3e9076fadbe77438343714ce617cca44d0
		if len(kb.InlineKeyboard) != 6 {
			t.Fatalf("expected 6 rows, got %d", len(kb.InlineKeyboard))
		}

		navRow := kb.InlineKeyboard[5]
<<<<<<< HEAD
=======
		// On last page: "⬅️ Назад" and "📄 3/3" (2 buttons)
>>>>>>> 0c59af3e9076fadbe77438343714ce617cca44d0
		if len(navRow) != 2 {
			t.Fatalf("expected 2 nav buttons on last page, got %d", len(navRow))
		}

		if navRow[0].Text != "⬅️ Назад" || navRow[0].CallbackData != "page_1" {
			t.Errorf("prev button mismatch: text=%q, callback=%q", navRow[0].Text, navRow[0].CallbackData)
		}
		if navRow[1].Text != "📄 3/3" || navRow[1].CallbackData != "noop_3" {
			t.Errorf("counter button mismatch: text=%q, callback=%q", navRow[1].Text, navRow[1].CallbackData)
		}
	})

	t.Run("clamps negative page index to 0", func(t *testing.T) {
		tracks := generateTracks(25)
		kb := BuildTracksKeyboard(tracks, -5)

		navRow := kb.InlineKeyboard[10]
		if navRow[0].Text != "📄 1/3" {
			t.Errorf("expected page clamped to 0 ('📄 1/3'), got %q", navRow[0].Text)
		}
	})

	t.Run("clamps oversized page index to totalPages-1", func(t *testing.T) {
		tracks := generateTracks(25)
		kb := BuildTracksKeyboard(tracks, 999)

		navRow := kb.InlineKeyboard[5]
		if navRow[1].Text != "📄 3/3" {
			t.Errorf("expected page clamped to last page ('📄 3/3'), got %q", navRow[1].Text)
		}
	})
}
