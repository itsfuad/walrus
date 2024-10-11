package typechecker

import (
	"fmt"
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
		val = Struct{
			DataType:   STRUCT_TYPE, // old: VALUE_TYPE(t.TypeName)
			StructName: node.UDTypeName,
			Elements:   props,
		}
	case ast.FunctionType:
		var params []FnParam

		fmt.Printf("Declaring type of function with UDTypeName: %s\n", node.UDTypeName)

		funcEnv := NewTypeENV(env, FUNCTION_SCOPE, node.UDTypeName, env.filePath)

		for _, param := range t.Parameters {
			typ := EvaluateTypeName(param.Type, funcEnv)

			params = append(params, FnParam{
				Name: param.Identifier.Name,
				IsOptional: param.IsOptional,
				Type: typ,
			})
		}

		var ret ValueTypeInterface
		typ := EvaluateTypeName(t.ReturnType, funcEnv)

		ret = typ

		val = Fn{
			DataType:      FUNCTION_TYPE,
			Params:        params,
			Returns:       ret,
			FunctionScope: *funcEnv,
		}

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
