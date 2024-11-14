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

	if err := checkLValue(Assignee, env); err != nil {
		errgen.AddError(env.filePath, Assignee.StartPos().Line, Assignee.EndPos().Line, Assignee.StartPos().Column, Assignee.EndPos().Column, "cannot assign to "+err.Error()).DisplayWithPanic()

	}

	expectedType := nodeType(Assignee, env)
	providedType := nodeType(valueToAssign, env)

	err := matchTypes(expectedType, providedType)
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

	varsToDecl := node.Variables

	for _, varToDecl := range varsToDecl {

		fmt.Print("Declaring variable ")
		utils.RED.Println(varToDecl.Identifier.Name)

		var expectedTypeInterface ValueTypeInterface

		// let a : int = 5;

		if varToDecl.ExplicitType != nil {
			expectedTypeInterface = evaluateTypeName(varToDecl.ExplicitType, env)
			fmt.Print("Explicit type: ")
			utils.PURPLE.Println(string(valueTypeInterfaceToString(expectedTypeInterface)))
		} else {
			expectedTypeInterface = nodeType(varToDecl.Value, env)
			fmt.Print("Auto detected type: ")
			utils.PURPLE.Println(string(valueTypeInterfaceToString(expectedTypeInterface)))
		}

		if varToDecl.Value != nil && varToDecl.ExplicitType != nil {
			//providedValue := CheckAST(node.Value, env)
			providedValue := nodeType(varToDecl.Value, env)
			err := matchTypes(expectedTypeInterface, providedValue)
			if err != nil {
				errgen.AddError(env.filePath, varToDecl.Value.StartPos().Line, varToDecl.Value.EndPos().Line, varToDecl.Value.StartPos().Column, varToDecl.Value.EndPos().Column, err.Error())
			}
		}

		err := env.DeclareVar(varToDecl.Identifier.Name, expectedTypeInterface, node.IsConst, false)
		if err != nil {

			errgen.AddError(env.filePath, varToDecl.Start.Line, varToDecl.End.Line, varToDecl.Start.Column, varToDecl.End.Column, err.Error())
		}

		if node.IsConst {
			utils.BLUE.Print("Declared constant variable ")
			utils.RED.Print(varToDecl.Identifier.Name)
			fmt.Print(" of type ")
			utils.PURPLE.Println(string(valueTypeInterfaceToString(expectedTypeInterface)))
		} else {
			utils.BLUE.Print("Declared variable ")
			utils.RED.Print(varToDecl.Identifier.Name)
			fmt.Print(" of type ")
			utils.PURPLE.Println(string(valueTypeInterfaceToString(expectedTypeInterface)))
		}

		//return the type of the variable
	}
	return NewVoid()
}
