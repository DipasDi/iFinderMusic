package ui

import (
	"fmt"

	"random/internal/model"
)

// WelcomeMessage returns the welcome greeting for the /start command.
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

// HelpMessage returns user instructions and help guide.
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

// SearchingMessage returns a notification indicating an active search in progress.
func SearchingMessage(query string) string {
	return fmt.Sprintf("🔎 *Ищу треки по запросу:* `%s`\n\n_Пожалуйста, подождите немного\\.\\.\\._ ⏳", EscapeCode(query))
}

// NotFoundMessage returns a message when no tracks match the query.
func NotFoundMessage(query string) string {
	return fmt.Sprintf("🔍 *Поиск по запросу:* `%s`\n\n"+
		"😔 *К сожалению, ничего не найдено\\!*\n\n"+
		"> 💡 *Советы:*\n"+
		"> • Проверьте правильность написания\n"+
		"> • Попробуйте указать только имя исполнителя\n"+
		"> • Попробуйте написать запрос латиницей",
		EscapeCode(query))
}

// SearchResultsMessage returns the search results header with pagination details.
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

// TrackLoadingMessage returns the download progress card for a track.
func TrackLoadingMessage(track model.Music) string {
	return fmt.Sprintf("⏳ *Загрузка трека\\.\\.\\.*\n\n"+
		"🎵 *Песня:* `%s`\n"+
		"👤 *Артист:* `%s`\n"+
		"⏱ *Длительность:* `%s`\n\n"+
		"_Файл скачивается и отправляется в чат\\.\\.\\._ 🚀",
		EscapeCode(track.Title), EscapeCode(track.Author), EscapeCode(track.Duration))
}

// AudioCaption returns the caption attached to the uploaded MP3 audio file.
func AudioCaption(track model.Music) string {
	return fmt.Sprintf("🎧 *%s* — *%s*\n⏱ _Длительность: %s_\n\n✨ _Приятного прослушивания\\! 🎵_",
		EscapeMD(track.Author), EscapeMD(track.Title), EscapeMD(track.Duration))
}

// ErrorMessage returns a formatted error notification message.
func ErrorMessage(title, description string) string {
	return fmt.Sprintf("⚠️ *%s*\n\n> %s", EscapeMD(title), EscapeMD(description))
}
