package typechecker

import (
	"fmt"
	"walrus/ast"
	"walrus/errgen"
)

func checkImplStmt(implStmt ast.ImplStmt, env *TypeEnvironment) ValueTypeInterface {
	// Resolve the type to implement
	env, err := env.ResolveType(implStmt.ImplFor.Name)
	if err != nil {
		errgen.MakeError(env.filePath, implStmt.Start.Line, implStmt.End.Line, implStmt.Start.Column, implStmt.End.Column, err.Error()).Display()
		return nil
	}

	// type must be a struct
	if env.types[implStmt.ImplFor.Name].(UserDefined).TypeDef.DType() != STRUCT_TYPE {
		errgen.MakeError(env.filePath, implStmt.Start.Line, implStmt.End.Line, implStmt.Start.Column, implStmt.End.Column, "can only implement for structs").Display()
		return nil
	}

	implForType := env.types[implStmt.ImplFor.Name].(UserDefined).TypeDef.(Struct)

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
	env.types[implStmt.ImplFor.Name] = UserDefined{
		DataType: USER_DEFINED_TYPE,
		TypeDef:  implForType,
	}

	return implForType
}