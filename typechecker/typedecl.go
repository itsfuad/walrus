package typechecker

import (
	"fmt"
	"walrus/ast"
	"walrus/builtins"
	"walrus/errgen"
)

func checkTypeDeclaration(node ast.TypeDeclStmt, env *TypeEnvironment) ValueTypeInterface {

	typeName := node.UDType

	fmt.Printf("declaring type %s\n", node.UDTypeName)

	var val ValueTypeInterface

	switch t := typeName.(type) {
	case ast.StructType:
		val = checkStructTypeDecl(node.UDTypeName, t, env)
	case ast.FunctionType:
		val = EvaluateTypeName(t, env)
	case ast.InterfaceType:
		val = checkInterfaceTypeDecl(node.UDTypeName, t, env)
	case ast.ArrayType:
		val = EvaluateTypeName(t.ArrayType, env)
		arr := Array{
			DataType:  builtins.ARRAY,
			ArrayType: val,
		}
		val = arr
	case nil:
		val = NewVoid()
	default:
		typ, err := stringToValueTypeInterface(VALUE_TYPE(t.Type()), env)
		if err != nil {
			errgen.MakeError(env.filePath, t.StartPos().Line, t.EndPos().Line, t.StartPos().Column, t.EndPos().Column, err.Error()).Display()
		}
		val = typ
	}

	typeVal := UserDefined{
		DataType: USER_DEFINED_TYPE,
		TypeName: node.UDTypeName,
		TypeDef:  val,
	}

	err := env.DeclareType(node.UDTypeName, typeVal)
	if err != nil {
		fmt.Printf("cannot declare type %s: %s\n", node.UDTypeName, err.Error())
		errgen.MakeError(env.filePath, node.Start.Line, node.End.Line, node.Start.Column, node.End.Column, err.Error()).Display()
	}

	return val
}