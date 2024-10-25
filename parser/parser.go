package parser

import (
	"fmt"
	"walrus/ast"
	"walrus/errgen"
	"walrus/lexer"
)

type Parser struct {
	tokens   []lexer.Token
	FilePath string
	Errors   []error
	index    int
}

func (p *Parser) currentToken() lexer.Token {
	return p.tokens[p.index]
}

func (p *Parser) currentTokenKind() lexer.TOKEN_KIND {
	return p.currentToken().Kind
}

func (p *Parser) nextToken() lexer.Token {
	if p.index+1 < len(p.tokens) {
		return p.tokens[p.index+1]
	}
	return lexer.Token{Kind: lexer.EOF_TOKEN}
}

func (p *Parser) nextTokenKind() lexer.TOKEN_KIND {
	return p.nextToken().Kind
}

func (p *Parser) hasToken() bool {
	return p.index < len(p.tokens) && p.currentTokenKind() != lexer.EOF_TOKEN
}

func (p *Parser) advance() lexer.Token {
	token := p.currentToken()
	p.index++
	return token
}

func (p *Parser) expectError(expectedKind lexer.TOKEN_KIND, err error) lexer.Token {
	token := p.currentToken()
	kind := token.Kind

	start := token.Start
	end := token.End

	if kind != expectedKind {
		if err != nil {
			errgen.MakeError(p.FilePath, start.Line, end.Line, start.Column, end.Column, err.Error()).DisplayWithPanic()
		} else {
			msg := fmt.Sprintf("parser:expected '%s' but got '%s'", expectedKind, kind)
			errgen.MakeError(p.FilePath, start.Line, end.Line, start.Column, end.Column, msg).DisplayWithPanic()
		}
	}
	return p.advance()
}

func (p *Parser) expect(expectedKind lexer.TOKEN_KIND) lexer.Token {
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

	p.expect(lexer.SEMI_COLON_TOKEN)

	return expr
}

func (p *Parser) Parse() ast.Node {

	var contents []ast.Node

	for p.hasToken() {
		stmt := parseNode(p)
		contents = append(contents, stmt)
	}

	return ast.ProgramStmt{
		Contents: contents,
	}
}

func NewParser(filePath string, tokens []lexer.Token) *Parser {

	bindLookupHandlers()
	bindTypeLookups()

	parser := &Parser{
		tokens:   tokens,
		FilePath: filePath,
		index:    0,
	}
	return parser
}
