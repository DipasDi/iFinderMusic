package scraper

import (
	"bytes"
	"image"
	"image/color"
	"image/jpeg"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"random/internal/model"

	"github.com/bogem/id3v2/v2"
)

// createDummyJPEG creates a synthetic JPEG image with given width and height for testing.
func createDummyJPEG(w, h int) []byte {
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			img.Set(x, y, color.RGBA{R: 200, G: 100, B: 50, A: 255})
		}
	}
	var buf bytes.Buffer
	_ = jpeg.Encode(&buf, img, &jpeg.Options{Quality: 80})
	return buf.Bytes()
}

func TestPrepareThumbnail_Valid(t *testing.T) {
	raw := createDummyJPEG(600, 400)
	thumb, err := PrepareThumbnail(raw)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	decoded, format, err := image.Decode(bytes.NewReader(thumb))
	if err != nil {
		t.Fatalf("decode thumbnail: %v", err)
	}
	if format != "jpeg" {
		t.Errorf("expected format jpeg, got %s", format)
	}

	bounds := decoded.Bounds()
	if bounds.Dx() > 320 || bounds.Dy() > 320 {
		t.Errorf("expected dimensions <= 320, got %dx%d", bounds.Dx(), bounds.Dy())
	}
}

func TestPrepareThumbnail_EmptyOrCorrupt(t *testing.T) {
	_, err := PrepareThumbnail(nil)
	if err == nil {
		t.Error("expected error for nil bytes, got nil")
	}

	_, err = PrepareThumbnail([]byte("not_an_image"))
	if err == nil {
		t.Error("expected error for invalid image bytes, got nil")
	}
}

func TestEmbedMetadata_WithRemoteCover(t *testing.T) {
	coverBytes := createDummyJPEG(200, 200)
	coverServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(coverBytes)
	}))
	defer coverServer.Close()

	tempFile, err := os.CreateTemp("", "meta-*.mp3")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tempFile.Name())
	_, _ = tempFile.Write([]byte("FAKE_MP3_STREAM_FOR_METADATA_TEST"))
	tempFile.Close()

	track := model.Music{
		Title:    "Test Song",
		Author:   "Test Artist",
		CoverURL: coverServer.URL + "/cover.jpg",
	}

	thumb, err := EmbedMetadata(tempFile.Name(), track)
	if err != nil {
		t.Fatalf("EmbedMetadata failed: %v", err)
	}
	if len(thumb) == 0 {
		t.Error("expected non-empty thumbnail bytes")
	}

	// Verify ID3 tags written
	tag, err := id3v2.Open(tempFile.Name(), id3v2.Options{Parse: true})
	if err != nil {
		t.Fatalf("failed to re-open tag: %v", err)
	}
	defer tag.Close()

	if tag.Artist() != "Test Artist" {
		t.Errorf("expected artist 'Test Artist', got %q", tag.Artist())
	}
	if tag.Title() != "Test Song" {
		t.Errorf("expected title 'Test Song', got %q", tag.Title())
	}

	frame := tag.GetLastFrame("APIC")
	if frame == nil {
		t.Fatal("expected APIC picture frame in tag, got nil")
	}
	pf, ok := frame.(id3v2.PictureFrame)
	if !ok || len(pf.Picture) == 0 {
		t.Fatal("expected non-empty PictureFrame")
	}
}

func TestEmbedMetadata_ExtractExistingCover(t *testing.T) {
	tempFile, err := os.CreateTemp("", "meta-existing-*.mp3")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tempFile.Name())
	_, _ = tempFile.Write([]byte("FAKE_MP3_WITH_PRESET_COVER"))
	tempFile.Close()

	presetJPEG := createDummyJPEG(100, 100)

	// Pre-seed with APIC
	tag, err := id3v2.Open(tempFile.Name(), id3v2.Options{Parse: true})
	if err != nil {
		t.Fatal(err)
	}
	tag.AddAttachedPicture(id3v2.PictureFrame{
		Encoding:    id3v2.EncodingUTF8,
		MimeType:    "image/jpeg",
		PictureType: id3v2.PTFrontCover,
		Picture:     presetJPEG,
	})
	_ = tag.Save()
	tag.Close()

	// EmbedMetadata without remote CoverURL
	track := model.Music{
		Title:  "Track With Local Cover",
		Author: "Local Artist",
	}

	thumb, err := EmbedMetadata(tempFile.Name(), track)
	if err != nil {
		t.Fatalf("EmbedMetadata failed: %v", err)
	}
	if len(thumb) == 0 {
		t.Error("expected thumbnail extracted from existing APIC frame, got empty")
	}
}
