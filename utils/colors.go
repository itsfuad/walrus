package utils

import (
	"fmt"
)

// ANSI color escape codes
type COLOR string

const (
	RESET        COLOR = "\033[0m"
	RED          COLOR = "\033[31m"
	BOLD_RED     COLOR = "\033[38;05;196m"
	GREEN        COLOR = "\033[32m"
	YELLOW       COLOR = "\033[33m"
	BOLD_YELLOW  COLOR = "\033[38;05;226m"
	ORANGE       COLOR = "\033[38;05;166m"
	BROWN        COLOR = "\033[38;05;130m"
	BRIGHT_BROWN COLOR = "\033[38;05;136m"
	BLUE         COLOR = "\033[34m"
	BOLD_BLUE    COLOR = "\033[38;05;27m"
	PURPLE       COLOR = "\033[35m"
	CYAN         COLOR = "\033[36m"
	WHITE        COLOR = "\033[37m"
	GREY         COLOR = "\033[90m"
	BOLD         COLOR = "\033[1m"
)

func (c COLOR) Printf(format string, args ...interface{}) {
	fmt.Printf(string(c)+format+string(RESET), args...)
}

func (c COLOR) Println(args ...interface{}) {
	fmt.Print(string(c))
	fmt.Println(args...)
	fmt.Print(string(RESET))
}

func (c COLOR) Print(args ...interface{}) {
	fmt.Print(string(c))
	fmt.Print(args...)
	fmt.Print(string(RESET))
}

func (c COLOR) Sprintf(format string, args ...interface{}) string {
	return string(c) + fmt.Sprintf(format, args...) + string(RESET)
}

func (c COLOR) Sprintln(args ...interface{}) string {
	return string(c) + fmt.Sprintln(args...) + string(RESET)
}

func (c COLOR) Sprint(args ...interface{}) string {
	return string(c) + fmt.Sprint(args...) + string(RESET)
}

func PrintWithColor(color COLOR, args ...interface{}) {
	color.Print(args...)
}

func SprintWithColor(color COLOR, args ...interface{}) string {
	return color.Sprint(args...)
}
