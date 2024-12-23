package ast

import "walrus/position"

type ProgramStmt struct {
	Contents []Node
	position.Location
}

func (a ProgramStmt) INode() {
	//empty method implements Node interface
}
func (a ProgramStmt) StartPos() position.Coordinate {
	return a.Location.Start
}
func (a ProgramStmt) EndPos() position.Coordinate {
	return a.Location.End
}

type VarDeclStmtVar struct {
	Identifier   IdentifierExpr
	Value        Node
	ExplicitType DataType
	position.Location
}

type VarDeclStmt struct {
	Variables []VarDeclStmtVar
	IsConst   bool
	position.Location
}

func (a VarDeclStmt) INode() {
	//empty method implements Node interface
}
func (a VarDeclStmt) StartPos() position.Coordinate {
	return a.Location.Start
}
func (a VarDeclStmt) EndPos() position.Coordinate {
	return a.Location.End
}

type TypeDeclStmt struct {
	UDTypeValue DataType
	UDTypeName  IdentifierExpr
	position.Location
}

func (a TypeDeclStmt) INode() {
	//empty method implements Node interface
}
func (a TypeDeclStmt) StartPos() position.Coordinate {
	return a.Location.Start
}
func (a TypeDeclStmt) EndPos() position.Coordinate {
	return a.Location.End
}

type BlockStmt struct {
	Contents []Node
	position.Location
}

func (a BlockStmt) INode() {
	//empty method implements Node interface
}
func (a BlockStmt) StartPos() position.Coordinate {
	return a.Location.Start
}
func (a BlockStmt) EndPos() position.Coordinate {
	return a.Location.End
}

type IfStmt struct {
	Condition      Node
	Block          BlockStmt
	AlternateBlock interface{}
	position.Location
}

func (a IfStmt) INode() {
	//empty method implements Node interface
}
func (a IfStmt) StartPos() position.Coordinate {
	return a.Location.Start
}
func (a IfStmt) EndPos() position.Coordinate {
	return a.Location.End
}

type ForStmt struct {
	Init      Node
	Condition Node
	Increment Node
	Block     BlockStmt
	position.Location
}

func (a ForStmt) INode() {
	//empty method implements Node interface
}
func (a ForStmt) StartPos() position.Coordinate {
	return a.Location.Start
}
func (a ForStmt) EndPos() position.Coordinate {
	return a.Location.End
}

type ForEachStmt struct {
	Key      Node
	Value    Node
	Iterable Node
	Block    BlockStmt
	position.Location
}

func (a ForEachStmt) INode() {
	//empty method implements Node interface
}

func (a ForEachStmt) StartPos() position.Coordinate {
	return a.Location.Start
}

func (a ForEachStmt) EndPos() position.Coordinate {
	return a.Location.End
}

type FunctionParam struct {
	Identifier   IdentifierExpr
	Type         DataType
	IsOptional   bool
	DefaultValue Node
	position.Location
}

func (a FunctionParam) INode() {
	//empty method implements Node interface
}

func (a FunctionParam) StartPos() position.Coordinate {
	return a.Location.Start
}

func (a FunctionParam) EndPos() position.Coordinate {
	return a.Location.End
}

type FunctionDeclStmt struct {
	Identifier IdentifierExpr
	FunctionLiteral
}

func (a FunctionDeclStmt) INode() {
	//empty method implements Node interface
}

func (a FunctionDeclStmt) StartPos() position.Coordinate {
	return a.Location.Start
}

func (a FunctionDeclStmt) EndPos() position.Coordinate {
	return a.Location.End
}

type ReturnStmt struct {
	Value Node
	position.Location
}

func (a ReturnStmt) INode() {
	//empty method implements Node interface
}

func (a ReturnStmt) StartPos() position.Coordinate {
	return a.Location.Start
}

func (a ReturnStmt) EndPos() position.Coordinate {
	return a.Location.End
}

type MethodToImplement struct {
	IsPrivate bool
	FunctionDeclStmt
}

type ImplStmt struct {
	ImplFor IdentifierExpr
	Methods []MethodToImplement
	position.Location
}

func (a ImplStmt) INode() {
	//empty method implements Node interface
}

func (a ImplStmt) StartPos() position.Coordinate {
	return a.Location.Start
}

func (a ImplStmt) EndPos() position.Coordinate {
	return a.Location.End
}

type SafeStmt struct {
	Value       IdentifierExpr
	SafeBlock   BlockStmt
	UnsafeBlock BlockStmt
	position.Location
}

func (a SafeStmt) INode() {
	//empty method implements Node interface
}

func (a SafeStmt) StartPos() position.Coordinate {
	return a.Location.Start
}

func (a SafeStmt) EndPos() position.Coordinate {
	return a.Location.End
}
