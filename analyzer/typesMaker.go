package analyzer

import (
	"fmt"
	"walrus/builtins"
)

func makeNumericType(isInt bool, bitSize uint8, isSigned bool) builtins.TC_TYPE {
	rawDataType := ""
	if isInt {
		if isSigned {
			rawDataType = "i"
		} else {
			rawDataType = "u"
		}
	} else {
		rawDataType = "f"
	}

	rawDataType += fmt.Sprintf("%d", bitSize)

	return builtins.TC_TYPE(rawDataType)
}

// helper type initialization functions
func NewInt(bitSize uint8, isSigned bool) Int {
	return Int{DataType: makeNumericType(true, bitSize, isSigned), BitSize: bitSize, IsSigned: isSigned}
}

func NewFloat(bitSize uint8) Float {
	return Float{DataType: makeNumericType(false, bitSize, false), BitSize: bitSize}
}

func NewStr() Str {
	return Str{DataType: STRING_TYPE}
}

func NewBool() Bool {
	return Bool{DataType: BOOLEAN_TYPE}
}

func NewNull() Null {
	return Null{DataType: NULL_TYPE}
}

func NewVoid() Void {
	return Void{DataType: VOID_TYPE}
}

func NewMap(keyType ExprType, valueType ExprType) Map {
	return Map{DataType: MAP_TYPE, KeyType: keyType, ValueType: valueType}
}

func NewMaybe(valueType ExprType) Maybe {
	return Maybe{DataType: MAYBE_TYPE, MaybeType: valueType}
}
