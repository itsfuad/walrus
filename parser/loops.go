package parser

import (
	"walrus/ast"
	"walrus/lexer"
	"walrus/errgen"
)

// parseForStmt parses a 'for' statement in the source code.
// It handles different types of 'for' loops:
// - Infinite loop: `for { }`
// - Condition-only loop: `for i < 10 { }`
// - Traditional for loop with initialization, condition, and increment: `for i := 0; i < 10; i++ { }`
// - For-each loop: `for v in arr { }` or `for i, v in arr { }`
//
// Parameters:
// - p: A pointer to the Parser instance.
//
// Returns:
// - An ast.Node representing the parsed 'for' statement.
// - If the syntax is invalid, it generates an error and returns nil.
func parseForStmt(p *Parser) ast.Node {
	start := p.advance().Start // eat for token
	//first token is either an identifier or open curly
	if p.currentTokenKind() == lexer.OPEN_CURLY {
		// infinite loop
		block := parseBlock(p)
		return ast.ForStmt{
			Start:     nil,
			Condition: nil,
			Increment: nil,
			Block:     block,
			Location: ast.Location{
				Start: start,
				End:   block.End,
			},
		}
	} else if p.currentTokenKind() == lexer.IDENTIFIER_TOKEN {

		identifier := p.currentToken()

		idententifierExpr := ast.IdentifierExpr{
			Name: identifier.Value,
			Location: ast.Location{
				Start: identifier.Start,
				End:   identifier.End,
			},
		}

		switch p.nextTokenKind() {
		case lexer.WALRUS_TOKEN: // for i := 0; i < 10; i++ { }
			p.advance() // eat the identifier
			p.advance() // eat the walrus token
			// value of the identifier
			value := parseExpr(p, ASSIGNMENT_BP)
			p.expect(lexer.SEMI_COLON_TOKEN)
			// condition
			condition := parseExpr(p, ASSIGNMENT_BP)
			p.expect(lexer.SEMI_COLON_TOKEN)
			// increment
			increment := parseExpr(p, ASSIGNMENT_BP)

			block := parseBlock(p)

			return ast.ForStmt{
				Start: ast.VarDeclStmt{
					Variable:     idententifierExpr,
					Value:        value,
					ExplicitType: nil,
					IsConst:      false,
					Location: ast.Location{
						Start: identifier.Start,
						End:   value.EndPos(),
					},
				},
				Condition: condition,
				Increment: increment,
				Block:     block,
				Location: ast.Location{
					Start: start,
					End:   block.End,
				},
			}
		case lexer.IN_TOKEN: // for v in arr { }
			p.advance() // eat the identifier
			p.advance() // eat the in token
			// value of the identifier
			value := parseExpr(p, ASSIGNMENT_BP)

			block := parseBlock(p)

			return ast.ForEachStmt{
				Key:      nil,
				Value:    idententifierExpr,
				Iterable: value,
				Block:    block,
				Location: ast.Location{
					Start: start,
					End:   block.End,
				},
			}
		case lexer.COMMA_TOKEN: // for i, v in arr { }
			p.advance() // eat the identifier
			p.advance() // eat the comma token
			// value of the identifier
			key := idententifierExpr
			value := ast.IdentifierExpr{
				Name: p.expect(lexer.IDENTIFIER_TOKEN).Value,
				Location: ast.Location{
					Start: identifier.Start,
					End:   identifier.End,
				},
			}

			p.expect(lexer.IN_TOKEN)

			// value of the identifier
			iterable := parseExpr(p, ASSIGNMENT_BP)

			block := parseBlock(p)

			return ast.ForEachStmt{
				Key:      key,
				Value:    value,
				Iterable: iterable,
				Block:    block,
				Location: ast.Location{
					Start: start,
					End:   block.End,
				},
			}

		default:
			// for i < 10 { } // condition only
			condition := parseExpr(p, ASSIGNMENT_BP)
			block := parseBlock(p)
			return ast.ForStmt{
				Start: ast.IdentifierExpr{
					Name: identifier.Value,
					Location: ast.Location{
						Start: identifier.Start,
						End:   identifier.End,
					},
				},
				Condition: condition,
				Increment: nil,
				Block:     block,
				Location: ast.Location{
					Start: start,
					End:   block.End,
				},
			}
		}

	} else {
		//error
		//msg := "invalid for loop syntax"
		errgen.MakeError(p.FilePath, start.Line, start.Line, start.Column, start.Column, "invalid for loop syntax").Display()
		return nil
	}
}