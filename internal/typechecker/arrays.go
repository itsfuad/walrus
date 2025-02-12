package typechecker

import (
	//Standard packages
	"fmt"
	//Walrus packages
	"walrus/internal/ast"
	"walrus/internal/report"
)

// evaluateIndexableAccess evaluates the access of an array element by its index.
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
func evaluateIndexableAccess(indexable ast.Indexable, e *TypeEnvironment) Tc {

	container := parseNodeValue(indexable.Container, e)
	index := parseNodeValue(indexable.Index, e)

	var indexedValueType Tc

	switch t := container.(type) {
	case Array:
		if !isIntType(index) {
			report.Add(e.filePath, indexable.Start.Line, indexable.End.Line, indexable.Index.StartPos().Column, indexable.Index.EndPos().Column, fmt.Sprintf("cannot use type '%s' to index array\n", tcToString(index))+report.TreeFormatString("type must be a valid signed integer")).SetLevel(report.NORMAL_ERROR)
		}
		indexedValueType = t.ArrayType
	case Str:
		if !isIntType(index) {
			report.Add(e.filePath, indexable.Start.Line, indexable.End.Line, indexable.Index.StartPos().Column, indexable.Index.EndPos().Column, fmt.Sprintf("cannot use type '%s' to index string\n", tcToString(index))+report.TreeFormatString("type must be a valid signed integer")).SetLevel(report.NORMAL_ERROR)
		}
		return NewInt(8, false)
	case Map:
		//if key is interface then error
		if _, ok := unwrapType(t.KeyType).(Interface); ok {
			report.Add(e.filePath, indexable.Start.Line, indexable.End.Line, indexable.Index.StartPos().Column, indexable.Index.EndPos().Column, fmt.Sprintf("cannot access index of type %s", INTERFACE_TYPE)).SetLevel(report.NORMAL_ERROR)
		}
		indexedValueType = t.ValueType
	default:
		report.Add(e.filePath, indexable.Start.Line, indexable.End.Line, indexable.Container.StartPos().Column, indexable.Container.EndPos().Column, fmt.Sprintf("cannot access index of type %s", tcToString(container))).SetLevel(report.CRITICAL_ERROR)
	}

	return indexedValueType
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
func evaluateArrayExpr(array ast.ArrayLiteral, env *TypeEnvironment) Tc {
	var expectedType Tc
	for i, value := range array.Values {
		v := parseNodeValue(value, env)
		if i == 0 {
			expectedType = v
		}
		//check every type is same or not
		err := validateTypeCompatibility(expectedType, v)
		if err != nil {
			report.Add(env.filePath, array.Start.Line, array.End.Line, array.Values[i].StartPos().Column, array.Values[i].EndPos().Column, err.Error()).SetLevel(report.NORMAL_ERROR)
		}
	}

	return Array{
		DataType:  ARRAY_TYPE,
		ArrayType: expectedType,
	}
}
