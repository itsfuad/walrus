package report

import (
	"testing"
)

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
