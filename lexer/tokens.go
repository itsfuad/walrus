package lexer

import (
	"fmt"
	"walrus/builtins"
)

type TOKEN_KIND string

const (
	//variable declaration
	LET_TOKEN           TOKEN_KIND = "let"
	CONST_TOKEN         TOKEN_KIND = "const"
	TYPE_TOKEN          TOKEN_KIND = "type"
	IF_TOKEN            TOKEN_KIND = "if"
	ELSE_TOKEN          TOKEN_KIND = "els"
	ELSEIF_TOKEN        TOKEN_KIND = "elf"
	IDENTIFIER_TOKEN    TOKEN_KIND = "identifier"
	PRIVATE_TOKEN       TOKEN_KIND = "priv"
	INT                 TOKEN_KIND = builtins.INT
	FLOAT               TOKEN_KIND = builtins.FLOAT
	CHR                 TOKEN_KIND = builtins.CHAR
	STR                 TOKEN_KIND = builtins.STRING
	BOOL                TOKEN_KIND = builtins.BOOL
	NULL                TOKEN_KIND = builtins.NULL
	NOT_TOKEN           TOKEN_KIND = "!"
	MINUS_TOKEN         TOKEN_KIND = "-"
	PLUS_TOKEN          TOKEN_KIND = "+"
	MUL_TOKEN           TOKEN_KIND = "*"
	DIV_TOKEN           TOKEN_KIND = "/"
	MOD_TOKEN           TOKEN_KIND = "%"
	EXP_TOKEN           TOKEN_KIND = "^"
	LESS_TOKEN          TOKEN_KIND = "<"
	GREATER_TOKEN       TOKEN_KIND = ">"
	LESS_EQUAL_TOKEN    TOKEN_KIND = "<="
	GREATER_EQUAL_TOKEN TOKEN_KIND = ">="
	NOT_EQUAL_TOKEN     TOKEN_KIND = "!="
	DOUBLE_EQUAL_TOKEN  TOKEN_KIND = "=="
	WALRUS_TOKEN        TOKEN_KIND = ":="
	COLON_TOKEN         TOKEN_KIND = ":"
	EQUALS_TOKEN        TOKEN_KIND = "="
	OPEN_PAREN          TOKEN_KIND = "("
	CLOSE_PAREN         TOKEN_KIND = ")"
	OPEN_BRACKET        TOKEN_KIND = "["
	CLOSE_BRACKET       TOKEN_KIND = "]"
	OPEN_CURLY          TOKEN_KIND = "{"
	CLOSE_CURLY         TOKEN_KIND = "}"
	COMMA_TOKEN         TOKEN_KIND = ","
	DOT_TOKEN           TOKEN_KIND = "."
	SEMI_COLON_TOKEN    TOKEN_KIND = ";"
	AT_TOKEN            TOKEN_KIND = "@"
	EOF_TOKEN           TOKEN_KIND = "eof"
)

var keyWordsMap map[string]TOKEN_KIND = map[string]TOKEN_KIND{
	"let":    LET_TOKEN,
	"const":  CONST_TOKEN,
	"if":     IF_TOKEN,
	"els":    ELSE_TOKEN,
	"elf":    ELSEIF_TOKEN,
	"type":   TYPE_TOKEN,
	"priv":   PRIVATE_TOKEN,
	"struct": builtins.STRUCT,
	"fn":     builtins.FUNCTION,
}

func IsKeyword(token string) bool {
	if _, ok := keyWordsMap[token]; ok {
		return true
	}
	return false
}

type Token struct {
	Kind  TOKEN_KIND
	Value string
	Start Position
	End   Position
}

func (t *Token) Debug(filename string) {
	if t.Value == string(t.Kind) {
		fmt.Printf("%s:%d:%d\t '%s'\n", filename, t.Start.Line, t.Start.Column, t.Value)
	} else {
		fmt.Printf("%s:%d:%d\t '%s'\t('%v')\n", filename, t.Start.Line, t.Start.Column, t.Value, t.Kind)
	}
}

func NewToken(kind TOKEN_KIND, value string, start Position, end Position) Token {
	return Token{
		Kind:  kind,
		Value: value,
		Start: start,
		End:   end,
	}
}
