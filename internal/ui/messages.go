package ui

import (
	"fmt"

	"random/internal/model"
)

// WelcomeMessage returns the MarkdownV2 welcome text for the /start command.
func WelcomeMessage() string {
	return "👋 *Добро пожаловать\\!*\n\n" +
		"Я — твой бот для быстрого поиска и скачивания музыки 🎧\n\n" +
		"📌 *Что я умею:*\n" +
		"> 🔎 Находить треки по названию и артисту\n" +
		"> ⚡️ Скачивать аудио высокого качества в MP3\n" +
		"> 📂 Быстро открывать популярные треки и новинки\n\n" +
		"💡 *Как начать:*\n" +
		"Просто напиши мне имя исполнителя или название песни, например:\n" +
		"`The Weeknd` или `Miyagi`\n\n" +
		"Либо воспользуйся быстрыми кнопками ниже 👇"
}

// HelpMessage returns the MarkdownV2 help guide for the /help command.
func HelpMessage() string {
	return "📖 *Справка по использованию бота*\n\n" +
		"🎧 *Поиск треков:*\n" +
		"> Напиши в чат название песни или имя исполнителя\\. Бот подберёт список подходящих вариантов с длительностью\\.\n\n" +
		"📥 *Скачивание:*\n" +
		"> Нажми на интерактивную кнопку с нужным треком, и бот сразу же отправит тебе аудиофайл\\.\n\n" +
		"⚡️ *Быстрые кнопки:*\n" +
		"> • *🔥 Популярные* — треки в тренде\n" +
		"> • *✨ Новинки* — свежие релизы\n\n" +
		"⚠️ _Если трек не скачивается, возможно, истёк срок временной сессии — просто повтори поиск\\._"
}

// SearchingMessage returns a status text displayed while search is in progress.
func SearchingMessage(query string) string {
	return fmt.Sprintf("🔎 *Ищу треки по запросу:* `%s`\n\n_Пожалуйста, подождите немного\\.\\.\\._ ⏳", EscapeCode(query))
}

// NotFoundMessage returns the text shown when no tracks match the query.
func NotFoundMessage(query string) string {
	return fmt.Sprintf("🔍 *Поиск по запросу:* `%s`\n\n"+
		"😔 *К сожалению, ничего не найдено\\!*\n\n"+
		"> 💡 *Советы:*\n"+
		"> • Проверьте правильность написания\n"+
		"> • Попробуйте указать только имя исполнителя\n"+
		"> • Попробуйте написать запрос латиницей",
		EscapeCode(query))
}

// SearchResultsMessage returns the header text for paginated track search results.
func SearchResultsMessage(query string, count int, page int, totalPages int) string {
	if totalPages <= 1 {
		return fmt.Sprintf("🎵 *Результаты поиска:* `%s`\n\n"+
			"Найдено треков: *%d*\n\n"+
			"> _Нажмите на кнопку с треком ниже, чтобы скачать аудиофайл:_ ",
			EscapeCode(query), count)
	}

	return fmt.Sprintf("🎵 *Результаты поиска:* `%s`\n\n"+
		"Найдено треков: *%d* \\(страница *%d* из *%d*\\)\n\n"+
		"> _Нажмите на кнопку с треком ниже, чтобы скачать аудиофайл:_ ",
		EscapeCode(query), count, page+1, totalPages)
}

// TrackLoadingMessage returns the status text shown when a track is being downloaded and prepared.
func TrackLoadingMessage(track model.Music) string {
	return fmt.Sprintf("⏳ *Загрузка трека\\.\\.\\.*\n\n"+
		"🎵 *Песня:* `%s`\n"+
		"👤 *Артист:* `%s`\n"+
		"⏱ *Длительность:* `%s`\n\n"+
		"_Файл скачивается и отправляется в чат\\.\\.\\._ 🚀",
		EscapeCode(track.Title), EscapeCode(track.Author), EscapeCode(track.Duration))
}

// ErrorMessage formats an error alert card in MarkdownV2.
func ErrorMessage(title, description string) string {
	return fmt.Sprintf("⚠️ *%s*\n\n> %s", EscapeMD(title), EscapeMD(description))
}
