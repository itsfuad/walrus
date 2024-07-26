package parser

import (
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
		errgen.MakeError(p.FilePath, assToken.Start.Line, assToken.Start.Column, assToken.End.Column, msg).Display()
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
		errgen.MakeError(p.FilePath, p.currentToken().Start.Line, p.currentToken().Start.Column, p.currentToken().End.Column, msg).Display()
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
		errgen.MakeError(p.FilePath, typeName.Start.Line, typeName.Start.Column, typeName.End.Column, "user defined types should start with capital letter").AddHint(fmt.Sprintf("type %s%s [your type]", strings.ToUpper(typeName.Value[:1]), typeName.Value[1:]), errgen.TEXT_HINT).Display()
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

func parseIfStmt(p *Parser) ast.Node {

	start := p.advance().Start // eat if token

	condition := parseExpr(p, ASSIGNMENT_BP)

	//parse block
	consequentBlock := parseBlock(p)

	var alternate ast.Node

	if p.hasToken() && p.currentTokenKind() == lexer.ELSE_TOKEN {
		p.advance() // eat else token
		alternate = parseBlock(p)
	} else if p.hasToken() && p.currentTokenKind() == lexer.ELSEIF_TOKEN {
		alternate = parseIfStmt(p)
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
