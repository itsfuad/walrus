package typechecker

import (
	"fmt"
	"walrus/ast"
	"walrus/errgen"
)

func checkMethodsImplementations(expected, provided ValueTypeInterface, filePath string, lineStart, lineEnd, colStart, colEnd int) {

	//check if the provided type implements the interface
	expectedMethods := expected.(Interface).Methods
	structType, ok := provided.(Struct)

	if !ok {
		//value must be a struct. display "type must be a struct" error
		errgen.MakeError(filePath, lineStart, lineEnd, colStart, colEnd, "type must be a struct").Display()
	}

	for name, method := range expectedMethods {
		//check if the method is implemented
		if _, ok := structType.Methods[name]; !ok {
			//method not implemented. display "method not implemented" error
			errgen.MakeError(filePath, lineStart, lineEnd, colStart, colEnd, fmt.Sprintf("method '%s' not implemented for interface '%s' on struct '%s'", name, expected.(Interface).InterfaceName, provided.(Struct).StructName)).Display()
		}
		// check the return type and parameters
		for i, param := range method.Params {
			expectedParam := valueTypeInterfaceToString(param.Type)
			providedParam := valueTypeInterfaceToString(structType.Methods[name].Params[i].Type)
			if expectedParam != providedParam {
				errgen.MakeError(filePath, lineStart, lineEnd, colStart, colEnd, fmt.Sprintf("method '%s' found for interface '%s' but parameter missmatch", name, expected.(Interface).InterfaceName)).Display()
			}
		}
		//check the return type
		expectedReturn := valueTypeInterfaceToString(method.Returns)
		providedReturn := valueTypeInterfaceToString(structType.Methods[name].Returns)
		if expectedReturn != providedReturn {
			errgen.MakeError(filePath, lineStart, lineEnd, colStart, colEnd, fmt.Sprintf("method '%s' found for interface '%s' but return type missmatch", name, expected.(UserDefined).TypeDef.(Interface).InterfaceName)).Display()
		}
	}
}

func checkImplStmt(implStmt ast.ImplStmt, env *TypeEnvironment) ValueTypeInterface {
	// Resolve the type to implement
	structEnv, err := env.ResolveType(implStmt.ImplFor.Name)
	if err != nil {
		errgen.MakeError(env.filePath, implStmt.Start.Line, implStmt.End.Line, implStmt.Start.Column, implStmt.End.Column, err.Error()).Display()
		return nil
	}

	//if scope is not global, throw an error
	if structEnv.scopeType != GLOBAL_SCOPE {
		errgen.MakeError(env.filePath, implStmt.Start.Line, implStmt.End.Line, implStmt.Start.Column, implStmt.End.Column, "implement statement must be at global scope").Display()
		return nil
	}

	// type must be a struct
	if structEnv.types[implStmt.ImplFor.Name].(UserDefined).TypeDef.DType() != STRUCT_TYPE {
		errgen.MakeError(env.filePath, implStmt.Start.Line, implStmt.End.Line, implStmt.Start.Column, implStmt.End.Column, "can only implement for structs").Display()
		return nil
	}

	implForType := structEnv.types[implStmt.ImplFor.Name].(UserDefined).TypeDef.(Struct)

	fmt.Printf("Implementing for type %s\n", valueTypeInterfaceToString(implForType))

	//add the methods to the struct
	for name, method := range implStmt.Methods {

		// if the method name is in the struct's elements, throw an error
		if _, ok := implForType.Elements[name]; ok {
			errgen.MakeError(structEnv.filePath, method.Start.Line, method.End.Line, method.Start.Column, method.End.Column, fmt.Sprintf("cannot use method name '%s'. same named property already exists in struct", name)).Display()
		}

		fnEnv := NewTypeENV(structEnv, FUNCTION_SCOPE, name, structEnv.filePath)

		implForType.Methods[name] = StructMethod{
			IsPrivate: method.IsPrivate,
			Fn: Fn{
				DataType:      FUNCTION_TYPE,
				Params:        checkParamaters(method.Params, fnEnv),
				Returns:       EvaluateTypeName(method.ReturnType, fnEnv),
				FunctionScope: *fnEnv,
			},
		}
	}

	fmt.Printf("Updated struct methods: %v\n", implForType.Methods)

	//update the struct in the environment
	structEnv.types[implStmt.ImplFor.Name] = UserDefined{
		DataType: USER_DEFINED_TYPE,
		TypeName: implStmt.ImplFor.Name,
		TypeDef:  implForType,
	}

	fmt.Printf("Updated struct at scope %s\n", structEnv.scopeName)

	return implForType
}
