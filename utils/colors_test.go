package utils

import (
	"bytes"
	"fmt"
	"os"
	"testing"
)

// TestColorize checks that the Colorize function adds the correct ANSI escape sequences around the text.
func TestColorize(t *testing.T) {
	tests := []struct {
		color    COLOR
		text     string
		expected string
	}{
		{RED, "hello", fmt.Sprintf("%shello%s", RED, RESET)},
		{GREEN, "world", fmt.Sprintf("%sworld%s", GREEN, RESET)},
		{BOLD_RED, "error", fmt.Sprintf("%serror%s", BOLD_RED, RESET)},
		{BLUE, "info", fmt.Sprintf("%sinfo%s", BLUE, RESET)},
		{ORANGE, "warning", fmt.Sprintf("%swarning%s", ORANGE, RESET)},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("Color: %s Text: %s", tt.color, tt.text), func(t *testing.T) {
			result := Colorize(tt.color, tt.text)
			if result != tt.expected {
				t.Errorf("got %q, want %q", result, tt.expected)
			}
		})
	}
}

// TestInvalidColor checks that Colorize panics if an invalid color code is passed.
func TestInvalidColor(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic for invalid color, but did not panic")
		}
	}()

	invalidColor := COLOR("\033[999m") // Invalid ANSI code
	Colorize(invalidColor, "Invalid")
}

// TestColorPrint captures stdout and checks the output of ColorPrint.
func TestColorPrint(t *testing.T) {
	// Create a pipe to capture os.Stdout output
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("Failed to create pipe: %v", err)
	}

	// Save the original os.Stdout
	stdout := os.Stdout
	defer func() { os.Stdout = stdout }() // Restore os.Stdout after the test

	// Redirect os.Stdout to the write end of the pipe
	os.Stdout = w

	// Call the function to be tested
	ColorPrint(GREEN, "output")

	// Close the writer and restore os.Stdout
	w.Close()

	// Read the captured output from the read end of the pipe
	var buf bytes.Buffer
	_, err = buf.ReadFrom(r)
	if err != nil {
		t.Fatalf("Failed to read from pipe: %v", err)
	}

	// Compare the output
	expected := fmt.Sprintf("%s%s%s\n", GREEN, "output", RESET) // Note: fmt.Println adds a newline
	actual := buf.String()

	if expected != actual {
		t.Errorf("got %q, want %q", actual, expected)
	}
}
