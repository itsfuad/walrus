package parser

import (
	"walrus/ast"
	"walrus/lexer"
)

func parseMapLiteral(p *Parser) ast.Node {

	mapType := parseMapType(p)

	//parse the opening curly brace
	p.expect(lexer.OPEN_CURLY)

	//parse the values
	values := make([]ast.MapProp, 0)

	for p.hasToken() && p.currentTokenKind() != lexer.CLOSE_CURLY {
		//parse the key
		key := parseExpr(p, DEFAULT_BP)
		//parse the colon
		p.expect(lexer.COLON_TOKEN)
		//parse the value
		value := parseExpr(p, DEFAULT_BP)

		prop := ast.MapProp{
			Key:   key,
			Value: value,
		}

		values = append(values, prop)

		if p.currentTokenKind() == lexer.COMMA_TOKEN {
			p.expect(lexer.COMMA_TOKEN)
		}
	}

	end := p.expect(lexer.CLOSE_CURLY).End

	return ast.MapLiteral{
		MapType: mapType.(ast.MapType),
		Values: values,
		Location: ast.Location{
			Start: mapType.StartPos(),
			End:   end,
		},
	}
}
