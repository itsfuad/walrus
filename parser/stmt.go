package parser

import (
	"walrus/ast"
	"walrus/lexer"
)

// parseUserDefinedTypeStmt parses a user-defined type statement in the source code.
// It expects the 'type' keyword followed by an identifier that starts with a capital letter.
// If the identifier does not start with a capital letter, an error is generated with a hint.
// The function then parses the user-defined type and expects a semicolon at the end.
// It returns an AST node representing the type declaration statement.
//
// Parameters:
// - p: A pointer to the Parser instance.
//
// Returns:
// - ast.Node: An AST node representing the type declaration statement.
func parseUserDefinedTypeStmt(p *Parser) ast.Node {

	start := p.advance().Start //eat type token

	typeName := p.expect(lexer.IDENTIFIER_TOKEN)

	udType := parseUDTType(p)

	p.expect(lexer.SEMI_COLON_TOKEN)

	return ast.TypeDeclStmt{
		UDType:     udType,
		UDTypeName: typeName.Value,
		Location: ast.Location{
			Start: start,
			End:   udType.EndPos(),
		},
	}
}

// parseBlock parses a block statement from the input tokens.
// It expects the block to start with an opening curly brace '{' and end with a closing curly brace '}'.
// The function iterates over the tokens within the braces, parsing each node and adding it to the block's body.
// It returns an ast.BlockStmt containing the parsed nodes and their location in the source code.
//
// Parameters:
//   - p: A pointer to the Parser instance.
//
// Returns:
//   - ast.BlockStmt: The parsed block statement, including its contents and location.
func parseBlock(p *Parser) ast.BlockStmt {

	start := p.expect(lexer.OPEN_CURLY).Start

	body := make([]ast.Node, 0)

	for p.hasToken() && p.currentTokenKind() != lexer.CLOSE_CURLY {
		body = append(body, parseNode(p))
	}

	end := p.expect(lexer.CLOSE_CURLY).End

	return ast.BlockStmt{
		Contents: body,
		Location: ast.Location{
			Start: start,
			End:   end,
		},
	}
}

// parseReturnStmt parses a return statement in the source code.
// It expects the current token to be a return token and advances the parser.
// If the next token is not a semicolon, it parses an expression for the return value.
// Finally, it expects a semicolon to end the return statement and returns an ast.ReturnStmt node.
//
// Parameters:
//   - p: A pointer to the Parser instance.
//
// Returns:
//   - An ast.Node representing the return statement.
func parseReturnStmt(p *Parser) ast.Node {

	start := p.advance().Start // eat return token

	var value ast.Node

	if p.currentTokenKind() != lexer.SEMI_COLON_TOKEN {
		value = parseExpr(p, ASSIGNMENT_BP)
	}

	end := p.expect(lexer.SEMI_COLON_TOKEN).End

	return ast.ReturnStmt{
		Value: value,
		Location: ast.Location{
			Start: start,
			End:   end,
		},
	}
}
