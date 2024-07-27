package typechecker

import (
	"fmt"
	"walrus/ast"
	"walrus/builtins"
	"walrus/errgen"
)

// generate interfaces from the type enum
func makeTypesInterface(typ VALUE_TYPE, env *TypeEnvironment) (ValueTypeInterface, error) {
	switch typ {
	case INT_TYPE:
		return Int{
			DataType: typ,
		}, nil
	case FLOAT_TYPE:
		return Float{
			DataType: typ,
		}, nil
	case CHAR_TYPE:
		return Chr{
			DataType: typ,
		}, nil
	case STRING_TYPE:
		return Str{
			DataType: typ,
		}, nil
	case BOOLEAN_TYPE:
		return Bool{
			DataType: typ,
		}, nil
	case NULL_TYPE:
		return Null{
			DataType: typ,
		}, nil
	case VOID_TYPE:
		return Void{
			DataType: typ,
		}, nil
	default:
		//search for the type
		e, err := env.ResolveType(string(typ))
		if err != nil {
			return nil, err
		}
		return UserDefined{
			DataType: typ,
			TypeDef:  e.types[string(typ)],
		}, nil
	}
}

func handleExplicitType(explicitType ast.DataType, env *TypeEnvironment) ValueTypeInterface {
	//Explicit type is defined
	var expectedTypeInterface ValueTypeInterface
	switch t := explicitType.(type) {
	case ast.ArrayType:
		val, err := EvaluateTypeName(t, env)
		if err != nil {
			errgen.MakeError(env.filePath, t.Start.Line, t.End.Line, t.Start.Column, t.End.Column, err.Error()).Display()
		}
		expectedTypeInterface = val
	default:
		val, err := makeTypesInterface(VALUE_TYPE(explicitType.Type()), env)
		if err != nil {
			errgen.MakeError(env.filePath, explicitType.StartPos().Line, explicitType.EndPos().Line, explicitType.StartPos().Column, explicitType.EndPos().Column, err.Error()).Display()
		}
		expectedTypeInterface = val
	}
	return expectedTypeInterface
}

func getTypename(typeName ValueTypeInterface) VALUE_TYPE {
	switch t := typeName.(type) {
	case Array:
		return "[]" + getTypename(t.ArrayType)
	case Struct:
		return VALUE_TYPE(t.StructName)
	case Fn:
		return VALUE_TYPE(t.DataType)
	case UserDefined:
		return getTypename(t.TypeDef)
	default:
		return t.DType()
	}
}

func MatchTypes(expected ValueTypeInterface, provided ValueTypeInterface, filePath string, lineStart, lineEnd, start, end int) {

	expectedType := getTypename(expected)
	gotType := getTypename(provided)

	if expectedType != gotType {
		errgen.MakeError(filePath, lineStart, lineEnd, start, end, fmt.Sprintf("typecheck:expected '%s', got '%s'", expectedType, gotType)).Display()
	}
}

func IsLValue(node ast.Node) bool {
	switch t := node.(type) {
	case ast.IdentifierExpr:
		return true
	case ast.ArrayIndexAccess:
		return IsLValue(t.Arrayvalue)
	case ast.StructPropertyAccessExpr:
		return IsLValue(t.Object)
	default:
		return false
	}
}

func IsNumberType(operand ValueTypeInterface) bool {
	switch operand.(type) {
	case Int, Float:
		return true
	default:
		return false
	}
}

func EvaluateTypeName(dtype ast.DataType, env *TypeEnvironment) (ValueTypeInterface, error) {
	switch t := dtype.(type) {
	case ast.ArrayType:
		val, err := EvaluateTypeName(t.ArrayType, env)
		if err != nil {
			return nil, err
		}
		arr := Array{
			DataType:  builtins.ARRAY,
			ArrayType: val,
		}
		return arr, nil
	default:
		return makeTypesInterface(VALUE_TYPE(t.Type()), env)
	}
}
