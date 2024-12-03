package utils

import (
	"fmt"
)

// ANSI color escape codes
type COLOR string

const (
    RESET  		COLOR 	=  	"\033[0m"
    RED    		COLOR 	= 	"\033[31m"
	BOLD_RED 	COLOR 	= 	"\033[38;05;196m"
    GREEN  		COLOR 	= 	"\033[32m"
    YELLOW 		COLOR 	= 	"\033[33m"
	ORANGE 		COLOR 	= 	"\033[38;05;166m"
    BLUE   		COLOR 	= 	"\033[34m"
    PURPLE 		COLOR 	= 	"\033[35m"
    CYAN   		COLOR 	= 	"\033[36m"
    WHITE  		COLOR 	= 	"\033[37m"
	GREY   		COLOR 	= 	"\033[90m"
	BOLD   		COLOR 	= 	"\033[1m"
)

func (c COLOR) Printf(format string, args ...interface{}) {
	fmt.Printf(string(c) + format + string(RESET), args...)
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