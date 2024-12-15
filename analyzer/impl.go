package analyzer

import (
	"fmt"
	"walrus/ast"
	"walrus/errgen"
)

func checkImplStmt(implStmt ast.ImplStmt, env *TypeEnvironment) ExprType {

	//scope must be global
	if env.scopeType != GLOBAL_SCOPE {
		errgen.Add(env.filePath, implStmt.Start.Line, implStmt.End.Line, implStmt.Start.Column, implStmt.End.Column, "implement statement must be at global scope").Level(errgen.CRITICAL_ERROR)
		return NewVoid()
	}

	// check if the type to implement exists
	structValue, err := getTypeDefinition(implStmt.ImplFor.Name)
	if err != nil {
		errgen.Add(env.filePath, implStmt.Start.Line, implStmt.End.Line, implStmt.Start.Column, implStmt.End.Column, err.Error()).Level(errgen.CRITICAL_ERROR)
	}

	// type must be a struct
	implForType, ok := structValue.(Struct)
	if !ok {
		errgen.Add(env.filePath, implStmt.Start.Line, implStmt.End.Line, implStmt.Start.Column, implStmt.End.Column, "only structs can be implemented").Level(errgen.CRITICAL_ERROR)
	}

	//add the methods to the struct's environment
	for _, method := range implStmt.Methods {
		name := method.Identifier.Name
		// if the method name is in the struct's elements, throw an error
		if _, ok := implForType.StructScope.variables[name]; ok {
			errgen.Add(env.filePath, method.Start.Line, method.End.Line, method.Start.Column, method.End.Column, fmt.Sprintf("'%s' is already defined in struct '%s'", name, implForType.StructName)).Level(errgen.CRITICAL_ERROR)
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
		err := implForType.StructScope.declareVar(name, methodToDeclare, false, false)
		if err != nil {
			errgen.Add(env.filePath, method.Start.Line, method.End.Line, method.Start.Column, method.End.Column, fmt.Sprintf("cannot declare method '%s'\n└── %s", method.Identifier.Name, err.Error())).Level(errgen.CRITICAL_ERROR)
		}

		//check the function body
		for _, stmt := range method.Body.Contents {
			CheckAST(stmt, fnEnv)
		}

	}

	return NewVoid()
}
