package position

import (
	"testing"
)

func TestAdvance(t *testing.T) {
	tests := []struct {
		name     string
		initial  Coordinate
		toSkip   string
		expected Coordinate
	}{
		{
			name:     "advance single byte",
			initial:  Coordinate{Line: 1, Column: 1, Index: 0},
			toSkip:   "a",
			expected: Coordinate{Line: 1, Column: 2, Index: 1},
		},
		{
			name:     "advance multiple bytes",
			initial:  Coordinate{Line: 1, Column: 1, Index: 0},
			toSkip:   "abc",
			expected: Coordinate{Line: 1, Column: 4, Index: 3},
		},
		{
			name:     "advance with newline",
			initial:  Coordinate{Line: 1, Column: 1, Index: 0},
			toSkip:   "a\nb",
			expected: Coordinate{Line: 2, Column: 2, Index: 3},
		},
		{
			name:     "advance with multiple newlines",
			initial:  Coordinate{Line: 1, Column: 1, Index: 0},
			toSkip:   "a\nb\nc",
			expected: Coordinate{Line: 3, Column: 2, Index: 5},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pos := &tt.initial
			pos.Advance(tt.toSkip)
			if pos.Line != tt.expected.Line || pos.Column != tt.expected.Column || pos.Index != tt.expected.Index {
				t.Errorf("Advance() = %v, want %v", pos, tt.expected)
			}
		})
	}
}

type TestNode struct {
	start Coordinate
	end   Coordinate
}

func (t TestNode) INode() {
	// IMPLEMENT
}

func (t TestNode) StartPos() Coordinate {
	return t.start
}

func (t TestNode) EndPos() Coordinate {
	return t.end
}

func TestNodeInterface(t *testing.T) {
	startPos := Coordinate{Line: 1, Column: 1}
	endPos := Coordinate{Line: 1, Column: 10}
	node := TestNode{start: startPos, end: endPos}

	if node.StartPos() != startPos {
		t.Errorf("expected start position %v, got %v", startPos, node.StartPos())
	}

	if node.EndPos() != endPos {
		t.Errorf("expected end position %v, got %v", endPos, node.EndPos())
	}
}

func TestLocation(t *testing.T) {
	startPos := Coordinate{Line: 1, Column: 1}
	endPos := Coordinate{Line: 1, Column: 10}
	location := Location{Start: startPos, End: endPos}

	if location.Start != startPos {
		t.Errorf("expected start position %v, got %v", startPos, location.Start)
	}

	if location.End != endPos {
		t.Errorf("expected end position %v, got %v", endPos, location.End)
	}
}
