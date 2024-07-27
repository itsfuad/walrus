package typechecker

import (
	"walrus/ast"
	"walrus/errgen"
)

func checkTypeDeclaration(node ast.TypeDeclStmt, env *TypeEnvironment) ValueTypeInterface {

	typeName := node.UDType

	var val ValueTypeInterface

	switch t := typeName.(type) {
	case ast.StructType:
		props := map[string]StructProperty{}
		for propname, propval := range t.Properties {
			propType, err := EvaluateTypeName(propval.PropType, env)
			if err != nil {
				errgen.MakeError(env.filePath, t.Start.Line, t.End.Line, t.Start.Column, t.End.Column, err.Error()).Display()
			}
			p := StructProperty{
				IsPrivate: propval.IsPrivate,
				Type:      propType,
			}
			props[propname] = p
		}
		val = Struct{
			DataType:   VALUE_TYPE(t.TypeName),
			StructName: node.UDTypeName,
			Elements:   props,
		}
	default:
		typ, err := EvaluateTypeName(node.UDType, env)
		if err != nil {
			errgen.MakeError(env.filePath, node.UDType.StartPos().Line, node.UDType.EndPos().Line, node.UDType.StartPos().Column, node.UDType.EndPos().Column, err.Error()).Display()
		}
		val = typ
	}

	typeVal := UserDefined{
		DataType: "user-defined",
		TypeDef:  val,
	}

	env.DeclareType(node.UDTypeName, typeVal)
	return nil
}