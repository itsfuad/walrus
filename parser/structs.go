package parser

import (
	"fmt"
	"walrus/ast"
	"walrus/errgen"
	"walrus/lexer"
)

// parseStructLiteral parses a struct literal from the input tokens.
// It expects the following sequence of tokens:
// - An '@' token indicating the start of a struct literal.
// - An identifier token representing the struct name.
// - An opening curly brace '{'.
// - A series of property definitions, each consisting of:
//   - An identifier token for the property name.
//   - A colon ':' token.
//   - An expression representing the property value.
//   - An optional comma ',' token if there are more properties.
//
// - A closing curly brace '}'.
//
// The function returns an ast.Node representing the parsed struct literal.
//
// Parameters:
// - p: A pointer to the Parser instance.
//
// Returns:
// - An ast.Node representing the parsed struct literal.
func parseStructLiteral(p *Parser) ast.Node {

	start := p.expect(lexer.AT_TOKEN).Start

	idetifierToken := p.expectError(lexer.IDENTIFIER_TOKEN, fmt.Errorf("expected a struct name"))

	identidier := ast.IdentifierExpr{
		Name: idetifierToken.Value,
		Location: ast.Location{
			Start: idetifierToken.Start,
			End:   idetifierToken.End,
		},
	}

	p.expect(lexer.OPEN_CURLY)

	//parse the values
	props := make(map[string]ast.StructProp)

	for p.hasToken() && p.currentTokenKind() != lexer.CLOSE_CURLY {
		//we expect an identifier
		iden := p.expect(lexer.IDENTIFIER_TOKEN)
		//then we expect colon
		p.expect(lexer.COLON_TOKEN)
		//now we expect value as expression
		val := parseExpr(p, DEFAULT_BP)

		if _, ok := props[iden.Value]; ok {
			errgen.AddError(p.FilePath, iden.Start.Line, iden.End.Line, iden.Start.Column, iden.End.Column, fmt.Sprintf("property '%s' was previously assigned", iden.Value), errgen.ERROR_CRITICAL)
		}

		props[iden.Value] = ast.StructProp{
			Prop: ast.IdentifierExpr{
				Name: iden.Value,
				Location: ast.Location{
					Start: iden.Start,
					End:   iden.End,
				},
			},
			Value: val,
		}

		//if the next token is not } then we have more values
		if p.currentTokenKind() != lexer.CLOSE_CURLY {
			//we expect comma
			p.expect(lexer.COMMA_TOKEN)
		}
	}

	end := p.expect(lexer.CLOSE_CURLY).End

	structVal := ast.StructLiteral{
		Identifier: identidier,
		Properties: props,
		Location: ast.Location{
			Start: start,
			End:   end,
		},
	}

	return structVal
}

// parsePropertyExpr parses a property access expression from the parser.
// It expects the current token to be a dot (.) followed by an identifier.
// The function constructs and returns an AST node representing the property access.
//
// Parameters:
// - p: A pointer to the Parser instance.
// - left: The left-hand side node of the property access expression.
// - bp: The binding power (precedence) of the expression.
//
// Returns:
// - An AST node representing the property access expression.
func parsePropertyExpr(p *Parser, left ast.Node, bp BINDING_POWER) ast.Node {

	p.expect(lexer.DOT_TOKEN)

	identifier := p.expect(lexer.IDENTIFIER_TOKEN)

	property := ast.IdentifierExpr{
		Name: identifier.Value,
		Location: ast.Location{
			Start: identifier.Start,
			End:   identifier.End,
		},
	}

	return ast.StructPropertyAccessExpr{
		Object:   left,
		Property: property,
		Location: ast.Location{
			Start: left.StartPos(),
			End:   property.End,
		},
	}
}
