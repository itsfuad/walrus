package parser

import (
	"walrus/ast"
	"walrus/errors"
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
		errors.MakeError(p.FilePath, assToken.Start.Line, assToken.Start.Column, assToken.End.Column, msg).Display()
	}

	if (p.currentTokenKind() != lexer.SEMI_COLON_TOKEN) {
		if assToken.Kind == lexer.COLON_TOKEN {
			p.expect(lexer.EQUALS_TOKEN)
		}
		value = parseExpr(p, ASSIGNMENT_BP)
	}

	//if const, we must have a value
	if isConst && value == nil {
		msg := "constants must have value when declared"
		errors.MakeError(p.FilePath, p.currentToken().Start.Line, p.currentToken().Start.Column, p.currentToken().End.Column, msg).Display()
	}

	end := p.expect(lexer.SEMI_COLON_TOKEN).End


	node := ast.VarDeclStmt{
		Variable: ast.IdentifierExpr{
			Name: identifier.Value,
			Location: ast.Location{
				Start: identifier.Start,
				End: identifier.End,
			},
		},
		Value: value,
		ExplicitType: explicitType,
		IsConst: isConst,
		IsAssigned: value != nil,
		Location: ast.Location{
			Start: declToken.Start,
			End: end,
		},
	}

	return node
}