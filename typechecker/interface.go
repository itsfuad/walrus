package typechecker

import (
	"walrus/ast"
	"walrus/errgen"
)

func checkInterfaceDeclaration(interfaceNode ast.InterfaceDeclStmt, env *TypeEnvironment) ValueTypeInterface {

	interfaceName := interfaceNode.Interface.Name

	methods := make(map[string]Fn)

	for name, method := range interfaceNode.Methods {

		fnEnv := NewTypeENV(env, FUNCTION_SCOPE, name, env.filePath)

		params := make([]FnParam, 0)

		for _, param := range method.Parameters {
			param := FnParam{
				Name: param.Identifier.Name,
				IsOptional: param.IsOptional,
				Type: EvaluateTypeName(param.Type, fnEnv),
			}
			params = append(params, param)
		}

		returns := EvaluateTypeName(method.ReturnType, fnEnv)

		method := Fn{
			DataType: FUNCTION_TYPE,
			Params:   params,
			Returns:  returns,
			FunctionScope: *fnEnv,
		}

		methods[name] = method
	}

	interfaceValue := Interface{
		DataType:      INTERFACE_TYPE,
		InterfaceName: interfaceName,
		Methods:       methods,
	}

	err := env.DeclareInterface(interfaceName, interfaceValue)
	if err != nil {
		errgen.MakeError(env.filePath, interfaceNode.Start.Line, interfaceNode.End.Line, interfaceNode.Start.Column, interfaceNode.End.Column, err.Error()).Display()
	}

	return interfaceValue
}