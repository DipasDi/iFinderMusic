package ui

import (
	"testing"
)

func TestEscapeMD(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "plain text without special chars",
			input:    "Hello World 123 Привет Мир",
			expected: "Hello World 123 Привет Мир",
		},
		{
			name:     "escapes dot and exclamation",
			input:    "Hello, World. Test!",
			expected: "Hello, World\\. Test\\!",
		},
		{
			name:     "escapes brackets and parentheses",
			input:    "track [remix] (feat. artist)",
			expected: "track \\[remix\\] \\(feat\\. artist\\)",
		},
		{
			name:     "escapes arithmetic and symbols",
			input:    "2 + 2 = 4 - 0, 100% > 50%",
			expected: "2 \\+ 2 \\= 4 \\- 0, 100% \\> 50%",
		},
		{
			name:     "escapes backslash and backtick",
			input:    "path\\to\\file and `code`",
			expected: "path\\\\to\\\\file and \\`code\\`",
		},
		{
			name:     "escapes markdown formatting symbols",
			input:    "_italic_ *bold* ~strike~ #tag |pipe| {brace}",
			expected: "\\_italic\\_ \\*bold\\* \\~strike\\~ \\#tag \\|pipe\\| \\{brace\\}",
		},
		{
			name:     "realistic music track title",
			input:    "Queen - Bohemian Rhapsody (Live at Wembley '86) [05:55]!",
			expected: "Queen \\- Bohemian Rhapsody \\(Live at Wembley '86\\) \\[05:55\\]\\!",
		},
		{
			name:     "all markdown v2 reserved characters",
			input:    "_*[]()~`>#+-={}.!\\",
			expected: "\\_\\*\\[\\]\\(\\)\\~\\`\\>\\#\\+\\-\\=\\{\\}\\.\\!\\\\",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := EscapeMD(tt.input)
			if actual != tt.expected {
				t.Errorf("EscapeMD(%q)\n got:  %q\n want: %q", tt.input, actual, tt.expected)
			}
		})
	}
}

func TestEscapeCode(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "plain text with regular characters",
			input:    "Queen - Bohemian Rhapsody (Live)",
			expected: "Queen - Bohemian Rhapsody (Live)",
		},
		{
			name:     "escapes backticks inside code",
			input:    "run `rm -rf` command",
			expected: "run \\`rm -rf\\` command",
		},
		{
			name:     "escapes backslashes",
			input:    `C:\Music\Track.mp3`,
			expected: `C:\\Music\\Track.mp3`,
		},
		{
			name:     "escapes combination of backticks and backslashes",
			input:    "`echo \\n`",
			expected: "\\`echo \\\\n\\`",
		},
		{
			name:     "does not escape other markdown symbols like * _ [ ]",
			input:    "*bold* _italic_ [link]",
			expected: "*bold* _italic_ [link]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := EscapeCode(tt.input)
			if actual != tt.expected {
				t.Errorf("EscapeCode(%q) = %q, want %q", tt.input, actual, tt.expected)
			}
		})
	}
}

func TestCleanMarkdownForFallback(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "string without escape slashes remains unchanged",
			input:    "Hello World 123",
			expected: "Hello World 123",
		},
		{
			name:     "removes escapes for markdown v2 symbols",
			input:    "\\_italic\\_ \\*bold\\* \\[link\\] \\(parens\\) \\~strike\\~ \\`code\\` \\>quote \\#tag \\+plus \\-minus \\=eq \\|pipe \\{brace\\} \\.dot \\!excl \\\\slash",
			expected: "_italic_ *bold* [link] (parens) ~strike~ `code` >quote #tag +plus -minus =eq |pipe {brace} .dot !excl \\slash",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := CleanMarkdownForFallback(tt.input)
			if actual != tt.expected {
				t.Errorf("CleanMarkdownForFallback(%q)\n got:  %q\n want: %q", tt.input, actual, tt.expected)
			}
		})
	}

	t.Run("round trip: CleanMarkdownForFallback reverses EscapeMD", func(t *testing.T) {
		samples := []string{
			"Hello, World. Test!",
			"track [remix] (feat. artist)",
			"2 + 2 = 4 - 0, 100% > 50%",
			"_italic_ *bold* ~strike~ #tag |pipe| {brace}",
			"Queen - Bohemian Rhapsody (Live at Wembley '86) [05:55]!",
		}

		for _, sample := range samples {
			escaped := EscapeMD(sample)
			cleaned := CleanMarkdownForFallback(escaped)
			if cleaned != sample {
				t.Errorf("round trip failed for %q:\n escaped: %q\n cleaned: %q", sample, escaped, cleaned)
			}
		}
	})
}
