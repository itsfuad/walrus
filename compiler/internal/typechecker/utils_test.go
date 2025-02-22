package typechecker

import (
	"testing"
	"walrus/compiler/internal/builtins"
)

var exp string = "expected %s, got %s"

func TestFunctionSignatureString(t *testing.T) {
	tests := []struct {
		name     string
		fn       Fn
		expected string
	}{
		{
			name: "simple function",
			fn: Fn{
				Params: []FnParam{
					{Name: "a", Type: Int{builtins.INT32, 32, true}},
					{Name: "b", Type: Float{builtins.FLOAT64, 64}},
				},
				Returns: Int{builtins.INT32, 32, true},
			},
			expected: "fn(a: i32, b: f64) -> i32",
		},
		{
			name: "function with no return type",
			fn: Fn{
				Params: []FnParam{
					{Name: "a", Type: Int{builtins.INT32, 32, true}},
				},
				Returns: Void{},
			},
			expected: "fn(a: i32)",
		},
		{
			name: "function with no parameters",
			fn: Fn{
				Params:  []FnParam{},
				Returns: Int{builtins.INT32, 32, true},
			},
			expected: "fn() -> i32",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := functionSignatureString(tt.fn)
			if result != tt.expected {
				t.Errorf(exp, tt.expected, result)
			}
		})
	}
}
func TestTcValueToString(t *testing.T) {
	tests := []struct {
		name     string
		val      Tc
		expected string
	}{
		{
			name: "array type tcValue",
			val: Array{
				ArrayType: Int{builtins.INT32, 32, true},
			},
			expected: "[]i32",
		},
		{
			name: "struct type tcValue",
			val: Struct{
				StructName: "MyStruct",
			},
			expected: "MyStruct",
		},
		{
			name: "interface type tcValue",
			val: Interface{
				InterfaceName: "MyInterface",
			},
			expected: "MyInterface",
		},
		{
			name: "function type tcValue",
			val: Fn{
				Params: []FnParam{
					{Name: "a", Type: Int{builtins.INT32, 32, true}},
					{Name: "b", Type: Int{builtins.INT16, 16, true}},
					{Name: "c", Type: Int{builtins.UINT8, 8, false}},
				},
				Returns: Int{builtins.INT32, 32, true},
			},
			expected: "fn(a: i32, b: i16, c: u8) -> i32",
		},
		{
			name: "map type tcValue",
			val: Map{
				KeyType:   Int{builtins.INT32, 32, true},
				ValueType: Float{builtins.FLOAT64, 64},
			},
			expected: "map[i32]f64",
		},
		{
			name: "user defined type tcValue",
			val: UserDefined{
				TypeDef: Int{builtins.INT32, 32, true},
			},
			expected: "i32",
		},
		{
			name:     "default type tcValue",
			val:      Int{builtins.INT32, 32, true},
			expected: "i32",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tcToString(tt.val)
			if result != tt.expected {
				t.Errorf(exp, tt.expected, result)
			}
		})
	}
}
func TestArrayToString(t *testing.T) {
	val := Array{
		ArrayType: Int{builtins.INT32, 32, true},
	}
	expected := "[]i32"
	result := tcToString(val)
	if result != expected {
		t.Errorf(exp, expected, result)
	}
}

func TestStructToString(t *testing.T) {
	val := Struct{
		StructName: "MyStruct",
	}
	expected := "MyStruct"
	result := tcToString(val)
	if result != expected {
		t.Errorf(exp, expected, result)
	}
}

func TestInterfaceToString(t *testing.T) {
	val := Interface{
		InterfaceName: "MyInterface",
	}
	expected := "MyInterface"
	result := tcToString(val)
	if result != expected {
		t.Errorf(exp, expected, result)
	}
}

func TestFunctionToString(t *testing.T) {
	val := Fn{
		Params: []FnParam{
			{Name: "a", Type: Int{builtins.INT32, 32, true}},
			{Name: "b", Type: Float{builtins.FLOAT64, 64}},
		},
		Returns: Int{builtins.INT32, 32, true},
	}
	expected := "fn(a: i32, b: f64) -> i32"
	result := tcToString(val)
	if result != expected {
		t.Errorf(exp, expected, result)
	}
}

func TestMapToString(t *testing.T) {
	val := Map{
		KeyType:   Int{builtins.INT32, 32, true},
		ValueType: Float{builtins.FLOAT64, 64},
	}
	expected := "map[i32]f64"
	result := tcToString(val)
	if result != expected {
		t.Errorf(exp, expected, result)
	}
}

func TestUserDefinedToString(t *testing.T) {
	val := UserDefined{
		TypeDef: Int{builtins.INT32, 32, true},
	}
	expected := "i32"
	result := tcToString(val)
	if result != expected {
		t.Errorf(exp, expected, result)
	}
}

func TestDefaultToString(t *testing.T) {
	val := Int{builtins.INT32, 32, true}
	expected := "i32"
	result := tcToString(val)
	if result != expected {
		t.Errorf(exp, expected, result)
	}
}
func TestIsNumberType(t *testing.T) {
	tests := []struct {
		name     string
		val      Tc
		expected bool
	}{
		{
			name:     "int type",
			val:      Int{builtins.INT32, 32, true},
			expected: true,
		},
		{
			name:     "float type",
			val:      Float{builtins.FLOAT64, 64},
			expected: true,
		},
		{
			name:     "string type",
			val:      Str{builtins.STRING},
			expected: false,
		},
		{
			name:     "struct type",
			val:      Struct{StructName: "MyStruct"},
			expected: false,
		},
		{
			name:     "array type",
			val:      Array{ArrayType: Int{builtins.INT32, 32, true}},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isNumberType(tt.val)
			if result != tt.expected {
				t.Errorf(exp, tt.expected, result)
			}
		})
	}
}
func TestIsIntType(t *testing.T) {
	tests := []struct {
		name     string
		val      Tc
		expected bool
	}{
		{
			name:     "int type",
			val:      Int{builtins.INT32, 32, true},
			expected: true,
		},
		{
			name:     "float type",
			val:      Float{builtins.FLOAT64, 64},
			expected: false,
		},
		{
			name:     "string type",
			val:      Str{builtins.STRING},
			expected: false,
		},
		{
			name:     "struct type",
			val:      Struct{StructName: "MyStruct"},
			expected: false,
		},
		{
			name:     "array type",
			val:      Array{ArrayType: Int{builtins.INT32, 32, true}},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isIntType(tt.val)
			if result != tt.expected {
				t.Errorf(exp, tt.expected, result)
			}
		})
	}
}
