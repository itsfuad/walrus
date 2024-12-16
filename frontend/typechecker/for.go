package typechecker

import (
	"walrus/errgen"
	"walrus/frontend/ast"
)

func checkForStmt(forStmt ast.ForStmt, env *TypeEnvironment) ExprType {

	// for loop can be infinite loop or have a start, end and step

	forLoopEnv := NewTypeENV(env, LOOP_SCOPE, "for loop", env.filePath)

	if forStmt.Init != nil || forStmt.Condition != nil || forStmt.Increment != nil {

		//must be a variable declaration, or an assignment
		switch t := forStmt.Init.(type) {
		case ast.VarDeclStmt:
			checkVariableDeclaration(t, forLoopEnv)
		case ast.VarAssignmentExpr:
			checkVariableAssignment(t, forLoopEnv)
		default:
			errgen.Add(env.filePath, forStmt.StartPos().Line, forStmt.EndPos().Line, forStmt.StartPos().Column, forStmt.EndPos().Column, "for loop initialization must be a variable declaration or assignment").Level(errgen.CRITICAL_ERROR)
		}

		cond := parseNodeValue(forStmt.Condition, forLoopEnv)

		//must be a boolean if !cond -> error, if !cond.Type == bool -> error
		if cond == nil || cond.DType() != BOOLEAN_TYPE {
			errgen.Add(env.filePath, forStmt.StartPos().Line, forStmt.EndPos().Line, forStmt.StartPos().Column, forStmt.EndPos().Column, "for loop condition must be a boolean expression").Level(errgen.CRITICAL_ERROR)
		}

		incr := parseNodeValue(forStmt.Increment, forLoopEnv)

		//must be assignment
		if _, ok := incr.(ast.IncrementalInterface); !ok {
			errgen.Add(env.filePath, forStmt.StartPos().Line, forStmt.EndPos().Line, forStmt.StartPos().Column, forStmt.EndPos().Column, "for loop increment must be incremental assignment").Level(errgen.CRITICAL_ERROR)
		}
	}

	//infinte loop
	for _, stmt := range forStmt.Block.Contents {
		CheckAST(stmt, forLoopEnv)
	}

	return NewVoid()
}
