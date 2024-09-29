package parser

import (
	"fmt"
	"walrus/ast"
	"walrus/builtins"
	"walrus/errgen"
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
	typeNUD(lexer.IDENTIFIER_TOKEN, parseBuiltinType)
	typeNUD(lexer.OPEN_BRACKET, parseArrayType)
}

func parseBuiltinType(p *Parser) ast.DataType {

	identifier := p.advance()

	switch identifier.Kind {
	case lexer.IDENTIFIER_TOKEN:
		break
	default:
		errgen.MakeError(p.FilePath, identifier.Start.Line, identifier.End.Line, identifier.Start.Column, identifier.End.Column, "invalid data type").Display()
	}

	value := identifier.Value

	loc := ast.Location{
		Start: identifier.Start,
		End:   identifier.End,
	}

	switch v := value; lexer.TOKEN_KIND(v) {
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
		return ast.UserDefinedType{
			TypeName: ast.DATA_TYPE(v),
			Location: loc,
		}
	}
}

func parseArrayType(p *Parser) ast.DataType {

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
		err := errgen.MakeError(p.FilePath, p.currentToken().Start.Line, p.currentToken().End.Line,  p.currentToken().Start.Column, p.currentToken().End.Column, fmt.Sprintf("Unexpected token %s\n", tokenKind))
		err.AddHint("Follow ", errgen.TEXT_HINT)
		err.AddHint("let x := 10", errgen.CODE_HINT)
		err.AddHint(" syntax or", errgen.TEXT_HINT)
		err.AddHint("Use primitive types like ", errgen.TEXT_HINT)
		err.AddHint("int, float, bool, char, str", errgen.CODE_HINT)
		err.AddHint(" or arrays of them", errgen.TEXT_HINT)
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

/*
Used to parse type for the type declaration with type keyword

Example:
type MyType struct {
	x: int,
	y: float,
};
*/
func parseUDTType(p *Parser) ast.DataType {

	identifier := p.currentToken()

	switch v := identifier.Value; lexer.TOKEN_KIND(v) {
	case builtins.STRUCT:
		p.advance()
		props := map[string]ast.StructPropType{}

		start := p.expect(lexer.OPEN_CURLY).Start

		for p.hasToken() && p.currentTokenKind() != lexer.CLOSE_CURLY {

			isPrivate := false

			if p.currentTokenKind() == lexer.PRIVATE_TOKEN {
				isPrivate = true
				p.advance()
			}

			iden := p.expect(lexer.IDENTIFIER_TOKEN)
			idenExpr := ast.IdentifierExpr{
				Name: iden.Value,
				Location: ast.Location{
					Start: iden.Start,
					End:   iden.End,
				},
			}
			p.expect(lexer.COLON_TOKEN)
			typeName := parseType(p, DEFAULT_BP)

			props[iden.Value] = ast.StructPropType{
				Prop:      idenExpr,
				PropType:  typeName,
				IsPrivate: isPrivate,
			}

			if p.currentTokenKind() != lexer.CLOSE_CURLY {
				p.expect(lexer.COMMA_TOKEN)
			}
		}

		end := p.expect(lexer.CLOSE_CURLY).End

		loc := ast.Location{
			Start: start,
			End:   end,
		}

		if len(props) == 0 {
			errgen.MakeError(p.FilePath, identifier.Start.Line, identifier.End.Line,  identifier.Start.Column, identifier.End.Column, "struct is empty").Display()
		}

		return ast.StructType{
			TypeName:   ast.DATA_TYPE(builtins.STRUCT),
			Properties: props,
			Location:   loc,
		}

	default:
		return parseType(p, DEFAULT_BP)
	}
}
