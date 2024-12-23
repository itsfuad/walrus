package typechecker

import (
	//Walrus packages
	"walrus/frontend/ast"
	"walrus/report"
)

func checkBlock(block ast.BlockStmt, env *TypeEnvironment) Block {

	var blockInfo Block

	blockInfo.ProblemLocation = block.Location

	for i, stmt := range block.Contents {
		val := CheckAST(stmt, env)
		if _, ok := val.(ReturnType); ok {
			//if has any more statements after return
			if i < len(block.Contents)-1 {
				report.Add(env.filePath, stmt.StartPos().Line, stmt.EndPos().Line, stmt.StartPos().Column, stmt.EndPos().Column, "unreachable code").Level(report.NORMAL_ERROR)
			}
			blockInfo.IsSatisfied = true
			return blockInfo
		} else if v, ok := val.(Block); ok {
			blockInfo.IsSatisfied = blockInfo.IsSatisfied || v.IsSatisfied
			if !v.IsSatisfied {
				blockInfo.ProblemLocation = v.ProblemLocation
			}
		}
	}

	return blockInfo
}

func checkIfStmt(ifNode ast.IfStmt, env *TypeEnvironment) Block {
	//condition
	cond := parseNodeValue(ifNode.Condition, env)
	if _, ok := cond.(Bool); !ok {
		report.Add(env.filePath, ifNode.Condition.StartPos().Line, ifNode.Condition.EndPos().Line, ifNode.Condition.StartPos().Column, ifNode.Condition.EndPos().Column, "Condition must be a boolean expression").Level(report.NORMAL_ERROR)
	}

	var block Block
	//then block
	ifBranchValue := checkBlock(ifNode.Block, env)

	block.ProblemLocation = ifBranchValue.ProblemLocation

	if ifNode.AlternateBlock != nil {
		var altBranchValue Block
		switch t := ifNode.AlternateBlock.(type) {
		case ast.IfStmt:
			altBranchValue = checkIfStmt(t, env)
		case ast.BlockStmt:
			altBranchValue = checkBlock(t, env)
		}

		block.IsSatisfied = ifBranchValue.IsSatisfied && altBranchValue.IsSatisfied

		if !altBranchValue.IsSatisfied {
			block.ProblemLocation = altBranchValue.ProblemLocation
		}



	} else {
		block.IsSatisfied = ifBranchValue.IsSatisfied
	}

	return block
}

// func checkElseifStmt(elseifNode ast.IfStmt, env *TypeEnvironment) Block {
// 	cond := parseNodeValue(elseifNode.Condition, env)
// 	if _, ok := cond.(Bool); !ok {
// 		report.Add(env.filePath, elseifNode.Condition.StartPos().Line, elseifNode.Condition.EndPos().Line, elseifNode.Condition.StartPos().Column, elseifNode.Condition.EndPos().Column, "Condition must be a boolean expression").Level(report.NORMAL_ERROR)
// 	}

// 	var block Block

// 	ifBranchValue := checkBlock(elseifNode.Block, env)

// 	if elseifNode.AlternateBlock != nil {
// 		var altBranchValue Block
// 		switch t := elseifNode.AlternateBlock.(type) {
// 		case ast.IfStmt:
// 			altBranchValue = checkElseifStmt(t, env)
// 		case ast.BlockStmt:
// 			altBranchValue = checkBlock(t, env)
// 		}

// 		if ifBranchValue.IsSatisfied && altBranchValue.IsSatisfied {
// 			block.IsSatisfied = true
// 		} else {
// 			block.IsSatisfied = false
// 		}
// 	} else {
// 		block.IsSatisfied = ifBranchValue.IsSatisfied
// 	}

// 	return block
// }
