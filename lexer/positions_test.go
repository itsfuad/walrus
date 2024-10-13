package lexer

import (
	"testing"
)

func TestAdvance(t *testing.T) {
	tests := []struct {
		name     string
		initial  Position
		toSkip   string
		expected Position
	}{
		{
			name:     "advance single character",
			initial:  Position{Line: 1, Column: 1, Index: 0},
			toSkip:   "a",
			expected: Position{Line: 1, Column: 2, Index: 1},
		},
		{
			name:     "advance multiple characters",
			initial:  Position{Line: 1, Column: 1, Index: 0},
			toSkip:   "abc",
			expected: Position{Line: 1, Column: 4, Index: 3},
		},
		{
			name:     "advance with newline",
			initial:  Position{Line: 1, Column: 1, Index: 0},
			toSkip:   "a\nb",
			expected: Position{Line: 2, Column: 2, Index: 3},
		},
		{
			name:     "advance with multiple newlines",
			initial:  Position{Line: 1, Column: 1, Index: 0},
			toSkip:   "a\nb\nc",
			expected: Position{Line: 3, Column: 2, Index: 5},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pos := &tt.initial
			pos.Advance(tt.toSkip)
			if pos.Line != tt.expected.Line || pos.Column != tt.expected.Column || pos.Index != tt.expected.Index {
				t.Errorf("Advance() = %v, want %v", pos, tt.expected)
			}
		})
	}
}
