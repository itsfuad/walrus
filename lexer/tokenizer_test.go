package lexer

import (
	"os"
	"testing"
)

func TestTokenize(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []Token
	}{
		{
			name:  "Empty file",
			input: "",
			expected: []Token{
				NewToken(EOF_TOKEN, "eof", Position{Line: 1, Column: 1, Index: 0}, Position{Line: 1, Column: 1, Index: 0}),
			},
		},
		{
			name:  "Single line comment",
			input: "// this is a comment",
			expected: []Token{
				NewToken(EOF_TOKEN, "eof", Position{Line: 1, Column: 21, Index: 20}, Position{Line: 1, Column: 21, Index: 20}),
			},
		},
		{
			name:  "Simple identifier",
			input: "var",
			expected: []Token{
				NewToken(IDENTIFIER_TOKEN, "var", Position{Line: 1, Column: 1, Index: 0}, Position{Line: 1, Column: 4, Index: 3}),
				NewToken(EOF_TOKEN, "eof", Position{Line: 1, Column: 4, Index: 3}, Position{Line: 1, Column: 4, Index: 3}),
			},
		},
		{
			name:  "String literal",
			input: `"hello"`,
			expected: []Token{
				NewToken(STR, "hello", Position{Line: 1, Column: 1, Index: 0}, Position{Line: 1, Column: 8, Index: 7}),
				NewToken(EOF_TOKEN, "eof", Position{Line: 1, Column: 8, Index: 7}, Position{Line: 1, Column: 8, Index: 7}),
			},
		},
		{
			name:  "Number literal",
			input: "123",
			expected: []Token{
				NewToken(INT32, "123", Position{Line: 1, Column: 1, Index: 0}, Position{Line: 1, Column: 4, Index: 3}),
				NewToken(EOF_TOKEN, "eof", Position{Line: 1, Column: 4, Index: 3}, Position{Line: 1, Column: 4, Index: 3}),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpfile := createTempFile(t, tt.input)
			defer os.Remove(tmpfile.Name())

			tokens := Tokenize(tmpfile.Name(), false)
			compareTokens(t, tokens, tt.expected)
		})
	}
}

func createTempFile(t *testing.T, content string) *os.File {
	tmpfile, err := os.CreateTemp("", "testfile")
	if err != nil {
		t.Fatal(err)
	}

	if _, err := tmpfile.Write([]byte(content)); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	return tmpfile
}

func compareTokens(t *testing.T, tokens, expected []Token) {
	if len(tokens) != len(expected) {
		t.Fatalf("expected %d tokens, got %d", len(expected), len(tokens))
	}

	for i, token := range tokens {
		if token.Kind != expected[i].Kind || token.Value != expected[i].Value {
			t.Errorf("expected token %v, got %v", expected[i], token)
		}
	}
}
