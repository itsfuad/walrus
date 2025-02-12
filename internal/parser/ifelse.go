package parser

import (
	"walrus/internal/ast"
	"walrus/internal/lexer"
)

// parseIfStmt parses an if statement from the input and returns an AST node representing the if statement.
// It expects the parser to be positioned at the 'if' token at the start of the if statement.
//
// The function performs the following steps:
// 1. Advances the parser to consume the 'if' token and records the start position.
// 2. Parses the condition expression of the if statement.
// 3. Parses the consequent block of the if statement.
// 4. Checks for the presence of an 'else' token and, if found, parses the alternate block or another if statement.
//
// The returned AST node includes the condition, the consequent block, and optionally the alternate block.
//
// Parameters:
// - p: A pointer to the Parser instance.
//
// Returns:
// - An AST node representing the parsed if statement.
func parseIfStmt(p *Parser) ast.Node {

	start := p.eat().Start // eat if token

	condition := parseExpr(p, ASSIGNMENT_BP)

	//parse block
	consequentBlock := parseBlock(p)

	var alternate ast.Node

	if p.hasToken() && p.currentTokenKind() == lexer.ELSE_TOKEN {
		p.eat() // eat else token
		if p.hasToken() && p.currentTokenKind() == lexer.IF_TOKEN {
			alternate = parseIfStmt(p)
		} else {
			alternate = parseBlock(p)
		}
	}

	return ast.IfStmt{
		Condition:      condition,
		Block:          consequentBlock,
		AlternateBlock: alternate,
		Location: ast.Location{
			Start: start,
			End:   consequentBlock.End,
		},
	}
}
