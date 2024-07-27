package typechecker

import (
	"fmt"
	"walrus/ast"
	"walrus/errgen"
)

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

	var returnType ValueTypeInterface

	for i, stmt := range node.Contents {
		typ := CheckAST(stmt, env)
		fmt.Println(typ)
		if typ.DType() == RETURN_TYPE {
			returnedExpr := typ.(ReturnType).Expression
			//if has more statements,
			if i + 1 < len(node.Contents) {
				s := node.Contents[i + 1]
				fmt.Println(s)
				errgen.MakeError(env.filePath, s.StartPos().Line, s.EndPos().Line, s.StartPos().Column, s.EndPos().Column, "remove this unreachable code after return").Display()
			}
			returnType = ReturnType{
				DataType: RETURN_TYPE,
				Expression: returnedExpr,
			}
			break
		}
	}

	if returnType == nil {
		returnType = Void{
			DataType: VOID_TYPE,
		}
	}

	return returnType
}