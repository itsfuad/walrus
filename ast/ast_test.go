package ast_test

import (
	"testing"
	"walrus/ast"
	"walrus/lexer"
)

type TestNode struct {
	start lexer.Position
	end   lexer.Position
}

func (t TestNode) INode() {
	// IMPLEMENT
}

func (t TestNode) StartPos() lexer.Position {
	return t.start
}

func (t TestNode) EndPos() lexer.Position {
	return t.end
}

func TestNodeInterface(t *testing.T) {
	startPos := lexer.Position{Line: 1, Column: 1}
	endPos := lexer.Position{Line: 1, Column: 10}
	node := TestNode{start: startPos, end: endPos}

	if node.StartPos() != startPos {
		t.Errorf("expected start position %v, got %v", startPos, node.StartPos())
	}

	if node.EndPos() != endPos {
		t.Errorf("expected end position %v, got %v", endPos, node.EndPos())
	}
}

func TestLocation(t *testing.T) {
	startPos := lexer.Position{Line: 1, Column: 1}
	endPos := lexer.Position{Line: 1, Column: 10}
	location := ast.Location{Start: startPos, End: endPos}

	if location.Start != startPos {
		t.Errorf("expected start position %v, got %v", startPos, location.Start)
	}

	if location.End != endPos {
		t.Errorf("expected end position %v, got %v", endPos, location.End)
	}
}
