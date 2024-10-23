package typechecker

import (
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


func stringToValueTypeInterface(typ VALUE_TYPE, env *TypeEnvironment) (ValueTypeInterface, error) {

	builtin, ok := env.builtins[string(typ)]
	if ok {
		return builtin, nil
	}

	//search for the type
	declaredEnv, err := env.ResolveType(string(typ))
	if err != nil {
		return nil, err
	}

	return declaredEnv.types[string(typ)].(UserDefined).TypeDef, nil
}


// valueTypeInterfaceToString converts a ValueTypeInterface to a string representation of VALUE_TYPE.
// It handles different types such as Array, Struct, Interface, and Fn by recursively converting
// their components to strings and formatting them appropriately.
//
// Parameters:
// - typeName: A ValueTypeInterface representing the type to be converted.
//
// Returns:
// - VALUE_TYPE: A string representation of the given ValueTypeInterface.
//
// Note: The function currently has a commented-out case for UserDefined types and a default case
// that prints the type and value of unhandled cases.
func valueTypeInterfaceToString(typeName ValueTypeInterface) VALUE_TYPE {
	switch t := typeName.(type) {
	case Array:
		return "[]" + valueTypeInterfaceToString(t.ArrayType)
	case Struct:
		return VALUE_TYPE(t.StructName)
	case Interface:
		return VALUE_TYPE(t.InterfaceName)
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
		return VALUE_TYPE(fmt.Sprintf("fn(%s)%s", ParamStrs, ReturnStr))
	//case UserDefined:
	//	return valueTypeInterfaceToString(t.TypeDef)
	default:
		fmt.Printf("Default case %T's value %v\n", t, t.DType())
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
func MatchTypes(expected, provided ValueTypeInterface, filePath string, lineStart, lineEnd, colStart, colEnd int) (error) {

	expectedType := valueTypeInterfaceToString(expected)
	gotType := valueTypeInterfaceToString(provided)

	if expectedType != gotType {
		if expected.DType() == INTERFACE_TYPE {
			checkMethodsImplementations(expected, provided, filePath, lineStart, lineEnd, colStart, colEnd)
			return nil
		}
		return fmt.Errorf("expected %s, got %s", expectedType, gotType)
	}
	return nil
}

func IsAssignable(node ast.Node, env *TypeEnvironment) (error) {
	//if not constant and is IdentifierExpr
	switch t := node.(type) {
	case ast.IdentifierExpr:
		//find the declaredEnv where the variable was declared
		declaredEnv, err := env.ResolveVar(t.Name)
		if err != nil {
			return err
		}
		if !declaredEnv.constants[t.Name] {
			return nil
		} else {
			return fmt.Errorf("identifier '%s' is constant", t.Name)
		}
	case ast.ArrayIndexAccess:
		return IsAssignable(t.Array, env)
	case ast.StructPropertyAccessExpr:
		return IsAssignable(t.Object, env)
	default:
		return fmt.Errorf("assignment expression must be a valid lvalue")
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
				Name: 		param.Identifier.Name,
				IsOptional: param.IsOptional,
				Type: 		paramType,
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
	case nil:
		return NewVoid()
	default:
		val, err := stringToValueTypeInterface(VALUE_TYPE(t.Type()), env)
		if err != nil {
			errgen.MakeError(env.filePath, dtype.StartPos().Line, dtype.EndPos().Line, dtype.StartPos().Column, dtype.EndPos().Column, err.Error()).Display()
		}
		return val
	}
}
