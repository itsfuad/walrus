package typechecker

import (
	"fmt"
	"walrus/ast"
	"walrus/errgen"
	"walrus/utils"
)

func EvaluateProgram(program ast.ProgramStmt, env *TypeEnvironment) TcValue {
	utils.PURPLE.Println("### Running type checker ###")
	for _, item := range program.Contents {
		CheckAST(item, env)
	}

	//print the file path
	utils.GREY.Printf("File path: %s\n", env.filePath)

	return NewVoid()
}

func CheckAST(node ast.Node, env *TypeEnvironment) TcValue {
	switch t := node.(type) {
	case ast.ProgramStmt:
		return EvaluateProgram(t, env)
	case ast.VarDeclStmt:
		return checkVariableDeclaration(t, env)
	case ast.VarAssignmentExpr:
		return checkVariableAssignment(t, env)
	case ast.TypeCastExpr:
		return checkTypeCast(t, env)
	case ast.IdentifierExpr:
		return checkIdentifier(t, env)
	case ast.IntegerLiteralExpr:
		return NewInt(t.BitSize, t.IsSigned)
	case ast.FloatLiteralExpr:
		return NewFloat(t.BitSize)
	case ast.StringLiteralExpr:
		return NewStr()
	case ast.ByteLiteralExpr:
		return NewInt(8, false)
	case ast.BinaryExpr:
		return checkBinaryExpr(t, env)
	case ast.UnaryExpr:
		return checkUnaryExpr(t, env)
	case ast.IncrementalInterface:
		return checkIncrementalExpr(t, env)
	case ast.ArrayLiteral:
		return evaluateArrayExpr(t, env)
	case ast.Indexable:
		return evaluateIndexableAccess(t, env)
	case ast.TypeDeclStmt:
		return checkTypeDeclaration(t, env)
	case ast.ImplStmt:
		return checkImplStmt(t, env)
	case ast.StructLiteral:
		return checkStructLiteral(t, env)
	case ast.StructPropertyAccessExpr:
		return checkPropertyAccess(t, env)
	case ast.MapLiteral:
		return checkMapLiteral(t, env)
	case ast.FunctionDeclStmt:
		return checkFunctionDeclStmt(t, env)
	case ast.FunctionLiteral:
		return checkFunctionExpr(t, env)
	case ast.FunctionCallExpr:
		return checkFunctionCall(t, env)
	case ast.IfStmt:
		return checkIfStmt(t, env)
	case ast.ForStmt:
		return checkForStmt(t, env)
	case ast.ReturnStmt:
		return checkReturnStmt(t, env)
	case nil:
		return NewNull()
	}
	errgen.AddError(env.filePath, node.StartPos().Line, node.EndPos().Line, node.StartPos().Column, node.EndPos().Column, fmt.Sprintf("<%T> node is not implemented yet to check", node)).DisplayWithPanic()
	return nil
}
