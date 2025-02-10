package parser

import (
	//Walrus packages
	"walrus/frontend/ast"
	"walrus/frontend/lexer"
	"walrus/report"
)

// parseLambdaFunction parses a lambda function expression from the input and returns an AST node representing the function.
// It expects the parser to be positioned at the start of the lambda function.
//
// The function performs the following steps:
// 1. Advances the parser to consume the 'fn' token and records the start position.
// 2. Parses the function signature, including parameters and return type.
// 3. Parses the function body block.
//
// Returns:
// - ast.Node: An AST node representing the parsed lambda function expression, including its parameters, return type, body, and location.
func parseLambdaFunction(p *Parser) ast.Node {
	//annonymous function
	start := p.eat().Start // eat fn token

	params, returnType := parseFunctionSignature(p)

	block := parseBlock(p)

	return ast.FunctionLiteral{
		Params:     params,
		ReturnType: returnType,
		Body:       block,
		Location: ast.Location{
			Start: start,
			End:   block.End,
		},
	}
}

// parseFunctionDeclStmt parses a function declaration statement from the input
// and returns an AST node representing the function declaration.
//
// The function expects the parser to be positioned at the start of the function
// declaration (i.e., the 'fn' token). It advances the parser, consumes the
// function name, parses the function signature (parameters and return type),
// and then parses the function body block.
//
// Parameters:
//   - p: A pointer to the Parser instance.
//
// Returns:
//   - ast.Node: An AST node representing the function declaration, which includes
//     the function's identifier, parameters, return type, and body block.
func parseFunctionDeclStmt(p *Parser) ast.Node {

	start := p.eat().Start // eat fn token

	nameToken := p.expect(lexer.IDENTIFIER_TOKEN)

	params, returnType := parseFunctionSignature(p)

	block := parseBlock(p)

	return ast.FunctionDeclStmt{
		Identifier: ast.IdentifierExpr{
			Name: nameToken.Value,
			Location: ast.Location{
				Start: nameToken.Start,
				End:   nameToken.End,
			},
		},
		FunctionLiteral: ast.FunctionLiteral{
			Params:     params,
			ReturnType: returnType,
			Body:       block,
			Location: ast.Location{
				Start: start,
				End:   block.End,
			},
		},
	}
}

// parseFunctionSignature parses the signature of a function, including its parameters and return type.
// It expects the function signature to start with an opening parenthesis and end with a closing parenthesis.
// Parameters can be either required or optional, and optional parameters must have a default value.
// The return type is optional and, if present, is indicated by an arrow (->) followed by the type.
//
// Parameters:
// - p (*Parser): The parser instance used to parse the function signature.
//
// Returns:
// - ([]ast.FunctionParam): A slice of FunctionParam representing the parameters of the function.
// - (ast.DataType): The return type of the function.
func parseFunctionSignature(p *Parser) ([]ast.FunctionParam, ast.DataType) {
	p.expect(lexer.OPEN_PAREN)

	//parse params
	var params []ast.FunctionParam

	for p.hasToken() && p.currentTokenKind() != lexer.CLOSE_PAREN {
		paramToken := p.expect(lexer.IDENTIFIER_TOKEN)
		param := ast.IdentifierExpr{
			Name: paramToken.Value,
			Location: ast.Location{
				Start: paramToken.Start,
				End:   paramToken.End,
			},
		}
		// if : then the param is not optional
		// if ?: then the param is optional
		currentToken := p.currentToken()

		if currentToken.Kind != lexer.COLON_TOKEN {
			report.Add(p.FilePath, currentToken.Start.Line, currentToken.End.Line, currentToken.Start.Column, currentToken.End.Column, "expected : ").Level(report.SYNTAX_ERROR)
		}

		p.eat()

		paramType := parseType(p, DEFAULT_BP)

		var defaultValue ast.Node

		params = append(params, ast.FunctionParam{
			Identifier:   param,
			Type:         paramType,
			DefaultValue: defaultValue,
			Location: ast.Location{
				Start: param.Start,
				End:   paramType.EndPos(),
			},
		})

		if p.currentTokenKind() != lexer.CLOSE_PAREN {
			p.expect(lexer.COMMA_TOKEN)
		}
	}

	p.expect(lexer.CLOSE_PAREN)

	var returnType ast.DataType

	//parse return type which is optional
	if p.currentTokenKind() != lexer.OPEN_CURLY {
		p.expect(lexer.ARROW_TOKEN)
		returnType = parseType(p, DEFAULT_BP)
	}

	return params, returnType
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

	p.eat() //eat the open paren
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
