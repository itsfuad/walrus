package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"walrus/errgen"
	"walrus/lexer"
	parseMachine "walrus/parser"
	"walrus/typechecker"
)

func main() {
	fmt.Println("Hello world!")
	filePath := "language/expressions.wal"
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

	tc := typechecker.NewTypeENV(nil, typechecker.GLOBAL_SCOPE, "global", filePath)
	tc.DeclareVar("null", typechecker.NewNull(), true, false)
	tc.DeclareVar("true", typechecker.NewBool(), true, false)
	tc.DeclareVar("false", typechecker.NewBool(), true, false)
	tc.DeclareVar("PI", typechecker.NewFloat(32), true, false)
	typechecker.CheckAST(tree, tc)
	errgen.DisplayErrors()
}