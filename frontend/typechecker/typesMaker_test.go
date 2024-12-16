package typechecker

import (
	"fmt"
	"testing"
	"walrus/frontend/builtins"
)

func TestMakeNumericType(t *testing.T) {
	tests := []struct {
		isInt    bool
		bitSize  uint8
		isSigned bool
		expected builtins.TC_TYPE
	}{
		{true, 8, true, "i8"},
		{true, 16, true, "i16"},
		{true, 32, true, "i32"},
		{true, 64, true, "i64"},
		{true, 8, false, "u8"},
		{true, 16, false, "u16"},
		{true, 32, false, "u32"},
		{true, 64, false, "u64"},
		{false, 32, false, "f32"},
		{false, 64, false, "f64"},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("isInt=%v,bitSize=%d,isSigned=%v", tt.isInt, tt.bitSize, tt.isSigned), func(t *testing.T) {
			result := makeNumericType(tt.isInt, tt.bitSize, tt.isSigned)
			if result != tt.expected {
				t.Errorf("makeNumericType(%v, %d, %v) = %v; want %v", tt.isInt, tt.bitSize, tt.isSigned, result, tt.expected)
			}
		})
	}
}
