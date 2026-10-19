package scraper

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"random/internal/model"

	"github.com/gocolly/colly/v2"
)

var defaultHeader = func(r *colly.Request) {
	r.Headers.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	r.Headers.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8")
	r.Headers.Set("Accept-Language", "ru-RU,ru;q=0.9,en-US;q=0.8,en;q=0.7")
}

// Search scrapes tracks matching the query or fetches tracks from the specified href URL.
func Search(query string, href string) []model.Music {
	var music []model.Music
	c := colly.NewCollector(
		colly.AllowedDomains("ru.hitmoz.org"),
	)

	c.OnRequest(defaultHeader)

	c.OnError(func(r *colly.Response, err error) {
		fmt.Printf("scrape error [%d]: %v\n", r.StatusCode, err)
	})

	c.OnHTML(".tracks__item", func(e *colly.HTMLElement) {
		title := strings.TrimSpace(e.ChildText(".track__title"))
		artist := strings.TrimSpace(e.ChildText(".track__desc"))
		duration := strings.TrimSpace(e.ChildText(".track__fulltime"))
		downloadLink := strings.TrimSpace(e.ChildAttr(".track__download-btn", "href"))

		var coverURL string
		meta := e.Attr("data-musmeta")
		if meta != "" {
			var obj struct {
				Img string `json:"img"`
			}
			if err := json.Unmarshal([]byte(meta), &obj); err == nil && obj.Img != "" {
				coverURL = strings.TrimSpace(obj.Img)
			}
		}
		if coverURL == "" {
			style := e.ChildAttr(".track__img", "style")
			if idx := strings.Index(style, "url('"); idx != -1 {
				sub := style[idx+5:]
				if end := strings.Index(sub, "')"); end != -1 {
					coverURL = strings.TrimSpace(sub[:end])
				}
			}
		}
		if coverURL != "" && !strings.HasPrefix(coverURL, "http") {
			coverURL = "https://ru.hitmoz.org" + coverURL
		}

		if title != "" && artist != "" {
			music = append(music, model.Music{
				Title:        title,
				Author:       artist,
				Duration:     duration,
				DownloadLink: downloadLink,
				CoverURL:     coverURL,
			})
		}
	})

	if href == "" {
		targetURL := "https://ru.hitmoz.org/search?q=" + url.QueryEscape(query)
		_ = c.Visit(targetURL)
	} else {
		_ = c.Visit(href)
	}

	return music
}

// Download streams the audio file for the given track into a temporary MP3 file on disk,
// embeds ID3 metadata & cover art, and returns the audio file along with thumbnail image bytes.
func Download(track model.Music) (*os.File, []byte, error) {
	rawURL := track.DownloadLink
	if !strings.HasPrefix(rawURL, "http") {
		rawURL = "https://ru.hitmoz.org" + rawURL
	}

	tempFile, err := os.CreateTemp("", "track-*.mp3")
	if err != nil {
		return nil, nil, err
	}

	req, err := http.NewRequest("GET", rawURL, nil)
	if err != nil {
		tempFile.Close()
		_ = os.Remove(tempFile.Name())
		return nil, nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Language", "ru-RU,ru;q=0.9,en-US;q=0.8,en;q=0.7")
	req.Header.Set("Referer", "https://ru.hitmoz.org/")

	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		tempFile.Close()
		_ = os.Remove(tempFile.Name())
		return nil, nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		tempFile.Close()
		_ = os.Remove(tempFile.Name())
		return nil, nil, fmt.Errorf("download server returned %s", resp.Status)
	}

	if _, err := io.Copy(tempFile, resp.Body); err != nil {
		tempFile.Close()
		_ = os.Remove(tempFile.Name())
		return nil, nil, err
	}

	tempPath := tempFile.Name()
	// Close write handle before tagging to avoid file locks on Windows
	_ = tempFile.Close()

	thumbBytes, err := EmbedMetadata(tempPath, track)
	if err != nil {
		log.Printf("[SCRAPER WARN] Ошибка записи ID3 метаданных: %v", err)
	}

	reopened, err := os.Open(tempPath)
	if err != nil {
		_ = os.Remove(tempPath)
		return nil, nil, err
	}

	return reopened, thumbBytes, nil
}
