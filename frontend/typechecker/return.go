package typechecker

import (
	//Standard packages
	"fmt"
	//Walrus packages
	"walrus/frontend/ast"
	"walrus/report"
)

func checkReturnStmt(returnNode ast.ReturnStmt, env *TypeEnvironment) Tc {
	//check if the function is declared
	if !env.isInFunctionScope() {
		report.Add(env.filePath, returnNode.StartPos().Line, returnNode.EndPos().Line, returnNode.StartPos().Column, returnNode.EndPos().Column, "return statement outside function").Level(report.NORMAL_ERROR)
	}

	//check if the return type matches the function return type
	returnType := parseNodeValue(returnNode.Value, env)

	fnReturns := getFunctionReturnValue(env, returnNode)

	err := validateTypeCompatibility(fnReturns, returnType)
	if err != nil {
		report.Add(env.filePath, returnNode.StartPos().Line, returnNode.EndPos().Line, returnNode.StartPos().Column, returnNode.EndPos().Column, fmt.Sprintf("cannot return '%s' from this scope\n"+report.TreeFormatString(fmt.Sprintf("function '%s' expects return type '%s'", env.scopeName, tcToString(fnReturns))), tcToString(returnType))).Level(report.NORMAL_ERROR)
	}

	return ReturnType{
		DataType:   RETURN_TYPE,
		Expression: returnType,
	}
}
