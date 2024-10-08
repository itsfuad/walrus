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
	lineStart 	int
	lineEnd		int
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
	utils.ColorPrint(utils.GREY, fmt.Sprintf("\nIn file: %s:%d:%d\n", e.filePath, e.lineStart, e.colStart))
	lines := strings.Split(string(fileData), "\n")
	line := lines[e.lineStart - 1]
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
	underLine := fmt.Sprintf("%s^%s", strings.Repeat(" ", e.colStart - 1), strings.Repeat("~", hLen))
	
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
	panic("")
}


func (e *ErrorType) AddHint(msg string, htype HINT) *ErrorType {
	e.hints = append(e.hints, Hint{
		message: msg,
		hintType: htype,
	})

	return e
}

func MakeError(filePath string, lineStart, lineEnd int, colStart, colEnd int, err string) *ErrorType {
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
	return &ErrorType{
		filePath: filePath,
		lineStart: lineStart,
		lineEnd: lineEnd,
		colStart: colStart,
		colEnd: colEnd,
		err: fmt.Errorf(err),
	}
}