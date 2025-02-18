package report

import (
	"fmt"
	"os"
	"strings"
	"walrus/compiler/colors"
	"walrus/compiler/internal/utils"
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

type Reports []*Report

// global errors are arrays of error pointers
var globalReports Reports

// Report represents a diagnostic report used both internally and by LSP.
type Report struct {
	FilePath  string
	LineStart int
	LineEnd   int
	ColStart  int
	ColEnd    int
	Message   string
	Hints     []string
	Level     REPORT_TYPE
}

// GetReports returns a slice of diagnostics converted from internal reports.
// It skips any reports that do not have a valid level.
func GetReports() Reports {
	var diags Reports
	for _, r := range globalReports {
		if r.Level == NULL {
			// Skip reports without valid level.
			continue
		}
		diags = append(diags, r)
	}

	return diags
}

func ClearReports() {
	globalReports = Reports{}
	colors.CYAN.Println("Reports cleared")
}

// printReport prints a formatted diagnostic report to stdout.
// It shows file location, a code snippet, underline highlighting, any hints,
// and panics if the diagnostic level is critical or indicates a syntax error.
func printReport(r *Report) {

	// Generate the code snippet and underline.
	// hLen is the padding length for hint messages.
	snippet, underline, hLen := makeParts(r)

	var reportMsgType string

	switch r.Level {
	case WARNING:
		reportMsgType = "Warning: "
	case INFO:
		reportMsgType = "Info: "
	case CRITICAL_ERROR:
		reportMsgType = "Critical Error: "
	case SYNTAX_ERROR:
		reportMsgType = "Syntax Error: "
	case NORMAL_ERROR:
		reportMsgType = "Error: "
	}

	reportColor := colorMap[r.Level]

	// The error message type and the message itself are printed in the same color.
	reportColor.Print(reportMsgType)
	reportColor.Print(r.Message + "\n")

	//numlen is the length of the line number
	numlen := len(fmt.Sprint(r.LineStart))

	// The file path is printed in grey.
	colors.GREY.Printf("%s> [%s:%d:%d]\n", strings.Repeat("=", numlen+1), r.FilePath, r.LineStart, r.ColStart)

	// The code snippet and underline are printed in the same color.
	fmt.Print(snippet)
	reportColor.Print(underline)

	showHints(r, hLen)
}

// makeParts reads the source file and generates a code snippet and underline
// indicating the location of the diagnostic. It returns the snippet, underline,
// and a padding value.
func makeParts(r *Report) (snippet, underline string, hLen int) {
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

// showHints outputs any associated hint messages for the diagnostic,
// using the provided padding for proper alignment.
func showHints(r *Report, padding int) {
	if len(r.Hints) > 0 {
		colors.YELLOW.Printf("%sHint:\n", strings.Repeat(" ", padding))
		for _, hint := range r.Hints {
			colors.YELLOW.Printf("%s- %s\n", strings.Repeat(" ", padding), hint)
		}
	} else {
		fmt.Println()
	}
}

// Hint appends a new hint message to the diagnostic and returns the updated diagnostic.
// It ignores empty hint messages.
func (r *Report) Hint(msg string) *Report {

	if msg == "" {
		return r
	}

	r.Hints = append(r.Hints, msg)
	return r
}

// Add creates and registers a new diagnostic report with basic position validation.
// It returns a pointer to the newly created Diagnostic.
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

// SetLevel assigns a diagnostic level to the report, increments its count,
// and triggers DisplayAll if the level is critical or denotes a syntax error.
func (e *Report) SetLevel(level REPORT_TYPE) {
	if level == NULL {
		panic("call SetLevel() method with valid Error level")
	}
	e.Level = level
	if level == CRITICAL_ERROR || level == SYNTAX_ERROR {
		panic("Critical error or syntax error detected")
	}
}

// DisplayAll outputs all the diagnostic reports. It recovers from panics,
// prints a summary status, and exits the process if errors are present.
func (r Reports) DisplayAll() {
	for _, err := range r {
		if err.Level == NULL {
			panic("call SetLevel() method with valid Error level")
		}
		printReport(err)
	}
	r.ShowStatus()
}

// ShowStatus displays a summary of compilation status along with counts of warnings and errors.
func (r Reports) ShowStatus() {
	warningCount := 0
	probCount := 0

	for _, report := range r {
		if report.Level == WARNING {
			warningCount++
		} else if report.Level == NORMAL_ERROR || report.Level == CRITICAL_ERROR || report.Level == SYNTAX_ERROR {
			probCount++
		}
	}

	var messageColor colors.COLOR

	if probCount > 0 {
		messageColor = colors.RED
		messageColor.Print("------------- failed with ")
	} else {
		messageColor = colors.GREEN
		messageColor.Print("------------- Passed")
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
