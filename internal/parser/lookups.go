package parser

import (
	"walrus/internal/ast"
	"walrus/internal/builtins"
	"walrus/internal/lexer"
)

type BINDING_POWER int

const (
	DEFAULT_BP BINDING_POWER = iota
	ASSIGNMENT_BP
	CASTING_BP
	LOGICAL_BP
	RELATIONAL_BP
	ADDITIVE_BP
	MULTIPLICATIVE_BP
	UNARY_BP
	CALL_BP
	PRIMARY_BP
	MEMBER_BP
)

type NUDHandler func(p *Parser) ast.Node
type STMTHandler func(p *Parser) ast.Node
type LEDHandler func(p *Parser, left ast.Node, bp BINDING_POWER) ast.Node

var NUDLookup = map[builtins.TOKEN_KIND]NUDHandler{}
var STMTLookup = map[builtins.TOKEN_KIND]STMTHandler{}
var LEDLookup = map[builtins.TOKEN_KIND]LEDHandler{}
var BPLookup = map[builtins.TOKEN_KIND]BINDING_POWER{}

func GetBP(token builtins.TOKEN_KIND) BINDING_POWER {
	if bp, ok := BPLookup[token]; ok {
		return bp
	} else {
		return DEFAULT_BP
	}
}

func led(tokenKind builtins.TOKEN_KIND, bp BINDING_POWER, handler LEDHandler) {
	BPLookup[tokenKind] = bp
	LEDLookup[tokenKind] = handler
}

func nud(tokenKind builtins.TOKEN_KIND, handler NUDHandler) {
	NUDLookup[tokenKind] = handler
}

func stmt(tokenKind builtins.TOKEN_KIND, handler STMTHandler) {
	STMTLookup[tokenKind] = handler
}

func bindLookupHandlers() {

	//Assignment
	led(lexer.EQUALS_TOKEN, ASSIGNMENT_BP, parseVarAssignmentExpr)
	led(lexer.PLUS_EQUALS_TOKEN, ASSIGNMENT_BP, parseVarAssignmentExpr)
	led(lexer.MINUS_EQUALS_TOKEN, ASSIGNMENT_BP, parseVarAssignmentExpr)
	led(lexer.MUL_EQUALS_TOKEN, ASSIGNMENT_BP, parseVarAssignmentExpr)
	led(lexer.DIV_EQUALS_TOKEN, ASSIGNMENT_BP, parseVarAssignmentExpr)
	led(lexer.MOD_EQUALS_TOKEN, ASSIGNMENT_BP, parseVarAssignmentExpr)
	led(lexer.EXP_EQUALS_TOKEN, ASSIGNMENT_BP, parseVarAssignmentExpr)

	led(lexer.OPEN_BRACKET, MEMBER_BP, parseIndexable)

	led(lexer.DOT_TOKEN, MEMBER_BP, parsePropertyExpr)
	led(lexer.OPEN_PAREN, CALL_BP, parseCallExpr)

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

	led(lexer.AS_TOKEN, CASTING_BP, parseTypeCastExpr)

	//Postfix
	led(lexer.PLUS_PLUS_TOKEN, UNARY_BP, parsePostfixExpr)   // a++
	led(lexer.MINUS_MINUS_TOKEN, UNARY_BP, parsePostfixExpr) // a--

	nud(lexer.IDENTIFIER_TOKEN, parsePrimaryExpr)  // identifier
	nud(lexer.INT8_TOKEN, parsePrimaryExpr)        // int literal, 8 bit
	nud(lexer.INT16_TOKEN, parsePrimaryExpr)       // int literal, 16 bit
	nud(lexer.INT32_TOKEN, parsePrimaryExpr)       // int literal, 32 bit
	nud(lexer.INT64_TOKEN, parsePrimaryExpr)       // int literal, 64 bit
	nud(lexer.FLOAT32_TOKEN, parsePrimaryExpr)     // float literal
	nud(lexer.FLOAT64_TOKEN, parsePrimaryExpr)     // float literal, 64 bit
	nud(lexer.UINT8_TOKEN, parsePrimaryExpr)       // uint literal, 8 bit
	nud(lexer.UINT16_TOKEN, parsePrimaryExpr)      // uint literal, 16 bit
	nud(lexer.UINT32_TOKEN, parsePrimaryExpr)      // uint literal, 32 bit
	nud(lexer.UINT64_TOKEN, parsePrimaryExpr)      // uint literal, 64 bit
	nud(lexer.STR_TOKEN, parsePrimaryExpr)         // string literal
	nud(lexer.OPEN_BRACKET, parseArrayExpr)        // array literal [1,2,3]
	nud(lexer.OPEN_PAREN, parseGroupingExpr)       // grouping expression a + (b+c)
	nud(lexer.FUNCTION_TOKEN, parseLambdaFunction) // anonymous function
	nud(lexer.STRUCT_TOKEN, parseStructLiteral)    // anonymous struct

	//nud(lexer.AT_TOKEN, parseStructLiteral)
	nud(lexer.MAP_TOKEN, parseMapLiteral)

	//Unary
	nud(lexer.MINUS_TOKEN, parseUnaryExpr) // unary minus : -a
	nud(lexer.NOT_TOKEN, parseUnaryExpr)   // unary not : !a
	//Increment and Decrement
	//Prefix
	nud(lexer.PLUS_PLUS_TOKEN, parsePrefixExpr)   // ++a
	nud(lexer.MINUS_MINUS_TOKEN, parsePrefixExpr) // --a

	//Statements
	stmt(lexer.LET_TOKEN, parseVarDeclStmt)       // variable declaration
	stmt(lexer.CONST_TOKEN, parseVarDeclStmt)     // constant declaration
	stmt(lexer.TYPE_TOKEN, parseUserDefinedTypes) // user defined type

	stmt(lexer.IF_TOKEN, parseIfStmt)                 // if statement
	stmt(lexer.FOR_TOKEN, parseForStmt)               // for statement
	stmt(lexer.FOREACH_TOKEN, parseForStmt)           // foreach statement
	stmt(lexer.FUNCTION_TOKEN, parseFunctionDeclStmt) // function declaration
	stmt(lexer.RETURN_TOKEN, parseReturnStmt)         // return statement
	stmt(lexer.IMPL_TOKEN, parseImplStmt)
}
