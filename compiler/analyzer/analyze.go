package analyzer

import (
	"fmt"
	"os"
	"path/filepath"
	"walrus/compiler/colors"
	"walrus/compiler/internal/parser"
	"walrus/compiler/internal/typechecker"
	"walrus/compiler/io"
	"walrus/compiler/report"
)

func Analyze(filePath string, displayErrors, debug, save2Json bool) {

	defer func() {
		if r := recover(); r != nil {
			if displayErrors {
				report.DisplayAll()
				colors.BOLD_RED.Println(r)
			}
		}
	}()

	//get the folder and file name
	folder, fileName := filepath.Split(filePath)

	tree, err := parser.NewParser(filePath, debug).Parse(save2Json)
	if err != nil {
		fmt.Println(report.TreeFormatString("compilation halted", "Error parsing file", err.Error()))
		os.Exit(-1)
	}
	//write the tree to a file named 'expressions.json' in 'code/ast' folder
	err = io.Serialize(&tree, folder, fileName)

	if err != nil {
		fmt.Println(report.TreeFormatString("compilation halted", "Error serializing AST", err.Error()))
		os.Exit(-1)
	}

	typechecker.Analyze(tree, filePath)
}
