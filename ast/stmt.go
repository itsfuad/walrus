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

type TypeDeclStmt struct {
	UDType     DataType
	UDTypeName string
	Location
}

func (a TypeDeclStmt) INode() {
	//empty method implements Node interface
}
func (a TypeDeclStmt) StartPos() lexer.Position {
	return a.Location.Start
}
func (a TypeDeclStmt) EndPos() lexer.Position {
	return a.Location.End
}

type BlockStmt struct {
	Contents []Node
	Location
}

func (a BlockStmt) INode() {
	//empty method implements Node interface
}
func (a BlockStmt) StartPos() lexer.Position {
	return a.Location.Start
}
func (a BlockStmt) EndPos() lexer.Position {
	return a.Location.End
}

type IfStmt struct {
	Condition      Node
	Block          BlockStmt
	AlternateBlock interface{}
	Location
}

func (a IfStmt) INode() {
	//empty method implements Node interface
}
func (a IfStmt) StartPos() lexer.Position {
	return a.Location.Start
}
func (a IfStmt) EndPos() lexer.Position {
	return a.Location.End
}

type FunctionParam struct {
	Identifier IdentifierExpr
	Type       DataType
	Location
}

func (a FunctionParam) INode() {
	//empty method implements Node interface
}

func (a FunctionParam) StartPos() lexer.Position {
	return a.Location.Start
}

func (a FunctionParam) EndPos() lexer.Position {
	return a.Location.End
}

type FunctionDeclStmt struct {
	Identifier IdentifierExpr
	FunctionExpr
}

func (a FunctionDeclStmt) INode() {
	//empty method implements Node interface
}

func (a FunctionDeclStmt) StartPos() lexer.Position {
	return a.Location.Start
}

func (a FunctionDeclStmt) EndPos() lexer.Position {
	return a.Location.End
}

type ReturnStmt struct {
	Value Node
	Location
}

func (a ReturnStmt) INode() {
	//empty method implements Node interface
}

func (a ReturnStmt) StartPos() lexer.Position {
	return a.Location.Start
}

func (a ReturnStmt) EndPos() lexer.Position {
	return a.Location.End
}
