package typechecker

import (
	"fmt"
	"walrus/ast"
	"walrus/errgen"
)

func evaluateArrayAccess(array ast.ArrayIndexAccess, env *TypeEnvironment) ValueTypeInterface {
	//Array must be evaluated to an array value
	arrType := CheckAST(array.Arrayvalue, env)
	if _, ok := arrType.(Array); !ok {
		line := array.Arrayvalue.StartPos().Line
		start := array.Arrayvalue.StartPos().Column
		end := array.Arrayvalue.EndPos().Column
		errgen.MakeError(env.filePath, line, start, end, fmt.Sprintf("cannot access index of type %s", arrType.DType())).AddHint("type must be an array", errgen.TEXT_HINT).Display()
	}
	//index must be evaluated to int
	indexType := CheckAST(array.Index, env)
	if _, ok := indexType.(Int); !ok {
		line := array.Index.StartPos().Line
		start := array.Index.StartPos().Column
		end := array.Index.EndPos().Column
		errgen.MakeError(env.filePath, line, start, end, fmt.Sprintf("cannot use index of type %s", indexType.DType())).AddHint("index must be valid integer", errgen.TEXT_HINT).Display()
	}
	return arrType.(Array).ArrayType
}

func evaluateArrayExpr(array ast.ArrayExpr, env *TypeEnvironment) ValueTypeInterface {
	var expectedType ValueTypeInterface
	for i, value := range array.Values {
		v := CheckAST(value, env)
		if i == 0 {
			expectedType = v
		}
		//check every type is same or not
		MatchTypes(expectedType, v, env.filePath, array.Start.Line, array.Values[i].StartPos().Column, array.Values[i].EndPos().Column)
	}

	return Array{
		DataType:  ARRAY_TYPE,
		ArrayType: expectedType,
	}
}
