package ast

import (
	"testing"
	"walrus/position"
)

func TestIdentifierExpr(t *testing.T) {
	expr := IdentifierExpr{
		Name: "test",
		Location: position.Location{
			Start: position.Coordinate{Line: 1, Column: 1},
			End:   position.Coordinate{Line: 1, Column: 5},
		},
	}

	if expr.Name != "test" {
		t.Errorf("expected Name to be 'test', got %s", expr.Name)
	}

	if expr.StartPos() != (position.Coordinate{Line: 1, Column: 1}) {
		t.Errorf("expected StartPos to be {Line: 1, Column: 1}, got %v", expr.StartPos())
	}

	if expr.EndPos() != (position.Coordinate{Line: 1, Column: 5}) {
		t.Errorf("expected EndPos to be {Line: 1, Column: 5}, got %v", expr.EndPos())
	}
}

func TestIntegerLiteralExpr(t *testing.T) {
	expr := IntegerLiteralExpr{
		Value: "123",
		Location: position.Location{
			Start: position.Coordinate{Line: 2, Column: 1},
			End:   position.Coordinate{Line: 2, Column: 4},
		},
	}

	if expr.Value != "123" {
		t.Errorf("expected Value to be '123', got %s", expr.Value)
	}

	if expr.StartPos() != (position.Coordinate{Line: 2, Column: 1}) {
		t.Errorf("expected StartPos to be {Line: 2, Column: 1}, got %v", expr.StartPos())
	}

	if expr.EndPos() != (position.Coordinate{Line: 2, Column: 4}) {
		t.Errorf("expected EndPos to be {Line: 2, Column: 4}, got %v", expr.EndPos())
	}
}
