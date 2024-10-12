package typechecker

import (
	"fmt"
	"walrus/ast"
	"walrus/errgen"
	"walrus/utils"
)

func EvaluateProgram(program ast.ProgramStmt, env *TypeEnvironment) ValueTypeInterface {
	utils.ColorPrint(utils.ORANGE, "### Evaluating program ###")
	for _, item := range program.Contents {
		CheckAST(item, env)
	}

	utils.ColorPrint(utils.GREEN, "--------- passed ---------")

	return Void{
		DataType: VOID_TYPE,
	}
}

func CheckAST(node ast.Node, env *TypeEnvironment) ValueTypeInterface {
	switch t := node.(type) {
	case ast.ProgramStmt:
		return EvaluateProgram(t, env)
	case ast.VarDeclStmt:
		return checkVariableDeclaration(t, env)
	case ast.VarAssignmentExpr:
		return checkVariableAssignment(t, env)
	case ast.IdentifierExpr:
		return checkIdentifier(t, env)
	case ast.IntegerLiteralExpr:
		return Int{
			DataType: INT_TYPE,
		}
	case ast.FloatLiteralExpr:
		return Float{
			DataType: FLOAT_TYPE,
		}
	case ast.StringLiteralExpr:
		return Str{
			DataType: STRING_TYPE,
		}
	case ast.CharLiteralExpr:
		return Chr{
			DataType: CHAR_TYPE,
		}
	case ast.BinaryExpr:
		return checkBinaryExpr(t, env)
	case ast.UnaryExpr:
		return checkUnaryExpr(t, env)
	case ast.IncrementalInterface:
		return checkIncrementalExpr(t, env)
	case ast.ArrayExpr:
		return evaluateArrayExpr(t, env)
	case ast.ArrayIndexAccess:
		return evaluateArrayAccess(t, env)
	case ast.TypeDeclStmt:
		return checkTypeDeclaration(t, env)
	case ast.StructLiteral:
		return checkStructLiteral(t, env)
	case ast.StructPropertyAccessExpr:
		return checkProperty(t, env)
	case ast.FunctionDeclStmt:
		return checkFunctionDeclStmt(t, env)
	case ast.FunctionExpr:
		return checkFunctionExpr(t, env)
	case ast.FunctionCallExpr:
		return checkFunctionCall(t, env)
	case ast.IfStmt:
		return checkIfStmt(t, env)
	case ast.ReturnStmt:
		return checkReturnStmt(t, env)
	}
	errgen.MakeError(env.filePath, node.StartPos().Line, node.EndPos().Line, node.StartPos().Column, node.EndPos().Column, fmt.Sprintf("<%T> node is not implemented yet to check", node)).Display()
	return nil
}

func checkFunctionExpr(funcNode ast.FunctionExpr, env *TypeEnvironment) ValueTypeInterface {
	name := fmt.Sprintf("_FN_%s", RandStringRunes(10))
	return analyzeFuntion(funcNode, name, env)
}

func analyzeFuntion(funcNode ast.FunctionExpr, name string, env *TypeEnvironment) Fn {

	fnEnv := NewTypeENV(env, FUNCTION_SCOPE, name, env.filePath)

	parameters := checkParamaters(funcNode.Params, fnEnv)

	//check return type
	returnType := EvaluateTypeName(funcNode.ReturnType, fnEnv)

	fn := Fn{
		DataType:      FUNCTION_TYPE,
		Params:        parameters,
		Returns:       returnType,
		FunctionScope: *fnEnv,
	}

	//declare the function
	env.DeclareVar(name, fn, true, false)

	//check the function body
	for _, stmt := range funcNode.Body.Contents {
		CheckAST(stmt, fnEnv)
	}

	return fn
}

func checkParamaters(params []ast.FunctionParam, fnEnv *TypeEnvironment) []FnParam {

	var parameters []FnParam

	for _, param := range params {

		if fnEnv.isDeclared(param.Identifier.Name) {
			errgen.MakeError(fnEnv.filePath, param.Identifier.Start.Line, param.Identifier.End.Line, param.Identifier.Start.Column, param.Identifier.End.Column, fmt.Sprintf("Parameter %s is already declared", param.Identifier.Name)).Display()
		}

		paramType := EvaluateTypeName(param.Type, fnEnv)

		if param.IsOptional {
			//default value type
			defaultValue := CheckAST(param.DefaultValue, fnEnv)
			MatchTypes(paramType, defaultValue, fnEnv.filePath, param.DefaultValue.StartPos().Line, param.DefaultValue.EndPos().Line, param.DefaultValue.StartPos().Column, param.DefaultValue.EndPos().Column)
		}

		fnEnv.DeclareVar(param.Identifier.Name, paramType, false, param.IsOptional)

		parameters = append(parameters, FnParam{
			Name:       param.Identifier.Name,
			IsOptional: param.IsOptional,
			Type:       paramType,
		})
	}
	return parameters
}

func checkFunctionCall(callNode ast.FunctionCallExpr, env *TypeEnvironment) ValueTypeInterface {
	//check if the function is declared
	if !env.isDeclared(callNode.Identifier.Name) {
		errgen.MakeError(env.filePath, callNode.Identifier.Start.Line, callNode.Identifier.End.Line, callNode.Identifier.Start.Column, callNode.Identifier.End.Column, fmt.Sprintf("Function %s is not declared", callNode.Identifier.Name)).Display()
	}

	//check if the function is a function
	fn, err := userDefinedToFn(env.variables[callNode.Identifier.Name])
	if err != nil {
		errgen.MakeError(env.filePath, callNode.Identifier.Start.Line, callNode.Identifier.End.Line, callNode.Identifier.Start.Column, callNode.Identifier.End.Column, fmt.Sprintf("'%s' is not a function", callNode.Identifier.Name)).Display()
	}

	//check if the number of arguments match the number of parameters
	fnParams := fn.Params
	if len(callNode.Arguments) != len(fnParams) {
		// exclude the optional parameters from the count
		optionalParams := 0
		for _, param := range fnParams {
			if param.IsOptional {
				optionalParams++
			}
		}
		if len(callNode.Arguments) < len(fnParams)-optionalParams {
			errgen.MakeError(env.filePath, callNode.Identifier.Start.Line, callNode.Identifier.End.Line, callNode.Identifier.Start.Column, callNode.Identifier.End.Column, fmt.Sprintf("Function %s expects at least %d arguments, got %d", callNode.Identifier.Name, len(fnParams)-optionalParams, len(callNode.Arguments))).Display()
		}
		if len(callNode.Arguments) > len(fnParams) {
			errgen.MakeError(env.filePath, callNode.Identifier.Start.Line, callNode.Identifier.End.Line, callNode.Identifier.Start.Column, callNode.Identifier.End.Column, fmt.Sprintf("Function %s expects at most %d arguments, got %d", callNode.Identifier.Name, len(fnParams), len(callNode.Arguments))).Display()
		}
	}

	//check if the arguments match the parameters
	for i := 0; i < len(callNode.Arguments); i++ {
		arg := CheckAST(callNode.Arguments[i], env)
		MatchTypes(fnParams[i].Type, arg, env.filePath, callNode.Arguments[i].StartPos().Line, callNode.Arguments[i].EndPos().Line, callNode.Arguments[i].StartPos().Column, callNode.Arguments[i].EndPos().Column)
	}

	return fn.Returns
}

func userDefinedToFn(ud ValueTypeInterface) (Fn, error) {
	// if UserDefined then chain until Fn or error
	switch t := ud.(type) {
	case Fn:
		return t, nil
	case UserDefined:
		return userDefinedToFn(t.TypeDef)
	default:
		return Fn{}, fmt.Errorf("'%s' is not a function", ud.DType())
	}
}

func checkFunctionDeclStmt(funcNode ast.FunctionDeclStmt, env *TypeEnvironment) ValueTypeInterface {

	// check if function is already declared
	funcName := funcNode.Identifier.Name

	if env.isDeclared(funcName) {
		errgen.MakeError(env.filePath, funcNode.Identifier.Start.Line, funcNode.Identifier.End.Line, funcNode.Identifier.Start.Column, funcNode.Identifier.End.Column, fmt.Sprintf("Function %s is already declared", funcName)).Display()
	}

	return analyzeFuntion(funcNode.FunctionExpr, funcName, env)
}

func checkIfStmt(ifNode ast.IfStmt, env *TypeEnvironment) ValueTypeInterface {
	//condition
	cond := CheckAST(ifNode.Condition, env)
	if cond.DType() != BOOLEAN_TYPE {
		errgen.MakeError(env.filePath, ifNode.Condition.StartPos().Line, ifNode.Condition.EndPos().Line, ifNode.Condition.StartPos().Column, ifNode.Condition.EndPos().Column, "Condition must be a boolean expression").Display()
	}

	//then block
	checkConditionBlock(ifNode.Block, env)

	if ifNode.AlternateBlock != nil {
		switch t := ifNode.AlternateBlock.(type) {
		case ast.IfStmt:
			checkIfStmt(t, env)
		case ast.BlockStmt:
			checkConditionBlock(t, env)
		}
	}

	return nil
}

func checkConditionBlock(block ast.BlockStmt, env *TypeEnvironment) ValueTypeInterface {
	for _, stmt := range block.Contents {
		CheckAST(stmt, env)
	}
	return nil
}

func checkReturnStmt(returnNode ast.ReturnStmt, env *TypeEnvironment) ValueTypeInterface {
	//check if the function is declared
	if env.scopeType != FUNCTION_SCOPE {
		errgen.MakeError(env.filePath, returnNode.StartPos().Line, returnNode.EndPos().Line, returnNode.StartPos().Column, returnNode.EndPos().Column, "Return statement must be inside a function").Display()
	}

	//check if the return type matches the function return type
	returnType := CheckAST(returnNode.Value, env)

	fn := getFunctionReturnValue(env, returnNode)

	MatchTypes(fn, returnType, env.filePath, returnNode.Start.Line, returnNode.End.Line, returnNode.Start.Column, returnNode.End.Column)

	return ReturnType{
		DataType:   RETURN_TYPE,
		Expression: returnType,
	}
}

func getFunctionReturnValue(env *TypeEnvironment, returnNode ast.Node) ValueTypeInterface {
	funcParent, err := env.ResolveFunctionEnv()
	if err != nil {
		errgen.MakeError(env.filePath, returnNode.StartPos().Line, returnNode.EndPos().Line, returnNode.StartPos().Column, returnNode.EndPos().Column, err.Error()).Display()
	}

	fnName := funcParent.scopeName
	fn := funcParent.parent.variables[fnName].(Fn)
	return fn.Returns
}
