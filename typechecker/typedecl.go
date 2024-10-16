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
			property := StructProperty{
				IsPrivate: propval.IsPrivate,
				Type:      propType,
			}
			props[propname] = property
		}

		val = Struct{
			DataType:   STRUCT_TYPE,
			StructName: node.UDTypeName,
			Elements:   props,
			Methods:   	map[string]StructMethod{},
		}
	case ast.FunctionType:
		val = checkFunctionSignature(node.UDTypeName, t, env)
	default:
		typ := EvaluateTypeName(node.UDType, env)
		val = typ
	}

	typeVal := UserDefined {
		DataType: USER_DEFINED_TYPE,
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

func checkImplStmt(implStmt ast.ImplStmt, env *TypeEnvironment) ValueTypeInterface {
	// Resolve the type to implement
	udType, _, err := env.ResolveType(implStmt.ImplFor.Name)
	if err != nil {
		errgen.MakeError(env.filePath, implStmt.Start.Line, implStmt.End.Line, implStmt.Start.Column, implStmt.End.Column, err.Error()).Display()
		return nil
	}

	implForType := udType.(Struct)

	fmt.Printf("Implementing for type %s\n", getTypename(implForType))

	//add the methods to the struct
	for name, method := range implStmt.Methods {

		fnEnv := NewTypeENV(env, FUNCTION_SCOPE, name, env.filePath)

		implForType.Methods[name] = StructMethod{
			IsPrivate: method.IsPrivate,
			Fn: Fn{
				DataType: FUNCTION_TYPE,
				Params:  checkParamaters(method.Params, fnEnv),
				Returns: EvaluateTypeName(method.ReturnType, fnEnv),
				FunctionScope: *fnEnv,
			},
		}
	}

	//update the struct in the environment
	env.types[implStmt.ImplFor.Name] = implForType

	return implForType
}