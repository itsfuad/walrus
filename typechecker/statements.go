package typechecker

import "walrus/ast"

func checkIfStmt(node ast.IfStmt, env *TypeEnvironment) ValueTypeInterface {
	//check the if condition
	CheckAST(node.Condition, env)
	//look inside the block
	CheckAST(node.Block, env)
	if node.AlternateBlock != nil {
		switch n := node.AlternateBlock.(type) {
		case ast.IfStmt:
			return checkIfStmt(n, env)
		case ast.BlockStmt:
			return checkBlock(n, env)
		}
	}
	return Void{
		DataType: VOID_TYPE,
	}
}

func checkBlock(node ast.BlockStmt, env *TypeEnvironment) ValueTypeInterface {

	for _, stmt := range node.Contents {
		CheckAST(stmt, env)
	}

	return Void{
		DataType: VOID_TYPE,
	}
}