package analyzer

import (
	"errors"
	"fmt"
	"math/rand"

	//"reflect"
	"time"

	"walrus/ast"
	"walrus/builtins"
	"walrus/errgen"
	"walrus/utils"
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
		declaredEnv, err := env.resolveVar(t.Name)
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
		return evalArray(t, env)
	case ast.FunctionType:
		return evalFn(t, env)
	case ast.MapType:
		return evalMap(t, env)
	case ast.MaybeType:
		return NewMaybe(evaluateTypeName(t.MaybeType, env))
	case ast.UserDefinedType:
		return evalUD(t, env)
	case nil:
		return NewVoid()
	default:
		return evalDefaultType(dtype, env)
	}
}

func evalDefaultType(defaultType ast.DataType, env *TypeEnvironment) TcValue {
	val, err := getTypeDefinition(string(defaultType.Type())) // need to get the most deep type
	if err != nil || val == nil {
		errgen.AddError(env.filePath, defaultType.StartPos().Line, defaultType.EndPos().Line, defaultType.StartPos().Column, defaultType.EndPos().Column, err.Error()).ErrorLevel(errgen.CRITICAL)
	}
	return val
}

func evalUD(analyzedUD ast.UserDefinedType, env *TypeEnvironment) TcValue {
	typename := analyzedUD.AliasName
	val, err := getTypeDefinition(typename) // need to get the most deep type
	if err != nil || val == nil {
		errgen.AddError(env.filePath, analyzedUD.StartPos().Line, analyzedUD.EndPos().Line, analyzedUD.StartPos().Column, analyzedUD.EndPos().Column, err.Error()).ErrorLevel(errgen.CRITICAL)
	}
	return val
}

func evalArray(analyzedArray ast.ArrayType, env *TypeEnvironment) TcValue {
	val := evaluateTypeName(analyzedArray.ArrayType, env)
	arr := Array{
		DataType:  builtins.ARRAY,
		ArrayType: val,
	}
	return arr
}

func evalFn(analyzedFunctionType ast.FunctionType, env *TypeEnvironment) TcValue {
	var params []FnParam
	for _, param := range analyzedFunctionType.Parameters {
		paramType := evaluateTypeName(param.Type, env)
		params = append(params, FnParam{
			Name:       param.Identifier.Name,
			IsOptional: param.IsOptional,
			Type:       paramType,
		})
	}

	returns := evaluateTypeName(analyzedFunctionType.ReturnType, env)

	scope := NewTypeENV(env, FUNCTION_SCOPE, fmt.Sprintf("_FN_%s", RandStringRunes(10)), env.filePath)

	return Fn{
		DataType:      builtins.FUNCTION,
		Params:        params,
		Returns:       returns,
		FunctionScope: *scope,
	}
}

func evalMap(analyzedMap ast.MapType, env *TypeEnvironment) TcValue {
	if analyzedMap.Map.Name == "map" {
		keyType := evaluateTypeName(analyzedMap.KeyType, env)
		valueType := evaluateTypeName(analyzedMap.ValueType, env)
		return NewMap(keyType, valueType)
	} else {
		//find the name in the type definition

		val, err := getTypeDefinition(analyzedMap.Map.Name) // need to get the most deep type
		if err != nil {
			errgen.AddError(env.filePath, analyzedMap.StartPos().Line, analyzedMap.EndPos().Line, analyzedMap.StartPos().Column, analyzedMap.EndPos().Column, err.Error()).ErrorLevel(errgen.NORMAL)
		}

		if mapVal, ok := val.(Map); ok {
			return NewMap(mapVal.KeyType, mapVal.ValueType)
		}

		errgen.AddError(env.filePath, analyzedMap.StartPos().Line, analyzedMap.EndPos().Line, analyzedMap.StartPos().Column, analyzedMap.EndPos().Column, fmt.Sprintf("'%s' is not a map", analyzedMap.Map.Name)).ErrorLevel(errgen.CRITICAL)

		return NewVoid()
	}
}

func matchTypes(expectedType, providedType TcValue) error {

	unwrappedExpected := unwrapType(expectedType)
	unwrappedProvided := unwrapType(providedType)

	switch t := unwrappedExpected.(type) {
	case Interface:
		errs := checkMethodsImplementations(unwrappedExpected, unwrappedProvided)
		if len(errs) > 0 {
			msgs := fmt.Sprintf("cannot use type '%s' as interface '%s'\n", tcValueToString(providedType), tcValueToString(expectedType))
			for _, err := range errs {
				msgs += utils.ORANGE.Sprintln(" - " + err.Error())
			}
			return errors.New(msgs)
		}
		return nil
	case Maybe:
		if unwrapType(t.MaybeType).DType() == unwrappedProvided.DType() || unwrappedProvided.DType() == builtins.NULL {
			return nil
		}
	}

	expectedStr := tcValueToString(unwrappedExpected)
	providedStr := tcValueToString(unwrappedProvided)

	if expectedStr != providedStr {
		return fmt.Errorf("cannot assign value of type '%s' to type '%s'", providedStr, expectedStr)
	}

	return nil
}

func tcValueToString(val TcValue) string {
	switch t := val.(type) {
	case Array:
		return fmt.Sprintf("[]%s", tcValueToString(t.ArrayType))
	case Struct:
		return t.StructName
	case Interface:
		return t.InterfaceName
	case Fn:
		return functionSignatureString(t)
	case Map:
		return fmt.Sprintf("map[%s]%s", tcValueToString(t.KeyType), tcValueToString(t.ValueType))
	case Maybe:
		return fmt.Sprintf("maybe{%s}", tcValueToString(t.MaybeType))
	case UserDefined:
		return tcValueToString(unwrapType(t.TypeDef))
	default:
		if t == nil {
			return "void"
		}
		return string(t.DType())
	}
}

func functionSignatureString(fn Fn) string {
	ParamStrs := ""
	for i, param := range fn.Params {
		ParamStrs += param.Name
		if param.IsOptional {
			ParamStrs += "?: "
		} else {
			ParamStrs += ": "
		}
		ParamStrs += string(tcValueToString(param.Type))
		if i != len(fn.Params)-1 {
			ParamStrs += ", "
		}
	}
	ReturnStr := string(tcValueToString(fn.Returns))
	if ReturnStr != "" {
		ReturnStr = " -> " + ReturnStr
	}
	return fmt.Sprintf("fn(%s)%s", ParamStrs, ReturnStr)
}

func checkMethodsImplementations(expected, provided TcValue) []error {

	//check if the provided type implements the interface
	errs := []error{}

	var interfaceType Interface
	interfaceType, ok := expected.(Interface)
	if !ok {
		return []error{fmt.Errorf("type must be an interface")}
	}

	var structType Struct
	structType, ok = provided.(Struct)
	if !ok {
		return []error{fmt.Errorf("type must be a struct")}
	}

	// check if all methods are implemented
	for methodName, method := range interfaceType.Methods {
		// check if method is present in the struct's variables
		methodVal, ok := structType.StructScope.variables[methodName]
		if !ok {
			//return fmt.Errorf("struct '%s' did not implement method '%s' of interface '%s'",
			//	structType.StructName, methodName, interfaceType.InterfaceName)
			errs = append(errs, fmt.Errorf("method '%s' is not implemented for interface '%s'", methodName, interfaceType.InterfaceName))
			continue
		}

		// check if the method is a function
		methodFn, ok := methodVal.(StructMethod)
		if !ok {
			//return fmt.Errorf("'%s' on struct '%s' is not a valid method for interface '%s'",
			//	methodName, structType.StructName, interfaceType.InterfaceName)
			errs = append(errs, fmt.Errorf("'%s' is expected to be a method", methodName))
			continue
		}

		// check the return type and parameters
		for i, param := range method.Params {
			expectedParam := tcValueToString(param.Type)
			providedParam := tcValueToString(methodFn.Fn.Params[i].Type)
			if expectedParam != providedParam {
				//return fmt.Errorf("method '%s' found for interface '%s' but parameter missmatch", methodName, interfaceType.InterfaceName)
				errs = append(errs, fmt.Errorf("method '%s' found for interface '%s' but parameter missmatch", methodName, interfaceType.InterfaceName))
			}
		}

		//check the return type
		expectedReturn := tcValueToString(method.Returns)
		providedReturn := tcValueToString(methodFn.Fn.Returns)
		if expectedReturn != providedReturn {
			//return fmt.Errorf("method '%s' found for interface '%s' but return type mismatched", methodName, interfaceType.InterfaceName)
			errs = append(errs, fmt.Errorf("method '%s' found for interface '%s' but return type mismatched", methodName, interfaceType.InterfaceName))
		}
	}

	return errs
}
