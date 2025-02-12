package report

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"walrus/colors"
	"walrus/internal/utils"
)

type REPORT_TYPE string

const (
	NULL           REPORT_TYPE = ""
	CRITICAL_ERROR REPORT_TYPE = "critical error" // Stops compilation immediately
	SYNTAX_ERROR   REPORT_TYPE = "syntax error"   // Syntax error, also stops compilation
	NORMAL_ERROR   REPORT_TYPE = "error"          // Regular error that doesn't halt compilation

	WARNING REPORT_TYPE = "warning" // Indicates potential issues
	INFO    REPORT_TYPE = "info"    // Informational message
)

// var colorMap = make(map[REPORT_TYPE]utils.COLOR)
var colorMap = map[REPORT_TYPE]colors.COLOR{
	CRITICAL_ERROR: colors.BOLD_RED,
	SYNTAX_ERROR:   colors.RED,
	NORMAL_ERROR:   colors.RED,
	WARNING:        colors.YELLOW,
	INFO:           colors.BLUE,
}

// global errors are arrays of error pointers
var globalReports []*Diagnostic
var reports = make(map[REPORT_TYPE]int)

// Add exported Diagnostic type to be used by the LSP server.
type Diagnostic struct {
	FilePath  string      `json:"filePath"`
	LineStart int         `json:"lineStart"`
	LineEnd   int         `json:"lineEnd"`
	ColStart  int         `json:"colStart"`
	ColEnd    int         `json:"colEnd"`
	Message   string      `json:"message"`
	Hints     []string    `json:"hints"`
	Level     REPORT_TYPE `json:"level"`
}

// GetDiagnostics converts internal reports to LSP diagnostics.
func GetDiagnostics() []Diagnostic {
	var diags []Diagnostic
	for _, r := range globalReports {
		if r.Level == NULL {
			// Skip reports without valid level.
			continue
		}
		diags = append(diags, *r)
	}
	return diags
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
func printReport(r *Diagnostic) {

	colors.GREY.Printf("%s:%d:%d: ", r.FilePath, r.LineStart, r.ColStart)

	snippet, underline, hLen := makeParts(r)

	var reportMsg string

	switch r.Level {
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

	reportColor := colorMap[r.Level]
	reportColor.Print(reportMsg)
	reportColor.Print(r.Message + "\n")

	fmt.Print(snippet)
	reportColor.Print(underline)

	showHints(r, hLen)

	if r.Level == CRITICAL_ERROR || r.Level == SYNTAX_ERROR {
		panic(fmt.Sprintf("Compilation halted due to %s\n", r.Level))
	}
}

func makeParts(r *Diagnostic) (snippet, underline string, hLen int) {
	fileData, err := os.ReadFile(r.FilePath)
	if err != nil {
		panic(err)
	}

	lines := strings.Split(string(fileData), "\n")
	line := lines[r.LineStart-1]

	hLen = 0

	if r.LineStart == r.LineEnd {
		hLen = (r.ColEnd - r.ColStart) - 1
	} else {
		//full line
		hLen = len(line) - 2
	}
	if hLen < 0 {
		hLen = 0
	}

	lineNumber := fmt.Sprintf("%d | ", r.LineStart)
	snippet = colors.GREY.Sprint(lineNumber) + line + "\n"
	underline = fmt.Sprintf("%s^%s\n", strings.Repeat(" ", (r.ColStart-1)+len(lineNumber)), strings.Repeat("~", hLen))

	return snippet, underline, hLen
}

func showHints(r *Diagnostic, padding int) {
	if len(r.Hints) > 0 {
		colors.YELLOW.Printf("%sHint:\n", strings.Repeat(" ", padding))
		for _, hint := range r.Hints {
			colors.YELLOW.Printf("%s- %s\n", strings.Repeat(" ", padding), hint)
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
func (r *Diagnostic) Hint(msg string) *Diagnostic {

	if msg == "" {
		return r
	}

	r.Hints = append(r.Hints, msg)
	return r
}

func Add(filePath string, lineStart, lineEnd int, colStart, colEnd int, msg string) *Diagnostic {
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

	report := &Diagnostic{
		FilePath:  filePath,
		LineStart: lineStart,
		LineEnd:   lineEnd,
		ColStart:  colStart,
		ColEnd:    colEnd,
		Message:   msg,
		Level:     NULL,
	}

	globalReports = append(globalReports, report)

	return report
}

func (e *Diagnostic) SetLevel(level REPORT_TYPE) {
	if level == NULL {
		panic("call SetLevel() method with valid Error level")
	}
	e.Level = level
	reports[level]++
	if level == CRITICAL_ERROR || level == SYNTAX_ERROR {
		DisplayAll()
	}
}

func DisplayAll() {
	//recover if panics
	defer func() {

		if reports[CRITICAL_ERROR] == 0 && reports[NORMAL_ERROR] == 0 && reports[SYNTAX_ERROR] == 0 {
			showStatus(true, "Compilation successful with")
			return
		}

		if r := recover(); r != nil {
			colors.BOLD_RED.Println(r)
		}

		showStatus(false, "Compilation failed with")
		os.Exit(-1)
	}()
	for _, err := range globalReports {
		if err.Level == NULL {
			panic("call SetLevel() method with valid Error level")
		}
		printReport(err)
	}
}

func showStatus(passed bool, msg string) {

	//show errors and warnings separately
	warningCount := reports[WARNING]
	probCount := reports[NORMAL_ERROR] + reports[CRITICAL_ERROR] + reports[SYNTAX_ERROR]

	var messageColor colors.COLOR

	if passed {
		messageColor = colors.GREEN
		messageColor.Printf("------------- %s ", msg)
	} else {
		messageColor = colors.RED
		messageColor.Printf("------------- %s ", msg)
	}

	totalProblemsString := ""

	if warningCount > 0 {
		totalProblemsString += colorMap[WARNING].Sprintf("%d %s", warningCount, utils.Plural("warning", "warnings ", warningCount))
		if probCount > 0 {
			totalProblemsString += colors.ORANGE.Sprintf(", ")
		}
	}

	if probCount > 0 {
		totalProblemsString += colorMap[NORMAL_ERROR].Sprintf("%d %s", probCount, utils.Plural("error", "errors", probCount))
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
			str += colors.GREY.Sprint("└── ") + colors.BROWN.Sprint(prop)
		} else {
			str += colors.GREY.Sprint("├── ") + colors.BROWN.Sprint(prop+"\n")
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
