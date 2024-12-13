package utils

import (
	"errors"
)

//func Tree print with one or more strings
func TreeFormatString(strings ...string) string {
	// use └, ├ as tree characters
	str := ""
	for _, prop := range strings {
		str += GREY.Sprint("├─── ") + BROWN.Sprint(prop + "\n")
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
