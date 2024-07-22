package parser

import (
	"fmt"
	"walrus/ast"
	"walrus/builtins"
	"walrus/errors"
	"walrus/lexer"
)

type typeNUDHandler func(p *Parser) ast.DataType
type typeLEDHandler func(p *Parser, left ast.DataType, bp BINDING_POWER) ast.DataType

var bpTypeLookups = map[lexer.TOKEN_KIND]BINDING_POWER{}
var typeNUDLookup = map[lexer.TOKEN_KIND]typeNUDHandler{}
var typeLEDLookup = map[lexer.TOKEN_KIND]typeLEDHandler{}

func typeNUD(kind lexer.TOKEN_KIND, handler typeNUDHandler) {
	typeNUDLookup[kind] = handler
}

func bindTypeLookups() {
	typeNUD(lexer.IDENTIFIER_TOKEN, parseDataType)
	typeNUD(lexer.OPEN_BRACKET, parseArrayType)
}

func parseDataType(p *Parser) ast.DataType {
	identifier := p.expect(lexer.IDENTIFIER_TOKEN)
	value := identifier.Value

	loc := ast.Location{
		Start: identifier.Start,
		End:   identifier.End,
	}

	switch v := value ; lexer.TOKEN_KIND(v) {
	case lexer.INT:
		return ast.IntegerType{
			TypeName: ast.DATA_TYPE(v),
			Location: loc,
		}
	case lexer.FLOAT:
		return ast.FloatType{
			TypeName: ast.DATA_TYPE(v),
			Location: loc,
		}
	case lexer.STR:
		return ast.StringType{
			TypeName: ast.DATA_TYPE(v),
			Location: loc,
		}
	case lexer.CHR:
		return ast.CharType{
			TypeName: ast.DATA_TYPE(v),
			Location: loc,
		}
	case lexer.BOOL:
		return ast.BooleanType{
			TypeName: ast.DATA_TYPE(v),
			Location: loc,
		}
	case lexer.NULL:
		return ast.NullType{
			TypeName: ast.DATA_TYPE(v),
			Location: loc,
		}
	default:
		return ast.StructType{
			TypeName: ast.DATA_TYPE(builtins.STRUCT),
			Location: loc,
		}
	}
}

func parseArrayType(p *Parser) ast.DataType {

	fmt.Println("Parsing array type")

	p.advance()
	p.expect(lexer.CLOSE_BRACKET)

	elemType := parseType(p, DEFAULT_BP)

	return ast.ArrayType{
		TypeName:  ast.DATA_TYPE(builtins.ARRAY),
		ArrayType: elemType,
		Location: ast.Location{
			Start: elemType.StartPos(),
			End:   elemType.EndPos(),
		},
	}
}

func parseType(p *Parser, bp BINDING_POWER) ast.DataType {
	// Fist parse the NUD
	tokenKind := p.currentTokenKind()
	nudFunction, exists := typeNUDLookup[tokenKind]

	if !exists {
		//panic(fmt.Sprintf("TYPE NUD handler expected for token %s\n", tokenKind))
		err := errors.MakeError(p.FilePath, p.currentToken().Start.Line, p.currentToken().Start.Column, p.currentToken().End.Column, fmt.Sprintf("Unexpected token %s\n", tokenKind))
		err.AddHint("Follow ", errors.TEXT_HINT)
		err.AddHint("let x := 10", errors.CODE_HINT)
		err.AddHint(" syntax or", errors.TEXT_HINT)
		err.AddHint("Use primitive types like ", errors.TEXT_HINT)
		err.AddHint("int, float, bool, char, str", errors.CODE_HINT)
		err.AddHint(" or arrays of them", errors.TEXT_HINT)
		err.Display()
		return nil
	}

	left := nudFunction(p)

	for bpTypeLookups[p.currentTokenKind()] > bp {

		tokenKind := p.currentTokenKind()

		ledFunction, exists := typeLEDLookup[tokenKind]

		if !exists {
			panic(fmt.Sprintf("TYPE LED handler expected for token %s\n", tokenKind))
		}

		left = ledFunction(p, left, bpTypeLookups[p.currentTokenKind()])
	}

	return left
}
