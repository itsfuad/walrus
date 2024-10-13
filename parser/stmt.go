package parser

import (
	"errors"
	"fmt"
	"strings"
	"walrus/ast"
	"walrus/errgen"
	"walrus/lexer"
)

// parseVarDeclStmt parses a variable declaration statement in the source code.
// It handles both `let` and `const` declarations, with optional type annotations
// and initial values.
//
// Parameters:
// - p: A pointer to the Parser instance.
//
// Returns:
// - An ast.Node representing the parsed variable declaration statement.
//
// The function performs the following steps:
// 1. Advances the parser to consume the `let` or `const` keyword.
// 2. Determines if the declaration is a constant.
// 3. Expects and consumes an identifier token for the variable name.
// 4. Optionally parses an explicit type if a colon `:` is present.
// 5. Parses the assignment operator `:=` or `=` and the initial value expression if present.
// 6. Ensures that constants have an initial value.
// 7. Expects and consumes a semicolon `;` to terminate the statement.
// 8. Constructs and returns an ast.VarDeclStmt node with the parsed information.
func parseVarDeclStmt(p *Parser) ast.Node {

	declToken := p.advance() // advance the let/const keyword

	// is it let or const?
	isConst := declToken.Kind == lexer.CONST_TOKEN

	// parse the variable name
	identifier := p.expect(lexer.IDENTIFIER_TOKEN)

	// parse the explicit type if present. This will be nil if no type is specified.
	var explicitType ast.DataType

	var value ast.Node

	assignmentToken := p.advance()

	if assignmentToken.Kind == lexer.COLON_TOKEN {
		// syntax is let a : <type>
		explicitType = parseType(p, DEFAULT_BP)
	} else if assignmentToken.Kind != lexer.WALRUS_TOKEN {
		msg := "Invalid variable declaration syntax"
		errgen.MakeError(p.FilePath, assignmentToken.Start.Line, assignmentToken.End.Line, assignmentToken.Start.Column, assignmentToken.End.Column, msg).AddHint("Maybe you want to use : or := instead of =", errgen.TEXT_HINT).Display()
	}

	if p.currentTokenKind() != lexer.SEMI_COLON_TOKEN {
		// then we have an assignment
		if assignmentToken.Kind == lexer.COLON_TOKEN {
			p.expect(lexer.EQUALS_TOKEN)
		}
		value = parseExpr(p, ASSIGNMENT_BP)
	}

	//if const, we must have a value
	if isConst && value == nil {
		msg := "constants must have value when declared"
		errgen.MakeError(p.FilePath, p.currentToken().Start.Line, p.currentToken().End.Line, p.currentToken().Start.Column, p.currentToken().End.Column, msg).Display()
	}

	end := p.expectSemicolon().End

	node := ast.VarDeclStmt{
		Variable: ast.IdentifierExpr{
			Name: identifier.Value,
			Location: ast.Location{
				Start: identifier.Start,
				End:   identifier.End,
			},
		},
		Value:        value,
		ExplicitType: explicitType,
		IsConst:      isConst,
		IsAssigned:   value != nil,
		Location: ast.Location{
			Start: declToken.Start,
			End:   end,
		},
	}

	return node
}


// parseUserDefinedTypeStmt parses a user-defined type statement in the source code.
// It expects the 'type' keyword followed by an identifier that starts with a capital letter.
// If the identifier does not start with a capital letter, an error is generated with a hint.
// The function then parses the user-defined type and expects a semicolon at the end.
// It returns an AST node representing the type declaration statement.
//
// Parameters:
// - p: A pointer to the Parser instance.
//
// Returns:
// - ast.Node: An AST node representing the type declaration statement.
func parseUserDefinedTypeStmt(p *Parser) ast.Node {

	start := p.advance().Start //eat type token

	typeName := p.expect(lexer.IDENTIFIER_TOKEN)

	if strings.ToUpper(typeName.Value[:1]) != typeName.Value[:1] {
		errgen.MakeError(p.FilePath, typeName.Start.Line, typeName.End.Line, typeName.Start.Column, typeName.End.Column, "user defined types should start with capital letter").AddHint(fmt.Sprintf("type %s%s [your type]", strings.ToUpper(typeName.Value[:1]), typeName.Value[1:]), errgen.TEXT_HINT).Display()
	}

	udType := parseUDTType(p)

	p.expectSemicolon()

	return ast.TypeDeclStmt{
		UDType:     udType,
		UDTypeName: typeName.Value,
		Location: ast.Location{
			Start: start,
			End:   udType.EndPos(),
		},
	}
}

// parseBlock parses a block statement from the input tokens.
// It expects the block to start with an opening curly brace '{' and end with a closing curly brace '}'.
// The function iterates over the tokens within the braces, parsing each node and adding it to the block's body.
// It returns an ast.BlockStmt containing the parsed nodes and their location in the source code.
//
// Parameters:
//   - p: A pointer to the Parser instance.
//
// Returns:
//   - ast.BlockStmt: The parsed block statement, including its contents and location.
func parseBlock(p *Parser) ast.BlockStmt {

	start := p.expect(lexer.OPEN_CURLY).Start

	body := make([]ast.Node, 0)

	for p.hasToken() && p.currentTokenKind() != lexer.CLOSE_CURLY {
		body = append(body, parseNode(p))
	}

	end := p.expect(lexer.CLOSE_CURLY).End

	return ast.BlockStmt{
		Contents: body,
		Location: ast.Location{
			Start: start,
			End:   end,
		},
	}
}

// parseForStmt parses a 'for' statement in the source code.
// It handles different types of 'for' loops:
// - Infinite loop: `for { }`
// - Condition-only loop: `for i < 10 { }`
// - Traditional for loop with initialization, condition, and increment: `for i := 0; i < 10; i++ { }`
// - For-each loop: `for i in arr { }`
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
		token := p.currentToken()
		idententifier := ast.IdentifierExpr{
			Name: token.Value,
			Location: ast.Location{
				Start: token.Start,
				End:   token.End,
			},
		}

		var startExpr ast.Node
		var condition ast.Node
		var increment ast.Node
		var iterable ast.Node

		switch p.currentTokenKind() {
		case lexer.WALRUS_TOKEN:
			// for i := 0; i < 10; i++ { }
			p.advance() // eat variable
			//parse start
			p.advance() // eat walrus token
			startExpr = parseExpr(p, ASSIGNMENT_BP)
			p.expect(lexer.SEMI_COLON_TOKEN)
			//parse condition
			condition = parseExpr(p, ASSIGNMENT_BP)
			p.expect(lexer.SEMI_COLON_TOKEN)
			//parse increment
			increment = parseExpr(p, ASSIGNMENT_BP)
		case lexer.IN_TOKEN:
			// for i in arr { }
			p.advance() // eat i token
			p.advance() // eat in token
			iterable = parseExpr(p, ASSIGNMENT_BP)
		default:
			// for i < 10 { } // condition only
			condition = parseExpr(p, ASSIGNMENT_BP)
		}

		block := parseBlock(p)

		return ast.ForStmt{
			Start:     idententifier,
			StartExpr: startExpr,
			Condition: condition,
			Increment: increment,
			Iterable:  iterable,
			Block:     block,
			Location: ast.Location{
				Start: start,
				End:   block.End,
			},
		}
	} else {
		//error
		//msg := "invalid for loop syntax"
		errgen.MakeError(p.FilePath, start.Line, start.Line, start.Column, start.Column, "invalid for loop syntax").Display()
		return nil
	}
}

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

	start := p.advance().Start // eat if token

	condition := parseExpr(p, ASSIGNMENT_BP)

	//parse block
	consequentBlock := parseBlock(p)

	var alternate ast.Node

	if p.hasToken() && p.currentTokenKind() == lexer.ELSE_TOKEN {
		p.advance() // eat else token
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
	start := p.advance().Start // eat fn token

	params, returnType := parseFunctionSignature(p)

	block := parseBlock(p)

	return ast.FunctionExpr{
		Params:     params,
		ReturnType: returnType,
		Body:      	block,
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

	start := p.advance().Start // eat fn token

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
		FunctionExpr: ast.FunctionExpr{
			Params:     params,
			ReturnType: returnType,
			Body:      block,
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

		if currentToken.Kind != lexer.COLON_TOKEN && currentToken.Kind != lexer.OPTIONAL_TOKEN {
			errgen.MakeError(p.FilePath, currentToken.Start.Line, currentToken.End.Line, currentToken.Start.Column, currentToken.End.Column, "expected : or ?:").Display()
		}

		p.advance()

		isOptional := currentToken.Kind == lexer.OPTIONAL_TOKEN

		paramType := parseType(p, DEFAULT_BP)

		var defaultValue ast.Node

		if isOptional {
			p.expectError(lexer.EQUALS_TOKEN, errors.New("expected default value for optional parameter. eg: param?: int = 10"))
			defaultValue = parseExpr(p, ASSIGNMENT_BP)
		}

		params = append(params, ast.FunctionParam{
			Identifier:   param,
			Type:         paramType,
			IsOptional:   isOptional,
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

// parseReturnStmt parses a return statement in the source code.
// It expects the current token to be a return token and advances the parser.
// If the next token is not a semicolon, it parses an expression for the return value.
// Finally, it expects a semicolon to end the return statement and returns an ast.ReturnStmt node.
//
// Parameters:
//   - p: A pointer to the Parser instance.
//
// Returns:
//   - An ast.Node representing the return statement.
func parseReturnStmt(p *Parser) ast.Node {

	start := p.advance().Start // eat return token

	var value ast.Node

	if p.currentTokenKind() != lexer.SEMI_COLON_TOKEN {
		value = parseExpr(p, ASSIGNMENT_BP)
	}

	end := p.expectSemicolon().End

	return ast.ReturnStmt{
		Value: value,
		Location: ast.Location{
			Start: start,
			End:   end,
		},
	}
}

// parseImplementStmt parses an implementation statement in the source code.
// The syntax for the implementation statement can be one of the following:
// - impl A, B, C for T { ... }
// - impl A for T { ... }
// - impl T { ... }
//
// It expects the parser to be positioned at the start of the 'impl' keyword.
//
// Parameters:
// - p: A pointer to the Parser instance.
//
// Returns:
// - An ast.Node representing the parsed implementation statement.
func parseImplementStmt(p *Parser) ast.Node {

	start := p.advance().Start // eat implement token

	// syntax: impl A, B, C for T { ... } or impl A for T { ... } or impl T { ... }
	var implFor ast.IdentifierExpr
	var trait ast.IdentifierExpr

	//parse the trait
	implForIdentifier := p.expect(lexer.IDENTIFIER_TOKEN)
	implFor = ast.IdentifierExpr{
		Name: implForIdentifier.Value,
		Location: ast.Location{
			Start: implForIdentifier.Start,
			End:   implForIdentifier.End,
		},
	}

	if p.currentTokenKind() != lexer.OPEN_BRACKET {
		p.expect(lexer.FOR_TOKEN)
		identifier := p.expect(lexer.IDENTIFIER_TOKEN)
		trait = implFor
		implFor = ast.IdentifierExpr{
			Name: identifier.Value,
			Location: ast.Location{
				Start: identifier.Start,
				End:   identifier.End,
			},
		}
	}

	//parse the type
	p.expect(lexer.OPEN_CURLY)

	methods := make(map[string]ast.FunctionDeclStmt, 0)

	for p.hasToken() && p.currentTokenKind() != lexer.CLOSE_CURLY {
		method := parseFunctionDeclStmt(p).(ast.FunctionDeclStmt)
		methods[method.Identifier.Name] = method
	}

	end := p.expect(lexer.CLOSE_CURLY).End

	return ast.ImplStmt{
		ImplFor: implFor,
		Trait:   trait,
		Methods: methods,
		Location: ast.Location{
			Start: start,
			End:   end,
		},
	}
}

// parseTraitDeclStmt parses a trait declaration statement from the provided parser.
// It expects the following structure:
// 
// trait <identifier> {
//     function <method_name>(<parameters>) -> <return_type>;
//     ...
// }
// 
// The function returns an ast.Node representing the trait declaration statement.
//
// Parameters:
// - p: A pointer to the Parser instance.
//
// Returns:
// - ast.Node: The parsed trait declaration statement node.
func parseTraitDeclStmt(p *Parser) ast.Node {

	start := p.advance().Start // eat trait token

	trait := p.expect(lexer.IDENTIFIER_TOKEN)

	p.expect(lexer.OPEN_CURLY)

	methods := make(map[string]ast.TraitMethod)

	for p.hasToken() && p.currentTokenKind() != lexer.CLOSE_CURLY {

		start := p.expect(lexer.FUNCTION).Start

		if p.currentTokenKind() != lexer.IDENTIFIER_TOKEN {
			errgen.MakeError(p.FilePath, p.currentToken().Start.Line, p.currentToken().End.Line, p.currentToken().Start.Column, p.currentToken().End.Column, "expected method name").Display()
		}

		name := p.expect(lexer.IDENTIFIER_TOKEN)

		dataType, params, returnType := getFunctionTypeSignature(p)

		if _, ok := methods[name.Value]; ok {
			msg := fmt.Sprintf("method %s already defined", name.Value)
			errgen.MakeError(p.FilePath, name.Start.Line, name.End.Line, name.Start.Column, name.End.Column, msg).Display()
		}

		methods[name.Value] = ast.TraitMethod{
			Identifier: ast.IdentifierExpr{
				Name: name.Value,
				Location: ast.Location{
					Start: name.Start,
					End:   name.End,
				},
			},
			FunctionType: ast.FunctionType{
				TypeName:   dataType,
				Parameters: params,
				ReturnType: returnType,
				Location: ast.Location{
					Start: start,
					End:   returnType.EndPos(),
				},
			},
		}

		if p.currentTokenKind() != lexer.CLOSE_CURLY {
			p.expect(lexer.SEMI_COLON_TOKEN)
		}
	}

	end := p.expect(lexer.CLOSE_CURLY).End

	return ast.TraitDeclStmt{
		Trait: ast.IdentifierExpr{
			Name: trait.Value,
			Location: ast.Location{
				Start: trait.Start,
				End:   trait.End,
			},
		},
		Methods: methods,
		Location: ast.Location{
			Start: start,
			End:   end,
		},
	}
}