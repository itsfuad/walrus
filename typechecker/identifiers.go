package typechecker

import (
	"walrus/ast"
	"walrus/errgen"
)

func checkIdentifier(node ast.IdentifierExpr, env *TypeEnvironment) TcValue {

	name := node.Name

	//identifier cannot be types or builtins
	if isTypeDefined(name) {
		errgen.AddError(env.filePath, node.StartPos().Line, node.EndPos().Line, node.StartPos().Column, node.EndPos().Column, "cannot use type as value", errgen.ERROR_CRITICAL)
	}

	//find the declaredEnv where the variable was declared
	declaredEnv, err := env.ResolveVar(name)
	if err != nil {
		errgen.AddError(env.filePath, node.StartPos().Line, node.EndPos().Line, node.StartPos().Column, node.EndPos().Column, err.Error(), errgen.ERROR_CRITICAL)

	}
	// if we found value on that scope, return the value. Else make error (though there is no change to reach the error)
	variable := declaredEnv.variables[name]

	val, err := unwrapType(variable)
	if err != nil {
		errgen.AddError(env.filePath, node.StartPos().Line, node.EndPos().Line, node.StartPos().Column, node.EndPos().Column, err.Error(), errgen.ERROR_CRITICAL)
	}

	return val
}
