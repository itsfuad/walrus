package typechecker

import (
	"walrus/ast"
)

func checkForStmt(forStmt ast.ForStmt, env *TypeEnvironment) ValueTypeInterface {

	// for loop can be infinite loop or have a start, end and step

	forLoopEnv := NewTypeENV(env, LOOP_SCOPE, "for loop", env.filePath)

	if forStmt.Init == nil {
		//infinte loop
		for _, stmt := range forStmt.Block.Contents {
			CheckAST(stmt, forLoopEnv)
		}
		return NewVoid()
	}

	//init: optional, condition: must be present, increment: optional
	CheckAST(forStmt.Init, forLoopEnv)

	CheckAST(forStmt.Condition, forLoopEnv)

	CheckAST(forStmt.Increment, forLoopEnv)

	return NewVoid()
}
