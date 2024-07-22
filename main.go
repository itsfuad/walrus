package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"walrus/lexer"
	parseMachine "walrus/parser"
	"walrus/typechecker"
)

func main() {
	fmt.Println("Hello world!")
	filePath := "language/arrayAccess.wal"
	tokens := lexer.Tokenize(filePath, true)
	parser := parseMachine.NewParser(filePath, tokens)
	tree := parser.Parse()

	file, err := os.Create(strings.TrimSuffix(filePath, filepath.Ext(filePath)) + ".json")
	if err != nil {
		panic(err)
	}

	//parse as string
	astString, err := json.MarshalIndent(tree, "", "  ")

	if err != nil {
		panic(err)
	}

	_, err = file.Write(astString)

	if err != nil {
		panic(err)
	}

	file.Close()

	tc := typechecker.NewTypeENV(nil, filePath)
	tc.DeclareVar("PI", &typechecker.Float{DataType: typechecker.FLOAT_TYPE, Name: "float"}, false)
	typechecker.EvaluateTypesOfNode(tree, tc)
}
