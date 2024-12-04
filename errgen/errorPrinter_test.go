package errgen

import (
	"testing"
)

func TestAddHint(t *testing.T) {
	err := &WalrusError{}
	err.AddHint("This is a hint", TEXT_HINT)

	if len(err.hints) != 1 {
		t.Errorf("Expected 1 hint, got %d", len(err.hints))
	}

	if err.hints[0].message != "This is a hint" {
		t.Errorf("Expected hint message 'This is a hint', got '%s'", err.hints[0].message)
	}

	if err.hints[0].hintType != TEXT_HINT {
		t.Errorf("Expected hint type TEXT_HINT, got %d", err.hints[0].hintType)
	}
}

func TestMakeError(t *testing.T) {
	err := makeError("test.go", 1, 1, 1, 1, "Test error", ERROR_CRITICAL)

	if err.filePath != "test.go" {
		t.Errorf("Expected filePath 'test.go', got '%s'", err.filePath)
	}

	if err.lineStart != 1 {
		t.Errorf("Expected lineStart 1, got %d", err.lineStart)
	}

	if err.lineEnd != 1 {
		t.Errorf("Expected lineEnd 1, got %d", err.lineEnd)
	}

	if err.colStart != 1 {
		t.Errorf("Expected colStart 1, got %d", err.colStart)
	}

	if err.colEnd != 1 {
		t.Errorf("Expected colEnd 1, got %d", err.colEnd)
	}

	if err.err.Error() != "Test error" {
		t.Errorf("Expected error message 'Test error', got '%s'", err.err.Error())
	}

	if err.level != ERROR_CRITICAL {
		t.Errorf("Expected error level ERROR_CRITICAL, got %d", err.level)
	}
}