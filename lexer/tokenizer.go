package lexer

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"walrus/errgen"
)

type regexHandler func(lex *Lexer, regex *regexp.Regexp)

type regexPattern struct {
	regex   *regexp.Regexp
	handler regexHandler
}

type Lexer struct {
	Errors     []error
	Tokens     []Token
	Position   Position
	sourceCode []byte
	patterns   []regexPattern
	FilePath   string
}

func (lex *Lexer) advance(match string) {
	lex.Position.Advance(match)
}

func (lex *Lexer) push(token Token) {
	lex.Tokens = append(lex.Tokens, token)
}

func (lex *Lexer) at() byte {
	return lex.sourceCode[lex.Position.Index]
}

func (lex *Lexer) remainder() string {
	return string(lex.sourceCode)[lex.Position.Index:]
}

func (lex *Lexer) atEOF() bool {
	return lex.Position.Index >= len(lex.sourceCode)
}

func createLexer(filePath *string) *Lexer {

	fileText, err := os.ReadFile(*filePath)
	if err != nil {
		fmt.Printf("cannot compile: %v", err)
		os.Exit(-1)
	}

	lex := &Lexer{
		sourceCode: fileText,
		Tokens:     make([]Token, 0),
		Position: Position{
			Line:   1,
			Column: 1,
			Index:  0,
		},

		patterns: []regexPattern{
			//{regexp.MustCompile(`\n`), skipHandler}, // newlines
			{regexp.MustCompile(`\s+`), skipHandler},                          // whitespace
			{regexp.MustCompile(`\/\/.*`), skipHandler},                       // single line comments
			{regexp.MustCompile(`\/\*[\s\S]*?\*\/`), skipHandler},             // multi line comments
			{regexp.MustCompile(`"[^"]*"`), stringHandler},                    // string literals
			{regexp.MustCompile(`'[^']'`), characterHandler},                  // character literals
			{regexp.MustCompile(`[0-9]+(?:\.[0-9]+)?`), numberHandler},        // decimal numbers
			{regexp.MustCompile(`[a-zA-Z_][a-zA-Z0-9_]*`), identifierHandler}, // identifiers
			{regexp.MustCompile(`\+\+`), defaultHandler(PLUS_PLUS_TOKEN, "++")},
			{regexp.MustCompile(`\-\-`), defaultHandler(MINUS_MINUS_TOKEN, "--")},
			{regexp.MustCompile(`\->`), defaultHandler(ARROW_TOKEN, "->")},
			{regexp.MustCompile(`@`), defaultHandler(AT_TOKEN, "@")},
			{regexp.MustCompile(`!=`), defaultHandler(NOT_EQUAL_TOKEN, "!=")},
			{regexp.MustCompile(`!`), defaultHandler(NOT_TOKEN, "!")},
			{regexp.MustCompile(`\-`), defaultHandler(MINUS_TOKEN, "-")},
			{regexp.MustCompile(`\+`), defaultHandler(PLUS_TOKEN, "+")},
			{regexp.MustCompile(`\*`), defaultHandler(MUL_TOKEN, "*")},
			{regexp.MustCompile(`/`), defaultHandler(DIV_TOKEN, "/")},
			{regexp.MustCompile(`%`), defaultHandler(MOD_TOKEN, "%")},
			{regexp.MustCompile(`\^`), defaultHandler(EXP_TOKEN, "^")},
			{regexp.MustCompile(`:=`), defaultHandler(WALRUS_TOKEN, ":=")},
			{regexp.MustCompile(`<`), defaultHandler(LESS_EQUAL_TOKEN, "<=")},
			{regexp.MustCompile(`<`), defaultHandler(LESS_TOKEN, "<")},
			{regexp.MustCompile(`>=`), defaultHandler(GREATER_EQUAL_TOKEN, ">=")},
			{regexp.MustCompile(`>`), defaultHandler(GREATER_TOKEN, ">")},
			{regexp.MustCompile(`==`), defaultHandler(DOUBLE_EQUAL_TOKEN, "==")},
			{regexp.MustCompile(`=`), defaultHandler(EQUALS_TOKEN, "=")},
			{regexp.MustCompile(`\?:`), defaultHandler(OPTIONAL_TOKEN, "?:")},
			{regexp.MustCompile(`:`), defaultHandler(COLON_TOKEN, ":")},
			{regexp.MustCompile(`;`), defaultHandler(SEMI_COLON_TOKEN, ";")},
			{regexp.MustCompile(`\(`), defaultHandler(OPEN_PAREN, "(")},
			{regexp.MustCompile(`\)`), defaultHandler(CLOSE_PAREN, ")")},
			{regexp.MustCompile(`\[`), defaultHandler(OPEN_BRACKET, "[")},
			{regexp.MustCompile(`\]`), defaultHandler(CLOSE_BRACKET, "]")},
			{regexp.MustCompile(`\{`), defaultHandler(OPEN_CURLY, "{")},
			{regexp.MustCompile(`\}`), defaultHandler(CLOSE_CURLY, "}")},
			{regexp.MustCompile(","), defaultHandler(COMMA_TOKEN, ",")},
			{regexp.MustCompile(`\.`), defaultHandler(DOT_TOKEN, ".")},
		},
	}
	return lex
}

// defaultHandler returns a regexHandler function that processes a token of the specified kind and value.
// The returned function advances the lexer's position by the length of the value and pushes a new token
// with the given kind and value to the lexer.
//
// Parameters:
// - kind: The type of the token to be processed.
// - value: The string value of the token.
//
// Returns:
// - A regexHandler function that processes the token and updates the lexer state.
func defaultHandler(kind TOKEN_KIND, value string) regexHandler {

	return func(lex *Lexer, _ *regexp.Regexp) {

		start := lex.Position
		lex.advance(value)
		end := lex.Position

		lex.push(NewToken(kind, value, start, end))
	}
}

// identifierHandler processes an identifier token in the lexer.
// It uses a regular expression to find the identifier in the lexer's remainder,
// advances the lexer's position, and then determines if the identifier is a keyword.
// If it is a keyword, it pushes a keyword token onto the lexer's token stack;
// otherwise, it pushes an identifier token.
//
// Parameters:
// - lex: A pointer to the Lexer instance.
// - regex: A regular expression used to find the identifier in the lexer's remainder.
func identifierHandler(lex *Lexer, regex *regexp.Regexp) {
	identifier := regex.FindString(lex.remainder())
	start := lex.Position
	lex.advance(identifier)
	end := lex.Position
	if IsKeyword(identifier) {
		lex.push((NewToken(TOKEN_KIND(identifier), identifier, start, end)))
	} else {
		lex.push(NewToken(IDENTIFIER_TOKEN, identifier, start, end))
	}
}

// numberHandler processes numeric tokens from the lexer input.
// It uses a regular expression to find a numeric match in the lexer's remainder,
// advances the lexer's position, and determines whether the number is a float or an integer.
// Depending on the type, it pushes the appropriate token (FLOAT or INT) to the lexer's token stack.
//
// Parameters:
//
//	lex - A pointer to the Lexer instance.
//	regex - A compiled regular expression used to find numeric matches in the lexer's input.
func numberHandler(lex *Lexer, regex *regexp.Regexp) {
	match := regex.FindString(lex.remainder())
	start := lex.Position
	lex.advance(match)
	end := lex.Position
	//find the number is a float or an integer
	if strings.Contains(match, ".") {
		lex.push(NewToken(FLOAT, match, start, end))
	} else {
		lex.push(NewToken(INT, match, start, end))
	}
}

// stringHandler processes a string literal in the lexer, using the provided regular expression to match the string.
// It excludes the quotes from the matched string, updates the lexer's position, and pushes a new token.
//
// Parameters:
//
//	lex   - A pointer to the Lexer instance.
//	regex - A regular expression used to find the string literal in the lexer's remainder.
func stringHandler(lex *Lexer, regex *regexp.Regexp) {
	match := regex.FindString(lex.remainder())
	//exclude the quotes
	stringLiteral := match[1 : len(match)-1]
	start := lex.Position
	lex.advance(match)
	end := lex.Position
	lex.push(NewToken(STR, stringLiteral, start, end))
}

// characterHandler processes a character literal in the lexer.
// It uses a regular expression to find the character literal in the remaining input,
// excludes the quotes from the matched string, and then creates a new BYTE token
// with the character literal's value and its position in the input.
//
// Parameters:
//
//	lex - The Lexer instance that contains the input and current position.
//	regex - The regular expression used to match the character literal.
func characterHandler(lex *Lexer, regex *regexp.Regexp) {
	match := regex.FindString(lex.remainder())
	//exclude the quotes
	characterLiteral := match[1 : len(match)-1]
	start := lex.Position
	lex.advance(match)
	end := lex.Position
	lex.push(NewToken(BYTE, characterLiteral, start, end))
}

// skipHandler processes a token that should be skipped by the lexer.
func skipHandler(lex *Lexer, regex *regexp.Regexp) {
	match := regex.FindString(lex.remainder())
	lex.advance(match)
}

// Tokenize reads the source code from the specified file and tokenizes it.
func Tokenize(filename string, debug bool) []Token {

	lex := createLexer(&filename)
	lex.FilePath = filename

	for !lex.atEOF() {

		matched := false

		for _, pattern := range lex.patterns {

			loc := pattern.regex.FindStringIndex(lex.remainder())

			if loc != nil && loc[0] == 0 {
				pattern.handler(lex, pattern.regex)
				matched = true
				break
			}
		}

		if !matched {
			errStr := fmt.Sprintf("lexer:unexpected character: '%c'", lex.at())
			errgen.MakeError(filename, lex.Position.Line, lex.Position.Line, lex.Position.Column, lex.Position.Column, errStr).DisplayWithPanic()
			return nil
		}
	}

	lex.push(NewToken(EOF_TOKEN, "eof", lex.Position, lex.Position))

	//litter.Dump(lex.Tokens)
	if debug {
		for _, token := range lex.Tokens {
			token.Debug(filename)
		}
	}

	return lex.Tokens
}
