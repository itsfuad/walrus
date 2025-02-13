package typechecker

import (
	"fmt"
	"walrus/ast"
	"walrus/builtins"
	"walrus/errgen"
	"walrus/lexer"
)

func checkIncrementalExpr(node ast.IncrementalInterface, env *TypeEnvironment) ValueTypeInterface {
	op := node.Op()
	arg := node.Arg()
	// the argument must be an identifier evaluated to a number
	typeVal := GetValueType(arg, env)
	if !IsNumberType(typeVal) {
		//errgen.MakeError(env.filePath, arg.StartPos().Line, arg.EndPos().Line, arg.StartPos().Column, arg.EndPos().Column, "invalid prefix operation with non-numeric type").DisplayWithPanic()
		errgen.AddError(env.filePath, arg.StartPos().Line, arg.EndPos().Line, arg.StartPos().Column, arg.EndPos().Column, "invalid prefix operation with non-numeric type")
	}
	if op.Kind != lexer.PLUS_PLUS_TOKEN && op.Kind != lexer.MINUS_MINUS_TOKEN {
		//errgen.MakeError(env.filePath, op.Start.Line, op.End.Line, op.Start.Column, op.End.Column, "invalid prefix operation").DisplayWithPanic()
		errgen.AddError(env.filePath, op.Start.Line, op.End.Line, op.Start.Column, op.End.Column, "invalid prefix operation")
	}
	return typeVal
}

func checkUnaryExpr(node ast.UnaryExpr, env *TypeEnvironment) ValueTypeInterface {
	op := node.Operator
	arg := node.Argument
	//evaluate argument. must be evaluated to number or boolean for ! (not)
	typeVal := GetValueType(arg, env)

	switch t := typeVal.(type) {
	case Int, Float:
		//allow - only
		if op.Kind != lexer.MINUS_TOKEN {
			//errgen.MakeError(env.filePath, op.Start.Line, op.End.Line, op.Start.Column, op.End.Column, "invalid unary operation with numeric types").DisplayWithPanic()
			errgen.AddError(env.filePath, op.Start.Line, op.End.Line, op.Start.Column, op.End.Column, "invalid unary operation with numeric types")
		}
	case Bool:
		if op.Kind != lexer.NOT_TOKEN {
			//errgen.MakeError(env.filePath, op.Start.Line, op.End.Line, op.Start.Column, op.End.Column, "invalid unary operation with boolean types").DisplayWithPanic()
			errgen.AddError(env.filePath, op.Start.Line, op.End.Line, op.Start.Column, op.End.Column, "invalid unary operation with boolean types")
		}
	default:
		//errgen.MakeError(env.filePath, op.Start.Line, op.End.Line, op.Start.Column, op.End.Column, fmt.Sprintf("this unary operation is not supported with %s types", t.DType())).DisplayWithPanic()
		errgen.AddError(env.filePath, op.Start.Line, op.End.Line, op.Start.Column, op.End.Column, fmt.Sprintf("this unary operation is not supported with %s types", t.DType()))
	}

	return typeVal
}

func checkBinaryExpr(node ast.BinaryExpr, env *TypeEnvironment) ValueTypeInterface {
	op := node.Operator

	left := GetValueType(node.Left, env)
	right := GetValueType(node.Right, env)

	leftType := left.DType()
	rightType := right.DType()

	var errLineStart, errLineEnd, errStart, errEnd int
	var errMsg string

	switch op.Kind {
	case lexer.PLUS_TOKEN:
		return checkAdditionAndConcat(node, left, right, env)
	case lexer.MINUS_TOKEN, lexer.MUL_TOKEN, lexer.DIV_TOKEN, lexer.MOD_TOKEN, lexer.EXP_TOKEN:
		//must have to be numeric type on both side
		if leftType != builtins.INT && leftType != builtins.FLOAT {
			errMsg = "cannot perform numeric operation. left hand side expression must be evaluated to a numeric type"
			errLineStart = node.Left.StartPos().Line
			errLineEnd = node.Left.EndPos().Line
			errStart = node.Left.StartPos().Column
			errEnd = node.Left.EndPos().Column
		} else if rightType != builtins.INT && rightType != builtins.FLOAT {
			errMsg = "cannot perform numeric operation. right hand side expression must be evaluated to a numeric type"
			errLineStart = node.Right.StartPos().Line
			errLineEnd = node.Right.EndPos().Line
			errStart = node.Right.StartPos().Column
			errEnd = node.Right.EndPos().Column
		} else {
			return left
		}
	case lexer.DOUBLE_EQUAL_TOKEN, lexer.NOT_EQUAL_TOKEN, lexer.LESS_EQUAL_TOKEN, lexer.LESS_TOKEN, lexer.GREATER_EQUAL_TOKEN, lexer.GREATER_TOKEN:
		return checkComparison(node, left, right, env)
	default:
		errMsg = "invalid operator"
		errLineStart = op.Start.Line
		errLineEnd = op.End.Line
		errStart = op.Start.Column
		errEnd = op.End.Column
	}

	//errgen.MakeError(env.filePath, errLineStart, errLineEnd, errStart, errEnd, errMsg).DisplayWithPanic()
	errgen.AddError(env.filePath, errLineStart, errLineEnd, errStart, errEnd, errMsg)
	return left
}

func checkComparison(node ast.BinaryExpr, left ValueTypeInterface, right ValueTypeInterface, env *TypeEnvironment) ValueTypeInterface {

	leftType := left.DType()
	rightType := right.DType()

	op := node.Operator

	boolean := NewBool()

	if op.Kind == lexer.DOUBLE_EQUAL_TOKEN || op.Kind == lexer.NOT_EQUAL_TOKEN {
		// ( ==, != ) allow every type
		if IsNumberType(left) && IsNumberType(right) {
			return boolean
		} else if leftType == rightType {
			return boolean
		}
	} else {
		// ( >=, >, <=, < ) allow only numeric types
		if IsNumberType(left) && IsNumberType(right) {
			return boolean
		}
	}
	errMsg := fmt.Sprintf("invalid compare operation between '%s' and '%s'", leftType, rightType)
	//errgen.MakeError(env.filePath, node.Start.Line, node.End.Line, node.Start.Column, node.End.Column, errMsg).DisplayWithPanic()
	errgen.AddError(env.filePath, node.Start.Line, node.End.Line, node.Start.Column, node.End.Column, errMsg)
	return left
}

func checkAdditionAndConcat(node ast.BinaryExpr, left ValueTypeInterface, right ValueTypeInterface, env *TypeEnvironment) ValueTypeInterface {

	leftType := left.DType()
	rightType := right.DType()

	var errLineStart, errLineEnd, errStart, errEnd int
	var errMsg string
	//only string concat, int and floats are allowed.
	if leftType == builtins.INT || leftType == builtins.FLOAT {
		//right has to be int or float
		if rightType != builtins.INT && rightType != builtins.FLOAT {
			errMsg = "cannot perform numeric operation. right hand side expression must be evaluated to a numeric type"
			errLineStart = node.Right.StartPos().Line
			errLineEnd = node.Right.EndPos().Line
			errStart = node.Right.StartPos().Column
			errEnd = node.Right.EndPos().Column
		} else {
			//return the type if left
			return left
		}
	} else if leftType == builtins.STRING {
		// we concat the type
		return left
	} else {
		errMsg = "invalid expression"
		errLineStart = node.Start.Line
		errLineEnd = node.End.Line
		errStart = node.StartPos().Column
		errEnd = node.EndPos().Column
	}

	//errgen.MakeError(env.filePath, errLineStart, errLineEnd, errStart, errEnd, errMsg).DisplayWithPanic()
	errgen.AddError(env.filePath, errLineStart, errLineEnd, errStart, errEnd, errMsg)
	return left
}
