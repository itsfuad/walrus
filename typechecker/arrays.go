package typechecker

import (
	"fmt"
	"walrus/ast"
	"walrus/errgen"
)

// evaluateArrayAccess evaluates the access of an array element by its index.
//
// Parameters:
// - array: The AST node representing the array index access.
// - env: The type environment in which the array access is evaluated.
//
// Returns:
// - ValueTypeInterface: The type of the elements contained in the array.
//
// The function performs the following checks:
// 1. Ensures that the array expression evaluates to an array type. If not, it generates an error with a hint that the type must be an array.
// 2. Ensures that the index expression evaluates to an integer type. If not, it generates an error with a hint that the index must be a valid integer.
//
// If both checks pass, the function returns the type of the elements contained in the array.
func evaluateArrayAccess(array ast.ArrayIndexAccess, env *TypeEnvironment) ValueTypeInterface {
	//Array must be evaluated to an array value
	arrType := CheckAST(array.Arrayvalue, env)
	if _, ok := arrType.(Array); !ok {
		lineStart := array.Arrayvalue.StartPos().Line
		lineEnd := array.Arrayvalue.EndPos().Line
		start := array.Arrayvalue.StartPos().Column
		end := array.Arrayvalue.EndPos().Column
		errgen.MakeError(env.filePath, lineStart, lineEnd, start, end, fmt.Sprintf("cannot access index of type %s", arrType.DType())).AddHint("type must be an array", errgen.TEXT_HINT).Display()
	}
	//index must be evaluated to int
	indexType := CheckAST(array.Index, env)
	if _, ok := indexType.(Int); !ok {
		lineStart := array.Index.StartPos().Line
		lineEnd := array.Index.EndPos().Line
		start := array.Index.StartPos().Column
		end := array.Index.EndPos().Column
		errgen.MakeError(env.filePath, lineStart, lineEnd, start, end, fmt.Sprintf("cannot use index of type %s", indexType.DType())).AddHint("index must be valid integer", errgen.TEXT_HINT).Display()
	}
	return arrType.(Array).ArrayType
}

// evaluateArrayExpr evaluates an array expression within a given type environment.
// It checks that all elements in the array are of the same type and returns an Array type.
//
// Parameters:
// - array: The array expression to evaluate.
// - env: The type environment in which the array expression is evaluated.
//
// Returns:
// - ValueTypeInterface: The type of the array, which includes the data type and the type of the array elements.
func evaluateArrayExpr(array ast.ArrayLiteral, env *TypeEnvironment) ValueTypeInterface {
	var expectedType ValueTypeInterface
	for i, value := range array.Values {
		v := CheckAST(value, env)
		if i == 0 {
			expectedType = v
		}
		//check every type is same or not
		MatchTypes(expectedType, v, env.filePath, array.Start.Line, array.End.Line, array.Values[i].StartPos().Column, array.Values[i].EndPos().Column)
	}

	return Array{
		DataType:  ARRAY_TYPE,
		ArrayType: expectedType,
	}
}
