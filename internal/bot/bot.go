package bot

import (
	"context"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"random/internal/config"
	"random/internal/model"
	"random/internal/scraper"
	"random/internal/storage"
	"random/internal/ui"

	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
)

// Bot encapsulates the Telego bot instance and application configuration.
type Bot struct {
	bot *telego.Bot
	cfg *config.Config
}

// New creates and initializes a new Bot instance.
func New(cfg *config.Config) (*Bot, error) {
	if cfg.BotToken == "" {
		return nil, fmt.Errorf("bot token is empty: please set BOT_TOKEN in .env or environment")
	}

	b, err := telego.NewBot(cfg.BotToken, telego.WithDefaultDebugLogger())
	if err != nil {
		return nil, err
	}
	return &Bot{bot: b, cfg: cfg}, nil
}

// Start begins listening for Telegram updates via long polling.
func (b *Bot) Start() error {
	updates, err := b.bot.UpdatesViaLongPolling(context.Background(), nil)
	if err != nil {
		return err
	}

	menuKeyboard := ui.MainMenuKeyboard()

	for update := range updates {
		if update.Message != nil {
			b.handleMessage(update.Message, menuKeyboard)
		}

		if update.CallbackQuery != nil {
			b.handleCallback(update.CallbackQuery)
		}
	}

	return nil
}

// SendCategoryMusic caches tracks and sends a paginated track list message to the chat.
func SendCategoryMusic(b *Bot, chatID telego.ChatID, music []model.Music, query string) {
	// Cache tracks and search query in Redis
	if err := storage.SaveUserTracks(chatID.ID, music); err != nil {
		fmt.Println("Failed to cache tracks in Redis:", err)
	}
	_ = storage.SaveUserQuery(chatID.ID, query)

	totalPages := (len(music) + ui.TracksPerPage - 1) / ui.TracksPerPage
	inlineKeyboard := ui.BuildTracksKeyboard(music, 0)
	listMsg, err := ui.SendMessageMD(b.bot, chatID, ui.SearchResultsMessage(query, len(music), 0, totalPages), inlineKeyboard)

	// Save track list message ID for subsequent chat cleanup
	if err == nil && listMsg != nil {
		_ = storage.SaveUserListMessageID(chatID.ID, listMsg.MessageID)
	}
}

// handleMessage processes incoming text messages from users.
func (b *Bot) handleMessage(msg *telego.Message, menuKeyboard *telego.ReplyKeyboardMarkup) {
	chatID := msg.Chat.ChatID()
	text := strings.TrimSpace(msg.Text)

	switch text {
	case "/start":
		_, _ = ui.SendMessageMD(b.bot, chatID, ui.WelcomeMessage(), menuKeyboard)
	case "/help", "❓ Помощь":
		_, _ = ui.SendMessageMD(b.bot, chatID, ui.HelpMessage(), nil)
	default:
		if text == "" {
			return
		}

		switch text {
		case "🔥 Популярные", "Популярные":
			// Fetch and send popular tracks
			music := scraper.Search("", "https://ru.hitmoz.org/songs/top-today")
			SendCategoryMusic(b, chatID, music, text)
			return
		case "✨ Новинки", "Новинки":
			// Fetch and send new releases
			music := scraper.Search("", "https://ru.hitmoz.org/songs/new")
			SendCategoryMusic(b, chatID, music, text)
			return
		}

		// Delete previous track list from chat if recorded in Redis
		oldMsgID := storage.GetUserListMessageID(chatID.ID)
		if oldMsgID != 0 {
			_ = b.bot.DeleteMessage(context.Background(), &telego.DeleteMessageParams{
				ChatID:    chatID,
				MessageID: oldMsgID,
			})
		}

		// Delete user message to maintain a clean chat interface
		_ = b.bot.DeleteMessage(context.Background(), &telego.DeleteMessageParams{
			ChatID:    chatID,
			MessageID: msg.MessageID,
		})

		// Send temporary searching indicator
		notifyMsg, _ := ui.SendMessageMD(b.bot, chatID, ui.SearchingMessage(text), nil)

		// Execute search query via scraper
		music := scraper.Search(text, "")

		// Remove temporary search indicator
		if notifyMsg != nil && notifyMsg.MessageID != 0 {
			_ = b.bot.DeleteMessage(context.Background(), &telego.DeleteMessageParams{
				ChatID:    chatID,
				MessageID: notifyMsg.MessageID,
			})
		}

		if len(music) == 0 {
			_, _ = ui.SendMessageMD(b.bot, chatID, ui.NotFoundMessage(text), nil)
			return
		}

		// Send search results
		SendCategoryMusic(b, chatID, music, text)
	}
}

// handleCallback processes inline button callback queries.
func (b *Bot) handleCallback(cq *telego.CallbackQuery) {
	chatID := cq.Message.GetChat().ChatID()

	// Page counter indicator click (displays a popup alert)
	if strings.HasPrefix(cq.Data, "noop") {
		p := strings.TrimPrefix(cq.Data, "noop_")
		toastText := "Текущая страница"
		if p != "" {
			toastText = fmt.Sprintf("Страница %s", p)
		}
		_ = b.bot.AnswerCallbackQuery(context.Background(), &telego.AnswerCallbackQueryParams{
			CallbackQueryID: cq.ID,
			Text:            toastText,
		})
		return
	}

	// Handle pagination page navigation
	if strings.HasPrefix(cq.Data, "page_") {
		_ = b.bot.AnswerCallbackQuery(context.Background(), &telego.AnswerCallbackQueryParams{
			CallbackQueryID: cq.ID,
		})

		targetPage, err := strconv.Atoi(strings.TrimPrefix(cq.Data, "page_"))
		if err != nil {
			return
		}

		tracks, err := storage.GetUserTracks(chatID.ID)
		if err != nil || len(tracks) == 0 {
			_, _ = ui.SendMessageMD(b.bot, chatID, ui.ErrorMessage("Треки устарели", "Срок действия кэша для этого поиска истёк. Пожалуйста, выполните поиск заново."), nil)
			return
		}

		totalPages := (len(tracks) + ui.TracksPerPage - 1) / ui.TracksPerPage
		if targetPage < 0 || targetPage >= totalPages {
			return
		}

		query := storage.GetUserQuery(chatID.ID)
		if query == "" {
			query = "Музыка"
		}

		newKeyboard := ui.BuildTracksKeyboard(tracks, targetPage)
		newText := ui.SearchResultsMessage(query, len(tracks), targetPage, totalPages)

		_, _ = ui.EditMessageMD(b.bot, chatID, cq.Message.GetMessageID(), newText, newKeyboard)
		return
	}

	// Acknowledge download button callback query
	_ = b.bot.AnswerCallbackQuery(context.Background(), &telego.AnswerCallbackQueryParams{
		CallbackQueryID: cq.ID,
	})

	// Handle track download action
	if strings.HasPrefix(cq.Data, "dl_") {
		idx, err := strconv.Atoi(strings.TrimPrefix(cq.Data, "dl_"))
		if err != nil {
			_, _ = ui.SendMessageMD(b.bot, chatID, ui.ErrorMessage("Ошибка обработки запроса", "Не удалось определить трек. Пожалуйста, выполните поиск заново."), nil)
			return
		}

		tracks, err := storage.GetUserTracks(chatID.ID)
		if err != nil || idx < 0 || idx >= len(tracks) {
			_, _ = ui.SendMessageMD(b.bot, chatID, ui.ErrorMessage("Трек устарел", "Срок действия кэша для этого поиска истёк. Пожалуйста, выполните поиск заново."), nil)
			return
		}

		track := tracks[idx]

		// Send download in-progress notification
		loadingMsg, _ := ui.SendMessageMD(b.bot, chatID, ui.TrackLoadingMessage(track), nil)

		loadingMsgID := 0
		if loadingMsg != nil {
			loadingMsgID = loadingMsg.MessageID
		}

		go func(cID telego.ChatID, t model.Music, lMsgID int) {
			err := b.sendTrackAudio(cID, t)
			if err != nil {
				_, _ = ui.SendMessageMD(b.bot, cID, ui.ErrorMessage("Ошибка при отправке трека", err.Error()+"\nПопробуйте выбрать другой вариант или повторите попытку позже."), nil)
			}
			if lMsgID != 0 {
				_ = b.bot.DeleteMessage(context.Background(), &telego.DeleteMessageParams{
					ChatID:    cID,
					MessageID: lMsgID,
				})
			}
		}(chatID, track, loadingMsgID)
	}
}

// sendTrackAudio downloads and sends the audio file to the user.
func (b *Bot) sendTrackAudio(chatID telego.ChatID, track model.Music) error {
	file, err := scraper.Download(track)
	if err != nil {
		return err
	}
	defer file.Close()
	defer os.Remove(file.Name())

	caption := ui.AudioCaption(track)
	audioParams := tu.Audio(chatID, tu.File(file)).
		WithTitle(track.Title).
		WithPerformer(track.Author).
		WithCaption(caption).
		WithParseMode(telego.ModeMarkdownV2)

	_, err = b.bot.SendAudio(context.Background(), audioParams)
	if err != nil {
		// Fallback: send with plain text caption if markdown parsing fails
		_, _ = file.Seek(0, io.SeekStart)
		plainCaption := fmt.Sprintf("🎧 %s — %s (%s)", track.Author, track.Title, track.Duration)
		fallbackParams := tu.Audio(chatID, tu.File(file)).
			WithTitle(track.Title).
			WithPerformer(track.Author).
			WithCaption(plainCaption)
		_, err = b.bot.SendAudio(context.Background(), fallbackParams)
	}

	return err
}
