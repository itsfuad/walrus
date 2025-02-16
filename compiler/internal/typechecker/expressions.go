package typechecker

import (
	//Standard packages
	"errors"
	"fmt"

	//Walrus packages
	"walrus/compiler/colors"
	"walrus/compiler/internal/ast"
	"walrus/compiler/internal/builtins"
	"walrus/compiler/internal/lexer"
	"walrus/compiler/report"
)

func checkIncrementalExpr(node ast.IncrementalInterface, env *TypeEnvironment) Tc {
	op := node.Op()
	arg := node.Arg()
	// the argument must be an identifier evaluated to a number
	typeVal := parseNodeValue(arg, env)
	if !isNumberType(typeVal) {

		report.Add(env.filePath, arg.StartPos().Line, arg.EndPos().Line, arg.StartPos().Column, arg.EndPos().Column, "invalid prefix operation with non-numeric type").SetLevel(report.NORMAL_ERROR)
	}
	if op.Kind != lexer.PLUS_PLUS_TOKEN && op.Kind != lexer.MINUS_MINUS_TOKEN {

		report.Add(env.filePath, op.Start.Line, op.End.Line, op.Start.Column, op.End.Column, "invalid prefix operation").SetLevel(report.NORMAL_ERROR)
	}
	return typeVal
}

func logCastSuccess(originalType Tc, toCast Tc) {
	colors.ORANGE.Print("casted type ")
	colors.PURPLE.Print(tcToString(originalType))
	fmt.Print(" to ")
	colors.PURPLE.Println(tcToString(toCast))
}

func checkTypeof(node ast.TypeofExpr, env *TypeEnvironment) Tc {
	parseNodeValue(node.Expression, env)
	return NewStr()
}

func checkTypeCast(node ast.TypeCastExpr, env *TypeEnvironment) Tc {

	originalType := parseNodeValue(node.Expression, env)
	toCast := evaluateTypeName(node.ToCast, env)

	if err := isCastable(originalType, toCast); err == nil {
		logCastSuccess(originalType, toCast)
		return toCast
	} else {
		report.Add(env.filePath, node.Start.Line, node.End.Line, node.Start.Column, node.End.Column, err.Error()).SetLevel(report.NORMAL_ERROR)
	}

	return originalType
}

func isCastable(src, dest Tc) error {
	srcStr := tcToString(src)
	destStr := tcToString(dest)

	var err error

	switch t := src.(type) {
	case Int, Float:
		colors.BLUE.Printf("checking cast from %s to %s\n", srcStr, destStr)
		if isNumberType(dest) {
			return nil
		}
	case Struct:
		err = isCastableStruct(t, dest)
		if err == nil {
			return nil
		}
	case Interface:
		err = checkMethodsImplementations(t, dest)
		if err == nil {
			return nil
		}
	default:
		if srcStr == destStr {
			return nil
		}
	}

	if err != nil {
		return fmt.Errorf("cannot cast '%s' to '%s'\n%s", srcStr, destStr, report.TreeFormatError(err))
	}

	return fmt.Errorf("cannot cast '%s' to '%s'", srcStr, destStr)
}

func isCastableStruct(src Struct, dest Tc) error {

	tName := tcToString(src)
	dName := tcToString(dest)

	fmt.Printf("checking cast from %s to %s\n", tName, dName)

	errMsg := fmt.Sprintf("cannot cast struct '%s' to '%s'", tName, dName)

	format := "%s\n%s"

	//dest must be a struct and have the same fields as src
	if destStruct, ok := dest.(Struct); ok {
		if err := checkMissingFields(src, destStruct); err != nil {
			return fmt.Errorf(format, errMsg, err.Error())
		}
		return nil
	} else if destInterface, ok := dest.(Interface); ok {
		//check if the struct implements the interface
		return checkMethodsImplementations(src, destInterface)
	} else {
		fmt.Printf("dest type: %T\n", dest)
		return errors.New(errMsg)
	}
}

func checkMissingFields(src Struct, dest Struct) error {
	//fmt.Printf("checking missing fields in %s\n", targetType)
	errs := make([]error, 0)
	for key, val := range dest.StructScope.variables {
		fmt.Printf("checking field %s\n", key)
		if dVal, ok := src.StructScope.variables[key]; !ok {
			errs = append(errs, fmt.Errorf("field '%s' is missing in struct '%s'", key, src.StructName))
		} else if err := validateTypeCompatibility(val, dVal); err != nil {
			fmt.Printf("uncompatible: %s\n", err.Error())
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return errors.New(report.TreeFormatError(errs...).Error())
	}

	return nil
}

func checkUnaryExpr(node ast.UnaryExpr, env *TypeEnvironment) Tc {
	op := node.Operator
	arg := node.Argument
	//evaluate argument. must be evaluated to number or boolean for ! (not)
	typeVal := parseNodeValue(arg, env)

	switch t := typeVal.(type) {
	case Int, Float:
		//allow - only
		if op.Kind != lexer.MINUS_TOKEN {

			report.Add(env.filePath, op.Start.Line, op.End.Line, op.Start.Column, op.End.Column, "invalid unary operation with numeric types").SetLevel(report.NORMAL_ERROR)
		}
	case Bool:
		if op.Kind != lexer.NOT_TOKEN {

			report.Add(env.filePath, op.Start.Line, op.End.Line, op.Start.Column, op.End.Column, "invalid unary operation with boolean types").SetLevel(report.NORMAL_ERROR)
		}
	default:
		report.Add(env.filePath, op.Start.Line, op.End.Line, op.Start.Column, op.End.Column, fmt.Sprintf("this unary operation is not supported with %s types", tcToString(t))).SetLevel(report.NORMAL_ERROR)
	}

	return typeVal
}

func checkBinaryExpr(node ast.BinaryExpr, env *TypeEnvironment) Tc {
	op := node.Binop

	left := parseNodeValue(node.Left, env)
	right := parseNodeValue(node.Right, env)

	leftType := tcToString(left)
	rightType := tcToString(right)

	var errLineStart, errLineEnd, errStart, errEnd int
	var errMsg string

	switch op.Kind {
	case lexer.PLUS_TOKEN:
		return checkAdditionAndConcat(node, left, right, env)
	case lexer.MINUS_TOKEN, lexer.MUL_TOKEN, lexer.DIV_TOKEN, lexer.MOD_TOKEN, lexer.EXP_TOKEN:
		//must have to be numeric type on both side
		if leftType != builtins.INT32 && leftType != builtins.FLOAT32 {
			errMsg = fmt.Sprintf("cannot perform numeric operation between type '%s' and '%s'. left hand side expression must be evaluated to a numeric type.", leftType, rightType)
			errLineStart = node.Left.StartPos().Line
			errLineEnd = node.Left.EndPos().Line
			errStart = node.Left.StartPos().Column
			errEnd = node.Left.EndPos().Column
		} else if rightType != builtins.INT32 && rightType != builtins.FLOAT32 {
			errMsg = fmt.Sprintf("cannot perform numeric operation between type '%s' and '%s'. right hand side expression must be evaluated to a numeric type.", leftType, rightType)
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

	report.Add(env.filePath, errLineStart, errLineEnd, errStart, errEnd, errMsg).SetLevel(report.NORMAL_ERROR)
	return left
}

func checkComparison(node ast.BinaryExpr, left Tc, right Tc, env *TypeEnvironment) Tc {

	leftType := tcToString(left)
	rightType := tcToString(right)

	op := node.Binop

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

	report.Add(env.filePath, node.Start.Line, node.End.Line, node.Start.Column, node.End.Column, errMsg).SetLevel(report.NORMAL_ERROR)
	return left
}

func checkAdditionAndConcat(node ast.BinaryExpr, left Tc, right Tc, env *TypeEnvironment) Tc {

	leftType := tcToString(left)
	rightType := tcToString(right)

	var errLineStart, errLineEnd, errStart, errEnd int
	var errMsg string
	//only string concat, int and floats are allowed.
	if leftType == builtins.INT32 || leftType == builtins.FLOAT32 {
		//right has to be int or float
		if rightType != builtins.INT32 && rightType != builtins.FLOAT32 {
			errMsg = fmt.Sprintf("cannot perform numeric operation between type '%s' and '%s'. right hand side expression must be evaluated to a numeric type.", leftType, rightType)
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

	report.Add(env.filePath, errLineStart, errLineEnd, errStart, errEnd, errMsg).SetLevel(report.NORMAL_ERROR)
	return left
}
