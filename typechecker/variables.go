package typechecker

import (
	"fmt"
	"walrus/ast"
	"walrus/errgen"
)

func checkVariableAssignment(node ast.VarAssignmentExpr, env *TypeEnvironment) ValueTypeInterface {

	Assignee := node.Assignee
	valueToAssign := node.Value

	expected := GetValueType(Assignee, env)
	provided := GetValueType(valueToAssign, env)

	MatchTypes(expected, provided, env.filePath, valueToAssign.StartPos().Line, valueToAssign.EndPos().Line, valueToAssign.StartPos().Column, valueToAssign.EndPos().Column)

	var varName string

	if !IsLValue(Assignee) {
		errgen.MakeError(env.filePath, Assignee.StartPos().Line, Assignee.EndPos().Line, Assignee.StartPos().Column, Assignee.EndPos().Column, "invalid assignment expression. the assignee must be a lvalue").AddHint("lvalue is something that has a memory address\nFor example: variables.\nso you cannot assign values something which does not exist in memory as an independent identifier.", errgen.TEXT_HINT).Display()
	}

	switch assignee := Assignee.(type) {
	case ast.IdentifierExpr:
		varName = assignee.Name
	case ast.ArrayIndexAccess:
		return nil
	case ast.StructPropertyAccessExpr:
		return nil
	default:
		panic("cannot assign to this type")
	}

	//get the stored type
	scope, err := env.ResolveVar(varName)
	if err != nil {
		errgen.MakeError(env.filePath, valueToAssign.StartPos().Line, valueToAssign.EndPos().Line, valueToAssign.StartPos().Column, valueToAssign.EndPos().Column, err.Error()).Display()
	}

	//if constant
	if scope.constants[varName] {
		errgen.MakeError(env.filePath, valueToAssign.StartPos().Line, valueToAssign.EndPos().Line, valueToAssign.StartPos().Column, valueToAssign.EndPos().Column, fmt.Sprintf("'%s' is constant", varName)).AddHint("cannot assign value to constant variables", errgen.TEXT_HINT).Display()
	}
	scope.variables[varName] = provided

	fmt.Printf("Assigned variable %s of type %T\n", varName, provided)

	return provided
}

func checkVariableDeclaration(node ast.VarDeclStmt, env *TypeEnvironment) ValueTypeInterface {

	varToDecl := node.Variable

	fmt.Printf("Declaring variable %s\n", varToDecl.Name)

	var expectedTypeInterface ValueTypeInterface

	if node.ExplicitType != nil {
		expectedTypeInterface = EvaluateTypeName(node.ExplicitType, env)
	} else {
		expectedTypeInterface = GetValueType(node.Value, env)
		//handleExplicitType(typestr, env)
		fmt.Printf("Auto detected type %T, %s\n", expectedTypeInterface, expectedTypeInterface.DType())
	}

	if node.Value != nil && node.ExplicitType != nil {
		//providedValue := CheckAST(node.Value, env)
		providedValue := GetValueType(node.Value, env)
		MatchTypes(expectedTypeInterface, providedValue, env.filePath, node.Value.StartPos().Line, node.Value.EndPos().Line, node.Value.StartPos().Column, node.Value.EndPos().Column)
	}

	err := env.DeclareVar(varToDecl.Name, expectedTypeInterface, node.IsConst, false)
	if err != nil {
		errgen.MakeError(env.filePath, node.Variable.StartPos().Line, node.Variable.EndPos().Line, node.Variable.StartPos().Column, node.Variable.EndPos().Column, err.Error()).Display()
	}

	if node.IsConst {
		fmt.Printf("Declared constant variable %s of type %s\n", varToDecl.Name, expectedTypeInterface.DType())
	} else {
		fmt.Printf("Declared variable %s of type %s\n", varToDecl.Name, expectedTypeInterface.DType())
	}

	return Void{
		DataType: VOID_TYPE,
	}
}
