package main

import (
	//Standard packages
	"fmt"
	"os"
	"path/filepath"

	//Walrus packages
	"walrus/compiler/colors"
	"walrus/compiler/internal/lexer"
	"walrus/compiler/internal/parser"
	"walrus/compiler/internal/report"
	"walrus/compiler/internal/typechecker"
	"walrus/compiler/io"
)

func main() {

	if len(os.Args) < 2 {
		colors.GREEN.Println("Usage: walrus <file>")
		os.Exit(-1)
	}

	filePath := os.Args[1]

	//must have .wal file
	if len(filePath) < 5 || filePath[len(filePath)-4:] != ".wal" {
		colors.RED.Println("Error: file must have .wal extension")
		os.Exit(-1)
	}

	//get the folder and file name
	folder, fileName := filepath.Split(filePath)

	tokens := lexer.Tokenize(filePath, true)
	tree := parser.NewParser(filePath, tokens).Parse(false)

	//write the tree to a file named 'expressions.json' in 'code/ast' folder
	err := io.Serialize(&tree, folder, fileName)
	if err != nil {
		fmt.Println(report.TreeFormatString("compilation halted", "Error serializing AST", err.Error()))
		os.Exit(-1)
	}

	typeCheckerEnv := typechecker.ProgramEnv(filePath)

	typechecker.CheckAST(tree, typeCheckerEnv)

	report.DisplayAll()
}
