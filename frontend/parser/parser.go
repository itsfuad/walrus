package parser

import (
	//Standard packages
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	//Walrus packages
	"walrus/frontend/ast"
	"walrus/frontend/builtins"
	"walrus/frontend/lexer"
	"walrus/report"
	"walrus/utils"
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

// rollback by N tokens
func (p *Parser) rollback(n int) {
	p.index -= n
}

func (p *Parser) expectError(expectedKind builtins.TOKEN_KIND, err error) lexer.Token {
	token := p.currentToken()
	kind := token.Kind

	start := token.Start
	end := token.End

	if kind != expectedKind {
		if err != nil {
			report.Add(p.FilePath, start.Line, end.Line, start.Column, end.Column, err.Error()).Level(report.SYNTAX_ERROR)
		} else {
			var msg string
			if lexer.IsKeyword(token.Value) {
				msg = fmt.Sprintf("unexpected keyword '%s' found. expected '%s'", token.Value, expectedKind)
			} else {
				msg = fmt.Sprintf("unexpected token '%s' found. expected '%s'", token.Value, expectedKind)
			}
			report.Add(p.FilePath, start.Line, end.Line, start.Column, end.Column, msg).Level(report.SYNTAX_ERROR)
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

	fmt.Printf("Parsing node %s\n", p.currentToken().Value)
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

func (p *Parser) Parse(saveJson bool) ast.Node {

	utils.GREEN.Printf("Parsing %s\n", p.FilePath)

	var contents []ast.Node

	for p.hasToken() {
		stmt := parseNode(p)
		contents = append(contents, stmt)
	}

	program := ast.ProgramStmt{
		Contents: contents,
	}

	if saveJson {
		file, err := os.Create(strings.TrimSuffix(p.FilePath, filepath.Ext(p.FilePath)) + ".json")
		if err != nil {
			panic(err)
		}

		//parse as string
		astString, err := json.MarshalIndent(program, "", "  ")

		if err != nil {
			panic(err)
		}

		_, err = file.Write(astString)

		if err != nil {
			panic(err)
		}

		file.Close()
	}

	utils.GREEN.Println("Parsing complete")

	return program
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
