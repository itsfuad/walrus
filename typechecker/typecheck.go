package typechecker

import (
	"fmt"
	"walrus/ast"
	"walrus/errgen"
	"walrus/utils"
)

func EvaluateProgram(program ast.ProgramStmt, env *TypeEnvironment) ValueTypeInterface {
	utils.ColorPrint(utils.ORANGE, "### Evaluating program ###")
	for _, item := range program.Contents {
		CheckAST(item, env)
	}

	utils.ColorPrint(utils.GREEN, "--------- passed ---------")

	return Void{
		DataType: VOID_TYPE,
	}
}

func CheckAST(node ast.Node, env *TypeEnvironment) ValueTypeInterface {
	switch t := node.(type) {
	case ast.ProgramStmt:
		return EvaluateProgram(t, env)
	case ast.VarDeclStmt:
		return checkVariableDeclaration(t, env)
	case ast.VarAssignmentExpr:
		return checkVariableAssignment(t, env)
	case ast.IdentifierExpr:
		return checkIdentifier(t, env)
	case ast.IntegerLiteralExpr:
		return Int{
			DataType: INT_TYPE,
		}
	case ast.FloatLiteralExpr:
		return Float{
			DataType: FLOAT_TYPE,
		}
	case ast.StringLiteralExpr:
		return Str{
			DataType: STRING_TYPE,
		}
	case ast.CharLiteralExpr:
		return Chr{
			DataType: CHAR_TYPE,
		}
	case ast.BinaryExpr:
		return checkBinaryExpr(t, env)
	case ast.UnaryExpr:
		return checkUnaryExpr(t, env)
	case ast.IncrementalInterface:
		return checkIncrementalExpr(t, env)
	case ast.ArrayExpr:
		return evaluateArrayExpr(t, env)
	case ast.ArrayIndexAccess:
		return evaluateArrayAccess(t, env)
	case ast.TypeDeclStmt:
		return checkTypeDeclaration(t, env)
	case ast.StructLiteral:
		return checkStructLiteral(t, env)
	case ast.StructPropertyAccessExpr:
		return checkPropertyAccess(t, env)
	case ast.TraitDeclStmt:
		return checkTraitDeclaration(t, env)
	case ast.FunctionDeclStmt:
		return checkFunctionDeclStmt(t, env)
	case ast.FunctionExpr:
		return checkFunctionExpr(t, env)
	case ast.FunctionCallExpr:
		return checkFunctionCall(t, env)
	case ast.IfStmt:
		return checkIfStmt(t, env)
	case ast.ReturnStmt:
		return checkReturnStmt(t, env)
	}
	errgen.MakeError(env.filePath, node.StartPos().Line, node.EndPos().Line, node.StartPos().Column, node.EndPos().Column, fmt.Sprintf("<%T> node is not implemented yet to check", node)).Display()
	return nil
}
