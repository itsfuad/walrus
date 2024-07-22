package ast

import "walrus/lexer"

type ProgramStmt struct {
	Contents []Node
	Location
}

func (a ProgramStmt) INode() {
	//empty method implements Node interface
}
func (a ProgramStmt) StartPos() lexer.Position {
	return a.Location.Start
}
func (a ProgramStmt) EndPos() lexer.Position {
	return a.Location.End
}

type VarDeclStmt struct {
	Variable     IdentifierExpr
	Value        Node
	ExplicitType DataType
	IsConst      bool
	IsAssigned   bool
	Location
}

func (a VarDeclStmt) INode() {
	//empty method implements Node interface
}
func (a VarDeclStmt) StartPos() lexer.Position {
	return a.Location.Start
}
func (a VarDeclStmt) EndPos() lexer.Position {
	return a.Location.End
}