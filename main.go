package main

import (
	"fmt"
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
	tree := parser.Parse(false)

	tc := typechecker.NewTypeENV(nil, typechecker.GLOBAL_SCOPE, "global", filePath)

	tc.DeclareVar("null", typechecker.NewNull(), true, false)
	tc.DeclareVar("true", typechecker.NewBool(), true, false)
	tc.DeclareVar("false", typechecker.NewBool(), true, false)
	tc.DeclareVar("PI", typechecker.NewFloat(32), true, false)
	/*
	tc.DeclareType("i8", typechecker.NewInt(8, true))
	tc.DeclareType("i16", typechecker.NewInt(16, true))
	tc.DeclareType("i32", typechecker.NewInt(32, true))
	tc.DeclareType("i64", typechecker.NewInt(64, true))
	tc.DeclareType("u8", typechecker.NewInt(8, false))
	tc.DeclareType("u16", typechecker.NewInt(16, false))
	tc.DeclareType("u32", typechecker.NewInt(32, false))
	tc.DeclareType("u64", typechecker.NewInt(64, false))
	tc.DeclareType("f32", typechecker.NewFloat(32))
	tc.DeclareType("f64", typechecker.NewFloat(64))
	tc.DeclareType("str", typechecker.NewStr())
	tc.DeclareType("byte", typechecker.NewInt(8, false))
	tc.DeclareType("bool", typechecker.NewBool())
	tc.DeclareType("void", typechecker.NewVoid())
	tc.DeclareType("null", typechecker.NewNull())
	*/
	typechecker.CheckAST(tree, tc)

	errgen.DisplayErrors()
}