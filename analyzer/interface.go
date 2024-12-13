package analyzer

import (
	"fmt"
	"walrus/ast"
	"walrus/errgen"
)

func checkInterfaceTypeDecl(interfaceName string, interfaceNode ast.InterfaceType, env *TypeEnvironment) Interface {

	methods := make([]InterfaceMethod, 0)

	for _, method := range interfaceNode.Methods {

		fnEnv := NewTypeENV(env, FUNCTION_SCOPE, method.Identifier.Name, env.filePath)

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

		//check if the method already exists
		for _, m := range methods {
			if m.Name == method.Identifier.Name {
				errgen.AddError(env.filePath, method.Identifier.Start.Line, method.Identifier.End.Line, method.Identifier.Start.Column, method.Identifier.End.Column,
					fmt.Sprintf("method '%s' already exists in interface '%s'", method.Identifier.Name, interfaceName)).ErrorLevel(errgen.CRITICAL)
			}
		}

		methods = append(methods, InterfaceMethod{
			Name:  method.Identifier.Name,
			Method: Fn{
				DataType:      FUNCTION_TYPE,
				Params:        params,
				Returns:       returns,
				FunctionScope: *fnEnv,
			},
		})
	}

	return Interface{
		DataType:      INTERFACE_TYPE,
		InterfaceName: interfaceName,
		Methods:       methods,
	}
}
