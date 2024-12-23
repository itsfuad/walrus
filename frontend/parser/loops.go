package parser

import (
	//Standard packages
	"fmt"

	//Walrus packages
	"walrus/frontend/ast"
	"walrus/frontend/lexer"
	"walrus/position"
	"walrus/report"
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
	loopType := p.advance() // advance past the 'for|foreach' keyword

	if loopType.Kind == lexer.FOR_TOKEN {
		// empty loop is an infinite loop for {} | detection method: no expressions are present
		// condition-only loop for condition { } | detection method: only one expression is present
		// traditional for loop for init; condition; increment { } | detection method: three expressions are present
		var init, cond, incr ast.Node

		// check if there is an opening brace
		if p.currentTokenKind() == lexer.OPEN_CURLY {
			// empty loop
			block := parseBlock(p)
			return ast.ForStmt{
				Init:      nil,
				Condition: nil,
				Increment: nil,
				Block:     block,
				Location: position.Location{
					Start: loopType.Start,
					End:   block.EndPos(),
				},
			}
		}

		// check if there is an init expression
		init = parseNode(p)

		// check if there no opening brace, then there is a condition and increment
		if p.currentTokenKind() != lexer.OPEN_CURLY {
			fmt.Printf("Current token: %v\n", p.currentTokenKind())
			cond = parseExpr(p, DEFAULT_BP)
			p.expect(lexer.SEMI_COLON_TOKEN)
			incr = parseExpr(p, DEFAULT_BP)
		} else {
			// condition-only loop
			cond = init
			init = nil
		}

		// parse the block
		block := parseBlock(p)

		return ast.ForStmt{
			Init:      init,
			Condition: cond,
			Increment: incr,
			Block:     block,
			Location: position.Location{
				Start: loopType.Start,
				End:   block.EndPos(),
			},
		}

	} else if loopType.Kind == lexer.FOREACH_TOKEN {
		// for-each loop foreach v in arr { } | detection method: one identifier 'v' is present
		// for-each loop foreach i, v in arr { } | detection method: two identifiers 'i' and 'v' are present
		var first, second ast.Node
		first = parseExpr(p, ASSIGNMENT_BP)
		if p.currentTokenKind() == lexer.COMMA_TOKEN {
			p.advance()
			second = parseExpr(p, ASSIGNMENT_BP)
		} else {
			second = first
			first = nil
		}

		p.expect(lexer.IN_TOKEN)

		// parse the array expression
		array := parseExpr(p, ASSIGNMENT_BP)

		p.expect(lexer.OPEN_CURLY)

		// parse the block
		block := parseBlock(p)

		return ast.ForEachStmt{
			Key:      first,
			Value:    second,
			Iterable: array,
			Block:    block,
			Location: position.Location{
				Start: loopType.Start,
				End:   block.EndPos(),
			},
		}
	} else {
		report.Add(p.FilePath, loopType.Start.Line, loopType.End.Line, loopType.Start.Column, loopType.End.Column, "Expected 'for' or 'foreach' keyword").Level(report.SYNTAX_ERROR)
	}

	return nil
}
