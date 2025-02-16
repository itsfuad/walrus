package analyzer

import (
	"testing"
)

func TestAnalyze(t *testing.T) {
	filePath := `d:\dev\Golang\walrus\compiler\code\start.wal`
	reports, err := Analyze(filePath, false, false, false)
	if err != nil {
		t.Error(err)
	}
	if reports != nil {
		reports.DisplayAll()
	}
}
