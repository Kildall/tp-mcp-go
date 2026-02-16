package query

import "testing"

func TestEscapeValue(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "normal string no change",
			input:    "hello world",
			expected: "hello world",
		},
		{
			name:     "single quote doubled",
			input:    "it's",
			expected: "it''s",
		},
		{
			name:     "multiple quotes",
			input:    "she's reading 'book'",
			expected: "she''s reading ''book''",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := EscapeValue(tt.input)
			if result != tt.expected {
				t.Errorf("EscapeValue(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestFormatStringCondition(t *testing.T) {
	result := FormatStringCondition("Name", "=", "O'Brien")
	expected := "Name = 'O''Brien'"
	if result != expected {
		t.Errorf("FormatStringCondition() = %q, want %q", result, expected)
	}
}

func TestFormatNumberCondition(t *testing.T) {
	result := FormatNumberCondition("Id", ">", 100)
	expected := "Id > 100"
	if result != expected {
		t.Errorf("FormatNumberCondition() = %q, want %q", result, expected)
	}
}
