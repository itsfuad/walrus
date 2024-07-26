package parser

import (
	"walrus/ast"
	"walrus/lexer"
)

type BINDING_POWER int

const (
	DEFAULT_BP BINDING_POWER = iota
	COMMA_BP
	ASSIGNMENT_BP
	LOGICAL_BP
	RELATIONAL_BP
	ADDITIVE_BP
	MULTIPLICATIVE_BP
	UNARY_BP
	CALL_BP
	MEMBER_BP
	PRIMARY_BP
)

type NUDHandler func(p *Parser) ast.Node
type STMTHandler func(p *Parser) ast.Node
type LEDHandler func(p *Parser, left ast.Node, bp BINDING_POWER) ast.Node

var NUDLookup = map[lexer.TOKEN_KIND]NUDHandler{}
var STMTLookup = map[lexer.TOKEN_KIND]STMTHandler{}
var LEDLookup = map[lexer.TOKEN_KIND]LEDHandler{}
var BPLookup = map[lexer.TOKEN_KIND]BINDING_POWER{}

func GetBP(token lexer.TOKEN_KIND) BINDING_POWER {
	if bp, ok := BPLookup[token]; ok {
		return bp
	} else {
		return DEFAULT_BP
	}
}

func led(tokenKind lexer.TOKEN_KIND, bp BINDING_POWER, handler LEDHandler) {
	BPLookup[tokenKind] = bp
	LEDLookup[tokenKind] = handler
}

func nud(tokenKind lexer.TOKEN_KIND, handler NUDHandler) {
	NUDLookup[tokenKind] = handler
}

func stmt(tokenKind lexer.TOKEN_KIND, handler STMTHandler) {
	STMTLookup[tokenKind] = handler
}

func bindLookupHandlers() {

	led(lexer.EQUALS_TOKEN, ASSIGNMENT_BP, parseVarAssignmentExpr)
	led(lexer.OPEN_BRACKET, MEMBER_BP, parseArrayAccess)
	nud(lexer.AT_TOKEN, parseStructLiteral)

	led(lexer.DOT_TOKEN, MEMBER_BP, parsePropertyExpr)

	//arithmetics
	led(lexer.PLUS_TOKEN, ADDITIVE_BP, parseBinaryExpr)
	led(lexer.MINUS_TOKEN, ADDITIVE_BP, parseBinaryExpr)
	led(lexer.MUL_TOKEN, MULTIPLICATIVE_BP, parseBinaryExpr)
	led(lexer.DIV_TOKEN, MULTIPLICATIVE_BP, parseBinaryExpr)
	led(lexer.MOD_TOKEN, MULTIPLICATIVE_BP, parseBinaryExpr)
	led(lexer.EXP_TOKEN, MULTIPLICATIVE_BP, parseBinaryExpr)

	led(lexer.DOUBLE_EQUAL_TOKEN, RELATIONAL_BP, parseBinaryExpr)
	led(lexer.NOT_EQUAL_TOKEN, RELATIONAL_BP, parseBinaryExpr)
	led(lexer.LESS_EQUAL_TOKEN, RELATIONAL_BP, parseBinaryExpr)
	led(lexer.LESS_TOKEN, RELATIONAL_BP, parseBinaryExpr)
	led(lexer.GREATER_EQUAL_TOKEN, RELATIONAL_BP, parseBinaryExpr)
	led(lexer.GREATER_TOKEN, RELATIONAL_BP, parseBinaryExpr)

	nud(lexer.IDENTIFIER_TOKEN, parsePrimaryExpr)
	nud(lexer.INT, parsePrimaryExpr)
	nud(lexer.FLOAT, parsePrimaryExpr)
	nud(lexer.STR, parsePrimaryExpr)
	nud(lexer.CHR, parsePrimaryExpr)
	nud(lexer.NULL, parsePrimaryExpr)
	nud(lexer.OPEN_BRACKET, parseArrayExpr)
	nud(lexer.OPEN_PAREN, parseGroupingExpr)

	stmt(lexer.LET_TOKEN, parseVarDeclStmt)
	stmt(lexer.CONST_TOKEN, parseVarDeclStmt)
	stmt(lexer.TYPE_TOKEN, parseUserDefinedTypeStmt)

	stmt(lexer.IF_TOKEN, parseIfStmt)

	//Unary
	nud(lexer.MINUS_TOKEN, parseUnaryExpr)
	nud(lexer.NOT_TOKEN, parseUnaryExpr)
}
