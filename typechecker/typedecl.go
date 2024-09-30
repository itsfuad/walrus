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
	case ast.FunctionType:
		params := map[string]ValueTypeInterface{}

		fmt.Printf("Declaring type of function with UDTypeName: %s\n", node.UDTypeName)

		funcEnv := NewTypeENV(env, FUNCTION_SCOPE, node.UDTypeName, env.filePath)

		for name, paramType := range t.Parameters {
			typ, err := EvaluateTypeName(paramType, funcEnv)
			if err != nil {
				errgen.MakeError(funcEnv.filePath, t.Start.Line, t.End.Line, t.Start.Column, t.End.Column, err.Error()).Display()
			}
			params[name] = typ
		}

		var ret ValueTypeInterface
		typ, err := EvaluateTypeName(t.ReturnType, funcEnv)
		if err != nil {
			errgen.MakeError(funcEnv.filePath, t.Start.Line, t.End.Line, t.Start.Column, t.End.Column, err.Error()).Display()
		}

		ret = typ

		val = Fn{
			DataType: 	FUNCTION_TYPE,
			Params:  	params,
			Returns: 	ret,
			FunctionScope: *funcEnv,
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

	err := env.DeclareType(node.UDTypeName, typeVal)
	if err != nil {
		errgen.MakeError(env.filePath, node.Start.Line, node.End.Line, node.Start.Column, node.End.Column, err.Error()).Display()
	}

	return nil
}