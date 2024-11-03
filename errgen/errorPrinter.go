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
	fmt.Printf("Total hints: %d\n", len(e.hints))
	DisplayErrors()
}

func PrintError(e *WalrusError, showFileName bool) {
	fileData, err := os.ReadFile(e.filePath)
	if err != nil {
		panic(err)
	}

	if showFileName {
		utils.ColorPrint(utils.BLUE, fmt.Sprintf("\nIn file: %s:%d:%d\n", e.filePath, e.lineStart, e.colStart))
	}

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

	lineNumber := fmt.Sprintf("%d | ", e.lineStart)
	utils.ColorPrint(utils.GREY, lineNumber)
	fmt.Println(line)
	underLine := fmt.Sprintf("%s^%s\n", strings.Repeat(" ", (e.colStart-1) + len(lineNumber)), strings.Repeat("~", hLen))

	utils.ColorPrint(utils.RED, underLine)
	utils.ColorPrint(utils.RED, e.err.Error() + "\n")

	if len(e.hints) > 0 {
		utils.ColorPrint(utils.GREEN, "Hint:\n")
		for _, hint := range e.hints {
			if hint.hintType == TEXT_HINT {
				utils.ColorPrint(utils.YELLOW, hint.message + "\n")
			} else {
				utils.ColorPrint(utils.ORANGE, hint.message + "\n")
			}
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

	fmt.Printf("Hint added. %d hints available\n", len(e.hints))

	return e
}

func makeError(filePath string, lineStart, lineEnd int, colStart, colEnd int, errMsg string) *WalrusError {
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

	err := &WalrusError{
		filePath:  filePath,
		lineStart: lineStart,
		lineEnd:   lineEnd,
		colStart:  colStart,
		colEnd:    colEnd,
		err:       errors.New(errMsg),
	}

	globalErrors = append(globalErrors, err)

	return err
}

//global errors are arrays of error pointers
var globalErrors []*WalrusError

// make an errorlist to add all errors and display later
func AddError(filePath string, lineStart, lineEnd int, colStart, colEnd int, err string) *WalrusError {
	errItem := makeError(filePath, lineStart, lineEnd, colStart, colEnd, err)
	utils.ColorPrint(utils.YELLOW, fmt.Sprintf("Error added. %d errors available\n", len(globalErrors)))
	return errItem
}

func DisplayErrors() {
	if len(globalErrors) == 0 {
		utils.ColorPrint(utils.GREEN, "------- No errors --------\n")
		return
	}
	for _, err := range globalErrors {
		PrintError(err, true)
	}
	utils.ColorPrint(utils.BOLD_RED, fmt.Sprintf("%d error(s) found\n", len(globalErrors)))
	os.Exit(-1)
}
