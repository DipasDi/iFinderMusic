package bot

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"random/internal/config"
	"random/internal/model"
	"random/internal/scraper"
	"random/internal/shazam"
	"random/internal/storage"
	"random/internal/ui"

	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
)

type Bot struct {
	bot       *telego.Bot
	cfg       *config.Config
	shazam    *shazam.ShazamClient
	userStore *storage.UserStore
}

// New constructs and initializes a new Bot instance with the provided configuration.
func New(cfg *config.Config) (*Bot, error) {
	b, err := telego.NewBot(cfg.BotToken)
	if err != nil {
		return nil, err
	}

	var userStore *storage.UserStore
	if cfg.TrackUsers {
		var err error
		userStore, err = storage.NewUserStore(cfg.UsersFilePath)
		if err != nil {
			log.Printf("warning: could not initialize user storage at %s: %v", cfg.UsersFilePath, err)
		} else {
			log.Printf("user tracking enabled, store file: %s", cfg.UsersFilePath)
		}
	} else {
		log.Printf("user tracking is disabled by configuration")
	}

	return &Bot{
		bot:       b,
		cfg:       cfg,
		shazam:    shazam.NewClient(),
		userStore: userStore,
	}, nil
}

// trackUser logs and persists the Telegram user info in the JSON store.
func (b *Bot) trackUser(u *telego.User) {
	if b == nil || b.cfg == nil || !b.cfg.TrackUsers || b.userStore == nil || u == nil {
		return
	}
	if err := b.userStore.TrackUser(u.ID, u.Username, u.FirstName, u.LastName); err != nil {
		log.Printf("user tracking error for %d (%s): %v", u.ID, u.Username, err)
	} else {
		log.Printf("user %d (%s) tracked successfully", u.ID, u.Username)
	}
}

// Start launches the Telegram long-polling update listener loop.
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

// SendCategoryMusic caches the track list in Redis and displays the first page with inline navigation buttons.
func SendCategoryMusic(b *Bot, chatID telego.ChatID, music []model.Music, query string) {
	if err := storage.SaveUserTracks(chatID.ID, music); err != nil {
		log.Printf("redis cache error: %v\n", err)
	}
	_ = storage.SaveUserQuery(chatID.ID, query)

	totalPages := (len(music) + ui.TracksPerPage - 1) / ui.TracksPerPage
	inlineKeyboard := ui.BuildTracksKeyboard(music, 0)
	listMsg, err := ui.SendMessageMD(b.bot, chatID, ui.SearchResultsMessage(query, len(music), 0, totalPages), inlineKeyboard)

	if err == nil && listMsg != nil {
		_ = storage.SaveUserListMessageID(chatID.ID, listMsg.MessageID)
	}
}

// handleMessage processes incoming text commands, search queries, and media messages.
func (b *Bot) handleMessage(msg *telego.Message, menuKeyboard *telego.ReplyKeyboardMarkup) {
	if msg.From != nil {
		b.trackUser(msg.From)
	}
	chatID := msg.Chat.ChatID()
	text := strings.TrimSpace(msg.Text)

	var fileID string
	var InfoMessage *telego.Message
	var TextTypeAudioFile string

	switch {
	case text != "":
		switch text {
		case "/start":
			_, _ = ui.SendMessageMD(b.bot, chatID, ui.WelcomeMessage(), menuKeyboard)
		case "/help", "❓ Помощь":
			_, _ = ui.SendMessageMD(b.bot, chatID, ui.HelpMessage(), nil)
		default:
			switch text {
			case "🔥 Популярные", "Популярные":
				music := scraper.Search("", "https://ru.hitmoz.org/songs/top-today")
				SendCategoryMusic(b, chatID, music, text)
				return

			case "✨ Новинки", "Новинки":
				music := scraper.Search("", "https://ru.hitmoz.org/songs/new")
				SendCategoryMusic(b, chatID, music, text)
				return
			}

			// delete previous search results to avoid chat clutter
			oldMsgID := storage.GetUserListMessageID(chatID.ID)
			if oldMsgID != 0 {
				_ = b.bot.DeleteMessage(context.Background(), &telego.DeleteMessageParams{
					ChatID:    chatID,
					MessageID: oldMsgID,
				})
			}

			// delete the user prompt message for cleaner UI
			_ = b.bot.DeleteMessage(context.Background(), &telego.DeleteMessageParams{
				ChatID:    chatID,
				MessageID: msg.MessageID,
			})

			notifyMsg, _ := ui.SendMessageMD(b.bot, chatID, ui.SearchingMessage(text), nil)

			music := scraper.Search(text, "")

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

			SendCategoryMusic(b, chatID, music, text)
		}
		return

	case msg.Voice != nil:
		TextTypeAudioFile = "🔍 Распознаю трек из аудио... "
		fileID = msg.Voice.FileID
	case msg.Audio != nil:
		TextTypeAudioFile = "🔍 Распознаю трек из аудио... "
		fileID = msg.Audio.FileID
	case msg.Video != nil:
		TextTypeAudioFile = "🔍 Распознаю трек из видео... "
		fileID = msg.Video.FileID
	case msg.VideoNote != nil:
		TextTypeAudioFile = "🔍 Распознаю трек из видео... "
		fileID = msg.VideoNote.FileID
	case msg.Document != nil:
		TextTypeAudioFile = "🔍 Распознаю трек из документа... "
		fileID = msg.Document.FileID
	}

	if TextTypeAudioFile != "" {
		InfoMessage, _ = b.bot.SendMessage(context.Background(), &telego.SendMessageParams{
			ChatID: chatID,
			Text:   TextTypeAudioFile,
		})
	}

	if fileID != "" {
		go b.handleMedia(chatID, fileID, InfoMessage)
		return
	}
}

// handleMedia downloads incoming media, identifies the song via Shazam, and queries the catalog for MP3s.
func (b *Bot) handleMedia(chatID telego.ChatID, fileID string, InfoMessage *telego.Message) {
	data, err := b.downloadTelegramFile(fileID)
	if err != nil {
		_, _ = b.bot.SendMessage(context.Background(), &telego.SendMessageParams{
			ChatID: chatID,
			Text:   fmt.Sprintf("❌ Ошибка при скачивании файла: %v", err),
		})

		if InfoMessage != nil {
			_ = b.bot.DeleteMessage(context.Background(), &telego.DeleteMessageParams{
				ChatID:    chatID,
				MessageID: InfoMessage.MessageID,
			})
		}
		return
	}

	trackResp, err := b.shazam.Recognize(data)
	if err != nil {
		log.Printf("[SHAZAM ERROR] Ошибка распознавания: %v\n", err)
		_, _ = b.bot.SendMessage(context.Background(), &telego.SendMessageParams{
			ChatID: chatID,
			Text:   "❌ Ошибка при распознавании звука.",
		})

		if InfoMessage != nil {
			_ = b.bot.DeleteMessage(context.Background(), &telego.DeleteMessageParams{
				ChatID:    chatID,
				MessageID: InfoMessage.MessageID,
			})
		}
		return
	}

	if trackResp == nil || trackResp.Track == nil {
		matchesCount := 0
		if trackResp != nil {
			matchesCount = len(trackResp.Matches)
		}
		log.Printf("[SHAZAM] Совпадений в базе не найдено (matches: %d)\n", matchesCount)
		_, _ = b.bot.SendMessage(context.Background(), &telego.SendMessageParams{
			ChatID: chatID,
			Text:   "😔 Не удалось распознать трек. Возможно, отрывок слишком короткий (нужно от 5 секунд) или в видео много посторонних шумов.",
		})

		if InfoMessage != nil {
			_ = b.bot.DeleteMessage(context.Background(), &telego.DeleteMessageParams{
				ChatID:    chatID,
				MessageID: InfoMessage.MessageID,
			})
		}
		return
	}

	artist := trackResp.Track.Subtitle
	title := trackResp.Track.Title
	query := fmt.Sprintf("%s - %s", artist, title)

	if InfoMessage != nil {
		_ = b.bot.DeleteMessage(context.Background(), &telego.DeleteMessageParams{
			ChatID:    chatID,
			MessageID: InfoMessage.MessageID,
		})
	}

	InfoMessage, _ = b.bot.SendMessage(context.Background(), &telego.SendMessageParams{
		ChatID: chatID,
		Text:   fmt.Sprintf("🎵 Найдено: %s\n🔎 Ищу аудиозапись в каталоге...", query),
	})

	var music []model.Music
	queries := []string{
		fmt.Sprintf("%s %s", artist, title),
		title,
		artist,
	}

	for _, q := range queries {
		q = strings.TrimSpace(q)
		if q == "" {
			continue
		}
		music = scraper.Search(q, "")
		if len(music) > 0 {
			break
		}
	}

	if len(music) == 0 {
		_, _ = b.bot.SendMessage(context.Background(), &telego.SendMessageParams{
			ChatID: chatID,
			Text:   fmt.Sprintf("Музыка %s определена, но в базе для скачивания её не нашлось.", query),
		})

		if InfoMessage != nil {
			_ = b.bot.DeleteMessage(context.Background(), &telego.DeleteMessageParams{
				ChatID:    chatID,
				MessageID: InfoMessage.MessageID,
			})
		}
		return
	}

	if trackResp.Track != nil {
		shazamImage := ""
		if img, ok := trackResp.Track.Images["coverart"]; ok && img != "" {
			shazamImage = img
		} else if trackResp.Track.Share.Image != "" {
			shazamImage = trackResp.Track.Share.Image
		}
		if shazamImage != "" {
			for i := range music {
				if music[i].CoverURL == "" {
					music[i].CoverURL = shazamImage
				}
			}
		}
	}

	SendCategoryMusic(b, chatID, music, query)
}

// downloadTelegramFile downloads a file from Telegram servers into memory by its file ID.
func (b *Bot) downloadTelegramFile(fileID string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	fileInfo, err := b.bot.GetFile(ctx, &telego.GetFileParams{
		FileID: fileID,
	})
	if err != nil {
		return nil, fmt.Errorf("получение информации о файле: %w", err)
	}

	downloadURL := b.bot.FileDownloadURL(fileInfo.FilePath)
	log.Printf("[MEDIA URL] Ссылка: %s\n", downloadURL)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, downloadURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("запрос на скачивание: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP статус %s", resp.Status)
	}

	return io.ReadAll(resp.Body)
}

// handleCallback handles inline button interactions such as pagination and track download requests.
func (b *Bot) handleCallback(cq *telego.CallbackQuery) {
	b.trackUser(&cq.From)
	chatID := cq.Message.GetChat().ChatID()

	// page indicator toast
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

	// page switch
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

	_ = b.bot.AnswerCallbackQuery(context.Background(), &telego.AnswerCallbackQueryParams{
		CallbackQueryID: cq.ID,
	})

	// track download
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

// sendTrackAudio downloads the MP3 file for a track, embeds cover art & ID3 metadata, and sends it to the user.
func (b *Bot) sendTrackAudio(chatID telego.ChatID, track model.Music) error {
	file, thumbBytes, err := scraper.Download(track)
	if err != nil {
		return err
	}
	defer file.Close()
	defer os.Remove(file.Name())

	var thumbFile *os.File
	if len(thumbBytes) > 0 {
		var err error
		thumbFile, err = os.CreateTemp("", "thumb-*.jpg")
		if err == nil {
			defer thumbFile.Close()
			defer os.Remove(thumbFile.Name())
			_, _ = thumbFile.Write(thumbBytes)
			_, _ = thumbFile.Seek(0, io.SeekStart)
		}
	}

	replyMarkup := ui.YandexMusicKeyboard(track.Author, track.Title)

	audioParams := tu.Audio(chatID, tu.File(file)).
		WithTitle(track.Title).
		WithPerformer(track.Author).
		WithParseMode(telego.ModeMarkdownV2)

	if thumbFile != nil {
		thumb := tu.File(thumbFile)
		audioParams.WithThumbnail(&thumb)
	}

	if replyMarkup != nil {
		audioParams.WithReplyMarkup(replyMarkup)
	}

	_, err = b.bot.SendAudio(context.Background(), audioParams)
	if err != nil {
		log.Printf("[BOT ERROR] Ошибка SendAudio (MarkdownV2): %v, пробуем fallback", err)
		_, _ = file.Seek(0, io.SeekStart)
		plainCaption := fmt.Sprintf("🎧 %s — %s (%s)", track.Author, track.Title, track.Duration)
		fallbackParams := tu.Audio(chatID, tu.File(file)).
			WithTitle(track.Title).
			WithPerformer(track.Author).
			WithCaption(plainCaption)
		if thumbFile != nil {
			_, _ = thumbFile.Seek(0, io.SeekStart)
			thumb := tu.File(thumbFile)
			fallbackParams.WithThumbnail(&thumb)
		}
		if replyMarkup != nil {
			fallbackParams.WithReplyMarkup(replyMarkup)
		}
		_, err = b.bot.SendAudio(context.Background(), fallbackParams)
		if err != nil {
			log.Printf("[BOT ERROR] Ошибка SendAudio (fallback): %v", err)
		}
	}

	return err
}
