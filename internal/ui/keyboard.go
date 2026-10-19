package ui

import (
	"fmt"
	"net/url"
	"strings"

	"random/internal/model"

	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
)

const TracksPerPage = 10

// MainMenuKeyboard creates the persistent bottom reply keyboard with quick actions.
func MainMenuKeyboard() *telego.ReplyKeyboardMarkup {
	return tu.Keyboard(
		tu.KeyboardRow(
			tu.KeyboardButton("🔥 Популярные"),
			tu.KeyboardButton("✨ Новинки"),
		),
		tu.KeyboardRow(
			tu.KeyboardButton("❓ Помощь"),
		),
	).WithResizeKeyboard()
}

// ButtonTrackText formats artist, title and duration for inline buttons (max 60 runes).
func ButtonTrackText(m model.Music) string {
	text := fmt.Sprintf("%s — %s [%s]", m.Author, m.Title, m.Duration)
	runes := []rune(text)
	if len(runes) > 60 {
		return string(runes[:57]) + "..."
	}
	return text
}

// BuildTracksKeyboard generates an inline keyboard with paginated track selection buttons and navigation.
func BuildTracksKeyboard(music []model.Music, page int) *telego.InlineKeyboardMarkup {
	totalTracks := len(music)
	if totalTracks == 0 {
		return tu.InlineKeyboard()
	}

	totalPages := (totalTracks + TracksPerPage - 1) / TracksPerPage
	if page < 0 {
		page = 0
	}
	if page >= totalPages {
		page = totalPages - 1
	}

	start := page * TracksPerPage
	end := start + TracksPerPage
	if end > totalTracks {
		end = totalTracks
	}

	var rows [][]telego.InlineKeyboardButton
	for i := start; i < end; i++ {
		m := music[i]
		buttonText := ButtonTrackText(m)
		// pass global track index in callback data for download
		rows = append(rows, tu.InlineKeyboardRow(
			tu.InlineKeyboardButton(buttonText).WithCallbackData(fmt.Sprintf("dl_%d", i)),
		))
	}

	if totalPages > 1 {
		var navRow []telego.InlineKeyboardButton

		if page > 0 {
			navRow = append(navRow, tu.InlineKeyboardButton("⬅️ Назад").WithCallbackData(fmt.Sprintf("page_%d", page-1)))
		}

		navRow = append(navRow, tu.InlineKeyboardButton(fmt.Sprintf("📄 %d/%d", page+1, totalPages)).WithCallbackData(fmt.Sprintf("noop_%d", page+1)))

		if page < totalPages-1 {
			navRow = append(navRow, tu.InlineKeyboardButton("Вперёд ➡️").WithCallbackData(fmt.Sprintf("page_%d", page+1)))
		}

		rows = append(rows, navRow)
	}

	return tu.InlineKeyboard(rows...)
}

// YandexMusicKeyboard creates an inline button opening Yandex Music search for the track.
func YandexMusicKeyboard(author, title string) *telego.InlineKeyboardMarkup {
	query := strings.TrimSpace(author + " " + title)
	if query == "" {
		return nil
	}
	searchURL := "https://music.yandex.ru/search?text=" + url.QueryEscape(query)
	return tu.InlineKeyboard(
		tu.InlineKeyboardRow(
			tu.InlineKeyboardButton("🎧 Яндекс Музыка").WithURL(searchURL),
		),
	)
}
