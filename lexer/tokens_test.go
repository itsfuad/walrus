package lexer

import (
	"bytes"
	"log"
	"os"
	"testing"
)

func TestIsKeyword(t *testing.T) {
	tests := []struct {
		token    string
		expected bool
	}{
		{"let", true},
		{"const", true},
		{"if", true},
		{"else", true},
		{"for", true},
		{"type", true},
		{"priv", true},
		{"embed", true},
		{"trait", true},
		{"impl", true},
		{"struct", true},
		{"fn", true},
		{"ret", true},
		{"in", true},
		{"unknown", false},
	}

	for _, test := range tests {
		result := IsKeyword(test.token)
		if result != test.expected {
			t.Errorf("IsKeyword(%s) = %v; want %v", test.token, result, test.expected)
		}
	}
}

func TestTokenDebug(t *testing.T) {
	token := NewToken(LET_TOKEN, "let", Position{Line: 1, Column: 1}, Position{Line: 1, Column: 4})
	expectedOutput := "testfile:1:1\t 'let'\n"

	// Redirect stdout to capture the output
	var output bytes.Buffer
	log.SetOutput(&output)
	defer func() { log.SetOutput(os.Stderr) }()

	token.Debug("testfile")

	if output.String() != expectedOutput {
		t.Errorf("Token.Debug() output = %v; want %v", output, expectedOutput)
	}
}

func TestNewToken(t *testing.T) {
	kind := LET_TOKEN
	value := "let"
	start := Position{Line: 1, Column: 1}
	end := Position{Line: 1, Column: 4}

	token := NewToken(kind, value, start, end)

	if token.Kind != kind {
		t.Errorf("NewToken().Kind = %v; want %v", token.Kind, kind)
	}
	if token.Value != value {
		t.Errorf("NewToken().Value = %v; want %v", token.Value, value)
	}
	if token.Start != start {
		t.Errorf("NewToken().Start = %v; want %v", token.Start, start)
	}
	if token.End != end {
		t.Errorf("NewToken().End = %v; want %v", token.End, end)
	}
}
