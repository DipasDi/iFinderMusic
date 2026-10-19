package scraper

import (
	"bytes"
	"fmt"
	"image"
	_ "image/gif"
	"image/jpeg"
	_ "image/png"
	"io"
	"log"
	"net/http"
	"time"

	"random/internal/model"

	"github.com/bogem/id3v2/v2"
)

var coverHTTPClient = &http.Client{
	Timeout: 5 * time.Second,
}

// resizeThumbnail scales the image down so that its width and height do not exceed maxDim,
// preserving the original aspect ratio.
func resizeThumbnail(img image.Image, maxDim int) image.Image {
	bounds := img.Bounds()
	w, h := bounds.Dx(), bounds.Dy()
	if w <= 0 || h <= 0 {
		return img
	}
	if w <= maxDim && h <= maxDim {
		return img
	}

	var newW, newH int
	if w > h {
		newW = maxDim
		newH = (h * maxDim) / w
	} else {
		newH = maxDim
		newW = (w * maxDim) / h
	}
	if newW <= 0 {
		newW = 1
	}
	if newH <= 0 {
		newH = 1
	}

	dst := image.NewRGBA(image.Rect(0, 0, newW, newH))
	for y := 0; y < newH; y++ {
		for x := 0; x < newW; x++ {
			srcX := bounds.Min.X + (x*w)/newW
			srcY := bounds.Min.Y + (y*h)/newH
			dst.Set(x, y, img.At(srcX, srcY))
		}
	}
	return dst
}

// PrepareThumbnail converts the raw image bytes into a JPEG conforming to Telegram limits
// (JPEG format, width/height <= 320, size < 200 KB).
func PrepareThumbnail(raw []byte) ([]byte, error) {
	if len(raw) == 0 {
		return nil, fmt.Errorf("empty raw image data")
	}

	img, _, err := image.Decode(bytes.NewReader(raw))
	if err != nil {
		return nil, fmt.Errorf("decode image: %w", err)
	}

	scaled := resizeThumbnail(img, 320)
	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, scaled, &jpeg.Options{Quality: 85}); err != nil {
		return nil, fmt.Errorf("encode jpeg: %w", err)
	}

	return buf.Bytes(), nil
}

// fetchCoverBytes downloads the image bytes from the given URL.
func fetchCoverBytes(coverURL string) ([]byte, error) {
	if coverURL == "" {
		return nil, fmt.Errorf("empty cover url")
	}

	req, err := http.NewRequest(http.MethodGet, coverURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")
	req.Header.Set("Referer", "https://ru.hitmoz.org/")

	resp, err := coverHTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("cover download returned status %d", resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}

// EmbedMetadata tags the MP3 file at filePath with artist, title, and cover art (ID3v2 APIC),
// and returns the thumbnail bytes suitable for Telegram SendAudio.
func EmbedMetadata(filePath string, track model.Music) ([]byte, error) {
	tag, err := id3v2.Open(filePath, id3v2.Options{Parse: true})
	if err != nil {
		return nil, fmt.Errorf("id3 open: %w", err)
	}
	defer tag.Close()

	if track.Author != "" {
		tag.SetArtist(track.Author)
	}
	if track.Title != "" {
		tag.SetTitle(track.Title)
	}

	var thumbBytes []byte

	if track.CoverURL != "" {
		rawCover, err := fetchCoverBytes(track.CoverURL)
		if err != nil {
			log.Printf("[SCRAPER WARN] Не удалось скачать обложку по ссылке %q: %v", track.CoverURL, err)
		} else if len(rawCover) > 0 {
			tag.AddAttachedPicture(id3v2.PictureFrame{
				Encoding:    id3v2.EncodingUTF8,
				MimeType:    "image/jpeg",
				PictureType: id3v2.PTFrontCover,
				Description: "Front Cover",
				Picture:     rawCover,
			})

			tb, err := PrepareThumbnail(rawCover)
			if err == nil {
				thumbBytes = tb
			}
		}
	}

	// If no thumbnail was generated from CoverURL, check if MP3 already has an embedded APIC frame
	if len(thumbBytes) == 0 {
		frame := tag.GetLastFrame("APIC")
		if frame != nil {
			if pf, ok := frame.(id3v2.PictureFrame); ok && len(pf.Picture) > 0 {
				tb, err := PrepareThumbnail(pf.Picture)
				if err == nil {
					thumbBytes = tb
				}
			}
		}
	}

	if err := tag.Save(); err != nil {
		return thumbBytes, fmt.Errorf("id3 save: %w", err)
	}

	return thumbBytes, nil
}
