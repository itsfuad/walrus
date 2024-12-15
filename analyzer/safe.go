package analyzer

import (
	"fmt"
	"walrus/ast"
	"walrus/errgen"
)

// checkSafeStmt checks the safety of a given SafeStmt node within the provided type environment.
// It ensures that the value of the node is of 'maybe' type and validates the safe and unsafe blocks.
//
// Parameters:
// - node: The SafeStmt node to be checked.
// - env: The type environment in which the node is being checked.
//
// Returns:
// - TcValue: A type-checked value indicating the result of the safety check.
//
// Errors:
// - Adds a critical error if the node's value is not of 'maybe' type.
// - Adds normal errors if there are issues within the safe or unsafe blocks.
func checkSafeStmt(node ast.SafeStmt, env *TypeEnvironment) ExprType {

	maybeVar := parseNodeValue(node.Value, env)

	if maybeVar.DType() != MAYBE_TYPE {
		errgen.Add(env.filePath, node.Value.Start.Line, node.Value.End.Line, node.Value.Start.Column, node.Value.End.Column, "safe-otherwise can only be used with 'maybe' types").Level(errgen.CRITICAL_ERROR)
	}

	// check the safe block where the maybe type is the type of the defined type
	err := checkSafeBlock(env, node.Value.Name, node.SafeBlock, maybeVar.(Maybe))
	if err != nil {
		errgen.Add(env.filePath, node.SafeBlock.StartPos().Line, node.SafeBlock.EndPos().Line, node.SafeBlock.StartPos().Column, node.SafeBlock.EndPos().Column, err.Error()).Level(errgen.NORMAL_ERROR)
	}
	err = checkOtherwiseBlock(env, node.Value.Name, node.UnsafeBlock)
	if err != nil {
		errgen.Add(env.filePath, node.UnsafeBlock.StartPos().Line, node.UnsafeBlock.EndPos().Line, node.UnsafeBlock.StartPos().Column, node.UnsafeBlock.EndPos().Column, err.Error()).Level(errgen.NORMAL_ERROR)
	}

	return NewVoid()
}

func checkSafeBlock(env *TypeEnvironment, name string, block ast.BlockStmt, value Maybe) error {

	//new scope for the safe block
	safeScope := NewTypeENV(env, SAFE_SCOPE, "safe block", env.filePath)

	//declare the variable in the safe block
	err := safeScope.declareVar(name, value.MaybeType, false, false)
	if err != nil {
		return fmt.Errorf("error declaring variable '%s' in safe block. "+err.Error(), name)
	}

	//check the block
	for _, stmt := range block.Contents {
		CheckAST(stmt, safeScope)
	}

	return nil
}

func checkOtherwiseBlock(env *TypeEnvironment, name string, block ast.BlockStmt) error {

	//new scope for the unsafe block
	unsafeScope := NewTypeENV(env, OTHERWISE_SCOPE, "unsafe block", env.filePath)

	//declare the variable in the unsafe block
	err := unsafeScope.declareVar(name, NewNull(), false, false)
	if err != nil {
		return fmt.Errorf("error declaring variable '%s' in unsafe block. "+err.Error(), name)
	}

	//check the block
	for _, stmt := range block.Contents {
		CheckAST(stmt, unsafeScope)
	}

	return nil
}
