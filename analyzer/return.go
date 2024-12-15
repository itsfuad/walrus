package analyzer

import (
	"fmt"
	"walrus/ast"
	"walrus/errgen"
)

func checkReturnStmt(returnNode ast.ReturnStmt, env *TypeEnvironment) ExprType {
	//check if the function is declared
	if !env.isInFunctionScope() {
		errgen.Add(env.filePath, returnNode.StartPos().Line, returnNode.EndPos().Line, returnNode.StartPos().Column, returnNode.EndPos().Column, "return statement outside function").Level(errgen.NORMAL_ERROR)
	}

	//check if the return type matches the function return type
	returnType := parseNodeValue(returnNode.Value, env)

	fnReturns := getFunctionReturnValue(env, returnNode)

	err := matchTypes(fnReturns, returnType)
	if err != nil {
		errgen.Add(env.filePath, returnNode.StartPos().Line, returnNode.EndPos().Line, returnNode.StartPos().Column, returnNode.EndPos().Column, fmt.Sprintf("cannot return '%s' from this scope. function '%s' expects return type '%s'", tcValueToString(returnType), env.scopeName, tcValueToString(fnReturns))).Level(errgen.NORMAL_ERROR)
	}

	return ReturnType{
		DataType:   RETURN_TYPE,
		Expression: returnType,
	}
}
