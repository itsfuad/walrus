package errgen

import (
	"errors"
	"testing"
	"walrus/utils"
)

const (
	error1 = "error message 1"
	error2 = "error message 2"
)

var tests = []struct {
	name     string
	input    []string
	expected string
}{
	{
		name:     "Single string",
		input:    []string{error1},
		expected: utils.GREY.Sprint("├─── ") + utils.BROWN.Sprintln(error1),
	},
	{
		name:     "Multiple strings",
		input:    []string{error1, error2},
		expected: utils.GREY.Sprint("├─── ") + utils.BROWN.Sprintln(error1) + utils.GREY.Sprint("├─── ") + utils.BROWN.Sprintln(error2),
	},
	{
		name:     "No strings",
		input:    []string{},
		expected: "",
	},
}
func TestTreeFormatString(t *testing.T) {
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := TreeFormatString(tt.input...)
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}


func TestTreeFormatError(t *testing.T) {
	//use the strings array and use as error
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result error
			var errs []error
			//put all the strings in the error
			for _, str := range tt.input {
				errs = append(errs, errors.New(str))
			}
			result = TreeFormatError(errs...)
			expected := tt.expected
			if result.Error() != expected {
				t.Errorf("expected %q, got %q", expected, result.Error())
			}
		})
	}
}