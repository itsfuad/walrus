package wio

import (
	"path/filepath"
	"testing"
)

func TestUriToFilePath(t *testing.T) {
	tests := []struct {
		name      string
		uri       string
		wantPath  string
		wantError bool
	}{
		{
			name:     "Windows file URI",
			uri:      "file:///C:/path/to/file",
			wantPath: filepath.FromSlash("C:/path/to/file"),
		},
		{
			name:     "Unix file URI",
			uri:      "file:///home/user/file.txt",
			wantPath: filepath.FromSlash("/home/user/file.txt"),
		},
		{
			name:     "URI with URL-encoded characters",
			uri:      "file:///C:/Program%20Files/App",
			wantPath: filepath.FromSlash("C:/Program Files/App"),
		},
		{
			name:      "Invalid URI",
			uri:       "://invalid",
			wantError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := UriToFilePath(tc.uri)
			if tc.wantError {
				if err == nil {
					t.Errorf("Expected error for URI %q, got nil", tc.uri)
				}
				return
			}
			if err != nil {
				t.Fatalf("Unexpected error for URI %q: %v", tc.uri, err)
			}
			if got != tc.wantPath {
				t.Errorf("uriToFilePath(%q) = %q, want %q", tc.uri, got, tc.wantPath)
			}
		})
	}
}
