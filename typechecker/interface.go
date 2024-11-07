package typechecker

import (
	"walrus/ast"
)

func checkInterfaceTypeDecl(interfaceName string, interfaceNode ast.InterfaceType, env *TypeEnvironment) Interface {

	methods := make(map[string]Fn)

	for name, method := range interfaceNode.Methods {

		fnEnv := NewTypeENV(env, FUNCTION_SCOPE, name, env.filePath)

		params := make([]FnParam, 0)

		for _, param := range method.Parameters {
			param := FnParam{
				Name:       param.Identifier.Name,
				IsOptional: param.IsOptional,
				Type:       evaluateTypeName(param.Type, fnEnv),
			}
			params = append(params, param)
		}

		returns := evaluateTypeName(method.ReturnType, fnEnv)

		method := Fn{
			DataType:      FUNCTION_TYPE,
			Params:        params,
			Returns:       returns,
			FunctionScope: *fnEnv,
		}

		methods[name] = method
	}

	return Interface{
		DataType:      INTERFACE_TYPE,
		InterfaceName: interfaceName,
		Methods:       methods,
	}
}
