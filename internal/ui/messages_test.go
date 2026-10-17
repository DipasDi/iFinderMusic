package ui

import (
	"strings"
	"testing"

	"random/internal/model"
)

func TestWelcomeMessage(t *testing.T) {
	msg := WelcomeMessage()
	if msg == "" {
		t.Fatal("WelcomeMessage() returned empty string")
	}
	if !strings.Contains(msg, "Добро пожаловать") {
		t.Errorf("WelcomeMessage() expected to contain 'Добро пожаловать', got: %s", msg)
	}
}

func TestHelpMessage(t *testing.T) {
	msg := HelpMessage()
	if msg == "" {
		t.Fatal("HelpMessage() returned empty string")
	}
	if !strings.Contains(msg, "Справка") {
		t.Errorf("HelpMessage() expected to contain 'Справка', got: %s", msg)
	}
}

func TestSearchingMessage(t *testing.T) {
	t.Run("escapes query with backticks", func(t *testing.T) {
		query := "`Special` Song"
		msg := SearchingMessage(query)

		expectedEscaped := "\\`Special\\` Song"
		if !strings.Contains(msg, expectedEscaped) {
			t.Errorf("SearchingMessage(%q) did not properly escape backticks, got: %s", query, msg)
		}
	})
}

func TestNotFoundMessage(t *testing.T) {
	query := "NonExistingSong `xyz`"
	msg := NotFoundMessage(query)

	if !strings.Contains(msg, "ничего не найдено") {
		t.Errorf("NotFoundMessage() expected to mention 'ничего не найдено', got: %s", msg)
	}
	if !strings.Contains(msg, "\\`xyz\\`") {
		t.Errorf("NotFoundMessage() did not escape backticks in query, got: %s", msg)
	}
}

func TestSearchResultsMessage(t *testing.T) {
	t.Run("single page or zero pages omits page indicator", func(t *testing.T) {
		msg := SearchResultsMessage("Rock", 5, 0, 1)

		if strings.Contains(msg, "страница") {
			t.Errorf("expected no page indicator when totalPages <= 1, got: %s", msg)
		}
		if !strings.Contains(msg, "Найдено треков: *5*") {
			t.Errorf("expected count *5*, got: %s", msg)
		}
		if !strings.Contains(msg, "`Rock`") {
			t.Errorf("expected query in code quotes `Rock`, got: %s", msg)
		}
	})

	t.Run("multi-page results include page indicator", func(t *testing.T) {
		msg := SearchResultsMessage("Pop", 50, 1, 5)

		expectedPageInfo := `(страница *2* из *5*)`
		if !strings.Contains(msg, `\(страница *2* из *5*\)`) {
			t.Errorf("expected %q in message, got: %s", expectedPageInfo, msg)
		}
		if !strings.Contains(msg, "Найдено треков: *50*") {
			t.Errorf("expected count *50*, got: %s", msg)
		}
	})
}

func TestTrackLoadingMessage(t *testing.T) {
	track := model.Music{
		Title:    "Track `One`",
		Author:   `Band\Artist`,
		Duration: "03:45",
	}

	msg := TrackLoadingMessage(track)

	if !strings.Contains(msg, "Track \\`One\\`") {
		t.Errorf("expected title to escape backticks, got: %s", msg)
	}
	if !strings.Contains(msg, "Band\\\\Artist") {
		t.Errorf("expected author to escape backslashes, got: %s", msg)
	}
	if !strings.Contains(msg, "03:45") {
		t.Errorf("expected duration in message, got: %s", msg)
	}
}

func TestAudioCaption(t *testing.T) {
	track := model.Music{
		Title:    "Song (Live at Wembley)",
		Author:   "Queen & Friends",
		Duration: "04:15",
	}

	caption := AudioCaption(track)

	if !strings.Contains(caption, `\(Live at Wembley\)`) {
		t.Errorf("expected escaped parentheses in caption, got: %s", caption)
	}
	if !strings.Contains(caption, "Queen & Friends") {
		t.Errorf("expected author in caption, got: %s", caption)
	}
}

func TestErrorMessage(t *testing.T) {
	title := "Error: Not Found!"
	desc := "Could not fetch track [id=123]."

	msg := ErrorMessage(title, desc)

	if !strings.Contains(msg, "Not Found\\!") {
		t.Errorf("expected escaped exclamation mark in title, got: %s", msg)
	}
	if !strings.Contains(msg, "\\[id\\=123\\]") {
		t.Errorf("expected escaped brackets and equal sign in desc, got: %s", msg)
	}
}
