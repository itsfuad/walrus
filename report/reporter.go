package report

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"walrus/utils"
)

type REPORT_TYPE string

const (
	NULL           	REPORT_TYPE = ""
	CRITICAL_ERROR 	REPORT_TYPE = "critical error" // Stops compilation immediately
	SYNTAX_ERROR   	REPORT_TYPE = "syntax error"   // Syntax error, also stops compilation
	NORMAL_ERROR   	REPORT_TYPE = "error"          // Regular error that doesn't halt compilation

	WARNING 		REPORT_TYPE = "warning" // Indicates potential issues
	INFO   			REPORT_TYPE = "info"    // Informational message
)

//var colorMap = make(map[REPORT_TYPE]utils.COLOR)
var colorMap = map[REPORT_TYPE]utils.COLOR{
	CRITICAL_ERROR: utils.BOLD_RED,
	SYNTAX_ERROR:   utils.RED,
	NORMAL_ERROR:   utils.RED,
	WARNING:        utils.YELLOW,
	INFO:           utils.BLUE,
}

// global errors are arrays of error pointers
var globalReports []*Report
var reports = make(map[REPORT_TYPE]int)

type Report struct {
	filePath  string
	lineStart int
	lineEnd   int
	colStart  int
	colEnd    int
	msg       string
	hints     []string
	level     REPORT_TYPE
}

// printReport formats and displays error information for a WalrusError.
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
func printReport(e *Report) {

	utils.GREY.Printf("%s:%d:%d: ", e.filePath, e.lineStart, e.colStart)

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
	snippet := utils.GREY.Sprint(lineNumber) + fmt.Sprintln(line)
	underline := fmt.Sprintf("%s^%s\n", strings.Repeat(" ", (e.colStart-1)+len(lineNumber)), strings.Repeat("~", hLen))

	
	if e.level == WARNING {
		colorMap[WARNING].Print("Warning: ")
		colorMap[WARNING].Print(e.msg + "\n")
		fmt.Print(snippet)
		colorMap[WARNING].Print(underline)
	} else if e.level == INFO {
		colorMap[INFO].Print("Info: ")
		colorMap[INFO].Print(e.msg + "\n")
		fmt.Print(snippet)
		colorMap[INFO].Print(underline)
	} else {
		if e.level == CRITICAL_ERROR {
			//stop further execution
			colorMap[CRITICAL_ERROR].Print("Critical Error: ")
		} else if e.level == SYNTAX_ERROR {
			colorMap[SYNTAX_ERROR].Print("Syntax Error: ")
		} else {
			colorMap[NORMAL_ERROR].Print("Error: ")
		}
		utils.RED.Print(e.msg + "\n")
		fmt.Print(snippet)
		utils.RED.Print(underline)
	}

	showHints(e, hLen)

	if e.level == CRITICAL_ERROR || e.level == SYNTAX_ERROR {
		panic(fmt.Sprintf("Compilation halted due to %s\n", e.level))
	}
}

func showHints(e *Report, padding int) {
	if len(e.hints) > 0 {
		utils.YELLOW.Printf("%sHint:\n", strings.Repeat(" ", padding))
		for _, hint := range e.hints {
			utils.YELLOW.Printf("%s- %s\n", strings.Repeat(" ", padding), hint)
		}
	} else {
		fmt.Println()
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
func (e *Report) Hint(msg string) *Report {

	if msg == "" {
		return e
	}

	e.hints = append(e.hints, msg)
	return e
}

func Add(filePath string, lineStart, lineEnd int, colStart, colEnd int, msg string) *Report {
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

	report := &Report{
		filePath:  filePath,
		lineStart: lineStart,
		lineEnd:   lineEnd,
		colStart:  colStart,
		colEnd:    colEnd,
		msg:       msg,
		level:     NULL,
	}

	globalReports = append(globalReports, report)

	return report
}

func (e *Report) Level(level REPORT_TYPE) {
	if level == NULL {
		panic("call ErrorLevel() method with valid Error level")
	}
	e.level = level
	reports[level]++
	if level == CRITICAL_ERROR || level == SYNTAX_ERROR {
		DisplayAll()
	}
}

func DisplayAll() {
	//recover if panics
	defer func() {
		displayProblemCount()
		if reports[CRITICAL_ERROR] == 0 && reports[NORMAL_ERROR] == 0 {
			utils.GREEN.Println("------------ Passed ------------")
		}
		if r := recover(); r != nil {
			utils.BOLD_RED.Println(r)
			os.Exit(-1)
		}
	}()
	for _, err := range globalReports {
		if err.level == NULL {
			panic("call Level() method with valid Error level")
		}
		printReport(err)
	}
}

func displayProblemCount() {
	//show errors and warnings separately
	warningCount := reports[WARNING]
	probCount := reports[NORMAL_ERROR] + reports[CRITICAL_ERROR]

	if warningCount > 0 {
		colorMap[WARNING].Printf("%d %s", warningCount, utils.Plural("warning", "warnings ", warningCount))
		if probCount > 0 {
			utils.ORANGE.Printf(", ")
		}
	}

	if probCount > 0 {
		colorMap[NORMAL_ERROR].Printf("%d %s", probCount, utils.Plural("error", "errors", probCount))
	}

	println()
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
