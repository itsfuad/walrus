package typechecker

import (
	"fmt"
	"walrus/ast"
	"walrus/errgen"
	//"walrus/errgen"
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
	initValue := CheckAST(forStmt.Init, forLoopEnv)
	if initValue != nil {
		//init value must be a statement
		fmt.Printf("init value %T\n", initValue)
		if _, ok := forStmt.Init.(ast.VarDeclStmt); !ok {
			//errgen.MakeError(env.filePath, forStmt.Init.StartPos().Line, forStmt.Init.EndPos().Line, forStmt.Init.StartPos().Column, forStmt.Init.EndPos().Column, "init value must be a statement").DisplayWithPanic()
			errgen.AddError(env.filePath, forStmt.Init.StartPos().Line, forStmt.Init.EndPos().Line, forStmt.Init.StartPos().Column, forStmt.Init.EndPos().Column, "init value must be a statement")
		}
	}

	conditionValue := CheckAST(forStmt.Condition, forLoopEnv)
	if conditionValue == nil {
		//errgen.MakeError(env.filePath, forStmt.Condition.StartPos().Line, forStmt.Condition.EndPos().Line, forStmt.Condition.StartPos().Column, forStmt.Condition.EndPos().Column, "condition value must be a present").DisplayWithPanic()
		errgen.AddError(env.filePath, forStmt.Condition.StartPos().Line, forStmt.Condition.EndPos().Line, forStmt.Condition.StartPos().Column, forStmt.Condition.EndPos().Column, "condition value must be a present")
	}

	incrementValue := CheckAST(forStmt.Increment, forLoopEnv)
	if incrementValue != nil {
		//increment value must be incremental statement; i++, i += 3
		fmt.Printf("increment value %T\n", incrementValue)
		if _, ok := forStmt.Increment.(ast.UnaryExpr); !ok {
			//errgen.MakeError(env.filePath, forStmt.Increment.StartPos().Line, forStmt.Increment.EndPos().Line, forStmt.Increment.StartPos().Column, forStmt.Increment.EndPos().Column, "increment value must be a statement").DisplayWithPanic()
			errgen.AddError(env.filePath, forStmt.Increment.StartPos().Line, forStmt.Increment.EndPos().Line, forStmt.Increment.StartPos().Column, forStmt.Increment.EndPos().Column, "increment value must be a statement")
		}
	}

	return NewVoid()
}
