package ast

import "walrus/position"

type NODE_TYPE int

type Node interface {
	INode()
	StartPos() position.Coordinate
	EndPos() position.Coordinate
}
