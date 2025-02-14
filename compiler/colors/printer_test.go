package colors

import (
	"strings"
	"testing"
)

const resetText = "Expected RESET escape code at the end, got: %s"

func TestColorPrintf(t *testing.T) {
	color := RED
	output := color.Sprintf("Hello %s", "World")
	if !strings.Contains(output, string(RED)) {
		t.Errorf("Expected RED color escape code, got: %s", output)
	}
	if !strings.Contains(output, "Hello World") {
		t.Errorf("Expected formatted string 'Hello World', got: %s", output)
	}
	if !strings.HasSuffix(output, string(RESET)) {
		t.Errorf(resetText, output)
	}
}

func TestColorPrintln(t *testing.T) {
	color := GREEN
	output := color.Sprintln("Test message")
	if !strings.Contains(output, string(GREEN)) {
		t.Errorf("Expected GREEN color escape code, got: %s", output)
	}
	if !strings.Contains(output, "Test message\n") {
		t.Errorf("Expected message 'Test message' with newline, got: %s", output)
	}
	if !strings.HasSuffix(output, string(RESET)) {
		t.Errorf(resetText, output)
	}
}

func TestColorPrint(t *testing.T) {
	color := BLUE
	output := color.Sprint("No formatting")
	if !strings.Contains(output, string(BLUE)) {
		t.Errorf("Expected BLUE color escape code, got: %s", output)
	}
	if !strings.Contains(output, "No formatting") {
		t.Errorf("Expected string 'No formatting', got: %s", output)
	}
	if !strings.HasSuffix(output, string(RESET)) {
		t.Errorf(resetText, output)
	}
}

func TestColorCombination(t *testing.T) {
	color := BOLD_RED
	output := color.Sprintf("Error: %s", "Something went wrong")
	if !strings.Contains(output, string(BOLD_RED)) {
		t.Errorf("Expected BOLD_RED color escape code, got: %s", output)
	}
	if !strings.Contains(output, "Error: Something went wrong") {
		t.Errorf("Expected formatted string 'Error: Something went wrong', got: %s", output)
	}
	if !strings.HasSuffix(output, string(RESET)) {
		t.Errorf(resetText, output)
	}
}

func TestMultipleColors(t *testing.T) {
	red := RED.Sprint("Red Text")
	green := GREEN.Sprint("Green Text")
	combined := red + " " + green

	if !strings.Contains(combined, string(RED)) || !strings.Contains(combined, "Red Text") {
		t.Errorf("Expected RED text with escape code, got: %s", combined)
	}
	if !strings.Contains(combined, string(GREEN)) || !strings.Contains(combined, "Green Text") {
		t.Errorf("Expected GREEN text with escape code, got: %s", combined)
	}
	if !strings.Contains(combined, string(RESET)) {
		t.Errorf("Expected RESET escape code in combined text, got: %s", combined)
	}
}
