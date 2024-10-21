package typechecker

import (
	"fmt"
	"walrus/ast"
	//"walrus/errgen"
)

func checkForStmt(forStmt ast.ForStmt, env *TypeEnvironment) ValueTypeInterface {

	// for loop can be infinite loop or have a start, end and step

	forLoopEnv := NewTypeENV(env, LOOP_SCOPE, "for loop", env.filePath)

	if forStmt.Init != nil {
		val := CheckAST(forStmt.Init, forLoopEnv)
		fmt.Printf("init: %T\n", val)
	}

	return nil
}