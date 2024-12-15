package main

import (
	"fmt"
	"os"

	"walrus/analyzer"
	"walrus/errgen"
	"walrus/lexer"
	"walrus/parser"
	"walrus/helpers"
)

func main() {

	fileName := "userTypes"
	folder := "code"
	filePath := fmt.Sprintf("%s/%s.wal", folder, fileName)
	tokens := lexer.Tokenize(filePath, true)
	tree := parser.NewParser(filePath, tokens).Parse(false)

	//write the tree to a file named 'expressions.json' in 'code/ast' folder
	err := helpers.Serialize(&tree, folder, fileName)
	if err != nil {
		fmt.Println(errgen.TreeFormatString("compilation halted", "Error serializing AST", err.Error()))
		os.Exit(-1)
	}

	anz := analyzer.ProgramEnv(filePath)

	analyzer.CheckAST(tree, anz)

	errgen.DisplayAll()
}
