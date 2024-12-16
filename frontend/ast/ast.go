package ast

import "walrus/frontend/lexer"

type NODE_TYPE int

type Node interface {
	INode()
	StartPos() lexer.Position
	EndPos() lexer.Position
}

type Location struct {
	Start lexer.Position
	End   lexer.Position
}
