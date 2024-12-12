package analyzer

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
	case ast.TypeDeclStmt:
		return checkTypeDeclaration(t, env)
	case ast.ImplStmt:
		return checkImplStmt(t, env)
	case ast.FunctionDeclStmt:
		return checkFunctionDeclStmt(t, env)
	case ast.IfStmt:
		return checkIfStmt(t, env)
	case ast.ForStmt:
		return checkForStmt(t, env)
	case ast.SafeStmt:
		return checkSafeStmt(t, env)
	default:
		return parseNodeValue(node, env)
	}
}

func parseNodeValue(node ast.Node, env *TypeEnvironment) TcValue {
	switch t := node.(type) {
	case ast.VarAssignmentExpr:
		return checkVariableAssignment(t, env) // value
	case ast.TypeCastExpr:
		return checkTypeCast(t, env) // value
	case ast.IdentifierExpr:
		return checkIdentifier(t, env) // value
	case ast.IntegerLiteralExpr:
		return NewInt(t.BitSize, t.IsSigned) // value
	case ast.FloatLiteralExpr:
		return NewFloat(t.BitSize) // value
	case ast.StringLiteralExpr:
		return NewStr() // value
	case ast.ByteLiteralExpr:
		return NewInt(8, false) // value
	case ast.BinaryExpr:
		return checkBinaryExpr(t, env) // value
	case ast.UnaryExpr:
		return checkUnaryExpr(t, env) // value
	case ast.IncrementalInterface:
		return checkIncrementalExpr(t, env) // value
	case ast.ArrayLiteral:
		return evaluateArrayExpr(t, env) // value
	case ast.Indexable:
		return evaluateIndexableAccess(t, env) // value
	case ast.StructLiteral:
		return checkStructLiteral(t, env) // value
	case ast.StructPropertyAccessExpr:
		return checkPropertyAccess(t, env) // value
	case ast.MapLiteral:
		return checkMapLiteral(t, env) // value
	case ast.FunctionLiteral:
		return checkFunctionExpr(t, env) // value
	case ast.FunctionCallExpr:
		return checkFunctionCall(t, env) // value
	case ast.ReturnStmt:
		return checkReturnStmt(t, env) // value
	case nil:
		return NewNull() // value
	default:
		errgen.AddError(env.filePath, node.StartPos().Line, node.EndPos().Line, node.StartPos().Column, node.EndPos().Column, fmt.Sprintf("<%T> node is not implemented yet to check", node)).ErrorLevel(errgen.CRITICAL)
		return NewVoid()
	}
}
