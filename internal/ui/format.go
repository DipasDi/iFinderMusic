package ui

import (
	"context"
	"strings"

	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
)

// EscapeMD экранирует специальные символы Telegram MarkdownV2
func EscapeMD(s string) string {
	var sb strings.Builder
	for _, r := range s {
		switch r {
		case '_', '*', '[', ']', '(', ')', '~', '`', '>', '#', '+', '-', '=', '|', '{', '}', '.', '!', '\\':
			sb.WriteRune('\\')
		}
		sb.WriteRune(r)
	}
	return sb.String()
}

// EscapeCode экранирует символы внутри инлайн-кода `...` в MarkdownV2
func EscapeCode(s string) string {
	var sb strings.Builder
	for _, r := range s {
		switch r {
		case '`', '\\':
			sb.WriteRune('\\')
		}
		sb.WriteRune(r)
	}
	return sb.String()
}

// CleanMarkdownForFallback удаляет экранирование для fallback-отправки обычным текстом
func CleanMarkdownForFallback(s string) string {
	replacer := strings.NewReplacer(
		"\\_", "_",
		"\\*", "*",
		"\\[", "[",
		"\\]", "]",
		"\\(", "(",
		"\\)", ")",
		"\\~", "~",
		"\\`", "`",
		"\\>", ">",
		"\\#", "#",
		"\\+", "+",
		"\\-", "-",
		"\\=", "=",
		"\\|", "|",
		"\\{", "{",
		"\\}", "}",
		"\\.", ".",
		"\\!", "!",
		"\\\\", "\\",
	)
	return replacer.Replace(s)
}

// SendMessageMD отправляет сообщение с разметкой MarkdownV2 (с fallback к обычному тексту)
func SendMessageMD(bot *telego.Bot, chatID telego.ChatID, text string, replyMarkup telego.ReplyMarkup) (*telego.Message, error) {
	params := tu.Message(chatID, text).WithParseMode(telego.ModeMarkdownV2)
	if replyMarkup != nil {
		params.WithReplyMarkup(replyMarkup)
	}

	msg, err := bot.SendMessage(context.Background(), params)
	if err != nil {
		plainParams := tu.Message(chatID, CleanMarkdownForFallback(text))
		if replyMarkup != nil {
			plainParams.WithReplyMarkup(replyMarkup)
		}
		return bot.SendMessage(context.Background(), plainParams)
	}

	return msg, nil
}

// EditMessageMD редактирует сообщение с разметкой MarkdownV2 (с fallback к обычному тексту)
func EditMessageMD(bot *telego.Bot, chatID telego.ChatID, messageID int, text string, replyMarkup *telego.InlineKeyboardMarkup) (*telego.Message, error) {
	params := tu.EditMessageText(chatID, messageID, text).
		WithParseMode(telego.ModeMarkdownV2).
		WithReplyMarkup(replyMarkup)

	msg, err := bot.EditMessageText(context.Background(), params)
	if err != nil {
		if strings.Contains(err.Error(), "message is not modified") {
			return nil, nil
		}
		plainParams := tu.EditMessageText(chatID, messageID, CleanMarkdownForFallback(text)).
			WithReplyMarkup(replyMarkup)
		return bot.EditMessageText(context.Background(), plainParams)
	}

	return msg, nil
}
