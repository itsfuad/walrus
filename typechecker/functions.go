package typechecker

import (
	"walrus/ast"
	"walrus/errgen"
)

func checkFuncDecl(node ast.FunctionDeclStmt, env *TypeEnvironment) ValueTypeInterface {

	funcName := node.Name

	funcEnv := NewTypeENV(env, FUNCTION_SCOPE, env.filePath)

	//check if function already exists
	if _, err := env.ResolveVar(funcName.Name); err == nil {
		errgen.MakeError(env.filePath, funcName.StartPos().Line, funcName.EndPos().Line, funcName.StartPos().Column, funcName.EndPos().Column, "function already exists").Display()
	}

	//check if return type is defined
	var returnType ValueTypeInterface
	if node.ReturnType != nil {
		returnType = handleExplicitType(node.ReturnType, env)
	} else {
		returnType = Void{
			DataType: VOID_TYPE,
		}
	}

	// now check the parameter types
	paramTypes := map[string]ValueTypeInterface{}
	for _, param := range node.Params {
		paramType := handleExplicitType(param.Type, env) // used env instead of funcEnv because the function is not yet defined
		paramTypes[param.Name.Name] = paramType
		// add the parameter to the function environment
		err := funcEnv.DeclareVar(param.Name.Name, paramType, false)
		if err != nil {
			errgen.MakeError(env.filePath, param.Name.StartPos().Line, param.Name.EndPos().Line, param.Name.StartPos().Column, param.Name.EndPos().Column, "parameter already exists").Display()
		}
	}

	// check the body of the function

	var returnedExpr ValueTypeInterface
	var returnNode ast.Node = ast.NullLiteralExpr{}

	/* // This code is needed to check for debugging later
	for i, stmt := range node.Block.Contents {
		typ := CheckAST(stmt, funcEnv)
		if typ.DType() == RETURN_TYPE {
			returnedExpr = typ.(ReturnType).Expression
			returnNode = stmt
			//if has more statements,
			if i + 1 < len(node.Block.Contents) {
				s := node.Block.Contents[i + 1]
				errgen.MakeError(env.filePath, s.StartPos().Line, s.EndPos().Line, s.StartPos().Column, s.EndPos().Column, "remove this unreachable code after return").Display()
			}
			break
		}
	}
	*/

	returnedExpr = checkBlock(node.Block, funcEnv)

	returnErrLineStart := node.End.Line
	returnErrLineEnd := node.End.Line
	returnErrColStart := node.End.Column - 1
	returnErrColEnd  := node.End.Column - 1

	if returnedExpr == nil {
		if returnType.DType() != VOID_TYPE {
			errgen.MakeError(env.filePath, returnErrLineStart, returnErrLineEnd, returnErrColStart, returnErrColEnd, "no return expression in function").Display()
		}
		returnedExpr = Void{
			DataType: VOID_TYPE,
		}
	} else {
		returnErrLineStart = returnNode.StartPos().Line
		returnErrLineEnd = returnNode.EndPos().Line
		returnErrColStart = returnNode.StartPos().Column
		returnErrColEnd = returnNode.EndPos().Column
		if returnedExpr.DType() != returnType.DType() {
			errgen.MakeError(env.filePath, returnErrLineStart, returnErrLineEnd, returnErrColStart, returnErrColEnd, "return type did not match with function return type").Display()
		}
	}
	// add the function to the environment
	err := env.DeclareVar(funcName.Name, Fn{
		DataType: FUNCTION_TYPE,
		Params: paramTypes,
		Returns: returnType,
	}, true);

	if err != nil {
		errgen.MakeError(env.filePath, funcName.StartPos().Line, funcName.EndPos().Line, funcName.StartPos().Column, funcName.EndPos().Column, err.Error()).Display()
	}

	return Void{
		DataType: VOID_TYPE,
	}
}

func checkReturnStmt(node ast.ReturnStmt, env *TypeEnvironment) ValueTypeInterface {
	//if the env is function env only then return is allowed
	if env.scopeType == GLOBAL_SCOPE {
		errgen.MakeError(env.filePath, node.Start.Line, node.End.Line, node.Start.Column, node.End.Column, "cannot return from global scope").Display()
	}

	var rType ValueTypeInterface

	if node.Value == nil {
		rType = Void{
			DataType: VOID_TYPE,
		}
	} else {
		rType = CheckAST(node.Value, env)
	}

	return ReturnType{
		DataType: RETURN_TYPE,
		Expression: rType,
	}
}