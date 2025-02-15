package analyzer

import (
	"testing"
)

func TestAnalyze(t *testing.T) {
	filePath := "./../compiler/code/test.wal"
	reports, err := Analyze(filePath, false, false, false)
	if err != nil {
		t.Error(err)
	}
	if reports != nil {
		reports.DisplayAll()
	}
}
