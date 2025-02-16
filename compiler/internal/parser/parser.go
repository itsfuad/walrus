package parser

import (
	//Standard packages
	"errors"
	"fmt"
	//Walrus packages
	"walrus/compiler/colors"
	"walrus/compiler/internal/ast"
	"walrus/compiler/internal/builtins"
	"walrus/compiler/internal/lexer"
	"walrus/compiler/report"
)

type Parser struct {
	tokens   []lexer.Token
	FilePath string
	Errors   []error
	index    int
}

func (p *Parser) currentToken() lexer.Token {
	if p.index >= len(p.tokens) {
		return lexer.Token{} // for safety
	}
	return p.tokens[p.index]
}

func (p *Parser) currentTokenKind() builtins.TOKEN_KIND {
	return p.currentToken().Kind
}

func (p *Parser) hasToken() bool {
	return p.index < len(p.tokens) && p.currentTokenKind() != lexer.EOF_TOKEN
}

func (p *Parser) eat() lexer.Token {
	token := p.currentToken()
	p.index++
	return token
}

func (p *Parser) expectError(expectedKind builtins.TOKEN_KIND, err error) lexer.Token {
	token := p.currentToken()
	kind := token.Kind

	start := token.Start
	end := token.End

	if kind != expectedKind {
		if err != nil {
			report.Add(p.FilePath, start.Line, end.Line, start.Column, end.Column, err.Error()).SetLevel(report.SYNTAX_ERROR)
		} else {
			var msg string
			if lexer.IsKeyword(token.Value) {
				msg = fmt.Sprintf("unexpected keyword '%s' found. expected '%s'", token.Value, expectedKind)
			} else {
				msg = fmt.Sprintf("unexpected token '%s' found. expected '%s'", token.Value, expectedKind)
			}
			report.Add(p.FilePath, start.Line, end.Line, start.Column, end.Column, msg).SetLevel(report.SYNTAX_ERROR)
		}
	}
	return p.eat()
}

func (p *Parser) expect(expectedKind builtins.TOKEN_KIND) lexer.Token {
	return p.expectError(expectedKind, nil)
}

type I interface {
	Display()
}

func parseNode(p *Parser) ast.Node {
	// can be a statement or an expression
	stmt_fn, exists := STMTLookup[p.currentTokenKind()]

	if exists {
		return stmt_fn(p)
	}

	// if not a statement, then it must be an expression
	expr := parseExpr(p, DEFAULT_BP)

	p.expectError(lexer.SEMI_COLON_TOKEN, errors.New("expected a semicolon at the end of statement"))

	return expr
}

func (p *Parser) Parse() (ast.Node, error) {

	var contents []ast.Node

	for p.hasToken() {
		stmt := parseNode(p)
		contents = append(contents, stmt)
	}

	program := ast.ProgramStmt{
		Contents: contents,
	}

	colors.GREEN.Println("Parsing complete")

	return program, nil
}

func NewParser(filePath string, debug bool) *Parser {

	tokens := lexer.Tokenize(filePath, debug)

	bindLookupHandlers()
	bindTypeLookups()

	parser := &Parser{
		tokens:   tokens,
		FilePath: filePath,
		index:    0,
	}
	return parser
}
