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

type ERROR_LEVEL int

const (
	ERROR_CRITICAL ERROR_LEVEL = iota
	ERROR_NORMAL
	WARNING
	INFO
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
	level     ERROR_LEVEL
}

func PrintError(e *WalrusError, showFileName bool) {
	fileData, err := os.ReadFile(e.filePath)
	if err != nil {
		panic(err)
	}

	if showFileName {
		utils.BLUE.Print("\nIn file: ")
		utils.GREY.Printf("%s:%d:%d\n", e.filePath, e.lineStart, e.colStart)
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
	utils.GREY.Print(lineNumber)
	fmt.Println(line)
	underLine := fmt.Sprintf("%s^%s\n", strings.Repeat(" ", (e.colStart-1)+len(lineNumber)), strings.Repeat("~", hLen))

	utils.RED.Print(underLine)
	if e.level == ERROR_CRITICAL {
		//stop further execution
		utils.BOLD_RED.Print("Critical Error: ")
	}

	utils.RED.Println(e.err.Error())

	if len(e.hints) > 0 {
		utils.GREEN.Println("Hint:")
		for _, hint := range e.hints {
			if hint.hintType == TEXT_HINT {
				utils.YELLOW.Println(hint.message)
			} else {
				utils.ORANGE.Println(hint.message)
			}
		}
	}

	if e.level == ERROR_CRITICAL {
		utils.ORANGE.Println("Compilation stopped due to critical error. Resolve the critical error to continue compilation")
		os.Exit(-1)
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

func makeError(filePath string, lineStart, lineEnd int, colStart, colEnd int, errMsg string, level ERROR_LEVEL) *WalrusError {
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
		level:     level,
	}

	globalErrors = append(globalErrors, err)

	return err
}

// global errors are arrays of error pointers
var globalErrors []*WalrusError

// make an errorlist to add all errors and display later
func AddError(filePath string, lineStart, lineEnd int, colStart, colEnd int, err string, level ERROR_LEVEL) *WalrusError {
	errItem := makeError(filePath, lineStart, lineEnd, colStart, colEnd, err, level)
	utils.YELLOW.Printf("Error added on %s:%d:%d. %d errors available\n", filePath, lineStart, colStart, len(globalErrors))
	if level == ERROR_CRITICAL {
		DisplayErrors()
	}
	return errItem
}

func DisplayErrors() {
	if len(globalErrors) == 0 {
		utils.GREEN.Println("------- Passed --------")
		return
	} else {
		utils.BOLD_RED.Printf("%d error(s) found\n", len(globalErrors))
	}
	for _, err := range globalErrors {
		PrintError(err, true)
	}
}
