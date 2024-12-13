package ast

import (
	"walrus/lexer"
)

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

type VarDeclStmtVar struct {
	Identifier   IdentifierExpr
	Value        Node
	ExplicitType DataType
	Location
}

type VarDeclStmt struct {
	Variables []VarDeclStmtVar
	IsConst   bool
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

type ForStmt struct {
	Init      Node
	Condition Node
	Increment Node
	Block     BlockStmt
	Location
}

func (a ForStmt) INode() {
	//empty method implements Node interface
}
func (a ForStmt) StartPos() lexer.Position {
	return a.Location.Start
}
func (a ForStmt) EndPos() lexer.Position {
	return a.Location.End
}

type ForEachStmt struct {
	Key      Node
	Value    Node
	Iterable Node
	Block    BlockStmt
	Location
}

func (a ForEachStmt) INode() {
	//empty method implements Node interface
}

func (a ForEachStmt) StartPos() lexer.Position {
	return a.Location.Start
}

func (a ForEachStmt) EndPos() lexer.Position {
	return a.Location.End
}

type FunctionParam struct {
	Identifier   IdentifierExpr
	Type         DataType
	IsOptional   bool
	DefaultValue Node
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
	FunctionLiteral
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

type ImplMethod struct {
	IsPrivate bool
	FunctionDeclStmt
}

type ImplStmt struct {
	ImplFor IdentifierExpr
	Methods []ImplMethod
	Location
}

func (a ImplStmt) INode() {
	//empty method implements Node interface
}

func (a ImplStmt) StartPos() lexer.Position {
	return a.Location.Start
}

func (a ImplStmt) EndPos() lexer.Position {
	return a.Location.End
}

type SafeStmt struct {
	Value       IdentifierExpr
	SafeBlock   BlockStmt
	UnsafeBlock BlockStmt
	Location
}

func (a SafeStmt) INode() {
	//empty method implements Node interface
}

func (a SafeStmt) StartPos() lexer.Position {
	return a.Location.Start
}

func (a SafeStmt) EndPos() lexer.Position {
	return a.Location.End
}
