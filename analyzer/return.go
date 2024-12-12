package analyzer

import (
	"fmt"
	"walrus/ast"
	"walrus/errgen"
)

func checkReturnStmt(returnNode ast.ReturnStmt, env *TypeEnvironment) TcValue {
	//check if the function is declared
	if !env.isInFunctionScope() {
		errgen.AddError(env.filePath, returnNode.StartPos().Line, returnNode.EndPos().Line, returnNode.StartPos().Column, returnNode.EndPos().Column, "return statement outside function").ErrorLevel(errgen.NORMAL)
	}

	//check if the return type matches the function return type
	returnType := parseNodeValue(returnNode.Value, env)

	fnReturns := getFunctionReturnValue(env, returnNode)

	err := matchTypes(fnReturns, returnType)
	if err != nil {
		errgen.AddError(env.filePath, returnNode.StartPos().Line, returnNode.EndPos().Line, returnNode.StartPos().Column, returnNode.EndPos().Column, fmt.Sprintf("cannot return '%s' from this scope. function '%s' expects return type '%s'", tcValueToString(returnType), env.scopeName, tcValueToString(fnReturns))).ErrorLevel(errgen.NORMAL)
	}

	return ReturnType{
		DataType:   RETURN_TYPE,
		Expression: returnType,
	}
}
