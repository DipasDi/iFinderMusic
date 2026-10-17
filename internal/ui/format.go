package ui

import (
	"context"
	"strings"

	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
)

// EscapeMD escapes characters reserved by Telegram MarkdownV2
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

// EscapeCode escapes backticks and backslashes inside `...` code blocks
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

// CleanMarkdownForFallback strips escaping slashes if telegram rejects markdown and we need to retry as plain text
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

func SendMessageMD(bot *telego.Bot, chatID telego.ChatID, text string, replyMarkup telego.ReplyMarkup) (*telego.Message, error) {
	params := tu.Message(chatID, text).WithParseMode(telego.ModeMarkdownV2)
	if replyMarkup != nil {
		params.WithReplyMarkup(replyMarkup)
	}

	msg, err := bot.SendMessage(context.Background(), params)
	if err != nil {
		// fallback to plain text if telegram rejected markdown entities
		plainParams := tu.Message(chatID, CleanMarkdownForFallback(text))
		if replyMarkup != nil {
			plainParams.WithReplyMarkup(replyMarkup)
		}
		return bot.SendMessage(context.Background(), plainParams)
	}

	return msg, nil
}

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
