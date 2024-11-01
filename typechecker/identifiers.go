package typechecker

import (
	"walrus/ast"
	"walrus/errgen"
)

func checkIdentifier(node ast.IdentifierExpr, env *TypeEnvironment) ValueTypeInterface {

	name := node.Name

	//identifier cannot be types or builtins
	if _, ok := env.types[name]; ok {
		errgen.AddError(env.filePath, node.StartPos().Line, node.EndPos().Line, node.StartPos().Column, node.EndPos().Column, "cannot use type as value").DisplayWithPanic()
	}

	if _, ok := env.builtins[name]; ok {
		errgen.AddError(env.filePath, node.StartPos().Line, node.EndPos().Line, node.StartPos().Column, node.EndPos().Column, "cannot use builtin as value").DisplayWithPanic()
	}

	//find the declaredEnv where the variable was declared
	declaredEnv, err := env.ResolveVar(name)
	if err != nil {
		errgen.AddError(env.filePath, node.StartPos().Line, node.EndPos().Line, node.StartPos().Column, node.EndPos().Column, err.Error()).DisplayWithPanic()
		//errgen.AddError(env.filePath, node.StartPos().Line, node.EndPos().Line, node.StartPos().Column, node.EndPos().Column, err.Error())
	}
	// if we found value on that scope, return the value. Else make error (though there is no change to reach the error)
	variable := declaredEnv.variables[name]

	value, err := getValueTypeInterface(variable, env)
	if err != nil {
		//errgen.AddError(env.filePath, node.StartPos().Line, node.EndPos().Line, node.StartPos().Column, node.EndPos().Column, err.Error()).DisplayWithPanic()
		errgen.AddError(env.filePath, node.StartPos().Line, node.EndPos().Line, node.StartPos().Column, node.EndPos().Column, err.Error())
	}
	return value
}
