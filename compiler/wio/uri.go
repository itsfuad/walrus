package wio

import (
	"net/url"
	"path/filepath"
)

// UriToFilePath converts a file:// URI to a platform-specific file path.
func UriToFilePath(uri string) (string, error) {
	// Parse the URI
	parsed, err := url.Parse(uri)
	if err != nil {
		return "", err
	}

	// Handle Windows paths (e.g., file:///C:/path/to/file)
	if parsed.Path[0] == '/' && len(parsed.Path) > 3 && parsed.Path[2] == ':' {
		parsed.Path = parsed.Path[1:]
	}

	// Decode URL-encoded characters (e.g., %3A -> :)
	filePath, err := url.PathUnescape(parsed.Path)
	if err != nil {
		return "", err
	}

	// Convert to platform-specific file path
	return filepath.FromSlash(filePath), nil
}
