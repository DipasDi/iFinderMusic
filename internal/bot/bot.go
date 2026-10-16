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

// Bot инкапсулирует клиента Telego и конфигурацию
type Bot struct {
	bot *telego.Bot
	cfg *config.Config
}

// New создает новый экземпляр бота
func New(cfg *config.Config) (*Bot, error) {
	b, err := telego.NewBot(cfg.BotToken, telego.WithDefaultDebugLogger())
	if err != nil {
		return nil, err
	}
	return &Bot{bot: b, cfg: cfg}, nil
}

// Start запускает long polling обработчик обновлений
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

// Отправка списка треков по категории
func SendCategoryMusic(b *Bot, chatID telego.ChatID, music []model.Music, query string) {
	// Кэшируем треки и запрос в Redis
	if err := storage.SaveUserTracks(chatID.ID, music); err != nil {
		fmt.Println("Ошибка сохранения треков в Redis:", err)
	}
	_ = storage.SaveUserQuery(chatID.ID, query)

	totalPages := (len(music) + ui.TracksPerPage - 1) / ui.TracksPerPage
	inlineKeyboard := ui.BuildTracksKeyboard(music, 0)
	listMsg, err := ui.SendMessageMD(b.bot, chatID, ui.SearchResultsMessage(query, len(music), 0, totalPages), inlineKeyboard)

	// Запоминаем ID сообщения со списком треков для последующей очистки
	if err == nil && listMsg != nil {
		_ = storage.SaveUserListMessageID(chatID.ID, listMsg.MessageID)
	}
}

// handleMessage обрабатывает текстовые сообщения от пользователей
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
			//Отправляем найденые треки из категории популярные
			music := scraper.Search("", "https://ru.hitmoz.org/songs/top-today")

			// Отправляем треки
			SendCategoryMusic(b, chatID, music, text)

			return
		case "✨ Новинки", "Новинки":
			//Отправляем найденые треки из категории новинки
			music := scraper.Search("", "https://ru.hitmoz.org/songs/new")

			// Отправляем треки
			SendCategoryMusic(b, chatID, music, text)
			return
		}

		// Удаляем предыдущий список треков из чата, если он был сохранен в Redis
		oldMsgID := storage.GetUserListMessageID(chatID.ID)
		if oldMsgID != 0 {
			_ = b.bot.DeleteMessage(context.Background(), &telego.DeleteMessageParams{
				ChatID:    chatID,
				MessageID: oldMsgID,
			})
		}

		// Удаляем сообщение пользователя для чистоты чата
		_ = b.bot.DeleteMessage(context.Background(), &telego.DeleteMessageParams{
			ChatID:    chatID,
			MessageID: msg.MessageID,
		})

		// Отправляем уведомление о начале поиска
		notifyMsg, _ := ui.SendMessageMD(b.bot, chatID, ui.SearchingMessage(text), nil)

		// Выполняем поиск через парсер
		music := scraper.Search(text, "")

		// Удаляем временное уведомление о поиске
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

		// Отправляем треки
		SendCategoryMusic(b, chatID, music, text)
	}
}

// handleCallback обрабатывает нажатия инлайн-кнопок
func (b *Bot) handleCallback(cq *telego.CallbackQuery) {
	chatID := cq.Message.GetChat().ChatID()

	// Клик по индикатору страницы (показывает подсказку)
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

	// Перелистывание страниц (пагинация)
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

	// Ответ на нажатие кнопки скачивания трека
	_ = b.bot.AnswerCallbackQuery(context.Background(), &telego.AnswerCallbackQueryParams{
		CallbackQueryID: cq.ID,
	})

	// Скачивание трека
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

		// Уведомление о начале скачивания
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

// sendTrackAudio скачивает трек и отправляет аудиофайл в Telegram
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
		// Fallback: при ошибке парсинга разметки отправляем обычный текст
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
