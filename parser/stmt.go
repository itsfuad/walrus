package parser

import (
	"errors"
	"fmt"
	"strings"
	"walrus/ast"
	"walrus/errgen"
	"walrus/lexer"
)

func parseVarDeclStmt(p *Parser) ast.Node {

	declToken := p.advance() // advance the let/const keyword

	// is it let or const?
	isConst := declToken.Kind == lexer.CONST_TOKEN

	identifier := p.expect(lexer.IDENTIFIER_TOKEN)

	var explicitType ast.DataType

	var value ast.Node

	assToken := p.advance()

	if assToken.Kind == lexer.COLON_TOKEN {
		// syntax is let a : <type>
		explicitType = parseType(p, DEFAULT_BP)
	} else if assToken.Kind != lexer.WALRUS_TOKEN {
		msg := "expected : or :="
		errgen.MakeError(p.FilePath, assToken.Start.Line, assToken.End.Line, assToken.Start.Column, assToken.End.Column, msg).Display()
	}

	if p.currentTokenKind() != lexer.SEMI_COLON_TOKEN {
		if assToken.Kind == lexer.COLON_TOKEN {
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

func parseForStmt(p *Parser) ast.Node {
	start := p.advance().Start // eat for token
	// we can have `for i := 0; i < 10; i++ { }` or `for i < 10 { }` or `for { }` or `for i in arr { }`
	// get which type of for loop it is
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

		// p.expect(lexer.COLON_TOKEN)

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

func parseTraitDeclStmt(p *Parser) ast.Node {

	start := p.advance().Start // eat trait token

	trait := p.expect(lexer.IDENTIFIER_TOKEN)

	p.expect(lexer.OPEN_CURLY)

	methods := make(map[string]ast.TraitMethod, 0)

	for p.hasToken() && p.currentTokenKind() != lexer.CLOSE_CURLY {

		start := p.expect(lexer.FUNCTION).Start

		if p.currentTokenKind() != lexer.IDENTIFIER_TOKEN {
			errgen.MakeError(p.FilePath, p.currentToken().Start.Line, p.currentToken().End.Line, p.currentToken().Start.Column, p.currentToken().End.Column, "expected method name").Display()
		}

		name := p.expect(lexer.IDENTIFIER_TOKEN)

		dataType, params, returnType := getFunctionTypeSignature(p)

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
