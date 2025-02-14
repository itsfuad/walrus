package report

import (
	"errors"
	"testing"

	"walrus/compiler/colors"
)

const error1 = "error message 1"
const error2 = "error message 2"

var tests = []struct {
	name     string
	input    []string
	expected string
}{
	{
		name:     "Single string",
		input:    []string{error1},
		expected: colors.GREY.Sprint("└── ") + colors.BROWN.Sprint(error1),
	},
	{
		name:     "Multiple strings",
		input:    []string{error1, error2},
		expected: colors.GREY.Sprint("├── ") + colors.BROWN.Sprintln(error1) + colors.GREY.Sprint("└── ") + colors.BROWN.Sprint(error2),
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
				t.Errorf("expected %q\ngot %q", tt.expected, result)
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
				t.Errorf("expected %q\ngot %q", expected, result.Error())
			}
		})
	}
}

// test panic and recover
func TestPanicRecover(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Log("recovered from panic")
		} else {
			t.Error("did not panic")
		}
	}()

	// Trigger the panic by passing a non-primitive value (e.g., a struct or a slice)
	Panicable(1, []int{2, 3})
}

func Panicable(data ...any) {
	//panic if data is not primitive
	for _, d := range data {
		if !isPrimitive(d) {
			panic("Invalid data type")
		}
	}
}

func isPrimitive(data any) bool {
	switch data.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64, string, bool:
		return true
	}
	return false
}
