package main

import (
	"fmt"
	"os"

	"walrus/errgen"
	"walrus/frontend/helpers"
	"walrus/frontend/lexer"
	"walrus/frontend/parser"
	"walrus/frontend/typechecker"
)

func main() {

	fileName := "variables"
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

	typeCheckerEnv := typechecker.ProgramEnv(filePath)

	typechecker.CheckAST(tree, typeCheckerEnv)

	errgen.DisplayAll()
}
