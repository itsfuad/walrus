package ast

import (
	"testing"
	"walrus/lexer"
)

func TestIdentifierExpr(t *testing.T) {
	expr := IdentifierExpr{
		Name: "test",
		Location: Location{
			Start: lexer.Position{Line: 1, Column: 1},
			End:   lexer.Position{Line: 1, Column: 5},
		},
	}

	if expr.Name != "test" {
		t.Errorf("expected Name to be 'test', got %s", expr.Name)
	}

	if expr.StartPos() != (lexer.Position{Line: 1, Column: 1}) {
		t.Errorf("expected StartPos to be {Line: 1, Column: 1}, got %v", expr.StartPos())
	}

	if expr.EndPos() != (lexer.Position{Line: 1, Column: 5}) {
		t.Errorf("expected EndPos to be {Line: 1, Column: 5}, got %v", expr.EndPos())
	}
}

func TestIntegerLiteralExpr(t *testing.T) {
	expr := IntegerLiteralExpr{
		Value: "123",
		Location: Location{
			Start: lexer.Position{Line: 2, Column: 1},
			End:   lexer.Position{Line: 2, Column: 4},
		},
	}

	if expr.Value != "123" {
		t.Errorf("expected Value to be '123', got %s", expr.Value)
	}

	if expr.StartPos() != (lexer.Position{Line: 2, Column: 1}) {
		t.Errorf("expected StartPos to be {Line: 2, Column: 1}, got %v", expr.StartPos())
	}

	if expr.EndPos() != (lexer.Position{Line: 2, Column: 4}) {
		t.Errorf("expected EndPos to be {Line: 2, Column: 4}, got %v", expr.EndPos())
	}
}

