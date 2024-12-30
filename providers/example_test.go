package providers_example

import "testing"

func TestSayHello(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "with empty name",
			input:    "",
			expected: "Hello, World!",
		},
		{
			name:     "with provided name",
			input:    "Alice",
			expected: "Hello, Alice!",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SayHello(tt.input)
			if got != tt.expected {
				t.Errorf("SayHello(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}
