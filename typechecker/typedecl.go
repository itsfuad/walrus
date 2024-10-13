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
			propType := EvaluateTypeName(propval.PropType, env)
			p := StructProperty{
				IsPrivate: propval.IsPrivate,
				Type:      propType,
			}
			props[propname] = p
		}

		embeds := make([]Struct, 0)

		//get all embedded structs
		for _, embed := range t.Embeds {
			//check if defined
			scope, err := env.ResolveType(embed.Name)
			if err != nil {
				errgen.MakeError(env.filePath, embed.StartPos().Line, embed.EndPos().Line, embed.StartPos().Column, embed.EndPos().Column, err.Error()).Display()
			}

			//check if struct
			udType, _ := scope.types[embed.Name].(UserDefined)
			if udType.TypeDef.DType() != STRUCT_TYPE {
				errgen.MakeError(env.filePath, embed.StartPos().Line, embed.EndPos().Line, embed.StartPos().Column, embed.EndPos().Column, "only structs can be embedded").Display()
			}

			//get struct and add its props to the current struct
			structType := udType.TypeDef.(Struct)
			for propName, prop := range structType.Elements {
				props[propName] = prop
			}

			embeds = append(embeds, structType)
		}

		val = Struct{
			DataType:   STRUCT_TYPE, // old: VALUE_TYPE(t.TypeName)
			StructName: node.UDTypeName,
			Elements:   props,
			Embeds: embeds,
		}
	case ast.FunctionType:
		val = checkFunctionSignature(node.UDTypeName, t, env)
	default:
		typ := EvaluateTypeName(node.UDType, env)
		val = typ
	}

	typeVal := UserDefined{
		DataType: "user-defined",
		TypeDef:  val,
	}

	err := env.DeclareType(node.UDTypeName, typeVal)
	if err != nil {
		errgen.MakeError(env.filePath, node.Start.Line, node.End.Line, node.Start.Column, node.End.Column, err.Error()).Display()
	}

	return nil
}

func checkFunctionSignature(name string, method ast.FunctionType, env *TypeEnvironment) Fn {
	var params []FnParam

	funcEnv := NewTypeENV(env, FUNCTION_SCOPE, name, env.filePath)

	for _, param := range method.Parameters {
		typ := EvaluateTypeName(param.Type, funcEnv)

		params = append(params, FnParam{
			Name:       param.Identifier.Name,
			IsOptional: param.IsOptional,
			Type:       typ,
		})
	}

	ret := EvaluateTypeName(method.ReturnType, funcEnv)

	return Fn{
		DataType:      FUNCTION_TYPE,
		Params:        params,
		Returns:       ret,
		FunctionScope: *funcEnv,
	}
}
