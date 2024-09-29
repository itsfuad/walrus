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
		return checkFunctionDeclStmt(t, env);
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

func checkFunctionCall(callNode ast.FunctionCallExpr, env *TypeEnvironment) ValueTypeInterface {
	//check if the function is declared
	if !env.isDeclared(callNode.Name.Name) {
		errgen.MakeError(env.filePath, callNode.Name.Start.Line, callNode.Name.End.Line, callNode.Name.Start.Column, callNode.Name.End.Column, fmt.Sprintf("Function %s is not declared", callNode.Name.Name)).Display()
	}

	//check if the function is a function
	fn := env.variables[callNode.Name.Name]
	if fn.DType() != FUNCTION_TYPE {
		errgen.MakeError(env.filePath, callNode.Name.Start.Line, callNode.Name.End.Line, callNode.Name.Start.Column, callNode.Name.End.Column, fmt.Sprintf("%s is not a function", callNode.Name.Name)).Display()
	}

	//check if the number of arguments match the number of parameters
	fnParams := fn.(Fn).Params
	if len(callNode.Arguments) != len(fnParams) {
		errgen.MakeError(env.filePath, callNode.Name.Start.Line, callNode.Name.End.Line, callNode.Name.Start.Column, callNode.Name.End.Column, fmt.Sprintf("Function %s expects %d arguments, got %d", callNode.Name.Name, len(fnParams), len(callNode.Arguments))).Display()
	}

	//check if the arguments match the parameters
	i := 0
	for _, param := range fnParams {
		arg := CheckAST(callNode.Arguments[i], env)
		if arg.DType() != param.DType() {
			errgen.MakeError(env.filePath, callNode.Arguments[i].StartPos().Line, callNode.Arguments[i].EndPos().Line, callNode.Arguments[i].StartPos().Column, callNode.Arguments[i].EndPos().Column, fmt.Sprintf("Argument %d expects type %s, got %s", i, param.DType(), arg.DType())).Display()
		}
		i++
	}

	return fn.(Fn).Returns
}

func checkFunctionDeclStmt(funcNode ast.FunctionDeclStmt, env *TypeEnvironment) ValueTypeInterface {

	// check if function is already declared
	funcName := funcNode.Name.Name

	if env.isDeclared(funcName) {
		errgen.MakeError(env.filePath, funcNode.Name.Start.Line, funcNode.Name.End.Line, funcNode.Name.Start.Column, funcNode.Name.End.Column, fmt.Sprintf("Function %s is already declared", funcName)).Display()
	}

	parameterTypes := make(map[string]ValueTypeInterface)

	//create a new environment for the function
	fnEnv := NewTypeENV(env, FUNCTION_SCOPE, funcName, env.filePath)

	//check parameters
	for _, param := range funcNode.Params {
		//check if the parameter is already declared
		if fnEnv.isDeclared(param.Name.Name) {
			errgen.MakeError(fnEnv.filePath, param.Name.Start.Line, param.Name.End.Line, param.Name.Start.Column, param.Name.End.Column, fmt.Sprintf("Parameter %s is already declared", param.Name.Name)).Display()
		}
		paramType, err := EvaluateTypeName(param.Type, fnEnv)
		if err != nil {
			errgen.MakeError(fnEnv.filePath, param.Start.Line, param.End.Line, param.Start.Column, param.End.Column, err.Error()).Display()
		}
		fnEnv.DeclareVar(param.Name.Name, paramType, false)
		parameterTypes[param.Name.Name] = paramType
	}

	//check return type
	returnType, err := EvaluateTypeName(funcNode.ReturnType, fnEnv)
	if err != nil {
		errgen.MakeError(fnEnv.filePath, funcNode.ReturnType.StartPos().Column, funcNode.ReturnType.EndPos().Line, funcNode.ReturnType.StartPos().Column, funcNode.ReturnType.EndPos().Column, err.Error()).Display()
	}

	
	//declare the function
	fn := Fn{
		DataType: 		FUNCTION_TYPE,
		Params:   		parameterTypes,
		Returns: 		returnType,
		FunctionScope: 	*fnEnv,
	}

	env.DeclareVar(funcName, fn, true)

	//check the function body
	for _, stmt := range funcNode.Block.Contents {
		CheckAST(stmt, fnEnv)
	}

	if len(errgen.GetGlobalErrors()) > 0 {
		errgen.DisplayErrors()
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

	if returnType.DType() != fn.DType() {
		errgen.MakeError(env.filePath, returnNode.StartPos().Line, returnNode.EndPos().Line, returnNode.StartPos().Column, returnNode.EndPos().Column, fmt.Sprintf("Return type does not match function return type. Expected %s, got %s", fn.DType(), returnType.DType())).Display()
	}

	return ReturnType{
		DataType: RETURN_TYPE,
		Expression: returnType,
	}
}

func getFunctionReturnValue(env *TypeEnvironment, returnNode ast.Node) ValueTypeInterface {
	funcParent, err := env.ResolveFunctionParent()
	if err != nil {
		errgen.MakeError(env.filePath, returnNode.StartPos().Line, returnNode.EndPos().Line, returnNode.StartPos().Column, returnNode.EndPos().Column, err.Error()).Display()
	}

	fnName := funcParent.scopeName
	fn := funcParent.parent.variables[fnName].(Fn)
	return fn.Returns
}