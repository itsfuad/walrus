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
	var parameters []FnParam
	name := fmt.Sprintf("_FN_%s", RandStringRunes(10))
	fnEnv := NewTypeENV(env, FUNCTION_SCOPE, name, env.filePath)

	for _, param := range funcNode.Params {

		if fnEnv.isDeclared(param.Identifier.Name) {
			errgen.MakeError(fnEnv.filePath, param.Identifier.Start.Line, param.Identifier.End.Line, param.Identifier.Start.Column, param.Identifier.End.Column, fmt.Sprintf("Parameter %s is already declared", param.Identifier.Name)).Display()
		}

		paramType, err := EvaluateTypeName(param.Type, fnEnv)
		if err != nil {
			errgen.MakeError(fnEnv.filePath, param.Start.Line, param.End.Line, param.Start.Column, param.End.Column, err.Error()).Display()
		}
		fnEnv.DeclareVar(param.Identifier.Name, paramType, false)
		parameters = append(parameters, FnParam{
			Name: param.Identifier.Name,
			Type: paramType,
		})
	}

	//check return type
	returnType, err := EvaluateTypeName(funcNode.ReturnType, fnEnv)
	if err != nil {
		errgen.MakeError(fnEnv.filePath, funcNode.ReturnType.StartPos().Line, funcNode.ReturnType.EndPos().Line, funcNode.ReturnType.StartPos().Column, funcNode.ReturnType.EndPos().Column, err.Error()).Display()
	}

	fn := Fn{
		DataType:      FUNCTION_TYPE,
		Params:        parameters,
		Returns:       returnType,
		FunctionScope: *fnEnv,
	}

	//declare the function
	env.DeclareVar(name, fn, true)

	//check the function body
	for _, stmt := range funcNode.Block.Contents {
		CheckAST(stmt, fnEnv)
	}

	return fn
}

func checkFunctionCall(callNode ast.FunctionCallExpr, env *TypeEnvironment) ValueTypeInterface {
	//check if the function is declared
	if !env.isDeclared(callNode.Name.Name) {
		errgen.MakeError(env.filePath, callNode.Name.Start.Line, callNode.Name.End.Line, callNode.Name.Start.Column, callNode.Name.End.Column, fmt.Sprintf("Function %s is not declared", callNode.Name.Name)).Display()
	}

	//check if the function is a function
	fn, err := userDefinedToFn(env.variables[callNode.Name.Name])
	if err != nil {
		errgen.MakeError(env.filePath, callNode.Name.Start.Line, callNode.Name.End.Line, callNode.Name.Start.Column, callNode.Name.End.Column, fmt.Sprintf("'%s' is not a function", callNode.Name.Name)).Display()
	}

	//check if the number of arguments match the number of parameters
	fnParams := fn.Params
	if len(callNode.Arguments) != len(fnParams) {
		fmt.Printf("Value: %v, Length: %d\n", callNode.Arguments, len(callNode.Arguments))
		fmt.Printf("Value: %v, Length: %d\n", fnParams, len(fnParams))
		errgen.MakeError(env.filePath, callNode.Name.Start.Line, callNode.Name.End.Line, callNode.Name.Start.Column, callNode.Name.End.Column, fmt.Sprintf("Function %s expects %d arguments, got %d", callNode.Name.Name, len(fnParams), len(callNode.Arguments))).Display()
	}

	//check if the arguments match the parameters
	for i := 0; i < len(callNode.Arguments); i++ {
		arg := CheckAST(callNode.Arguments[i], env)
		if !matchTypes(fnParams[i].Type, arg) {
			errgen.MakeError(env.filePath, callNode.Arguments[i].StartPos().Line, callNode.Arguments[i].EndPos().Line, callNode.Arguments[i].StartPos().Column, callNode.Arguments[i].EndPos().Column, fmt.Sprintf("Argument %s expects type %s, got %s", fnParams[i].Name, fnParams[i].Type.DType(), arg.DType())).Display()
		}
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

func matchTypes(expected, provided ValueTypeInterface) bool {
	if expected.DType() == provided.DType() {
		return true
	}

	// the typed can be user defined which wraps the actual type but may be the types are same
	if _, ok := expected.(UserDefined); ok {
		return matchTypes(expected.(UserDefined).TypeDef, provided)
	}

	if _, ok := provided.(UserDefined); ok {
		return matchTypes(expected, provided.(UserDefined).TypeDef)
	}

	return false
}

func checkFunctionDeclStmt(funcNode ast.FunctionDeclStmt, env *TypeEnvironment) ValueTypeInterface {

	// check if function is already declared
	funcName := funcNode.Identifier.Name

	if env.isDeclared(funcName) {
		errgen.MakeError(env.filePath, funcNode.Identifier.Start.Line, funcNode.Identifier.End.Line, funcNode.Identifier.Start.Column, funcNode.Identifier.End.Column, fmt.Sprintf("Function %s is already declared", funcName)).Display()
	}

	var parameterTypes []FnParam

	//create a new environment for the function
	fnEnv := NewTypeENV(env, FUNCTION_SCOPE, funcName, env.filePath)

	//check parameters
	for _, param := range funcNode.Params {
		//check if the parameter is already declared
		if fnEnv.isDeclared(param.Identifier.Name) {
			errgen.MakeError(fnEnv.filePath, param.Identifier.Start.Line, param.Identifier.End.Line, param.Identifier.Start.Column, param.Identifier.End.Column, fmt.Sprintf("Parameter %s is already declared", param.Identifier.Name)).Display()
		}
		paramType, err := EvaluateTypeName(param.Type, fnEnv)
		if err != nil {
			errgen.MakeError(fnEnv.filePath, param.Start.Line, param.End.Line, param.Start.Column, param.End.Column, err.Error()).Display()
		}
		fnEnv.DeclareVar(param.Identifier.Name, paramType, false)
		parameterTypes = append(parameterTypes, FnParam{
			Name: param.Identifier.Name,
			Type: paramType,
		})
	}

	//check return type
	returnType, err := EvaluateTypeName(funcNode.ReturnType, fnEnv)
	if err != nil {
		errgen.MakeError(fnEnv.filePath, funcNode.ReturnType.StartPos().Column, funcNode.ReturnType.EndPos().Line, funcNode.ReturnType.StartPos().Column, funcNode.ReturnType.EndPos().Column, err.Error()).Display()
	}

	//declare the function
	fn := Fn{
		DataType:      FUNCTION_TYPE,
		Params:        parameterTypes,
		Returns:       returnType,
		FunctionScope: *fnEnv,
	}

	env.DeclareVar(funcName, fn, true)

	//check the function body
	for _, stmt := range funcNode.Block.Contents {
		CheckAST(stmt, fnEnv)
	}

	return fn
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

	if !matchTypes(fn, returnType) {
		errgen.MakeError(env.filePath, returnNode.StartPos().Line, returnNode.EndPos().Line, returnNode.StartPos().Column, returnNode.EndPos().Column, fmt.Sprintf("Return type does not match function return type. Expected %s, got %s", fn.DType(), returnType.DType())).Display()
	}

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
