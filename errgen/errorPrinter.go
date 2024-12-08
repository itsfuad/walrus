// Package errgen provides error handling and reporting functionality for the Walrus compiler.
//
// The package implements a custom error type WalrusError and various utility functions
// for error management, including error creation, display, and collection.
//
// Error Levels:
//   - ERROR_CRITICAL: Stops compilation immediately
//   - ERROR_NORMAL: Regular error that doesn't halt compilation
//   - WARNING: Indicates potential issues
//   - INFO: Informational messages
//
// Example usage:
//
//	err := AddError("main.go", 1, 1, 1, 10, "unexpected token", ERROR_NORMAL)
//	err.AddHint("Did you forget a semicolon?")
//	DisplayErrors()
//
// The package provides colored output for better error visualization and supports
// features like:
//   - Line highlighting with ^ and ~ characters
//   - File location reporting
//   - Multiple hints per error
//   - Global error collection
//   - Critical error handling with immediate program termination
//
// Global error tracking allows accumulating multiple errors before displaying them,
// unless a critical error is encountered, which triggers immediate display and program exit.
package errgen

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"walrus/utils"
)

type ERROR_LEVEL int

const (
	ERROR_CRITICAL ERROR_LEVEL = iota			// Stops compilation immediately	
	ERROR_NORMAL								// Regular error that doesn't halt compilation
	WARNING										// Indicates potential issues
	INFO										// Informational messages
)

type WalrusError struct {
	filePath  string
	lineStart int
	lineEnd   int
	colStart  int
	colEnd    int
	err       error
	hints     []string
	level     ERROR_LEVEL
}


// PrintError formats and displays error information for a WalrusError.
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
			utils.GREEN.Printf("  %s\n", hint)
		}
	}

	if e.level == ERROR_CRITICAL {
		utils.ORANGE.Println("Compilation stopped due to critical error. Resolve the critical error to continue compilation")
		os.Exit(-1)
	}
}

// AddHint appends a hint message to the error's hints slice.
// If the provided message is empty, it returns the error without modification.
// Each hint provides additional context or suggestions about the error.
//
// Parameters:
//   - msg: The hint message to add
//
// Returns:
//   - *WalrusError: Returns the error instance to allow for method chaining
func (e *WalrusError) AddHint(msg string) *WalrusError {

	if msg == "" {
		return e
	}

	e.hints = append(e.hints, msg)

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
