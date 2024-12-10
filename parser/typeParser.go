package parser

import (
	"errors"
	"fmt"
	"walrus/ast"
	"walrus/builtins"
	"walrus/errgen"
	"walrus/lexer"
)

type typeNUDHandler func(p *Parser) ast.DataType
type typeLEDHandler func(p *Parser, left ast.DataType, bp BINDING_POWER) ast.DataType

var bpTypeLookups = map[builtins.TOKEN_KIND]BINDING_POWER{}
var typeNUDLookup = map[builtins.TOKEN_KIND]typeNUDHandler{}
var typeLEDLookup = map[builtins.TOKEN_KIND]typeLEDHandler{}

func typeNUD(kind builtins.TOKEN_KIND, handler typeNUDHandler) {
	typeNUDLookup[kind] = handler
}

func bindTypeLookups() {
	typeNUD(lexer.IDENTIFIER_TOKEN, parseDataType)
	typeNUD(lexer.OPEN_BRACKET, parseArrayType)
	typeNUD(lexer.FUNCTION_TOKEN, parseFunctionType)
	typeNUD(lexer.MAP_TOKEN, parseMapType)
}

func parseMapType(p *Parser) ast.DataType {

	// map[<keyType>]<valueType>
	// or
	// type UserDefinedType map[<keyType>]<valueType>
	// UserDefinedType

	var mapToken lexer.Token

	if p.currentTokenKind() == lexer.MAP_TOKEN {
		mapToken = p.advance()
	} else {
		//we expect an identifier
		mapToken = p.expectError(lexer.IDENTIFIER_TOKEN, errors.New("expected 'map' keyword or the map type"))
		return ast.MapType{
			TypeName: builtins.PARSER_TYPE(builtins.MAP),
			Map: ast.IdentifierExpr{
				Name: mapToken.Value,
				Location: ast.Location{
					Start: mapToken.Start,
					End:   mapToken.End,
				},
			},
			KeyType:   nil,
			ValueType: nil,
			Location: ast.Location{
				Start: mapToken.Start,
				End:   mapToken.End,
			},
		}
	}

	p.expect(lexer.OPEN_BRACKET)

	keyType := parseType(p, DEFAULT_BP)

	p.expect(lexer.CLOSE_BRACKET)

	valueType := parseType(p, DEFAULT_BP)

	return ast.MapType{
		TypeName: builtins.PARSER_TYPE(builtins.MAP),
		Map: ast.IdentifierExpr{
			Name: mapToken.Value,
			Location: ast.Location{
				Start: mapToken.Start,
				End:   mapToken.End,
			},
		},
		KeyType:   keyType,
		ValueType: valueType,
		Location: ast.Location{
			Start: mapToken.Start,
			End:   valueType.EndPos(),
		},
	}
}

func parseFunctionType(p *Parser) ast.DataType {

	start := p.advance().Start // eat function token

	typeName, params, returnType := getFunctionTypeSignature(p)

	loc := ast.Location{
		Start: start,
		End:   returnType.EndPos(),
	}

	return ast.FunctionType{
		TypeName:   typeName,
		Parameters: params,
		ReturnType: returnType,
		Location:   loc,
	}
}

//parseType is the entry point for parsing types
/*
Returns
 - ast.DataType : The parsed type
 - []ast.FunctionTypeParam : The parameters of the function type
 - ast.DataType : The return type of the function type
*/
func getFunctionTypeSignature(p *Parser) (builtins.PARSER_TYPE, []ast.FunctionTypeParam, ast.DataType) {
	p.expect(lexer.OPEN_PAREN)
	var params []ast.FunctionTypeParam
	for p.hasToken() && p.currentTokenKind() != lexer.CLOSE_PAREN {
		iden := p.expect(lexer.IDENTIFIER_TOKEN)
		//if exists, then it is a duplicate
		for _, param := range params {
			if param.Identifier.Name == iden.Value {
				errgen.AddError(p.FilePath, iden.Start.Line, iden.End.Line, iden.Start.Column, iden.End.Column, fmt.Sprintf("parameter '%s' already defined", iden.Value)).ErrorLevel(errgen.CRITICAL)
			}
		}

		curentToken := p.currentToken()

		if curentToken.Kind != lexer.COLON_TOKEN && curentToken.Kind != lexer.OPTIONAL_TOKEN {
			errgen.AddError(p.FilePath, curentToken.Start.Line, curentToken.End.Line, curentToken.Start.Column, curentToken.End.Column, "expected : or ?:").ErrorLevel(errgen.CRITICAL)
		}

		isOptional := p.advance().Kind == lexer.OPTIONAL_TOKEN

		typeName := parseType(p, DEFAULT_BP)

		params = append(params, ast.FunctionTypeParam{
			Identifier: ast.IdentifierExpr{
				Name: iden.Value,
				Location: ast.Location{
					Start: iden.Start,
					End:   iden.End,
				},
			},
			Type:       typeName,
			IsOptional: isOptional,
			Location: ast.Location{
				Start: iden.Start,
				End:   typeName.EndPos(),
			},
		})

		if p.currentTokenKind() != lexer.CLOSE_PAREN {
			p.expect(lexer.COMMA_TOKEN)
		}
	}

	p.expect(lexer.CLOSE_PAREN)

	var returnType ast.DataType

	if p.currentTokenKind() == lexer.ARROW_TOKEN {
		p.advance()
		returnType = parseType(p, DEFAULT_BP)
	} else {
		returnType = ast.VoidType{
			TypeName: builtins.PARSER_TYPE(builtins.VOID),
			Location: ast.Location{
				Start: p.currentToken().Start,
				End:   p.currentToken().End,
			},
		}
	}

	return builtins.PARSER_TYPE(builtins.FUNCTION), params, returnType
}

// Parses the builtin types like int, float, bool, char, str, null.
// If the type is not a builtin type, then it is a user defined type
// Type must be a single token identifier
func parseDataType(p *Parser) ast.DataType {

	identifier := p.advance()

	switch identifier.Kind {
	case lexer.IDENTIFIER_TOKEN:
		break
	default:
		errgen.AddError(p.FilePath, identifier.Start.Line, identifier.End.Line, identifier.Start.Column, identifier.End.Column, "invalid data type").ErrorLevel(errgen.CRITICAL)
	}

	value := identifier.Value

	loc := ast.Location{
		Start: identifier.Start,
		End:   identifier.End,
	}

	switch v := value; builtins.TOKEN_KIND(v) {
	case lexer.INT8_TOKEN, lexer.INT16_TOKEN, lexer.INT32_TOKEN, lexer.INT64_TOKEN, lexer.UINT8_TOKEN, lexer.UINT16_TOKEN, lexer.UINT32_TOKEN, lexer.UINT64_TOKEN:
		return ast.IntegerType{
			TypeName: builtins.PARSER_TYPE(v),
			BitSize:  builtins.GetBitSize(builtins.PARSER_TYPE(v)),
			IsSigned: builtins.IsSigned(builtins.PARSER_TYPE(v)),
			Location: loc,
		}
	case lexer.FLOAT32_TOKEN, lexer.FLOAT64_TOKEN:
		return ast.FloatType{
			TypeName: builtins.PARSER_TYPE(v),
			BitSize:  builtins.GetBitSize(builtins.PARSER_TYPE(v)),
			Location: loc,
		}
	case lexer.STR_TOKEN:
		return ast.StringType{
			TypeName: builtins.PARSER_TYPE(v),
			Location: loc,
		}
	case lexer.BOOL_TOKEN:
		return ast.BooleanType{
			TypeName: builtins.PARSER_TYPE(v),
			Location: loc,
		}
	case lexer.NULL_TOKEN:
		return ast.NullType{
			TypeName: builtins.PARSER_TYPE(v),
			Location: loc,
		}
	default:
		return ast.UserDefinedType{
			TypeName:  builtins.PARSER_TYPE(builtins.USER_DEFINED),
			AliasName: value,
			Location:  loc,
		}
	}
}

// parseArrayType parses an array type from the input and returns an ast.DataType
// representing the array type.
//
// The function expects the parser to be positioned at the opening bracket of the array type.
// It advances the parser, expects a closing bracket, and then parses the element type of the array.
//
// Parameters:
// - p: A pointer to the Parser instance.
//
// Returns:
// - ast.DataType: An instance of ast.ArrayType representing the parsed array type.
func parseArrayType(p *Parser) ast.DataType {

	p.advance()
	p.expect(lexer.CLOSE_BRACKET)

	elemType := parseType(p, DEFAULT_BP)

	return ast.ArrayType{
		TypeName:  builtins.PARSER_TYPE(builtins.ARRAY),
		ArrayType: elemType,
		Location: ast.Location{
			Start: elemType.StartPos(),
			End:   elemType.EndPos(),
		},
	}
}

// parseType parses a data type from the given parser instance, respecting the specified binding power.
// It first attempts to parse a null denotation (NUD) based on the current token kind.
// If no NUD handler is found for the token, it generates an error with hints and displays it.
// If a NUD handler is found, it proceeds to parse left denotations (LED) while the binding power of the current token kind is greater than the specified binding power.
// The function returns the parsed data type.
//
// Parameters:
// - p: A pointer to the Parser instance from which to parse the data type.
// - bp: The binding power to respect during parsing.
//
// Returns:
// - An ast.DataType representing the parsed data type, or nil if an error occurs.
func parseType(p *Parser, bp BINDING_POWER) ast.DataType {
	// Fist parse the NUD
	tokenKind := p.currentTokenKind()
	nudFunction, exists := typeNUDLookup[tokenKind]

	if !exists {
		//panic(fmt.Sprintf("TYPE NUD handler expected for token %s\n", tokenKind))
		errgen.AddError(p.FilePath, p.currentToken().Start.Line, p.currentToken().End.Line, p.currentToken().Start.Column, p.currentToken().End.Column, fmt.Sprintf("unexpected token '%s'\n", tokenKind)).ErrorLevel(errgen.CRITICAL)
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
	switch v := p.currentToken().Value; builtins.TOKEN_KIND(v) {
	case builtins.STRUCT:
		return parseStructType(p)
	case builtins.INTERFACE:
		return parseInterfaceType(p)
	default:
		return parseType(p, DEFAULT_BP)
	}
}

// parseStructType parses a struct type definition from the provided parser.
// It expects the parser to be positioned at the start of the struct definition.
//
// The function handles the following:
// - Parsing the struct identifier.
// - Parsing properties of the struct, including their types and visibility (public/private).
// - Parsing embedded structs.
//
// It returns an ast.DataType representing the parsed struct type.
//
// Parameters:
// - p: A pointer to the Parser instance.
//
// Returns:
// - ast.DataType: The parsed struct type.
//
// Errors:
// - If the struct is empty, an error is generated and displayed.
func parseStructType(p *Parser) ast.DataType {

	identifier := p.advance() // eat struct token

	props := make(map[string]ast.StructPropType)

	start := p.expect(lexer.OPEN_CURLY).Start

	for p.hasToken() && p.currentTokenKind() != lexer.CLOSE_CURLY {

		isPrivate := false

		if p.currentTokenKind() == lexer.PRIVATE_TOKEN {
			isPrivate = true
			p.advance()
		}

		iden := p.expect(lexer.IDENTIFIER_TOKEN)

		if _, ok := props[iden.Value]; ok {
			errgen.AddError(p.FilePath, iden.Start.Line, iden.End.Line, iden.Start.Column, iden.End.Column, fmt.Sprintf("property '%s' already defined", iden.Value)).ErrorLevel(errgen.CRITICAL)
		}

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
		errgen.AddError(p.FilePath, identifier.Start.Line, identifier.End.Line, identifier.Start.Column, identifier.End.Column, "struct is empty").ErrorLevel(errgen.CRITICAL)
	}

	return ast.StructType{
		TypeName:   builtins.PARSER_TYPE(builtins.STRUCT),
		Properties: props,
		Location:   loc,
	}
}

func parseInterfaceType(p *Parser) ast.DataType {

	start := p.advance().Start

	p.expect(lexer.OPEN_CURLY)

	methods := make(map[string]ast.InterfaceMethod)

	for p.hasToken() && p.currentTokenKind() != lexer.CLOSE_CURLY {

		start := p.expect(lexer.FUNCTION_TOKEN).Start

		if p.currentTokenKind() != lexer.IDENTIFIER_TOKEN {
			errgen.AddError(p.FilePath, p.currentToken().Start.Line, p.currentToken().End.Line, p.currentToken().Start.Column, p.currentToken().End.Column, "expected method name").ErrorLevel(errgen.CRITICAL)
		}

		name := p.expect(lexer.IDENTIFIER_TOKEN)

		dataType, params, returnType := getFunctionTypeSignature(p)

		if _, ok := methods[name.Value]; ok {
			msg := fmt.Sprintf("method %s already defined", name.Value)
			errgen.AddError(p.FilePath, name.Start.Line, name.End.Line, name.Start.Column, name.End.Column, msg).ErrorLevel(errgen.CRITICAL)
		}

		methods[name.Value] = ast.InterfaceMethod{
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

	return ast.InterfaceType{
		TypeName: builtins.PARSER_TYPE(builtins.INTERFACE),
		Methods:  methods,
		Location: ast.Location{
			Start: start,
			End:   end,
		},
	}
}
