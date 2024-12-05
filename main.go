package main

import (
	"walrus/errgen"
	"walrus/lexer"
	"walrus/parser"
	"walrus/typechecker"
)

func main() {

	filePath := "language/types.wal"
	tokens := lexer.Tokenize(filePath, true)
	tree := parser.NewParser(filePath, tokens).Parse(false)

	tc := typechecker.ProgramEnv(filePath)

	typechecker.CheckAST(tree, tc)

	errgen.DisplayErrors()
}