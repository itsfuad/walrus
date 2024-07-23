package errgen

import (
	"fmt"
	"os"
	"strings"
	"walrus/utils"
)

type HINT int

const (
	TEXT_HINT	HINT = iota
	CODE_HINT
)

type Hint struct {
	message		string
	hintType	HINT
}

type ErrorType struct {
	filePath	string
	line 		int
	colStart	int
	colEnd		int
	err 		error
	hints		[]Hint
}

func (e *ErrorType) Display() {
	fileData, err := os.ReadFile(e.filePath)
	if err != nil {
		panic(err)
	}
	utils.ColorPrint(utils.GREY, fmt.Sprintf("\nIn file: %s:%d:%d\n", e.filePath, e.line, e.colStart))
	lines := strings.Split(string(fileData), "\n")
	line := lines[e.line - 1]
	hLen := (e.colEnd - e.colStart) - 1
	if hLen < 0 {
		hLen = 0
	}
	fmt.Println(line)
	underLine := fmt.Sprintf("%s^%s", strings.Repeat(" ", e.colStart - 1), strings.Repeat("~", hLen))
	
	utils.ColorPrint(utils.RED, underLine)
	utils.ColorPrint(utils.PURPLE, e.err.Error())
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
	panic("")
	os.Exit(1)
}


func (e *ErrorType) AddHint(msg string, htype HINT) *ErrorType {
	e.hints = append(e.hints, Hint{
		message: msg,
		hintType: htype,
	})

	return e
}

func MakeError(filePath string, line int, colStart, colEnd int, err string) *ErrorType {
	return &ErrorType{
		filePath: filePath,
		line: line,
		colStart: colStart,
		colEnd: colEnd,
		err: fmt.Errorf(err),
	}
}