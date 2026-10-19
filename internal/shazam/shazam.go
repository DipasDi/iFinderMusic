package shazam

import (
	"context"
	"log"

	"github.com/alexgorbatchev/goshazam"
	"github.com/alexgorbatchev/goshazam/pkg/models"
)

// ShazamClient wraps the goshazam music recognition service.
type ShazamClient struct {
	client *goshazam.Shazam
}

// NewClient creates and returns a new ShazamClient instance.
func NewClient() *ShazamClient {
	return &ShazamClient{
		client: goshazam.New(),
	}
}

// Recognize processes raw audio or video bytes and queries Shazam for song identification.
func (c *ShazamClient) Recognize(audio []byte) (*models.ResponseTrack, error) {
	log.Printf("[SHAZAM] Начало распознавания (входной размер: %d байт)\n", len(audio))
	resp, err := c.client.Recognize(context.Background(), audio)
	if err != nil {
		log.Printf("[SHAZAM] Ошибка Recognize: %v\n", err)
		return nil, err
	}
	if resp != nil {
		trackName := "nil"
		if resp.Track != nil {
			trackName = resp.Track.Subtitle + " - " + resp.Track.Title
		}
		log.Printf("[SHAZAM] Результат: matches=%d, track=%s\n", len(resp.Matches), trackName)
	}
	return resp, nil
}

