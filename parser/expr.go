package parser

import (
	"fmt"
	"walrus/ast"
	"walrus/builtins"
	"walrus/errgen"
	"walrus/lexer"
)

// parseExpr parses an expression with the given binding power.
// It first parses the NUD (Null Denotation) of the expression,
// then continues to parse the LED (Left Denotation) of the expression
// until the binding power of the current token is less than or equal to the given binding power.
//
// The parsed expression is returned as an ast.Node.
//
// bp parameter is the limit.
// parser will go down the BINDING_POWER table until it reaches the limit.
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
		errgen.AddError(p.FilePath, p.currentToken().Start.Line, p.currentToken().End.Line, p.currentToken().Start.Column, p.currentToken().End.Column, msg).ErrorLevel(errgen.CRITICAL)
	}

	expr := nudFunction(p)

	for GetBP(p.currentTokenKind()) > bp {

		tokenKind = p.currentTokenKind()

		ledFunction, exists := LEDLookup[tokenKind]

		if !exists {
			msg := fmt.Sprintf("parser:led:unexpected token %s\n", tokenKind)
			errgen.AddError(p.FilePath, p.currentToken().Start.Line, p.currentToken().End.Line, p.currentToken().Start.Column, p.currentToken().End.Column, msg).ErrorLevel(errgen.CRITICAL)
		}

		expr = ledFunction(p, expr, GetBP(p.currentTokenKind()))
	}

	return expr
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
	case lexer.INT8_TOKEN, lexer.INT16_TOKEN, lexer.INT32_TOKEN, lexer.INT64_TOKEN, lexer.UINT8_TOKEN, lexer.UINT16_TOKEN, lexer.UINT32_TOKEN, lexer.UINT64_TOKEN:
		return ast.IntegerLiteralExpr{
			Value:    rawValue,
			BitSize:  builtins.GetBitSize(builtins.PARSER_TYPE(primaryToken.Kind)),
			IsSigned: builtins.IsSigned(builtins.PARSER_TYPE(primaryToken.Kind)),
			Location: loc,
		}
	case lexer.FLOAT32_TOKEN, lexer.FLOAT64_TOKEN:

		return ast.FloatLiteralExpr{
			Value:    rawValue,
			BitSize:  builtins.GetBitSize(builtins.PARSER_TYPE(primaryToken.Kind)),
			Location: loc,
		}

	case lexer.STR_TOKEN:
		return ast.StringLiteralExpr{
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
		errgen.AddError(p.FilePath, p.currentToken().Start.Line, p.currentToken().End.Line, p.currentToken().Start.Column, p.currentToken().End.Column, msg).ErrorLevel(errgen.CRITICAL)
	}

	return nil
}

// parseGroupingExpr parses a grouping expression enclosed in parentheses.
// It expects an opening parenthesis, followed by an expression, and a closing parenthesis.
// Returns the parsed expression node.
//
// Parameters:
// - p: A pointer to the Parser instance.
//
// Returns:
// - ast.Node: The parsed expression node.
func parseGroupingExpr(p *Parser) ast.Node {
	p.expect(lexer.OPEN_PAREN)
	expr := parseExpr(p, DEFAULT_BP)
	p.expect(lexer.CLOSE_PAREN)
	return expr
}

// parsePostfixExpr parses a postfix expression, which consists of an identifier
// followed by a postfix operator (e.g., increment or decrement).
//
// Parameters:
// - p: A pointer to the Parser instance.
// - left: The left-hand side node, which must be an identifier.
// - bp: The binding power of the operator.
//
// Returns:
// - An ast.Node representing the parsed postfix expression.
//
// Errors:
// - If the left-hand side node is not an identifier, an error is generated and displayed.
func parsePostfixExpr(p *Parser, left ast.Node, bp BINDING_POWER) ast.Node {
	start := left.StartPos()
	// left must be an identifier
	if _, ok := left.(ast.IdentifierExpr); !ok {
		errgen.AddError(p.FilePath, left.StartPos().Line, left.EndPos().Line, left.StartPos().Column, left.EndPos().Column, "only identifiers can be incremented or decremented").ErrorLevel(errgen.CRITICAL)
	}
	operator := p.advance()
	return ast.PostfixExpr{
		Operator: operator,
		Argument: left.(ast.IdentifierExpr),
		Location: ast.Location{
			Start: start,
			End:   operator.End,
		},
	}
}

// parsePrefixExpr parses a prefix expression from the input tokens.
// It expects the current token to be the start of a prefix expression,
// advances to the operator token, and then expects an identifier token
// as the argument. It returns an ast.PrefixExpr node representing the
// parsed prefix expression.
//
// Parameters:
//   - p: A pointer to the Parser instance.
//
// Returns:
//   - ast.Node: An abstract syntax tree node representing the parsed
//     prefix expression.
func parsePrefixExpr(p *Parser) ast.Node {
	start := p.currentToken().Start
	operator := p.advance()
	argument := p.expect(lexer.IDENTIFIER_TOKEN)
	return ast.PrefixExpr{
		Operator: operator,
		Argument: ast.IdentifierExpr{
			Name: argument.Value,
			Location: ast.Location{
				Start: argument.Start,
				End:   argument.End,
			},
		},
		Location: ast.Location{
			Start: start,
			End:   argument.End,
		},
	}
}

// parseUnaryExpr parses a unary expression from the input tokens.
// It expects the current token to be a unary operator (e.g., '-', '!').
// If the operator is valid, it advances the parser and parses the operand expression.
// If the operator is invalid, it generates and displays an error.
// The function returns an ast.UnaryExpr node representing the parsed unary expression.
//
// Parameters:
// - p: A pointer to the Parser instance.
//
// Returns:
// - An ast.Node representing the parsed unary expression.
func parseUnaryExpr(p *Parser) ast.Node {

	start := p.currentToken().Start

	operator := p.advance()

	switch operator.Kind {
	case lexer.MINUS_TOKEN, lexer.NOT_TOKEN:
		break
	default:
		errgen.AddError(p.FilePath, operator.Start.Line, operator.End.Line, operator.Start.Column, operator.End.Column, fmt.Sprintf("invalid unary operator '%s'", operator.Value)).ErrorLevel(errgen.CRITICAL)
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

// parseBinaryExpr parses a binary expression from the input.
// It takes a parser, a left-hand side node, and a binding power as arguments.
// It returns an AST node representing the binary expression.
//
// Parameters:
//   - p: A pointer to the Parser instance.
//   - left: The left-hand side node of the binary expression.
//   - bp: The binding power of the binary operator.
//
// Returns:
//   - An AST node representing the binary expression.
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

func parseTypeCastExpr(p *Parser, left ast.Node, bp BINDING_POWER) ast.Node {
	start := left.StartPos()
	p.expect(lexer.AS_TOKEN)
	castType := parseType(p, bp)

	return ast.TypeCastExpr{
		Expression: left,
		ToCast:     castType,
		Location: ast.Location{
			Start: start,
			End:   castType.EndPos(),
		},
	}
}
