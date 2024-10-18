package parser

import (
	"walrus/ast"
	"walrus/lexer"
)

// parseArrayExpr parses an array expression from the input tokens.
// It expects the current token to be the opening bracket '[' of the array.
// The function will consume tokens until it reaches the closing bracket ']'.
// Elements within the array are expected to be separated by commas.
// Returns an ast.ArrayExpr node representing the parsed array expression.
//
// Parameters:
//   - p: A pointer to the Parser instance.
//
// Returns:
//   - ast.Node: An AST node representing the parsed array expression.
func parseArrayExpr(p *Parser) ast.Node {

	start := p.advance().Start //eat the [ token

	var values []ast.Node

	for p.currentTokenKind() != lexer.CLOSE_BRACKET {
		value := parseExpr(p, PRIMARY_BP)
		values = append(values, value)
		if p.currentTokenKind() != lexer.CLOSE_BRACKET {
			p.expect(lexer.COMMA_TOKEN)
		}
	}

	end := p.expect(lexer.CLOSE_BRACKET).End

	return ast.ArrayLiteral{
		Values: values,
		Location: ast.Location{
			Start: start,
			End:   end,
		},
	}
}

// parseArrayAccess parses an array access expression from the input.
// It expects the current token to be an opening bracket '[' and parses
// the index expression inside the brackets. The function returns an
// ast.ArrayIndexAccess node representing the array access.
//
// Parameters:
// - p: A pointer to the Parser instance.
// - left: The left-hand side node, representing the array being accessed.
// - bp: The binding power for the expression parsing.
//
// Returns:
// - An ast.Node representing the array index access.
func parseArrayAccess(p *Parser, left ast.Node, bp BINDING_POWER) ast.Node {
	p.expect(lexer.OPEN_BRACKET)
	index := parseExpr(p, bp)
	end := p.expect(lexer.CLOSE_BRACKET).End
	return ast.ArrayIndexAccess{
		Arrayvalue: left,
		Index:      index,
		Location: ast.Location{
			Start: left.StartPos(),
			End:   end,
		},
	}
}