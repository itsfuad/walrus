package typechecker

import (
	"fmt"
	"walrus/ast"
	"walrus/errgen"
	"walrus/lexer"
)

func checkUnaryExpr(node ast.UnaryExpr, env *TypeEnvironment) ValueTypeInterface {
	op := node.Operator
	arg := node.Argument
	//evaluate argument. must be evaluated to number or boolean for ! (not)

	typeVal := CheckAST(arg, env)

	switch t := typeVal.(type) {
	case Int, Float:
		//allow - only
		if op.Kind != lexer.MINUS_TOKEN {
			errgen.MakeError(env.filePath, op.Start.Line, op.Start.Column, op.End.Column, "invalid unary operation with numeric types").Display()
		}
	case Bool:
		if op.Kind != lexer.NOT_TOKEN {
			errgen.MakeError(env.filePath, op.Start.Line, op.Start.Column, op.End.Column, "invalid unary operation with boolean types").Display()
		}
	default:
		errgen.MakeError(env.filePath, op.Start.Line, op.Start.Column, op.End.Column, fmt.Sprintf("this unary operation is not supported with %s types", t.DType())).Display()
	}

	return typeVal
}