package parser

import (
	"walrus/ast"
	"walrus/errgen"
	"walrus/lexer"
)

// parseVarDeclStmt parses a variable declaration statement in the source code.
// It handles both `let` and `const` declarations, with optional type annotations
// and initial values.
//
// Parameters:
// - p: A pointer to the Parser instance.
//
// Returns:
// - An ast.Node representing the parsed variable declaration statement.
//
// The function performs the following steps:
// 1. Advances the parser to consume the `let` or `const` keyword.
// 2. Determines if the declaration is a constant.
// 3. Expects and consumes an identifier token for the variable name.
// 4. Optionally parses an explicit type if a colon `:` is present.
// 5. Parses the assignment operator `:=` or `=` and the initial value expression if present.
// 6. Ensures that constants have an initial value.
// 7. Expects and consumes a semicolon `;` to terminate the statement.
// 8. Constructs and returns an ast.VarDeclStmt node with the parsed information.
func parseVarDeclStmt(p *Parser) ast.Node {

	declToken := p.advance() // advance the let/const keyword

	// is it let or const?
	isConst := declToken.Kind == lexer.CONST_TOKEN

	// parse the variable name
	identifier := p.expect(lexer.IDENTIFIER_TOKEN)

	// parse the explicit type if present. This will be nil if no type is specified.
	var explicitType ast.DataType

	var value ast.Node

	assignmentToken := p.advance()

	if assignmentToken.Kind == lexer.COLON_TOKEN {
		// syntax is let a : <type>
		explicitType = parseType(p, DEFAULT_BP)
	} else if assignmentToken.Kind != lexer.WALRUS_TOKEN {
		msg := "Invalid variable declaration syntax"
		errgen.MakeError(p.FilePath, assignmentToken.Start.Line, assignmentToken.End.Line, assignmentToken.Start.Column, assignmentToken.End.Column, msg).AddHint("Maybe you want to use : or := instead of =", errgen.TEXT_HINT).DisplayWithPanic()
	}

	if p.currentTokenKind() != lexer.SEMI_COLON_TOKEN {
		// then we have an assignment
		if assignmentToken.Kind == lexer.COLON_TOKEN {
			p.expect(lexer.EQUALS_TOKEN)
		}
		value = parseExpr(p, ASSIGNMENT_BP)
	}

	//if const, we must have a value
	if isConst && value == nil {
		msg := "constants must have value when declared"
		errgen.MakeError(p.FilePath, p.currentToken().Start.Line, p.currentToken().End.Line, p.currentToken().Start.Column, p.currentToken().End.Column, msg).DisplayWithPanic()
	}

	end := p.expect(lexer.SEMI_COLON_TOKEN).End

	node := ast.VarDeclStmt{
		Variable: ast.IdentifierExpr{
			Name: identifier.Value,
			Location: ast.Location{
				Start: identifier.Start,
				End:   identifier.End,
			},
		},
		Value:        value,
		ExplicitType: explicitType,
		IsConst:      isConst,
		Location: ast.Location{
			Start: declToken.Start,
			End:   end,
		},
	}

	return node
}

// parseVarAssignmentExpr parses a variable assignment expression in the source code.
// It takes a parser instance, the left-hand side node, and the binding power as arguments.
// The function ensures that the left-hand side of the assignment is a valid identifier,
// array index access, or struct property access. If not, it generates an error message.
// It then advances the parser to the operator, parses the right-hand side expression,
// and constructs a VarAssignmentExpr node with the appropriate location information.
//
// Parameters:
//   - p: *Parser - The parser instance.
//   - left: ast.Node - The left-hand side node of the assignment.
//   - bp: BINDING_POWER - The binding power for the expression.
//
// Returns:
//   - ast.Node - The parsed variable assignment expression node.
func parseVarAssignmentExpr(p *Parser, left ast.Node, bp BINDING_POWER) ast.Node {

	operator := p.advance()

	right := parseExpr(p, bp)

	endPos := right.EndPos()

	return ast.VarAssignmentExpr{
		Assignee: left,
		Value:    right,
		Operator: operator,
		Location: ast.Location{
			Start: left.StartPos(),
			End:   endPos,
		},
	}
}
