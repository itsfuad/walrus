package typechecker

import (
	"fmt"
	"walrus/ast"
	"walrus/errgen"
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
func evaluateIndexableAccess(indexable ast.Indexable, env *TypeEnvironment) ValueTypeInterface {
	
	container := nodeType(indexable.Container, env)
	index := nodeType(indexable.Index, env)

	var indexedValueType ValueTypeInterface

	switch t := container.(type) {
	case Array:
		if !isIntType(index) {
			//return t.ArrayType, fmt.Errorf("index must be a valid integer")
			errgen.AddError(env.filePath, indexable.Start.Line, indexable.End.Line, indexable.Index.StartPos().Column, indexable.Index.EndPos().Column, "index must be a valid integer")
		}
		indexedValueType = t.ArrayType
	case Str:
		if !isIntType(index) {
			//fmt.Errorf("index must be a valid integer")
			errgen.AddError(env.filePath, indexable.Start.Line, indexable.End.Line, indexable.Index.StartPos().Column, indexable.Index.EndPos().Column, "index must be a valid integer")
		}
		return NewInt(8, false)
	case Map:
		//if key is interface then error
		if t.KeyType.DType() == INTERFACE_TYPE {
			//return t.ValueType, fmt.Errorf("cannot access index of type %s", INTERFACE_TYPE)
			errgen.AddError(env.filePath, indexable.Start.Line, indexable.End.Line, indexable.Index.StartPos().Column, indexable.Index.EndPos().Column, fmt.Sprintf("cannot access index of type %s", INTERFACE_TYPE))
		}
		indexedValueType = t.ValueType
	default:
		//return nil, fmt.Errorf("cannot access index of type %s", container.DType())
		errgen.AddError(env.filePath, indexable.Start.Line, indexable.End.Line, indexable.Container.StartPos().Column, indexable.Container.EndPos().Column, fmt.Sprintf("cannot access index of type %s", container.DType())).DisplayWithPanic()
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
func evaluateArrayExpr(array ast.ArrayLiteral, env *TypeEnvironment) ValueTypeInterface {
	var expectedType ValueTypeInterface
	for i, value := range array.Values {
		v := nodeType(value, env)
		if i == 0 {
			expectedType = v
		}
		//check every type is same or not
		err := matchTypes(expectedType, v)
		if err != nil {

			errgen.AddError(env.filePath, array.Start.Line, array.End.Line, array.Values[i].StartPos().Column, array.Values[i].EndPos().Column, err.Error())
		}
	}

	return Array{
		DataType:  ARRAY_TYPE,
		ArrayType: expectedType,
	}
}
