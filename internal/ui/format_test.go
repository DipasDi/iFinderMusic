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
