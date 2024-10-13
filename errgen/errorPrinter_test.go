package errgen

import (
	"errors"
	"testing"
)

func TestMakeError(t *testing.T) {
	tests := []struct {
		filePath  string
		lineStart int
		lineEnd   int
		colStart  int
		colEnd    int
		errMsg    string
		expected  *WalrusError
	}{
		{
			filePath:  "testfile.go",
			lineStart: 10,
			lineEnd:   10,
			colStart:  5,
			colEnd:    15,
			errMsg:    "test error",
			expected: &WalrusError{
				filePath:  "testfile.go",
				lineStart: 10,
				lineEnd:   10,
				colStart:  5,
				colEnd:    15,
				err:       errors.New("test error"),
			},
		},
		{
			filePath:  "testfile.go",
			lineStart: -1,
			lineEnd:   -1,
			colStart:  -1,
			colEnd:    -1,
			errMsg:    "test error",
			expected: &WalrusError{
				filePath:  "testfile.go",
				lineStart: 1,
				lineEnd:   1,
				colStart:  1,
				colEnd:    1,
				err:       errors.New("test error"),
			},
		},
	}

	for _, tt := range tests {
		result := MakeError(tt.filePath, tt.lineStart, tt.lineEnd, tt.colStart, tt.colEnd, tt.errMsg)
		if result.filePath != tt.expected.filePath ||
			result.lineStart != tt.expected.lineStart ||
			result.lineEnd != tt.expected.lineEnd ||
			result.colStart != tt.expected.colStart ||
			result.colEnd != tt.expected.colEnd ||
			result.err.Error() != tt.expected.err.Error() {
			t.Errorf("MakeError(%s, %d, %d, %d, %d, %s) = %+v; expected %+v",
				tt.filePath, tt.lineStart, tt.lineEnd, tt.colStart, tt.colEnd, tt.errMsg, result, tt.expected)
		}
	}
}
