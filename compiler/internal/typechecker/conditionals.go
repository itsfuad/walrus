package typechecker

import (
	//Walrus packages
	"walrus/compiler/internal/ast"
	"walrus/compiler/internal/report"
)

func checkConditionBlock(block ast.BlockStmt, env *TypeEnvironment) Tc {
	for _, stmt := range block.Contents {
		CheckAST(stmt, env)
	}
	return NewVoid()
}

func checkIfStmt(ifNode ast.IfStmt, env *TypeEnvironment) Tc {
	//condition
	cond := parseNodeValue(ifNode.Condition, env)
	if _, ok := cond.(Bool); !ok {
		report.Add(env.filePath, ifNode.Condition.StartPos().Line, ifNode.Condition.EndPos().Line, ifNode.Condition.StartPos().Column, ifNode.Condition.EndPos().Column, "Condition must be a boolean expression").SetLevel(report.NORMAL_ERROR)
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
