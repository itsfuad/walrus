package errgen

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"walrus/utils"
)

type PROBLEM_TYPE string

const (
	NULL           	PROBLEM_TYPE = ""
	CRITICAL_ERROR 	PROBLEM_TYPE = "critical error" // Stops compilation immediately
	SYNTAX_ERROR   	PROBLEM_TYPE = "syntax error"   // Syntax error, also stops compilation
	NORMAL_ERROR   	PROBLEM_TYPE = "error"          // Regular error that doesn't halt compilation

	WARNING 		PROBLEM_TYPE = "warning" // Indicates potential issues
)

// global errors are arrays of error pointers
var globalProblems []*Problem
var problems = make(map[PROBLEM_TYPE]int)

type Problem struct {
	filePath  string
	lineStart int
	lineEnd   int
	colStart  int
	colEnd    int
	err       error
	hints     []string
	level     PROBLEM_TYPE
}

// printProblem formats and displays error information for a WalrusError.
// It prints the error location, the relevant code line, and visual indicators
// showing where the error occurred. For critical errors, it will terminate
// program execution.
//
// Parameters:
//   - e: Pointer to a WalrusError containing error details
//   - showFileName: Boolean flag to control whether the file name is displayed
//
// The function:
//   - Reads the source file
//   - Displays file location (if showFileName is true)
//   - Shows the problematic line of code
//   - Highlights the error position with ^ and ~ characters
//   - Prints the error message
//   - Shows hints if available
//   - Exits program if error is critical
//
// If file reading fails, the function will panic.
func printProblem(e *Problem) {
	fileData, err := os.ReadFile(e.filePath)
	if err != nil {
		panic(err)
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

	markUnderline(e, underLine)

	utils.GREY.Printf("at: %s:%d:%d\n", e.filePath, e.lineStart, e.colStart)

	if len(e.hints) > 0 {
		utils.YELLOW.Println("Hint:")
		for _, hint := range e.hints {
			utils.YELLOW.Printf("- %s\n", hint)
		}
	} else {
		fmt.Println()
	}

	if e.level == CRITICAL_ERROR || e.level == SYNTAX_ERROR {
		panic(fmt.Sprintf("Compilation halted due to %s\n", e.level))
	}
}

func markUnderline(e *Problem, underLine string) {
	if e.level == WARNING {
		utils.YELLOW.Print(underLine)
		utils.YELLOW.Print("Warning: ")
		utils.YELLOW.Print(e.err.Error() + "\n")
	} else {
		utils.RED.Print(underLine)
		if e.level == CRITICAL_ERROR {
			//stop further execution
			utils.BOLD_RED.Print("Critical Error: ")
		} else if e.level == SYNTAX_ERROR {
			utils.BOLD_RED.Print("Syntax Error: ")
		} else {
			utils.RED.Print("Error: ")
		}
		utils.RED.Print(e.err.Error() + "\n")
	}
}


// Hint appends a hint message to the error's hints slice.
// If the provided message is empty, it returns the error without modification.
// Each hint provides additional context or suggestions about the error.
//
// Parameters:
//   - msg: The hint message to add
//
// Returns:
//   - *WalrusError: Returns the error instance to allow for method chaining
func (e *Problem) Hint(msg string) *Problem {

	if msg == "" {
		return e
	}

	e.hints = append(e.hints, msg)
	return e
}

func Add(filePath string, lineStart, lineEnd int, colStart, colEnd int, errMsg string) *Problem {
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

	err := &Problem{
		filePath:  filePath,
		lineStart: lineStart,
		lineEnd:   lineEnd,
		colStart:  colStart,
		colEnd:    colEnd,
		err:       errors.New(errMsg),
		level:     NULL,
	}

	globalProblems = append(globalProblems, err)

	return err
}

func (e *Problem) Level(level PROBLEM_TYPE) {
	if level == NULL {
		panic("call ErrorLevel() method with valid Error level")
	}
	e.level = level
	problems[level]++
	if level == CRITICAL_ERROR || level == SYNTAX_ERROR {
		DisplayAll()
	}
}

func DisplayAll() {
	//recover if panics
	defer func() {
		if problems[CRITICAL_ERROR] == 0 && problems[NORMAL_ERROR] == 0 {
			utils.GREEN.Println("-------- Passed --------")
		}
		displayProblemCount()
		if r := recover(); r != nil {
			utils.BOLD_RED.Println(r)
			os.Exit(-1)
		}
	}()
	for _, err := range globalProblems {
		if err.level == NULL {
			panic("call Level() method with valid Error level")
		}
		printProblem(err)
	}
}

func displayProblemCount() {
	//show errors and warnings separately

	if problems[WARNING] > 0 {
		str := fmt.Sprintf("%d warning", problems[WARNING])
		if problems[WARNING] > 1 {
			str += "s"
		}
		utils.YELLOW.Print(str)
	}

	if problems[NORMAL_ERROR]+problems[CRITICAL_ERROR] > 0 {
		if problems[WARNING] > 0 {
			utils.YELLOW.Print(", ")
		}
		str := fmt.Sprintf("%d error", problems[NORMAL_ERROR] + problems[CRITICAL_ERROR])
		if problems[NORMAL_ERROR] > 1 {
			str += "s"
		}
		utils.RED.Println(str)
	}
}

// func Tree print with one or more strings
func TreeFormatString(strings ...string) string {
	// use └, ├ as tree characters
	str := ""
	for i, prop := range strings {
		if i == len(strings)-1 {
			str += utils.GREY.Sprint("└── ") + utils.BROWN.Sprint(prop)
		} else {
			str += utils.GREY.Sprint("├── ") + utils.BROWN.Sprint(prop+"\n")
		}
	}
	return str
}

func TreeFormatError(errs ...error) error {
	strs := []string{}
	for _, err := range errs {
		strs = append(strs, err.Error())
	}
	return errors.New(TreeFormatString(strs...))
}
