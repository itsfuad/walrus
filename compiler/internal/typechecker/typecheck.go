package typechecker

import (
	//Standard packages
	"fmt"
	//Walrus packages
	"walrus/compiler/colors"
	"walrus/compiler/internal/ast"
	"walrus/compiler/report"
)

func evaluateProgram(program ast.ProgramStmt, env *TypeEnvironment) Tc {
	colors.PURPLE.Println("### Running type checker ###")
	for _, item := range program.Contents {
		checkAST(item, env)
	}

	//print the file path
	colors.BLUE.Printf("Evaluated File: %s\n", env.filePath)

	return NewVoid()
}

func Analyze(tree ast.Node, filePath string) {

	colors.PURPLE.Println("### Running type checker ###")

	env := ProgramEnv(filePath)

	checkAST(tree, env)
}

func checkAST(node ast.Node, env *TypeEnvironment) Tc {
	switch t := node.(type) {
	case ast.ProgramStmt:
		return evaluateProgram(t, env)
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
	default:
		return parseNodeValue(node, env)
	}
}

func parseNodeValue(node ast.Node, env *TypeEnvironment) Tc {
	switch t := node.(type) {
	case ast.VarAssignmentExpr:
		return checkVariableAssignment(t, env) // value
	case ast.TypeofExpr:
		return checkTypeof(t, env)
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
	default:
		report.Add(env.filePath, node.StartPos().Line, node.EndPos().Line, node.StartPos().Column, node.EndPos().Column, fmt.Sprintf("<%T> node is not implemented yet to check", node)).SetLevel(report.CRITICAL_ERROR)
		return NewVoid()
	}
}
