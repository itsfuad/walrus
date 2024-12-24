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
	reportType     REPORT_TYPE
}

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
func printReport(r *Report) {

	utils.GREY.Printf("%s:%d:%d: ", r.filePath, r.lineStart, r.colStart)

	var reportMsg string

	switch r.reportType {
	case WARNING:
		reportMsg = "Warning: "
	case INFO:
		reportMsg = "Info: "
	case CRITICAL_ERROR:
		reportMsg = "Critical Error: "
	case SYNTAX_ERROR:
		reportMsg = "Syntax Error: "
	case NORMAL_ERROR:
		reportMsg = "Error: "
	}

	reportColor := colorMap[r.reportType]
	reportColor.Print(reportMsg)
	reportColor.Print(r.msg + "\n")

	makeParts(r, reportColor)

	if r.reportType == CRITICAL_ERROR || r.reportType == SYNTAX_ERROR {
		panic(fmt.Sprintf("Compilation halted due to %s\n", r.reportType))
	}
}

func makeParts(r *Report, color utils.COLOR){
	fileData, err := os.ReadFile(r.filePath)
	if err != nil {
		panic(err)
	}

	lines := strings.Split(string(fileData), "\n")

	if r.lineStart == r.lineEnd {
		underlineLength := (r.colEnd - r.colStart) - 1
		handleSingleLineError(lines[r.lineStart - 1], r.lineStart, r, underlineLength, color, true)
	} else {
		handleMultiLineError(lines, r, color)
	}
}

func handleSingleLineError(line string, lineNo int, r *Report, underlineLength int, color utils.COLOR, underline bool) {
	if underlineLength < 0 {
		underlineLength = 0
	}
	lineNumber := fmt.Sprintf("%d | ", lineNo)
	utils.GREY.Print(lineNumber)
	fmt.Println(line)

	if !underline {
		return
	}

	color.Printf("%s^%s\n", strings.Repeat(" ", (r.colStart-1)+len(lineNumber)), strings.Repeat("~", underlineLength))
}

func handleMultiLineError(lines []string, r *Report, color utils.COLOR) {
	codeLines := lines[r.lineStart-1 : r.lineEnd]
	for i, line := range codeLines {
		underlineLength := len(line) - r.colStart + 1
		handleSingleLineError(line, r.lineStart + i, r, underlineLength, color, i == len(codeLines)-1)
	}
}

func Add(filePath string, lineStart, lineEnd int, colStart, colEnd int, msg string, reportType REPORT_TYPE) *Report {
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
		reportType:     reportType,
	}

	globalReports = append(globalReports, report)

	reports[reportType]++

	return report
}

func DisplayAll() {
	//recover if panics
	defer func() {

		if reports[CRITICAL_ERROR] == 0 && reports[NORMAL_ERROR] == 0 && reports[SYNTAX_ERROR] == 0 {
			showStatus(true, "Compilation successful with")
			return
		}

		if r := recover(); r != nil {
			utils.BOLD_RED.Println(r)
		}

		showStatus(false, "Compilation failed with")
		os.Exit(-1)
	}()
	for _, err := range globalReports {
		printReport(err)
	}
}

func showStatus(passed bool, msg string) {

	//show errors and warnings separately
	warningCount := reports[WARNING]
	probCount := reports[NORMAL_ERROR] + reports[CRITICAL_ERROR] + reports[SYNTAX_ERROR]

	var messageColor utils.COLOR

	if passed {
		messageColor = utils.GREEN
		messageColor.Printf("------------- %s ", msg)
	} else {
		messageColor = utils.RED
		messageColor.Printf("------------- %s ", msg)
	}

	totalProblemsString := ""

	if warningCount > 0 {
		totalProblemsString += colorMap[WARNING].Sprintf("%d %s", warningCount, utils.Plural("warning", "warnings ", warningCount))
		if probCount > 0 {
			totalProblemsString += utils.ORANGE.Sprintf(", ")
		}
	}

	if probCount > 0 {
		totalProblemsString += colorMap[NORMAL_ERROR].Sprintf("%d %s", probCount, utils.Plural("error", "errors", probCount))
	} else {
		totalProblemsString += messageColor.Sprint("no issues found")
	}

	messageColor.Print(totalProblemsString)
	messageColor.Println(" -------------")
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
