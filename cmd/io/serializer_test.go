package io

import (
	"os"
	"path/filepath"
	"testing"

	"walrus/internal/ast"
	"walrus/internal/lexer"
)

func TestSerialize(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := "temp"
	if _, err := os.Stat(tempDir); os.IsNotExist(err) {
		err := os.Mkdir(tempDir, os.ModePerm)
		if err != nil {
			t.Fatal(err)
		}
	}
	defer os.RemoveAll(tempDir)
	// Create a temporary file for testing
	tempFile := filepath.Join(tempDir, "temp.json")
	file, err := os.Create(tempFile)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	defer os.Remove(tempFile)

	// Create a sample AST node
	var program ast.Node = ast.ProgramStmt{
		Contents: []ast.Node{
			ast.VarDeclStmt{
				Variables: []ast.VarDeclStmtVar{
					{
						Identifier: ast.IdentifierExpr{
							Name: "x",
							Location: ast.Location{
								Start: lexer.Position{Line: 1, Column: 1},
								End:   lexer.Position{Line: 1, Column: 2},
							},
						},
						Value: ast.IntegerLiteralExpr{
							Value:   "10",
							BitSize: 64,
							Location: ast.Location{
								Start: lexer.Position{Line: 1, Column: 5},
								End:   lexer.Position{Line: 1, Column: 6},
							},
							IsSigned: true,
						},
						ExplicitType: nil,
						Location: ast.Location{
							Start: lexer.Position{Line: 1, Column: 1},
							End:   lexer.Position{Line: 1, Column: 6},
						},
					},
				},
				IsConst: false,
				Location: ast.Location{
					Start: lexer.Position{Line: 1, Column: 1},
					End:   lexer.Position{Line: 1, Column: 6},
				},
			},
		},
	}

	// Serialize the AST node to a file
	err = Serialize(&program, tempDir, "temp")
	if err != nil {
		t.Fatal(err)
	}
}
