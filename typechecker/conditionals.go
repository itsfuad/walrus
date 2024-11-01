package typechecker

import (
	"walrus/ast"
	"walrus/errgen"
)

func checkConditionBlock(block ast.BlockStmt, env *TypeEnvironment) ValueTypeInterface {
	for _, stmt := range block.Contents {
		CheckAST(stmt, env)
	}
	return NewVoid()
}

func checkIfStmt(ifNode ast.IfStmt, env *TypeEnvironment) ValueTypeInterface {
	//condition
	cond := GetValueType(ifNode.Condition, env)
	if cond.DType() != BOOLEAN_TYPE {
		//errgen.AddError(env.filePath, ifNode.Condition.StartPos().Line, ifNode.Condition.EndPos().Line, ifNode.Condition.StartPos().Column, ifNode.Condition.EndPos().Column, "Condition must be a boolean expression").DisplayWithPanic()
		errgen.AddError(env.filePath, ifNode.Condition.StartPos().Line, ifNode.Condition.EndPos().Line, ifNode.Condition.StartPos().Column, ifNode.Condition.EndPos().Column, "Condition must be a boolean expression")
	}

	//then block
	checkConditionBlock(ifNode.Block, env)

	if ifNode.AlternateBlock != nil {
		switch t := ifNode.AlternateBlock.(type) {
		case ast.IfStmt:
			checkIfStmt(t, env)
		case ast.BlockStmt:
			checkConditionBlock(t, env)
		}
	}

	return NewVoid()
}
