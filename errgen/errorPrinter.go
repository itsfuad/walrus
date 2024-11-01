package errgen

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"walrus/utils"
)

type HINT int

const (
	TEXT_HINT HINT = iota
	CODE_HINT
)

type Hint struct {
	message  string
	hintType HINT
}

type WalrusError struct {
	filePath  string
	lineStart int
	lineEnd   int
	colStart  int
	colEnd    int
	err       error
	hints     []Hint
}

func (e *WalrusError) DisplayWithPanic() {
	PrintError(e)
	panic("")
}

func PrintError(e *WalrusError) {
	fileData, err := os.ReadFile(e.filePath)
	if err != nil {
		panic(err)
	}
	utils.ColorPrint(utils.GREY, fmt.Sprintf("\nIn file: %s:%d:%d\n", e.filePath, e.lineStart, e.colStart))
	lines := strings.Split(string(fileData), "\n")
	line := lines[e.lineStart-1]
	hLen := 0
	if e.lineStart == e.lineEnd {
		hLen = (e.colEnd - e.colStart) - 1
	} else {
		//full line
		hLen = len(line) - 2
	}
	if hLen < 0 {
		hLen = 0
	}
	fmt.Println(line)
	underLine := fmt.Sprintf("%s^%s", strings.Repeat(" ", e.colStart-1), strings.Repeat("~", hLen))

	utils.ColorPrint(utils.RED, underLine)
	utils.ColorPrint(utils.RED, e.err.Error())
	for i, hint := range e.hints {
		if i == 0 {
			utils.ColorPrint(utils.GREEN, "Hint:")
		}
		if hint.hintType == TEXT_HINT {
			utils.ColorPrint(utils.YELLOW, hint.message)
		} else {
			utils.ColorPrint(utils.ORANGE, hint.message)
		}
	}
}

func (e *WalrusError) AddHint(msg string, htype HINT) *WalrusError {

	if msg == "" {
		return e
	}

	e.hints = append(e.hints, Hint{
		message:  msg,
		hintType: htype,
	})

	return e
}

func MakeError(filePath string, lineStart, lineEnd int, colStart, colEnd int, err string) *WalrusError {
	if lineStart < 1 {
		lineStart = 1
	}
	if lineEnd < 1 {
		lineEnd = 1
	}
	if colStart < 1 {
		colStart = 1
	}
	if colEnd < 1 {
		colEnd = 1
	}
	return &WalrusError{
		filePath:  filePath,
		lineStart: lineStart,
		lineEnd:   lineEnd,
		colStart:  colStart,
		colEnd:    colEnd,
		err:       errors.New(err),
	}
}

var globalErrors = make([]WalrusError, 0)

// make an errorlist to add all errors and display later
func AddError(filePath string, lineStart, lineEnd int, colStart, colEnd int, err string) *WalrusError {
	errItem := MakeError(filePath, lineStart, lineEnd, colStart, colEnd, err)
	globalErrors = append(globalErrors, *errItem)
	utils.ColorPrint(utils.YELLOW, fmt.Sprintf("Error added. %d errors available\n", len(globalErrors)))
	return errItem
}

func DisplayErrors() {
	if len(globalErrors) == 0 {
		utils.ColorPrint(utils.GREEN, "------- No errors --------")
		return
	}

	
	for _, err := range globalErrors {
		PrintError(&err)
	}
	utils.ColorPrint(utils.BOLD_RED, fmt.Sprintf("%d error(s) found\n", len(globalErrors)))
	os.Exit(-1)
}
