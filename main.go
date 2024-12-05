package main

import (
	"walrus/errgen"
	"walrus/lexer"
	"walrus/parser"
	"walrus/analyzer"
)

func main() {

	filePath := "language/types.wal"
	tokens := lexer.Tokenize(filePath, true)
	tree := parser.NewParser(filePath, tokens).Parse(false)

	anz := analyzer.ProgramEnv(filePath)

	analyzer.CheckAST(tree, anz)

	errgen.DisplayErrors()
}