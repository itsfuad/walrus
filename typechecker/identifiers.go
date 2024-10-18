package typechecker

import (
	"walrus/ast"
	"walrus/errgen"
)

func checkIdentifier(node ast.IdentifierExpr, env *TypeEnvironment) ValueTypeInterface {

	name := node.Name
	//find the scope where the variable was declared
	scope, err := env.ResolveVar(name)
	if err != nil {
		errgen.MakeError(env.filePath, node.StartPos().Line, node.EndPos().Line, node.StartPos().Column, node.EndPos().Column, err.Error()).Display()
	}
	// if we found value on that scope, return the value. Else make error (though there is no change to reach the error)
	if val, ok := scope.variables[name]; ok {
		val, err := getValueTypeInterface(val, env)
		if err != nil {
			errgen.MakeError(env.filePath, node.StartPos().Line, node.EndPos().Line, node.StartPos().Column, node.EndPos().Column, err.Error()).Display()
		}
		return val
	}
	errgen.MakeError(env.filePath, node.StartPos().Line, node.EndPos().Line, node.StartPos().Column, node.EndPos().Column, "failed to check type. not found in environment").Display()
	return nil
}