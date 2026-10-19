package scraper

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"random/internal/model"
)

<<<<<<< HEAD
// TestDownload_Success verifies successful audio download from a mock HTTP server.
=======
>>>>>>> 0c59af3e9076fadbe77438343714ce617cca44d0
func TestDownload_Success(t *testing.T) {
	fakeAudioBytes := []byte("ID3_FAKE_MP3_STREAM_CONTENT_123456789")

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Referer") != "https://ru.hitmoz.org/" {
			http.Error(w, "bad referer", http.StatusBadRequest)
			return
		}
		if r.Header.Get("User-Agent") == "" {
			http.Error(w, "missing user agent", http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(fakeAudioBytes)
	}))
	defer server.Close()

	track := model.Music{
		Title:        "Test Title",
		Author:       "Test Artist",
		DownloadLink: server.URL + "/track.mp3",
	}

	tempFile, err := Download(track)
	if err != nil {
		t.Fatalf("Download failed: %v", err)
	}
	defer func() {
		_ = tempFile.Close()
		_ = os.Remove(tempFile.Name())
	}()

	data, err := io.ReadAll(tempFile)
	if err != nil {
		t.Fatalf("read failed: %v", err)
	}

	if string(data) != string(fakeAudioBytes) {
		t.Errorf("got %q, want %q", string(data), string(fakeAudioBytes))
	}
}

<<<<<<< HEAD
// TestDownload_HttpErrors verifies proper error handling for 404 and 500 HTTP responses.
=======
>>>>>>> 0c59af3e9076fadbe77438343714ce617cca44d0
func TestDownload_HttpErrors(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/notfound" {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		http.Error(w, "server error", http.StatusInternalServerError)
	}))
	defer server.Close()

	t.Run("404 error", func(t *testing.T) {
		track := model.Music{
			DownloadLink: server.URL + "/notfound",
		}

		file, err := Download(track)
		if err == nil {
			if file != nil {
				_ = file.Close()
				_ = os.Remove(file.Name())
			}
			t.Fatal("expected error on 404, got nil")
		}
	})

	t.Run("500 error", func(t *testing.T) {
		track := model.Music{
			DownloadLink: server.URL + "/error",
		}

		file, err := Download(track)
		if err == nil {
			if file != nil {
				_ = file.Close()
				_ = os.Remove(file.Name())
			}
			t.Fatal("expected error on 500, got nil")
		}
	})
}

<<<<<<< HEAD
// TestDownload_InvalidUrl verifies that attempting to download with an invalid URL fails gracefully.
=======
>>>>>>> 0c59af3e9076fadbe77438343714ce617cca44d0
func TestDownload_InvalidUrl(t *testing.T) {
	track := model.Music{
		DownloadLink: "http://invalid-url\x7f/file.mp3",
	}

	file, err := Download(track)
	if err == nil {
		if file != nil {
			_ = file.Close()
			_ = os.Remove(file.Name())
		}
		t.Fatal("expected error for invalid url, got nil")
	}
}
