package parser

import (
	"walrus/frontend/ast"
	"walrus/frontend/lexer"
)

// parseSafeStmt parses a safe statement in the source code.
// A safe statement consists of a safe block and an unsafe block.
// The function expects the parser to be positioned at the start of the safe statement.
//
// The structure of a safe statement is as follows:
//
//	safe <identifier> {
//	    // safe block
//	} otherwise {
//
//	    // unsafe block
//	}
//
// The function advances the parser, extracts the identifier, parses the safe block,
// expects an 'otherwise' token, and then parses the unsafe block.
//
// Parameters:
// - p: A pointer to the Parser instance.
//
// Returns:
// - An ast.Node representing the parsed safe statement.
func parseSafeStmt(p *Parser) ast.Node {
	// safe block

	start := p.eat().Start // eat safe token

	varName := p.expect(lexer.IDENTIFIER_TOKEN)

	//now we are in the safe block
	safeBody := parseBlock(p)

	//unsafe block
	p.expect(lexer.OTHERWISE_TOKEN) // eat unsafe token

	//now we are in the unsafe block
	unsafeBody := parseBlock(p)

	return ast.SafeStmt{
		Value: ast.IdentifierExpr{
			Name: varName.Value,
			Location: ast.Location{
				Start: varName.Start,
				End:   varName.End,
			},
		},
		SafeBlock:   safeBody,
		UnsafeBlock: unsafeBody,
		Location: ast.Location{
			Start: start,
			End:   unsafeBody.End,
		},
	}
}
