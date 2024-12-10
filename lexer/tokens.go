package lexer

import (
	"fmt"
	"walrus/builtins"
	"walrus/utils"
)

const (
	//keywords
	LET_TOKEN        		builtins.TOKEN_KIND = "let"
	CONST_TOKEN      		builtins.TOKEN_KIND = "const"
	TYPE_TOKEN       		builtins.TOKEN_KIND = "type"
	IF_TOKEN         		builtins.TOKEN_KIND = "if"
	ELSE_TOKEN       		builtins.TOKEN_KIND = "else"
	FOR_TOKEN        		builtins.TOKEN_KIND = "for"
	FOREACH_TOKEN    		builtins.TOKEN_KIND = "foreach"
	IDENTIFIER_TOKEN 		builtins.TOKEN_KIND = "identifier"
	PRIVATE_TOKEN    		builtins.TOKEN_KIND = "priv"
	IMPLEMENT_TOKEN  		builtins.TOKEN_KIND = "impl"
	RETURN_TOKEN     		builtins.TOKEN_KIND = "ret"
	IN_TOKEN         		builtins.TOKEN_KIND = "in"
	AS_TOKEN         		builtins.TOKEN_KIND = "as"
	SAFE_TOKEN	   			builtins.TOKEN_KIND = "safe"
	OTHERWISE_TOKEN			builtins.TOKEN_KIND = "otherwise"
	//data types
	INT8_TOKEN      		builtins.TOKEN_KIND = builtins.INT8
	INT16_TOKEN     		builtins.TOKEN_KIND = builtins.INT16
	INT32_TOKEN     		builtins.TOKEN_KIND = builtins.INT32
	INT64_TOKEN     		builtins.TOKEN_KIND = builtins.INT64
	UINT8_TOKEN     		builtins.TOKEN_KIND = builtins.UINT8
	UINT16_TOKEN    		builtins.TOKEN_KIND = builtins.UINT16
	UINT32_TOKEN    		builtins.TOKEN_KIND = builtins.UINT32
	UINT64_TOKEN    		builtins.TOKEN_KIND = builtins.UINT64
	FLOAT32_TOKEN   		builtins.TOKEN_KIND = builtins.FLOAT32
	FLOAT64_TOKEN   		builtins.TOKEN_KIND = builtins.FLOAT64
	STR_TOKEN       		builtins.TOKEN_KIND = builtins.STRING
	BYTE_TOKEN      		builtins.TOKEN_KIND = builtins.BYTE
	BOOL_TOKEN      		builtins.TOKEN_KIND = builtins.BOOL
	NULL_TOKEN      		builtins.TOKEN_KIND = builtins.NULL
	STRUCT_TOKEN    		builtins.TOKEN_KIND = builtins.STRUCT
	FUNCTION_TOKEN  		builtins.TOKEN_KIND = builtins.FUNCTION
	INTERFACE_TOKEN 		builtins.TOKEN_KIND = builtins.INTERFACE
	MAYBE_TOKEN				builtins.TOKEN_KIND = builtins.MAYBE
	MAP_TOKEN       		builtins.TOKEN_KIND = builtins.MAP

	//increment and decrement
	PLUS_PLUS_TOKEN   		builtins.TOKEN_KIND = "++"
	MINUS_MINUS_TOKEN 		builtins.TOKEN_KIND = "--"
	//logical operators
	NOT_TOKEN           	builtins.TOKEN_KIND = "!"
	MINUS_TOKEN         	builtins.TOKEN_KIND = "-"
	PLUS_TOKEN          	builtins.TOKEN_KIND = "+"
	MUL_TOKEN           	builtins.TOKEN_KIND = "*"
	DIV_TOKEN           	builtins.TOKEN_KIND = "/"
	MOD_TOKEN           	builtins.TOKEN_KIND = "%"
	EXP_TOKEN           	builtins.TOKEN_KIND = "^"
	LESS_TOKEN          	builtins.TOKEN_KIND = "<"
	GREATER_TOKEN       	builtins.TOKEN_KIND = ">"
	LESS_EQUAL_TOKEN    	builtins.TOKEN_KIND = "<="
	GREATER_EQUAL_TOKEN 	builtins.TOKEN_KIND = ">="
	NOT_EQUAL_TOKEN     	builtins.TOKEN_KIND = "!="
	DOUBLE_EQUAL_TOKEN  	builtins.TOKEN_KIND = "=="
	//assignment
	WALRUS_TOKEN       		builtins.TOKEN_KIND = ":="
	COLON_TOKEN        		builtins.TOKEN_KIND = ":"
	EQUALS_TOKEN       		builtins.TOKEN_KIND = "="
	PLUS_EQUALS_TOKEN  		builtins.TOKEN_KIND = "+="
	MINUS_EQUALS_TOKEN 		builtins.TOKEN_KIND = "-="
	MUL_EQUALS_TOKEN   		builtins.TOKEN_KIND = "*="
	DIV_EQUALS_TOKEN   		builtins.TOKEN_KIND = "/="
	MOD_EQUALS_TOKEN   		builtins.TOKEN_KIND = "%="
	EXP_EQUALS_TOKEN   		builtins.TOKEN_KIND = "^="
	//delimiters
	OPEN_PAREN       		builtins.TOKEN_KIND = "("
	CLOSE_PAREN      		builtins.TOKEN_KIND = ")"
	OPEN_BRACKET     		builtins.TOKEN_KIND = "["
	CLOSE_BRACKET    		builtins.TOKEN_KIND = "]"
	OPEN_CURLY       		builtins.TOKEN_KIND = "{"
	CLOSE_CURLY      		builtins.TOKEN_KIND = "}"
	COMMA_TOKEN      		builtins.TOKEN_KIND = ","
	DOT_TOKEN        		builtins.TOKEN_KIND = "."
	SEMI_COLON_TOKEN 		builtins.TOKEN_KIND = ";"
	ARROW_TOKEN      		builtins.TOKEN_KIND = "->"
	FAT_ARROW_TOKEN  		builtins.TOKEN_KIND = "=>"
	OPTIONAL_TOKEN   		builtins.TOKEN_KIND = "?:"
	AT_TOKEN         		builtins.TOKEN_KIND = "@"
	DOLLAR_TOKEN     		builtins.TOKEN_KIND = "$"
	EOF_TOKEN        		builtins.TOKEN_KIND = "eof"
)

var keyWordsMap map[string]builtins.TOKEN_KIND = map[string]builtins.TOKEN_KIND{
	"let":       LET_TOKEN,
	"const":     CONST_TOKEN,
	"if":        IF_TOKEN,
	"else":      ELSE_TOKEN,
	"for":       FOR_TOKEN,
	"foreach":   FOREACH_TOKEN,
	"type":      TYPE_TOKEN,
	"priv":      PRIVATE_TOKEN,
	"interface": INTERFACE_TOKEN,
	"impl":      IMPLEMENT_TOKEN,
	"struct":    STRUCT_TOKEN,
	"fn":        FUNCTION_TOKEN,
	"map":       MAP_TOKEN,
	"maybe":	 MAYBE_TOKEN,
	"safe":		 SAFE_TOKEN,
	"otherwise": OTHERWISE_TOKEN,
	"ret":       RETURN_TOKEN,
	"in":        IN_TOKEN,
	"as":        AS_TOKEN,
}

func IsKeyword(token string) bool {
	if _, ok := keyWordsMap[token]; ok {
		return true
	}
	return false
}

type Token struct {
	Kind  builtins.TOKEN_KIND
	Value string
	Start Position
	End   Position
}

func (t *Token) Debug(filename string) {
	utils.GREY.Printf("%s:%d:%d ", filename, t.Start.Line, t.Start.Column)
	if t.Value == string(t.Kind) {
		fmt.Printf("'%s'\n", t.Value)
	} else {
		fmt.Printf("'%s' ('%v')\n", t.Value, t.Kind)
	}
}

func NewToken(kind builtins.TOKEN_KIND, value string, start Position, end Position) Token {
	return Token{
		Kind:  kind,
		Value: value,
		Start: start,
		End:   end,
	}
}
