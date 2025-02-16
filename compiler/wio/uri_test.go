package wio

import (
	"path/filepath"
	"testing"
)

func TestUriToFilePath(t *testing.T) {
	tests := []struct {
		name     string
		uri      string
		expected string
	}{
		{
			name:     "Windows file path",
			uri:      "file:///d:/dev/Golang/walrus/compiler/code/start.wal",
			expected: "d:/dev/Golang/walrus/compiler/code/start.wal",
		},
		{
			name:     "Unix file path",
			uri:      "file:///usr/local/bin",
			expected: "/usr/local/bin",
		},
		{
			name:     "Relative file URI without host",
			uri:      "file://start.wal",
			expected: filepath.FromSlash(""), // parsed.Path is empty
		},
		{
			name:     "Non-file scheme",
			uri:      "http:///usr/local/bin",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := UriToFilePath(tt.uri)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			if result != tt.expected {
				t.Errorf("UriToFilePath(%q) returned %q; expected %q", tt.uri, result, tt.expected)
			}
		})
	}
}
