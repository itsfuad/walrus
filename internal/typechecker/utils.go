package typechecker

import (
	//Standard packages
	"errors"
	"fmt"
	"math/rand"
	"time"

	//Walrus packages
	"walrus/internal/ast"
	"walrus/internal/builtins"
	"walrus/internal/report"
	"walrus/internal/utils"
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

func isNumberType(operand Tc) bool {
	switch operand.(type) {
	case Int, Float:
		return true
	default:
		return false
	}
}

func isIntType(operand Tc) bool {
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
func evaluateTypeName(dtype ast.DataType, env *TypeEnvironment) Tc {
	switch t := dtype.(type) {
	case ast.ArrayType:
		return evalArray(t, env)
	case ast.FunctionType:
		return evalFn(t, env)
	case ast.MapType:
		return evalMap(t, env)
	case ast.UserDefinedType:
		return evalUD(t, env)
	case ast.StructType:
		return evalStruct(t, env)
	case nil:
		return NewVoid()
	default:
		return evalDefaultType(dtype, env)
	}
}

func evalStruct(s ast.StructType, env *TypeEnvironment) Struct {

	name := "struct { "
	for i, prop := range s.Properties {

		propType := evaluateTypeName(prop.PropType, env)

		name += prop.Prop.Name + ": " + tcToString(propType)
		if i != len(s.Properties)-1 {
			name += ", "
		}
	}

	name += " }"

	typ := checkStructTypeDecl(name, s, env)

	// declare the struct type
	declareType(name, typ) // if the struct is already declared, it will be skipped

	return typ
}

func evalDefaultType(defaultType ast.DataType, env *TypeEnvironment) Tc {
	val, err := getTypeDefinition(string(defaultType.Type())) // need to get the most deep type
	if err != nil || val == nil {
		report.Add(env.filePath, defaultType.StartPos().Line, defaultType.EndPos().Line, defaultType.StartPos().Column, defaultType.EndPos().Column, err.Error()).SetLevel(report.CRITICAL_ERROR)
	}
	return val
}

// evalUD evaluates a user-defined type and returns a type-checked user-defined type.
// It takes an analyzed user-defined type of type ast.UserDefinedType and a type environment.
// It returns a type-checked user-defined type (Tc) with the evaluated user-defined type.
//
// Parameters:
// - analyzedUD: the AST representation of the user-defined type to be evaluated.
// - env: the type environment used for type evaluation.
//
// Returns:
// - Tc: a type-checked user-defined type with the evaluated user-defined type.
func evalUD(analyzedUD ast.UserDefinedType, env *TypeEnvironment) Tc {
	typename := analyzedUD.AliasName
	val, err := getTypeDefinition(typename) // need to get the most deep type
	if err != nil || val == nil {
		report.Add(env.filePath, analyzedUD.StartPos().Line, analyzedUD.EndPos().Line, analyzedUD.StartPos().Column, analyzedUD.EndPos().Column, err.Error()).SetLevel(report.CRITICAL_ERROR)
	}
	return val
}

// evalArray evaluates an AST array type and returns a type-checked array.
// It takes an analyzed array of type ast.ArrayType and a type environment.
// It returns a type-checked array (Tc) with the evaluated array type.
//
// Parameters:
// - analyzedArray: the AST representation of the array type to be evaluated.
// - env: the type environment used for type evaluation.
//
// Returns:
// - Tc: a type-checked array with the evaluated array type.
func evalArray(analyzedArray ast.ArrayType, env *TypeEnvironment) Tc {
	val := evaluateTypeName(analyzedArray.ArrayType, env)
	arr := Array{
		DataType:  builtins.ARRAY,
		ArrayType: val,
	}
	return arr
}

// evalFn evaluates a given function type within a specified type environment and returns a type-checked function.
//
// Parameters:
// - analyzedFunctionType: The AST representation of the function type to be analyzed.
// - env: The type environment in which the function type is evaluated.
//
// Returns:
// - Tc: A type-checked function containing the function's data type, parameters, return type, and function scope.
//
// The function performs the following steps:
// 1. Creates a new type environment for the function scope.
// 2. Iterates over the function parameters and checks for duplicate parameter names.
// 3. Evaluates the type of each parameter and adds it to the list of function parameters.
// 4. Evaluates the return type of the function.
// 5. Constructs and returns a type-checked function with the evaluated parameters and return type.
func evalFn(analyzedFunctionType ast.FunctionType, env *TypeEnvironment) Tc {

	scope := NewTypeENV(env, FUNCTION_SCOPE, fmt.Sprintf("_FN_%s", RandStringRunes(10)), env.filePath)

	var params []FnParam
	for _, param := range analyzedFunctionType.Parameters {
		//check if the parameter is already declared
		if _, found := utils.Some(params, func(p FnParam) bool {
			return p.Name == param.Identifier.Name
		}); found {
			report.Add(scope.filePath, param.Identifier.Start.Line, param.Identifier.End.Line, param.Identifier.Start.Column, param.Identifier.End.Column,
				fmt.Sprintf("parameter '%s' is already defined", param.Identifier.Name)).SetLevel(report.CRITICAL_ERROR)
		}

		paramType := evaluateTypeName(param.Type, scope)
		params = append(params, FnParam{
			Name: param.Identifier.Name,
			Type: paramType,
		})
	}

	returns := evaluateTypeName(analyzedFunctionType.ReturnType, scope)

	return Fn{
		DataType:      builtins.FUNCTION,
		Params:        params,
		Returns:       returns,
		FunctionScope: *scope,
	}
}

// evalMap evaluates an AST representation of a map type and returns a type-checked
// representation of the map. It checks if the map type is a built-in "map" type or
// a user-defined type.
//
// Parameters:
// - analyzedMap: The AST representation of the map type to be evaluated.
// - env: The type environment containing type information and definitions.
//
// Returns:
// - Tc: The type-checked representation of the map.
//
// If the map type is a built-in "map", it evaluates the key and value types and
// returns a new map type. If the map type is user-defined, it retrieves the type
// definition and checks if it is a map type. If it is, it returns a new map type
// based on the retrieved definition. If the type is not a map, it reports an error
// and returns a void type.
func evalMap(analyzedMap ast.MapType, env *TypeEnvironment) Tc {
	if analyzedMap.Map.Name == "map" {
		keyType := evaluateTypeName(analyzedMap.KeyType, env)
		valueType := evaluateTypeName(analyzedMap.ValueType, env)
		return NewMap(keyType, valueType)
	} else {
		//find the name in the type definition
		val, err := getTypeDefinition(analyzedMap.Map.Name) // need to get the most deep type
		if err != nil {
			report.Add(env.filePath, analyzedMap.StartPos().Line, analyzedMap.EndPos().Line, analyzedMap.StartPos().Column, analyzedMap.EndPos().Column, err.Error()).SetLevel(report.NORMAL_ERROR)
		}

		if mapVal, ok := val.(Map); ok {
			return NewMap(mapVal.KeyType, mapVal.ValueType)
		}

		report.Add(env.filePath, analyzedMap.StartPos().Line, analyzedMap.EndPos().Line, analyzedMap.StartPos().Column, analyzedMap.EndPos().Column, fmt.Sprintf("'%s' is not a map", analyzedMap.Map.Name)).SetLevel(report.CRITICAL_ERROR)

		return NewVoid()
	}
}

// validateTypeCompatibility checks if the provided type is compatible with the expected type.
// It first unwraps both types and then performs a type switch on the unwrapped expected type.
// If the expected type is an Interface, it checks if the provided type implements the required methods.
// Otherwise, it converts both types to their string representations and compares them.
// If the types are not compatible, it returns an error indicating the type mismatch.
//
// Parameters:
//   - expectedType: The type that is expected.
//   - providedType: The type that is provided.
//
// Returns:
//   - error: An error indicating the type incompatibility, or nil if the types are compatible.
func validateTypeCompatibility(expectedType, providedType Tc) error {

	unwrappedExpected := unwrapType(expectedType)
	unwrappedProvided := unwrapType(providedType)

	switch unwrappedExpected.(type) {
	case Interface:
		return checkMethodsImplementations(unwrappedProvided, unwrappedExpected)
	}

	expectedStr := tcToString(unwrappedExpected)
	providedStr := tcToString(unwrappedProvided)

	if expectedStr != providedStr {
		return fmt.Errorf("cannot assign value of type '%s' to type '%s'", providedStr, expectedStr)
	}

	return nil
}

// tcToString converts a type-checker (Tc) value to its string representation.
// It handles various types including Array, Struct, Interface, Fn, Map, Maybe,
// and UserDefined. For each type, it returns a formatted string that represents
// the type. If the type is nil, it returns "void". For other types, it returns
// the string representation of the type's DType.
//
// Parameters:
// - val: The type-checker (Tc) value to be converted to a string.
//
// Returns:
// - A string representation of the type-checker value.
func tcToString(val Tc) string {
	switch t := val.(type) {
	case Array:
		return fmt.Sprintf("[]%s", tcToString(t.ArrayType))
	case Struct:
		return t.StructName
	case Interface:
		return t.InterfaceName
	case Fn:
		return functionSignatureString(t)
	case Map:
		return fmt.Sprintf("map[%s]%s", tcToString(t.KeyType), tcToString(t.ValueType))
	case Maybe:
		return fmt.Sprintf("maybe{%s}", tcToString(t.MaybeType))
	case UserDefined:
		return tcToString(unwrapType(t.TypeDef))
	default:
		if t == nil {
			return "void"
		}
		return string(t.DType())
	}
}

// functionSignatureString generates a string representation of a function's signature.
// It takes a function `fn` of type `Fn` as input and returns a string that describes
// the function's parameters and return type in the format: "fn(param1: type1, param2: type2) -> returnType".
// If the function has no return type, the return type part is omitted.
//
// Parameters:
// - fn: The function whose signature is to be represented as a string.
//
// Returns:
// - A string representing the function's signature.
func functionSignatureString(fn Fn) string {
	ParamStrs := ""
	for i, param := range fn.Params {
		ParamStrs += param.Name
		ParamStrs += ": "
		ParamStrs += string(tcToString(param.Type))
		if i != len(fn.Params)-1 {
			ParamStrs += ", "
		}
	}
	ReturnStr := string(tcToString(fn.Returns))
	if ReturnStr != "" {
		ReturnStr = " -> " + ReturnStr
	}
	return fmt.Sprintf("fn(%s)%s", ParamStrs, ReturnStr)
}

// checkMethodsImplementations checks if the provided type `src` implements the interface `dest`.
// It returns an error if `src` does not implement `dest` or if `dest` is not an interface.
//
// Parameters:
// - src: The type to check for interface implementation. It can be a struct or an interface.
// - dest: The interface type that `src` should implement.
//
// Returns:
//   - error: An error if `src` does not implement `dest` or if `dest` is not an interface.
//     The error message includes details about the mismatch or invalid type.
func checkMethodsImplementations(src, dest Tc) error {

	//check if the provided type implements the interface
	expectedTypeName := tcToString(dest)
	errMsg := fmt.Sprintf("cannot use type '%s' as interface '%s'\n", tcToString(src), expectedTypeName)
	errs := make([]error, 0)

	var interfaceType Interface
	interfaceType, ok := dest.(Interface)
	if !ok {
		return errors.New(errMsg + report.TreeFormatString(fmt.Sprintf("type '%s' must be an interface", expectedTypeName)))
	}

	//check if the provided type is a struct or interface
	switch t := src.(type) {
	case Struct:
		handleStructDest(t, interfaceType, errs)
	case Interface:
		handleInterfaceDest(t, interfaceType, errs)
	default:
		return errors.New(errMsg + report.TreeFormatString("type must be a struct or interface"))
	}

	if len(errs) > 0 {
		return errors.New(errMsg + report.TreeFormatError(errs...).Error())
	}

	return nil
}

// handleStructDest checks if a given struct implements all methods of a specified interface.
// It verifies the presence, parameter types, and return types of each method.
// If any method is missing or has a mismatch in parameters or return types, an error is appended to the errs slice.
//
// Parameters:
// - src: The struct to be checked.
// - destInterface: The interface that the struct should implement.
// - errs: A slice to which any errors found during the check will be appended.
func handleStructDest(src Struct, destInterface Interface, errs []error) {

	// check if all methods are present
	for _, interfaceMethod := range destInterface.Methods {
		// check if method is present in the struct's variables
		methodVal, ok := src.StructScope.variables[interfaceMethod.Name]
		if !ok {
			errs = append(errs, fmt.Errorf("missing method '%s' on '%s'", interfaceMethod.Name, src.StructName))
			continue
		}

		// check if the method is a function
		methodFn, ok := methodVal.(StructMethod)
		if !ok {
			errs = append(errs, fmt.Errorf("'%s' is expected to be a method", interfaceMethod.Name))
			continue
		}

		// check the return type and parameters
		for i, param := range interfaceMethod.Method.Params {
			expectedParam := tcToString(param.Type)
			providedParam := tcToString(methodFn.Fn.Params[i].Type)
			if expectedParam != providedParam {
				//return fmt.Errorf("method '%s' found for interface '%s' but parameter missmatch", methodName, interfaceType.InterfaceName)
				errs = append(errs, fmt.Errorf("method '%s', but parameter missmatch", interfaceMethod.Name))
			}
		}

		//check the return type
		expectedReturn := tcToString(interfaceMethod.Method.Returns)
		providedReturn := tcToString(methodFn.Fn.Returns)
		if expectedReturn != providedReturn {
			//return fmt.Errorf("method '%s' found for interface '%s' but return type mismatched", methodName, interfaceType.InterfaceName)
			errs = append(errs, fmt.Errorf("method '%s' found, but return type mismatched", interfaceMethod.Name))
		}
	}
}

// handleInterfaceDest checks if all methods in the destination interface are present and compatible with the source interface.
// It validates the parameters and return types of each method.
//
// Parameters:
// - src: The source interface to be checked.
// - destInterface: The destination interface to be checked against.
// - errs: A slice to collect any errors found during the validation process.
//
// The function performs the following checks:
// 1. Ensures that each method in the destination interface is present in the source interface.
// 2. Validates the compatibility of method parameters between the source and destination interfaces.
// 3. Validates the compatibility of return types between the source and destination interfaces.
//
// If any method is missing or has incompatible parameters or return types, an error is appended to the errs slice.
func handleInterfaceDest(src Interface, destInterface Interface, errs []error) {
	// both are interfaces, check if all methods are present and compatible
	for _, interfaceMethod := range destInterface.Methods {
		// check if method is present in the struct's variables
		if method, found := utils.Some(src.Methods, func(m InterfaceMethodType) bool {
			return m.Name == interfaceMethod.Name
		}); found {
			fmt.Printf("checking method %s\n", method.Name)
			//check parameters
			for i, param := range interfaceMethod.Method.Params {
				if err := validateTypeCompatibility(param.Type, method.Method.Params[i].Type); err != nil {
					errs = append(errs, fmt.Errorf("method '%s' found, but parameter missmatch", interfaceMethod.Name))
				}
			}

			//check return type
			if err := validateTypeCompatibility(interfaceMethod.Method.Returns, method.Method.Returns); err != nil {
				errs = append(errs, fmt.Errorf("method '%s' found, but return type mismatched", interfaceMethod.Name))
			}

		} else {
			errs = append(errs, fmt.Errorf("missing method '%s' on '%s'", interfaceMethod.Name, src.InterfaceName))
		}
	}
}
