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
		//errgen.MakeError(filePath, lineStart, lineEnd, colStart, colEnd, "type must be a struct").DisplayWithPanic()
		errgen.AddError(filePath, lineStart, lineEnd, colStart, colEnd, "type must be a struct")
	}

	for name, method := range expectedMethods {
		// check if method is present in the struct's variables and is a function
		methodVal, ok := structType.StructScope.variables[name]
		if !ok {
			//errgen.MakeError(filePath, lineStart, lineEnd, colStart, colEnd,
			//	fmt.Sprintf("method '%s' not found in struct '%s' (expected implementation for interface '%s')",
			//		name, provided.(Struct).StructName, expected.(Interface).InterfaceName)).DisplayWithPanic()
			errgen.AddError(filePath, lineStart, lineEnd, colStart, colEnd,
				fmt.Sprintf("method '%s' not found in struct '%s' (expected implementation for interface '%s')",
					name, provided.(Struct).StructName, expected.(Interface).InterfaceName))
		}
		methodFn, ok := methodVal.(StructMethod)
		if !ok {
			//errgen.MakeError(filePath, lineStart, lineEnd, colStart, colEnd,
			//	fmt.Sprintf("'%s' in struct '%s' is not a valid method (expected implementation for interface '%s')",
			//		name, provided.(Struct).StructName, expected.(Interface).InterfaceName)).DisplayWithPanic()
			errgen.AddError(filePath, lineStart, lineEnd, colStart, colEnd,
				fmt.Sprintf("'%s' in struct '%s' is not a valid method (expected implementation for interface '%s')",
					name, provided.(Struct).StructName, expected.(Interface).InterfaceName))
		}

		// check the return type and parameters
		for i, param := range method.Params {
			expectedParam := valueTypeInterfaceToString(param.Type)
			providedParam := valueTypeInterfaceToString(methodFn.Fn.Params[i].Type)
			if expectedParam != providedParam {
				//errgen.MakeError(filePath, lineStart, lineEnd, colStart, colEnd, fmt.Sprintf("method '%s' found for interface '%s' but parameter missmatch", name, expected.(Interface).InterfaceName)).DisplayWithPanic()
				errgen.AddError(filePath, lineStart, lineEnd, colStart, colEnd, fmt.Sprintf("method '%s' found for interface '%s' but parameter missmatch", name, expected.(Interface).InterfaceName))
			}
		}

		//check the return type
		expectedReturn := valueTypeInterfaceToString(method.Returns)
		providedReturn := valueTypeInterfaceToString(methodFn.Fn.Returns)
		if expectedReturn != providedReturn {
			//errgen.MakeError(filePath, lineStart, lineEnd, colStart, colEnd, fmt.Sprintf("method '%s' found for interface '%s' but return type missmatch", name, expected.(Interface).InterfaceName)).DisplayWithPanic()
			errgen.AddError(filePath, lineStart, lineEnd, colStart, colEnd, fmt.Sprintf("method '%s' found for interface '%s' but return type missmatch", name, expected.(Interface).InterfaceName))
		}
	}
}

func checkImplStmt(implStmt ast.ImplStmt, env *TypeEnvironment) ValueTypeInterface {

	//scope must be global
	if env.scopeType != GLOBAL_SCOPE {
		//errgen.MakeError(env.filePath, implStmt.Start.Line, implStmt.End.Line, implStmt.Start.Column, implStmt.End.Column, "implement statement must be at global scope").DisplayWithPanic()
		errgen.AddError(env.filePath, implStmt.Start.Line, implStmt.End.Line, implStmt.Start.Column, implStmt.End.Column, "implement statement must be at global scope")
		return NewVoid()
	}

	// check if the type to implement exists
	structValue, err := env.GetTypeFromEnv(implStmt.ImplFor.Name)
	if err != nil {
		//errgen.MakeError(env.filePath, implStmt.Start.Line, implStmt.End.Line, implStmt.Start.Column, implStmt.End.Column, err.Error()).DisplayWithPanic()
		errgen.AddError(env.filePath, implStmt.Start.Line, implStmt.End.Line, implStmt.Start.Column, implStmt.End.Column, err.Error())
	}

	// type must be a struct
	implForType, ok := structValue.(Struct)
	if !ok {
		//fmt.Printf("Type %v\n", structValue)
		//errgen.MakeError(env.filePath, implStmt.Start.Line, implStmt.End.Line, implStmt.Start.Column, implStmt.End.Column, "only structs can be implemented").DisplayWithPanic()
		errgen.AddError(env.filePath, implStmt.Start.Line, implStmt.End.Line, implStmt.Start.Column, implStmt.End.Column, "only structs can be implemented")
	}

	fmt.Printf("Implementing for type %s\n", valueTypeInterfaceToString(implForType))

	//add the methods to the struct's environment
	for name, method := range implStmt.Methods {

		// if the method name is in the struct's elements, throw an error
		if _, ok := implForType.StructScope.variables[name]; ok {
			//errgen.MakeError(env.filePath, method.Start.Line, method.End.Line, method.Start.Column, method.End.Column, fmt.Sprintf("name '%s' already exists in struct", name)).DisplayWithPanic()
			errgen.AddError(env.filePath, method.Start.Line, method.End.Line, method.Start.Column, method.End.Column, fmt.Sprintf("name '%s' already exists in struct", name))
		}

		fnEnv := NewTypeENV(&implForType.StructScope, FUNCTION_SCOPE, name, implForType.StructScope.filePath)

		//check the parameters and declare them
		params := checkandDeclareParamaters(method.Params, fnEnv)

		//check the return type
		returnType := EvaluateTypeName(method.ReturnType, fnEnv)

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
			//errgen.MakeError(env.filePath, method.Start.Line, method.End.Line, method.Start.Column, method.End.Column, err.Error()).DisplayWithPanic()
			errgen.AddError(env.filePath, method.Start.Line, method.End.Line, method.Start.Column, method.End.Column, err.Error())
		}

		//check the function body
		for _, stmt := range method.Body.Contents {
			CheckAST(stmt, fnEnv)
		}

	}

	return NewVoid()
}
