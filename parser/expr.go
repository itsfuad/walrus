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

// parseVarAssignmentExpr parses a variable assignment expression in the source code.
// It takes a parser instance, the left-hand side node, and the binding power as arguments.
// The function ensures that the left-hand side of the assignment is a valid identifier,
// array index access, or struct property access. If not, it generates an error message.
// It then advances the parser to the operator, parses the right-hand side expression,
// and constructs a VarAssignmentExpr node with the appropriate location information.
//
// Parameters:
//   - p: *Parser - The parser instance.
//   - left: ast.Node - The left-hand side node of the assignment.
//   - bp: BINDING_POWER - The binding power for the expression.
//
// Returns:
//   - ast.Node - The parsed variable assignment expression node.
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
	props := map[string]ast.Node{}

	for p.hasToken() && p.currentTokenKind() != lexer.CLOSE_CURLY {
		//we expect an identifier
		iden := p.expect(lexer.IDENTIFIER_TOKEN)
		//then we expect colon
		p.expect(lexer.COLON_TOKEN)
		//now we expect value as expression
		val := parseExpr(p, DEFAULT_BP)

		if _, ok := props[iden.Value]; ok {
			errgen.MakeError(p.FilePath, iden.Start.Line, iden.End.Line, iden.Start.Column, iden.End.Column, fmt.Sprintf("property '%s' was previously assigned", iden.Value)).Display()
		}

		props[iden.Value] = val

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
		errgen.MakeError(p.FilePath, left.StartPos().Line, left.EndPos().Line, left.StartPos().Column, left.EndPos().Column, "only identifiers can be incremented or decremented").Display()
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

// parseCallExpr parses a function call expression in the source code.
// It expects the current token to be an open parenthesis and will parse
// all arguments until it encounters a closing parenthesis.
//
// Parameters:
// - p: A pointer to the Parser instance.
// - left: The left-hand side node, typically the function being called.
// - bp: The binding power, which determines the precedence of the expression.
//
// Returns:
//   - An ast.Node representing the function call expression, including the
//     caller, arguments, and their location in the source code.
func parseCallExpr(p *Parser, left ast.Node, bp BINDING_POWER) ast.Node {

	p.advance() //eat the open paren
	startPos := left.StartPos()
	var args []ast.Node
	// parse the arguments
	for p.currentTokenKind() != lexer.CLOSE_PAREN {
		arg := parseExpr(p, DEFAULT_BP)
		args = append(args, arg)
		if p.currentTokenKind() != lexer.CLOSE_PAREN {
			p.expect(lexer.COMMA_TOKEN)
		}
	}

	endPos := p.expect(lexer.CLOSE_PAREN).End

	return ast.FunctionCallExpr{
		Caller:    left,
		Arguments: args,
		Location: ast.Location{
			Start: startPos,
			End:   endPos,
		},
	}
}
