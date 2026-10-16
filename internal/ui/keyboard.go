package ui

import (
	"fmt"

	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
	"random/internal/model"
)

const TracksPerPage = 10

// MainMenuKeyboard возвращает постоянное меню внизу чата
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

// ButtonTrackText форматирует текст для инлайн-кнопки трека
func ButtonTrackText(m model.Music) string {
	text := fmt.Sprintf("%s — %s [%s]", m.Author, m.Title, m.Duration)
	runes := []rune(text)
	if len(runes) > 60 {
		return string(runes[:57]) + "..."
	}
	return text
}

// BuildTracksKeyboard формирует инлайн-клавиатуру со списком треков и кнопками перелистывания страниц (пагинации)
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
	// Треки для выбранной страницы
	for i := start; i < end; i++ {
		m := music[i]
		buttonText := ButtonTrackText(m)
		// В callback_data передаем глобальный индекс трека для скачивания
		rows = append(rows, tu.InlineKeyboardRow(
			tu.InlineKeyboardButton(buttonText).WithCallbackData(fmt.Sprintf("dl_%d", i)),
		))
	}

	// Строка переключения страниц
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
