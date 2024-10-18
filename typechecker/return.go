package typechecker

import (
	"walrus/ast"
	"walrus/errgen"
)

func checkReturnStmt(returnNode ast.ReturnStmt, env *TypeEnvironment) ValueTypeInterface {
	//check if the function is declared
	if env.scopeType != FUNCTION_SCOPE {
		errgen.MakeError(env.filePath, returnNode.StartPos().Line, returnNode.EndPos().Line, returnNode.StartPos().Column, returnNode.EndPos().Column, "Return statement must be inside a function").Display()
	}

	//check if the return type matches the function return type
	returnType := GetValueType(returnNode.Value, env)

	fn := getFunctionReturnValue(env, returnNode)

	MatchTypes(fn, returnType, env.filePath, returnNode.Start.Line, returnNode.End.Line, returnNode.Start.Column, returnNode.End.Column)

	return ReturnType{
		DataType:   RETURN_TYPE,
		Expression: returnType,
	}
}