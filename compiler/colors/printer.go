package colors

import (
	"fmt"
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
