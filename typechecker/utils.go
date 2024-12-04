package typechecker

import (
	"errors"
	"fmt"
	"math/rand"

	//"reflect"
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

func checkLValue(node ast.Node, env *TypeEnvironment) error {
	//if not constant and is IdentifierExpr
	switch t := node.(type) {
	case ast.IdentifierExpr:
		if isTypeDefined(t.Name) {
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
		return checkLValue(t.Container, env)
	case ast.StructPropertyAccessExpr:
		return checkLValue(t.Object, env)
	default:
		return fmt.Errorf("invalid lvalue")
	}
}

func isNumberType(operand TcValue) bool {
	switch operand.(type) {
	case Int, Float:
		return true
	default:
		return false
	}
}

func isIntType(operand TcValue) bool {
	switch operand.(type) {
	case Int:
		return true
	default:
		return false
	}
}

// evaluateTypeName evaluates the given DataType and returns a corresponding ValueTypeInterface.
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
func evaluateTypeName(dtype ast.DataType, env *TypeEnvironment) TcValue {
	switch t := dtype.(type) {
	case ast.ArrayType:
		val := evaluateTypeName(t.ArrayType, env)
		arr := Array{
			DataType:  builtins.ARRAY,
			ArrayType: val,
		}
		return arr
	case ast.FunctionType:
		var params []FnParam
		for _, param := range t.Parameters {
			paramType := evaluateTypeName(param.Type, env)
			params = append(params, FnParam{
				Name:       param.Identifier.Name,
				IsOptional: param.IsOptional,
				Type:       paramType,
			})
		}

		returns := evaluateTypeName(t.ReturnType, env)

		scope := NewTypeENV(env, FUNCTION_SCOPE, fmt.Sprintf("_FN_%s", RandStringRunes(10)), env.filePath)

		return Fn{
			DataType:      builtins.FUNCTION,
			Params:        params,
			Returns:       returns,
			FunctionScope: *scope,
		}
	case ast.MapType:
		keyType := evaluateTypeName(t.KeyType, env)
		valueType := evaluateTypeName(t.ValueType, env)
		return NewMap(keyType, valueType)
	case ast.UserDefinedType:
		typename := t.AliasName
		val, err := getTypeDefinition(typename) // need to get the most deep type
		if err != nil || val == nil {
			errgen.AddError(env.filePath, dtype.StartPos().Line, dtype.EndPos().Line, dtype.StartPos().Column, dtype.EndPos().Column, err.Error(), errgen.ERROR_CRITICAL)
		}
		return val
	case nil:
		return NewVoid()
	default:
		val, err := getTypeDefinition(string(t.Type())) // need to get the most deep type
		if err != nil || val == nil {
			errgen.AddError(env.filePath, dtype.StartPos().Line, dtype.EndPos().Line, dtype.StartPos().Column, dtype.EndPos().Column, err.Error(), errgen.ERROR_CRITICAL)
		}
		return val
	}
}

func matchTypes(expectedType, providedType TcValue) error {

	if expectedType.DType() == builtins.INTERFACE {
		return checkMethodsImplementations(expectedType, providedType)
	}

	unwrappedExpected, err := unwrapType(expectedType)
	if err != nil {
		return err
	}

	unwrappedProvided, err := unwrapType(providedType)
	if err != nil {
		return err
	}

	expectedStr := tcValueToString(unwrappedExpected)
	providedStr := tcValueToString(unwrappedProvided)

	if expectedStr != providedStr {
		return fmt.Errorf("expected type '%s', got '%s'", expectedStr, providedStr)
	}

	return nil
}

func tcValueToString(val TcValue) string {
	switch t := val.(type) {
	case Array:
		return "[]" + tcValueToString(t.ArrayType)
	case Struct:
		return t.StructName
	case Interface:
		return t.InterfaceName
	case Fn:
		ParamStrs := ""
		for i, param := range t.Params {
			ParamStrs += param.Name
			if param.IsOptional {
				ParamStrs += "?: "
			} else {
				ParamStrs += ": "
			}
			ParamStrs += string(tcValueToString(param.Type))
			if i != len(t.Params)-1 {
				ParamStrs += ", "
			}
		}
		ReturnStr := string(tcValueToString(t.Returns))
		if ReturnStr != "" {
			ReturnStr = " -> " + ReturnStr
		}
		return fmt.Sprintf("fn(%s)%s", ParamStrs, ReturnStr)
	case Map:
		return fmt.Sprintf("map[%s]%s", tcValueToString(t.KeyType), tcValueToString(t.ValueType))
	case UserDefined:
		return tcValueToString(t.TypeDef)
	default:
		return string(t.DType())
	}
}
