package parser

import (
	"fmt"
	"walrus/ast"
	"walrus/lexer"
)

func parseMapLiteral(p *Parser) ast.Node {

	fmt.Println("Parsing map literal")
	
	start := p.expect(lexer.MAP_TOKEN).Start

	//parse the opening curly brace
	p.expect(lexer.OPEN_CURLY)

	//parse the values
	values := make(map[ast.Node]ast.Node)

	for p.hasToken() && p.currentTokenKind() != lexer.CLOSE_CURLY {
		//parse the key
		key := parseExpr(p, DEFAULT_BP)
		//parse the colon
		p.expect(lexer.COLON_TOKEN)
		//parse the value
		value := parseExpr(p, DEFAULT_BP)

		values[key] = value

		if p.currentTokenKind() == lexer.COMMA_TOKEN {
			p.expect(lexer.COMMA_TOKEN)
		}
	}

	end := p.expect(lexer.CLOSE_CURLY).End

	return ast.MapLiteral{
		Values: values,
		Location: ast.Location{
			Start: start,
			End:   end,
		},
	}
}