package typechecker

import (
	"fmt"
	"walrus/ast"
	"walrus/errgen"
)

func checkVariableAssignment(node ast.VarAssignmentExpr, env *TypeEnvironment) ValueTypeInterface {

	Assignee := node.Assignee
	valueToAssign := node.Value

	if err := IsAssignable(Assignee, env); err != nil {
		errgen.MakeError(env.filePath, Assignee.StartPos().Line, Assignee.EndPos().Line, Assignee.StartPos().Column, Assignee.EndPos().Column, err.Error()).Display()
	}

	expected := GetValueType(Assignee, env)
	provided := GetValueType(valueToAssign, env)

	err := MatchTypes(expected, provided, env.filePath, valueToAssign.StartPos().Line, valueToAssign.EndPos().Line, valueToAssign.StartPos().Column, valueToAssign.EndPos().Column)
	if err != nil {
		errgen.MakeError(env.filePath, valueToAssign.StartPos().Line, valueToAssign.EndPos().Line, valueToAssign.StartPos().Column, valueToAssign.EndPos().Column, err.Error()).Display()
	}

	return provided
}

func checkVariableDeclaration(node ast.VarDeclStmt, env *TypeEnvironment) ValueTypeInterface {

	varToDecl := node.Variable

	fmt.Printf("Declaring variable %s\n", varToDecl.Name)

	var expectedTypeInterface ValueTypeInterface

	if node.ExplicitType != nil {
		expectedTypeInterface = EvaluateTypeName(node.ExplicitType, env)
		fmt.Printf("Explicit type %T, %s\n", expectedTypeInterface, expectedTypeInterface.DType())
	} else {
		expectedTypeInterface = GetValueType(node.Value, env)
		//handleExplicitType(typestr, env)
		fmt.Printf("Auto detected type %T, %s\n", expectedTypeInterface, expectedTypeInterface.DType())
	}

	if node.Value != nil && node.ExplicitType != nil {
		//providedValue := CheckAST(node.Value, env)
		providedValue := GetValueType(node.Value, env)
		err := MatchTypes(expectedTypeInterface, providedValue, env.filePath, node.Value.StartPos().Line, node.Value.EndPos().Line, node.Value.StartPos().Column, node.Value.EndPos().Column)
		if err != nil {
			errgen.MakeError(env.filePath, node.Value.StartPos().Line, node.Value.EndPos().Line, node.Value.StartPos().Column, node.Value.EndPos().Column, err.Error()).Display()
		}
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
