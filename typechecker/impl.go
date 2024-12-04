package typechecker

import (
	"fmt"
	"walrus/ast"
	"walrus/errgen"
)

func checkMethodsImplementations(expected, provided TcValue) error {

	//check if the provided type implements the interface
	expectedMethods := expected.(Interface).Methods
	structType, ok := provided.(Struct)
	if !ok {
		return fmt.Errorf("type must be a struct")
	}

	for methodName, method := range expectedMethods {
		// check if method is present in the struct's variables and is a function
		methodVal, ok := structType.StructScope.variables[methodName]
		if !ok {
			return fmt.Errorf("struct '%s' did not implement method '%s' of interface '%s'",
				provided.(Struct).StructName, methodName, expected.(Interface).InterfaceName)
		}
		methodFn, ok := methodVal.(StructMethod)
		if !ok {
			return fmt.Errorf("'%s' on struct '%s' is not a valid method for interface '%s'",
				methodName, provided.(Struct).StructName, expected.(Interface).InterfaceName)
		}

		// check the return type and parameters
		for i, param := range method.Params {
			expectedParam := tcValueToString(param.Type)
			providedParam := tcValueToString(methodFn.Fn.Params[i].Type)
			if expectedParam != providedParam {
				return fmt.Errorf("method '%s' found for interface '%s' but parameter missmatch", methodName, expected.(Interface).InterfaceName)
			}
		}

		//check the return type
		expectedReturn := tcValueToString(method.Returns)
		providedReturn := tcValueToString(methodFn.Fn.Returns)
		if expectedReturn != providedReturn {
			return fmt.Errorf("method '%s' found for interface '%s' but return type mismatched", methodName, expected.(Interface).InterfaceName)
		}
	}

	return nil
}

func checkImplStmt(implStmt ast.ImplStmt, env *TypeEnvironment) TcValue {

	//scope must be global
	if env.scopeType != GLOBAL_SCOPE {
		errgen.AddError(env.filePath, implStmt.Start.Line, implStmt.End.Line, implStmt.Start.Column, implStmt.End.Column, "implement statement must be at global scope", errgen.ERROR_CRITICAL)
		return NewVoid()
	}

	// check if the type to implement exists
	structValue, err := getTypeDefinition(implStmt.ImplFor.Name)
	if err != nil {
		errgen.AddError(env.filePath, implStmt.Start.Line, implStmt.End.Line, implStmt.Start.Column, implStmt.End.Column, err.Error(), errgen.ERROR_CRITICAL)
	}

	// type must be a struct
	implForType, ok := structValue.(Struct)
	if !ok {
		errgen.AddError(env.filePath, implStmt.Start.Line, implStmt.End.Line, implStmt.Start.Column, implStmt.End.Column, "only structs can be implemented", errgen.ERROR_CRITICAL)
	}

	//fmt.Printf("Implementing type %s\n", valueTypeInterfaceToString(implForType))

	//add the methods to the struct's environment
	for name, method := range implStmt.Methods {

		// if the method name is in the struct's elements, throw an error
		if _, ok := implForType.StructScope.variables[name]; ok {

			errgen.AddError(env.filePath, method.Start.Line, method.End.Line, method.Start.Column, method.End.Column, fmt.Sprintf("name '%s' already exists in struct", name), errgen.ERROR_CRITICAL)
		}

		fnEnv := NewTypeENV(&implForType.StructScope, FUNCTION_SCOPE, name, implForType.StructScope.filePath)

		//check the parameters and declare them
		params := checkandDeclareParamaters(method.Params, fnEnv)

		//check the return type
		returnType := evaluateTypeName(method.ReturnType, fnEnv)

		methodToDeclare := StructMethod{
			IsPrivate: method.IsPrivate,
			Fn: Fn{
				DataType:      FUNCTION_TYPE,
				Params:        params,
				Returns:       returnType,
				FunctionScope: *fnEnv,
			},
		}

		//declare the method on the struct's environment and then check the body
		err := implForType.StructScope.DeclareVar(name, methodToDeclare, false, false)
		if err != nil {

			errgen.AddError(env.filePath, method.Start.Line, method.End.Line, method.Start.Column, method.End.Column, err.Error(), errgen.ERROR_CRITICAL)
		}

		//check the function body
		for _, stmt := range method.Body.Contents {
			CheckAST(stmt, fnEnv)
		}

	}

	return NewVoid()
}
