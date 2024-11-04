package typechecker

import (
	"fmt"
	"walrus/ast"
	"walrus/errgen"
	"walrus/utils"
)

// checkVariableAssignment checks the assignment of a value to a variable in the given type environment.
// It verifies if the assignee is assignable and if the types of the assignee and the value to be assigned match.
// If any errors are encountered during these checks, they are displayed using the error generation utility.
//
// Parameters:
// - node: an AST node representing the variable assignment expression.
// - env: a pointer to the TypeEnvironment which holds type information and the file path.
//
// Returns:
// - ValueTypeInterface: the type of the value being assigned.
func checkVariableAssignment(node ast.VarAssignmentExpr, env *TypeEnvironment) ValueTypeInterface {

	Assignee := node.Assignee
	valueToAssign := node.Value

	if err := CheckLValue(Assignee, env); err != nil {
		errgen.AddError(env.filePath, Assignee.StartPos().Line, Assignee.EndPos().Line, Assignee.StartPos().Column, Assignee.EndPos().Column, "cannot assign to "+err.Error()).DisplayWithPanic()

	}

	expectedType := GetValueType(Assignee, env)
	providedType := GetValueType(valueToAssign, env)

	err := MatchTypes(expectedType, providedType)
	if err != nil {

		errgen.AddError(env.filePath, valueToAssign.StartPos().Line, valueToAssign.EndPos().Line, valueToAssign.StartPos().Column, valueToAssign.EndPos().Column, err.Error())
	}

	return providedType
}

// checkVariableDeclaration checks the declaration of a variable in the given AST node
// and updates the type environment accordingly.
//
// Parameters:
// - node: The AST node representing the variable declaration statement.
// - env: The type environment in which the variable is being declared.
//
// Returns:
// - ValueTypeInterface: The type of the declared variable.
//
// The function performs the following steps:
//  1. Retrieves the variable to be declared from the AST node.
//  2. Prints a message indicating the variable being declared.
//  3. Determines the expected type of the variable, either from an explicit type
//     specified in the declaration or by inferring it from the assigned value.
//  4. If both an explicit type and a value are provided, checks that the value's
//     type matches the expected type and reports an error if they do not match.
//  5. Declares the variable in the type environment and reports any errors that occur.
//  6. Prints a message indicating whether the variable is a constant and its type.
//  7. Returns a Void type indicating the end of the declaration process.
func checkVariableDeclaration(node ast.VarDeclStmt, env *TypeEnvironment) ValueTypeInterface {

	varToDecl := node.Variable

	fmt.Print("Declaring variable ")
	utils.ColorPrint(utils.RED, varToDecl.Name + "\n")

	var expectedTypeInterface ValueTypeInterface

	if node.ExplicitType != nil {
		expectedTypeInterface = EvaluateTypeName(node.ExplicitType, env)
		fmt.Print("Explicit type: ")
		utils.ColorPrint(utils.PURPLE, string(valueTypeInterfaceToString(expectedTypeInterface)) + "\n")
	} else {
		expectedTypeInterface = GetValueType(node.Value, env)
		fmt.Print("Auto detected type: ")
		utils.ColorPrint(utils.PURPLE, string(valueTypeInterfaceToString(expectedTypeInterface)) + "\n")
	}

	if node.Value != nil && node.ExplicitType != nil {
		//providedValue := CheckAST(node.Value, env)
		providedValue := GetValueType(node.Value, env)
		err := MatchTypes(expectedTypeInterface, providedValue)
		if err != nil {
			errgen.AddError(env.filePath, node.Value.StartPos().Line, node.Value.EndPos().Line, node.Value.StartPos().Column, node.Value.EndPos().Column, err.Error())
		}
	}
	err := env.DeclareVar(varToDecl.Name, expectedTypeInterface, node.IsConst, false)
	if err != nil {

		errgen.AddError(env.filePath, node.Variable.StartPos().Line, node.Variable.EndPos().Line, node.Variable.StartPos().Column, node.Variable.EndPos().Column, err.Error())
	}

	if node.IsConst {
		fmt.Print("Declared constant variable ")
		utils.ColorPrint(utils.RED, varToDecl.Name + "\n")
	} else {
		fmt.Print("Declared variable ")
		utils.ColorPrint(utils.RED, varToDecl.Name + "\n")
	}

	//return the type of the variable
	return expectedTypeInterface
}
