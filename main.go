package main

import (
	"walrus/errgen"
	"walrus/lexer"
	parseMachine "walrus/parser"
	"walrus/typechecker"
)

func main() {

	filePath := "language/types.wal"
	tokens := lexer.Tokenize(filePath, true)
	parser := parseMachine.NewParser(filePath, tokens)
	tree := parser.Parse(false)

	tc := typechecker.ProgramEnv(filePath)

	typechecker.CheckAST(tree, tc)

	errgen.DisplayErrors()
}