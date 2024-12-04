package typechecker

import (
	"fmt"
	"walrus/ast"
	"walrus/builtins"
	"walrus/errgen"
	"walrus/lexer"
	"walrus/utils"
)

func checkIncrementalExpr(node ast.IncrementalInterface, env *TypeEnvironment) TcValue {
	op := node.Op()
	arg := node.Arg()
	// the argument must be an identifier evaluated to a number
	typeVal := CheckAST(arg, env)
	if !isNumberType(typeVal) {

		errgen.AddError(env.filePath, arg.StartPos().Line, arg.EndPos().Line, arg.StartPos().Column, arg.EndPos().Column, "invalid prefix operation with non-numeric type", errgen.ERROR_NORMAL)
	}
	if op.Kind != lexer.PLUS_PLUS_TOKEN && op.Kind != lexer.MINUS_MINUS_TOKEN {

		errgen.AddError(env.filePath, op.Start.Line, op.End.Line, op.Start.Column, op.End.Column, "invalid prefix operation", errgen.ERROR_NORMAL)
	}
	return typeVal
}

func logCastSuccess(originalType TcValue, toCast TcValue) {
	utils.ORANGE.Print("casted type ")
	utils.PURPLE.Print(string(originalType.DType()))
	fmt.Print(" to ")
	utils.PURPLE.Println(string(toCast.DType()))
}

func checkTypeCast(node ast.TypeCastExpr, env *TypeEnvironment) TcValue {

	originalType := CheckAST(node.Expression, env)
	toCast := evaluateTypeName(node.ToCast, env)

	if originalType.DType() == toCast.DType() {
		logCastSuccess(originalType, toCast)
		return originalType
	}

	if isNumberType(originalType) && isNumberType(toCast) {
		logCastSuccess(originalType, toCast)
		return toCast
	}

	errgen.AddError(env.filePath, node.Start.Line, node.End.Line, node.Start.Column, node.End.Column, fmt.Sprintf("cannot cast '%s' to '%s'", originalType.DType(), toCast.DType()), errgen.ERROR_NORMAL)

	return originalType
}

func checkUnaryExpr(node ast.UnaryExpr, env *TypeEnvironment) TcValue {
	op := node.Operator
	arg := node.Argument
	//evaluate argument. must be evaluated to number or boolean for ! (not)
	typeVal := CheckAST(arg, env)

	switch t := typeVal.(type) {
	case Int, Float:
		//allow - only
		if op.Kind != lexer.MINUS_TOKEN {

			errgen.AddError(env.filePath, op.Start.Line, op.End.Line, op.Start.Column, op.End.Column, "invalid unary operation with numeric types", errgen.ERROR_NORMAL)
		}
	case Bool:
		if op.Kind != lexer.NOT_TOKEN {

			errgen.AddError(env.filePath, op.Start.Line, op.End.Line, op.Start.Column, op.End.Column, "invalid unary operation with boolean types", errgen.ERROR_NORMAL)
		}
	default:

		errgen.AddError(env.filePath, op.Start.Line, op.End.Line, op.Start.Column, op.End.Column, fmt.Sprintf("this unary operation is not supported with %s types", t.DType()), errgen.ERROR_NORMAL)
	}

	return typeVal
}

func checkBinaryExpr(node ast.BinaryExpr, env *TypeEnvironment) TcValue {
	op := node.Operator

	left := CheckAST(node.Left, env)
	right := CheckAST(node.Right, env)

	leftType := left.DType()
	rightType := right.DType()

	var errLineStart, errLineEnd, errStart, errEnd int
	var errMsg string

	switch op.Kind {
	case lexer.PLUS_TOKEN:
		return checkAdditionAndConcat(node, left, right, env)
	case lexer.MINUS_TOKEN, lexer.MUL_TOKEN, lexer.DIV_TOKEN, lexer.MOD_TOKEN, lexer.EXP_TOKEN:
		//must have to be numeric type on both side
		if leftType != builtins.INT32 && leftType != builtins.FLOAT32 {
			errMsg = "cannot perform numeric operation. left hand side expression must be evaluated to a numeric type"
			errLineStart = node.Left.StartPos().Line
			errLineEnd = node.Left.EndPos().Line
			errStart = node.Left.StartPos().Column
			errEnd = node.Left.EndPos().Column
		} else if rightType != builtins.INT32 && rightType != builtins.FLOAT32 {
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

	errgen.AddError(env.filePath, errLineStart, errLineEnd, errStart, errEnd, errMsg, errgen.ERROR_NORMAL)
	return left
}

func checkComparison(node ast.BinaryExpr, left TcValue, right TcValue, env *TypeEnvironment) TcValue {

	leftType := left.DType()
	rightType := right.DType()

	op := node.Operator

	boolean := NewBool()

	if op.Kind == lexer.DOUBLE_EQUAL_TOKEN || op.Kind == lexer.NOT_EQUAL_TOKEN {
		// ( ==, != ) allow every type
		if isNumberType(left) && isNumberType(right) {
			return boolean
		} else if leftType == rightType {
			return boolean
		}
	} else {
		// ( >=, >, <=, < ) allow only numeric types
		if isNumberType(left) && isNumberType(right) {
			return boolean
		}
	}
	errMsg := fmt.Sprintf("invalid compare operation between '%s' and '%s'", leftType, rightType)

	errgen.AddError(env.filePath, node.Start.Line, node.End.Line, node.Start.Column, node.End.Column, errMsg, errgen.ERROR_NORMAL)
	return left
}

func checkAdditionAndConcat(node ast.BinaryExpr, left TcValue, right TcValue, env *TypeEnvironment) TcValue {

	leftType := left.DType()
	rightType := right.DType()

	var errLineStart, errLineEnd, errStart, errEnd int
	var errMsg string
	//only string concat, int and floats are allowed.
	if leftType == builtins.INT32 || leftType == builtins.FLOAT32 {
		//right has to be int or float
		if rightType != builtins.INT32 && rightType != builtins.FLOAT32 {
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

	errgen.AddError(env.filePath, errLineStart, errLineEnd, errStart, errEnd, errMsg, errgen.ERROR_NORMAL)
	return left
}
