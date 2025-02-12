package typechecker

import (
	//Walrus packages
	"walrus/internal/ast"
	"walrus/internal/report"
)

func checkIdentifier(node ast.IdentifierExpr, env *TypeEnvironment) Tc {

	name := node.Name

	//identifier cannot be types or builtins
	if isTypeDefined(name) && (name != "null" && name != "void") {
		report.Add(env.filePath, node.StartPos().Line, node.EndPos().Line, node.StartPos().Column, node.EndPos().Column, "cannot use type as value").Level(report.CRITICAL_ERROR)
	}

	//find the declaredEnv where the variable was declared
	declaredEnv, err := env.resolveVar(name)
	if err != nil {
		report.Add(env.filePath, node.StartPos().Line, node.EndPos().Line, node.StartPos().Column, node.EndPos().Column, err.Error()).Level(report.CRITICAL_ERROR)
	}

	// if we found value on that scope, return the value. Else make error (though there is no change to reach the error)
	variable := declaredEnv.variables[name]

	return unwrapType(variable)
}
