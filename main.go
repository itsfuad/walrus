package main

import (
	"encoding/json"
	"fmt"
	"os"

	"walrus/analyzer"
	"walrus/ast"
	"walrus/errgen"
	"walrus/lexer"
	"walrus/parser"
)

func serialize(root *ast.Node, folder, filename string) error {

	//create the folder if it does not exist
	if _, err := os.Stat(folder + "/ast"); os.IsNotExist(err) {
		os.Mkdir(folder+"/ast", os.ModePerm)
	}

	file, err := os.Create(fmt.Sprintf("%s/ast/%s.json", folder, filename))
	if err != nil {
		fmt.Printf("Error creating file: %s", err)
	}
	defer file.Close()

	//write the tree to a file named 'expressions.json' in 'code/ast' folder
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "    ")
	err = encoder.Encode(root)
	if err != nil {
		fmt.Printf("Error encoding JSON: %s", err)
		return err
	}
	return nil
}

func main() {

	fileName := "userTypes"
	folder := "code"
	filePath := fmt.Sprintf("%s/%s.wal", folder, fileName)
	tokens := lexer.Tokenize(filePath, true)
	tree := parser.NewParser(filePath, tokens).Parse(false)

	//write the tree to a file named 'expressions.json' in 'code/ast' folder
	serialize(&tree, folder, fileName)

	anz := analyzer.ProgramEnv(filePath)

	analyzer.CheckAST(tree, anz)

	errgen.DisplayErrors()
}
