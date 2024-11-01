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

	var retval ValueTypeInterface
	//Array must be evaluated to an array value
	arrType := GetValueType(array.Array, env)
	arr, ok1 := arrType.(Array)
	if !ok1 {
		_, ok2 := arrType.(Str)
		if !ok2 {
			lineStart := array.Array.StartPos().Line
			lineEnd := array.Array.EndPos().Line
			start := array.Array.StartPos().Column
			end := array.Array.EndPos().Column
			errgen.AddError(env.filePath, lineStart, lineEnd, start, end, fmt.Sprintf("cannot access index of type %s", arrType.DType())).AddHint("type must be an array", errgen.TEXT_HINT).DisplayWithPanic()
			//errgen.AddError(env.filePath, lineStart, lineEnd, start, end, fmt.Sprintf("cannot access index of type %s", arrType.DType())).AddHint("type must be an array", errgen.TEXT_HINT)
		}
		retval = NewInt(8, false)
	} else {
		retval = arr.ArrayType
	}
	//index must be evaluated to int
	indexType := GetValueType(array.Index, env)
	if _, ok := indexType.(Int); !ok {
		lineStart := array.Index.StartPos().Line
		lineEnd := array.Index.EndPos().Line
		start := array.Index.StartPos().Column
		end := array.Index.EndPos().Column
		//errgen.AddError(env.filePath, lineStart, lineEnd, start, end, fmt.Sprintf("cannot use index of type %s", indexType.DType())).AddHint("index must be valid integer", errgen.TEXT_HINT).DisplayWithPanic()
		errgen.AddError(env.filePath, lineStart, lineEnd, start, end, fmt.Sprintf("cannot use index of type %s", indexType.DType())).AddHint("index must be valid integer", errgen.TEXT_HINT)
	}

	return retval
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
		v := GetValueType(value, env)
		if i == 0 {
			expectedType = v
		}
		//check every type is same or not
		err := MatchTypes(expectedType, v, env.filePath, array.Start.Line, array.End.Line, array.Values[i].StartPos().Column, array.Values[i].EndPos().Column)
		if err != nil {
			//errgen.AddError(env.filePath, array.Start.Line, array.End.Line, array.Values[i].StartPos().Column, array.Values[i].EndPos().Column, err.Error()).DisplayWithPanic()
			errgen.AddError(env.filePath, array.Start.Line, array.End.Line, array.Values[i].StartPos().Column, array.Values[i].EndPos().Column, err.Error())
		}
	}

	return Array{
		DataType:  ARRAY_TYPE,
		ArrayType: expectedType,
	}
}
