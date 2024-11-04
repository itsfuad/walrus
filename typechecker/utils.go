package typechecker

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"walrus/ast"
	"walrus/builtins"
	"walrus/errgen"
)

func init() {
	rand.New(rand.NewSource(time.Now().UnixNano())).Intn(len(letterRunes))
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func stringToValueTypeInterface(typ builtins.TC_TYPE, env *TypeEnvironment) (ValueTypeInterface, error) {
	switch typ {
	case INT8_TYPE:
		return NewInt(8, true), nil
	case INT16_TYPE:
		return NewInt(16, true), nil
	case INT32_TYPE:
		return NewInt(32, true), nil
	case INT64_TYPE:
		return NewInt(64, true), nil
	case UINT8_TYPE:
		return NewInt(8, false), nil
	case UINT16_TYPE:
		return NewInt(16, false), nil
	case UINT32_TYPE:
		return NewInt(32, false), nil
	case UINT64_TYPE:
		return NewInt(64, false), nil
	case BYTE_TYPE:
		return NewInt(8, false), nil
	case STRING_TYPE:
		return NewStr(), nil
	case FLOAT32_TYPE:
		return NewFloat(32), nil
	case FLOAT64_TYPE:
		return NewFloat(64), nil
	case NULL_TYPE:
		return NewNull(), nil
	case VOID_TYPE:
		return NewVoid(), nil
	default:
		//search for the type
		declaredEnv, err := env.ResolveType(string(typ))
		if err != nil {
			return nil, err
		}
		return declaredEnv.types[string(typ)], nil
	}
}

// valueTypeInterfaceToString converts a ValueTypeInterface to a string representation of builtins.VALUE_TYPE.
// It handles different types such as Array, Struct, Interface, and Fn by recursively converting
// their components to strings and formatting them appropriately.
//
// Parameters:
// - typeName: A ValueTypeInterface representing the type to be converted.
//
// Returns:
// - builtins.VALUE_TYPE: A string representation of the given ValueTypeInterface.
//
// Note: The function currently has a commented-out case for UserDefined types and a default case
// that prints the type and value of unhandled cases.
func valueTypeInterfaceToString(typeName ValueTypeInterface) builtins.TC_TYPE {
	switch t := typeName.(type) {
	case Array:
		return "[]" + valueTypeInterfaceToString(t.ArrayType)
	case Struct:
		return builtins.TC_TYPE(t.StructName)
	case Interface:
		return builtins.TC_TYPE(t.InterfaceName)
	case Fn:
		ParamStrs := ""
		for i, param := range t.Params {
			ParamStrs += param.Name
			if param.IsOptional {
				ParamStrs += "?: "
			} else {
				ParamStrs += ": "
			}
			ParamStrs += string(valueTypeInterfaceToString(param.Type))
			if i != len(t.Params)-1 {
				ParamStrs += ", "
			}
		}
		ReturnStr := string(valueTypeInterfaceToString(t.Returns))
		if ReturnStr != "" {
			ReturnStr = " -> " + ReturnStr
		}
		return builtins.TC_TYPE(fmt.Sprintf("fn(%s)%s", ParamStrs, ReturnStr))
	//case UserDefined:
	//	return valueTypeInterfaceToString(t.TypeDef)
	case Map:
		return builtins.TC_TYPE(fmt.Sprintf("map[%s]%s", valueTypeInterfaceToString(t.KeyType), valueTypeInterfaceToString(t.ValueType)))
	default:
		return t.DType()
	}
}

// MatchTypes compares the expected and provided ValueTypeInterface types.
// If the types do not match, it returns an error indicating the mismatch.
// If the expected type is an interface, it checks the method implementations
// of the provided type against the expected type.
//
// Parameters:
//   - expected: The expected ValueTypeInterface type.
//   - provided: The provided ValueTypeInterface type.
//   - filePath: The file path where the type check is being performed.
//   - lineStart: The starting line number of the type check.
//   - lineEnd: The ending line number of the type check.
//   - colStart: The starting column number of the type check.
//   - colEnd: The ending column number of the type check.
//
// Returns:
//   - error: An error if the types do not match, otherwise nil.
func MatchTypes(expected, provided ValueTypeInterface) error {

	expectedType := valueTypeInterfaceToString(expected)
	gotType := valueTypeInterfaceToString(provided)

	if expectedType != gotType {
		if expected.DType() == INTERFACE_TYPE {
			return checkMethodsImplementations(expected, provided)
		}
		return fmt.Errorf("expected %s, got %s", expectedType, gotType)
	}
	return nil
}

func CheckLValue(node ast.Node, env *TypeEnvironment) error {
	//if not constant and is IdentifierExpr
	switch t := node.(type) {
	case ast.IdentifierExpr:

		if _, ok := env.types[t.Name]; ok {
			return errors.New("type")
		}

		//find the declaredEnv where the variable was declared
		declaredEnv, err := env.ResolveVar(t.Name)
		if err != nil {
			return err
		}
		if !declaredEnv.constants[t.Name] {
			return nil
		} else {
			return errors.New("constant")
		}
	case ast.Indexable:
		return CheckLValue(t.Container, env)
	case ast.StructPropertyAccessExpr:
		return CheckLValue(t.Object, env)
	default:
		return fmt.Errorf("invalid lvalue")
	}
}

func IsNumberType(operand ValueTypeInterface) bool {
	switch operand.(type) {
	case Int, Float:
		return true
	default:
		return false
	}
}

// EvaluateTypeName evaluates the given DataType and returns a corresponding ValueTypeInterface.
// It handles different types of DataType such as ArrayType, FunctionType, and others.
//
// Parameters:
//   - dtype: The DataType to be evaluated.
//   - env: The TypeEnvironment in which the DataType is evaluated.
//
// Returns:
//   - A ValueTypeInterface representing the evaluated type.
//
// The function performs the following steps
//  1. If the dtype is an ArrayType, it recursively evaluates the element type and returns an Array.
//  2. If the dtype is a FunctionType, it evaluates the parameter types and return type, creates a new function scope, and returns a Fn.
//  3. If the dtype is nil, it returns a Void type.
//  4. For other types, it attempts to create a ValueTypeInterface and handles any errors that occur.
func EvaluateTypeName(dtype ast.DataType, env *TypeEnvironment) ValueTypeInterface {
	switch t := dtype.(type) {
	case ast.ArrayType:
		val := EvaluateTypeName(t.ArrayType, env)
		arr := Array{
			DataType:  builtins.ARRAY,
			ArrayType: val,
		}
		return arr
	case ast.FunctionType:
		var params []FnParam
		for _, param := range t.Parameters {
			paramType := EvaluateTypeName(param.Type, env)
			params = append(params, FnParam{
				Name:       param.Identifier.Name,
				IsOptional: param.IsOptional,
				Type:       paramType,
			})
		}

		returns := EvaluateTypeName(t.ReturnType, env)

		scope := NewTypeENV(env, FUNCTION_SCOPE, fmt.Sprintf("_FN_%s", RandStringRunes(10)), env.filePath)

		return Fn{
			DataType:      builtins.FUNCTION,
			Params:        params,
			Returns:       returns,
			FunctionScope: *scope,
		}
	case ast.MapType:
		keyType := EvaluateTypeName(t.KeyType, env)
		valueType := EvaluateTypeName(t.ValueType, env)
		return NewMap(keyType, valueType)
	case ast.UserDefinedType:
		typename := t.AliasName
		val, err := stringToValueTypeInterface(builtins.TC_TYPE(typename), env)
		if err != nil {
			errgen.AddError(env.filePath, dtype.StartPos().Line, dtype.EndPos().Line, dtype.StartPos().Column, dtype.EndPos().Column, err.Error())
		}
		if val == nil {
			errgen.DisplayErrors()
		}
		return val
	case nil:
		return NewVoid()
	default:
		val, err := stringToValueTypeInterface(builtins.TC_TYPE(t.Type()), env)
		if err != nil {
			errgen.AddError(env.filePath, dtype.StartPos().Line, dtype.EndPos().Line, dtype.StartPos().Column, dtype.EndPos().Column, err.Error())
		}
		if val == nil {
			errgen.DisplayErrors()
		}
		return val
	}
}
