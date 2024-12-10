package analyzer

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
func checkVariableAssignment(node ast.VarAssignmentExpr, env *TypeEnvironment) TcValue {

	Assignee := node.Assignee
	valueToAssign := node.Value

	if err := checkLValue(Assignee, env); err != nil {
		errgen.AddError(env.filePath, Assignee.StartPos().Line, Assignee.EndPos().Line, Assignee.StartPos().Column, Assignee.EndPos().Column, "cannot assign to "+err.Error()).ErrorLevel(errgen.CRITICAL)

	}

	expectedType := CheckAST(Assignee, env)
	providedType := CheckAST(valueToAssign, env)

	err := matchTypes(expectedType, providedType)
	if err != nil {

		errgen.AddError(env.filePath, valueToAssign.StartPos().Line, valueToAssign.EndPos().Line, valueToAssign.StartPos().Column, valueToAssign.EndPos().Column, err.Error()).ErrorLevel(errgen.NORMAL)
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
func checkVariableDeclaration(node ast.VarDeclStmt, env *TypeEnvironment) TcValue {

	varsToDecl := node.Variables

	for _, varToDecl := range varsToDecl {

		utils.BLUE.Print("Declaring variable ")
		utils.RED.Println(varToDecl.Identifier.Name)

		var expectedTypeInterface TcValue

		// let a : int = 5;

		if varToDecl.ExplicitType != nil {
			expectedTypeInterface = evaluateTypeName(varToDecl.ExplicitType, env)
			fmt.Print("Explicit type: ")
			utils.PURPLE.Println(tcValueToString(expectedTypeInterface))
		} else {
			expectedTypeInterface = CheckAST(varToDecl.Value, env)
			utils.ORANGE.Print("Auto detected type: ")
			utils.PURPLE.Println(tcValueToString(expectedTypeInterface))
		}

		if varToDecl.Value != nil && varToDecl.ExplicitType != nil {
			//providedValue := CheckAST(node.Value, env)
			providedValue := CheckAST(varToDecl.Value, env)
			err := matchTypes(expectedTypeInterface, providedValue)
			if err != nil {
				errgen.AddError(env.filePath, varToDecl.Value.StartPos().Line, varToDecl.Value.EndPos().Line, varToDecl.Value.StartPos().Column, varToDecl.Value.EndPos().Column, fmt.Sprintf("error declaring variable '%s'. %s", varToDecl.Identifier.Name, err.Error())).ErrorLevel(errgen.NORMAL)
			}
		}

		err := env.declareVar(varToDecl.Identifier.Name, expectedTypeInterface, node.IsConst, false)
		if err != nil {
			errgen.AddError(env.filePath, varToDecl.Identifier.StartPos().Line, varToDecl.Identifier.EndPos().Line, varToDecl.Identifier.StartPos().Column, varToDecl.Identifier.EndPos().Column, err.Error()).ErrorLevel(errgen.CRITICAL)
		}

		if node.IsConst {
			utils.GREEN.Print("Declared constant variable ")
			utils.RED.Print(varToDecl.Identifier.Name)
			fmt.Print(" of type ")
			utils.PURPLE.Println(tcValueToString(expectedTypeInterface))
		} else {
			utils.GREEN.Print("Declared variable ")
			utils.RED.Print(varToDecl.Identifier.Name)
			fmt.Print(" of type ")
			utils.PURPLE.Println(tcValueToString(expectedTypeInterface))
		}
	}
	return NewVoid()
}
