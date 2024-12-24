package main

import (
	//Standard packages
	"fmt"
	"os"

	//Walrus packages
	"walrus/frontend/helpers"
	"walrus/frontend/lexer"
	"walrus/frontend/parser"
	"walrus/frontend/typechecker"
	"walrus/report"
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
		fmt.Println("compilation halted" + report.TreeFormatString("Error serializing AST", err.Error()))
		os.Exit(-1)
	}

	typeCheckerEnv := typechecker.ProgramEnv(filePath)

	typechecker.CheckAST(tree, typeCheckerEnv)

	report.DisplayAll()
}
