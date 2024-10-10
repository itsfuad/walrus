package parser

import (
	"fmt"
	"walrus/ast"
	"walrus/errgen"
	"walrus/lexer"
)

// parseExpr parses an expression with the given binding power.
// It first parses the NUD (Null Denotation) of the expression,
// then continues to parse the LED (Left Denotation) of the expression
// until the binding power of the current token is less than or equal to the given binding power.
// The parsed expression is returned as an ast.Node.
func parseExpr(p *Parser, bp BINDING_POWER) ast.Node {

	// Fist parse the NUD
	token := p.currentToken()

	tokenKind := token.Kind

	nudFunction, exists := NUDLookup[tokenKind]

	if !exists {
		var msg string
		if lexer.IsKeyword(string(tokenKind)) {
			msg = fmt.Sprintf("parser:nud:unexpected keyword '%s'\n", tokenKind)
		} else {
			msg = fmt.Sprintf("parser:nud:unexpected token '%s'\n", tokenKind)
		}
		errgen.MakeError(p.FilePath, p.currentToken().Start.Line, p.currentToken().End.Line, p.currentToken().Start.Column, p.currentToken().End.Column, msg).Display()
	}

	left := nudFunction(p)

	for GetBP(p.currentTokenKind()) > bp {

		tokenKind = p.currentTokenKind()

		ledFunction, exists := LEDLookup[tokenKind]

		if !exists {
			msg := fmt.Sprintf("parser:led:unexpected token %s\n", tokenKind)
			errgen.MakeError(p.FilePath, p.currentToken().Start.Line, p.currentToken().End.Line, p.currentToken().Start.Column, p.currentToken().End.Column, msg).Display()
		}

		left = ledFunction(p, left, GetBP(p.currentTokenKind()))
	}

	return left
}

// parsePrimaryExpr parses a primary expression in the input stream.
// It handles numeric literals, string literals, identifiers, boolean literals, and null literals.
// If the current token does not match any of these types, it panics with an error message.
func parsePrimaryExpr(p *Parser) ast.Node {

	startpos := p.currentToken().Start

	endpos := p.currentToken().End

	primaryToken := p.advance()

	rawValue := primaryToken.Value

	loc := ast.Location{
		Start: startpos,
		End:   endpos,
	}

	switch primaryToken.Kind {
	case lexer.INT:
		return ast.IntegerLiteralExpr{
			Value:    rawValue,
			Location: loc,
		}
	case lexer.FLOAT:

		return ast.FloatLiteralExpr{
			Value:    rawValue,
			Location: loc,
		}

	case lexer.STR:
		return ast.StringLiteralExpr{
			Value:    rawValue,
			Location: loc,
		}
	case lexer.BYTE:
		return ast.CharLiteralExpr{
			Value:    rawValue,
			Location: loc,
		}
	case lexer.IDENTIFIER_TOKEN:
		return ast.IdentifierExpr{
			Name:     rawValue,
			Location: loc,
		}
	default:
		msg := fmt.Sprintf("Cannot create primary expression from %s\n", primaryToken.Value)
		errgen.MakeError(p.FilePath, p.currentToken().Start.Line, p.currentToken().End.Line, p.currentToken().Start.Column, p.currentToken().End.Column, msg).Display()
	}

	return nil
}

func parseGroupingExpr(p *Parser) ast.Node {
	p.expect(lexer.OPEN_PAREN)
	expr := parseExpr(p, DEFAULT_BP)
	p.expect(lexer.CLOSE_PAREN)
	return expr
}

func parseVarAssignmentExpr(p *Parser, left ast.Node, bp BINDING_POWER) ast.Node {

	switch left.(type) {
	case ast.IdentifierExpr:
		break
	case ast.ArrayIndexAccess:
		break
	case ast.StructPropertyAccessExpr:
		break
	default:
		errMsg := "cannot assign to a non-identifier\n"
		errgen.MakeError(p.FilePath, left.StartPos().Line, left.EndPos().Line, left.StartPos().Column, left.EndPos().Column, errMsg).Display()
	}

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

	return ast.ArrayExpr{
		Values: values,
		Location: ast.Location{
			Start: start,
			End:   end,
		},
	}
}

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
	props := map[string]ast.Node{}

	for p.hasToken() && p.currentTokenKind() != lexer.CLOSE_CURLY {
		//we expect an identifier
		iden := p.expect(lexer.IDENTIFIER_TOKEN)
		//then we expect colon
		p.expect(lexer.COLON_TOKEN)
		//now we expect value as expression
		val := parseExpr(p, DEFAULT_BP)

		props[iden.Value] = val

		//if the next token is not } then we have more values
		if p.currentTokenKind() != lexer.CLOSE_CURLY {
			//we expect comma
			p.expect(lexer.COMMA_TOKEN)
		}
	}

	end := p.expect(lexer.CLOSE_CURLY).End

	structVal := ast.StructLiteral{
		Name:       identidier,
		Properties: props,
		Location: ast.Location{
			Start: start,
			End:   end,
		},
	}

	return structVal
}

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

func parseUnaryExpr(p *Parser) ast.Node {

	start := p.currentToken().Start

	operator := p.advance()

	switch operator.Kind {
	case lexer.MINUS_TOKEN, lexer.NOT_TOKEN:
		break
	default:
		errgen.MakeError(p.FilePath, operator.Start.Line, operator.End.Line, operator.Start.Column, operator.End.Column, fmt.Sprintf("invalid unary operator '%s'", operator.Value)).Display()
	}

	argument := parseExpr(p, UNARY_BP)

	return ast.UnaryExpr{
		Operator: operator,
		Argument: argument,
		Location: ast.Location{
			Start: start,
			End:   argument.EndPos(),
		},
	}
}

func parseBinaryExpr(p *Parser, left ast.Node, bp BINDING_POWER) ast.Node {

	op := p.advance()

	right := parseExpr(p, bp)

	return ast.BinaryExpr{
		Operator: op,
		Left:     left,
		Right:    right,
		Location: ast.Location{
			Start: left.StartPos(),
			End:   right.EndPos(),
		},
	}
}

func parseCallExpr(p *Parser, left ast.Node, bp BINDING_POWER) ast.Node {

	if _, ok := left.(ast.IdentifierExpr); !ok {
		errgen.MakeError(p.FilePath, left.StartPos().Line, left.EndPos().Line, left.StartPos().Column, left.EndPos().Column, "only identifiers can be called").Display()
	}

	name := left.(ast.IdentifierExpr)

	p.advance() //eat the open paren
	startPos := left.StartPos()
	var args []ast.Node
	for p.currentTokenKind() != lexer.CLOSE_PAREN {
		arg := parseExpr(p, DEFAULT_BP)
		args = append(args, arg)
		if p.currentTokenKind() != lexer.CLOSE_PAREN {
			p.expect(lexer.COMMA_TOKEN)
		}
	}

	endPos := p.expect(lexer.CLOSE_PAREN).End

	return ast.FunctionCallExpr{
		Name:      name,
		Arguments: args,
		Location: ast.Location{
			Start: startPos,
			End:   endPos,
		},
	}
}
