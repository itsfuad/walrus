package main

import (
	"fmt"
	"walrus/errgen"
	"walrus/lexer"
	parseMachine "walrus/parser"
	"walrus/typechecker"
)

func main() {
	fmt.Println("Hello world!")
	filePath := "language/expressions.wal"
	tokens := lexer.Tokenize(filePath, true)
	parser := parseMachine.NewParser(filePath, tokens)
	tree := parser.Parse(false)

	tc := typechecker.ProgramEnv(filePath)

	typechecker.CheckAST(tree, tc)

	errgen.DisplayErrors()
}