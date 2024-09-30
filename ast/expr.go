package ast

import "walrus/lexer"

// Any word which is not a keyword or literal
type IdentifierExpr struct {
	Name string
	Location
}

func (a IdentifierExpr) INode() {
	//empty method implements Node interface
}
func (a IdentifierExpr) StartPos() lexer.Position {
	return a.Location.Start
}
func (a IdentifierExpr) EndPos() lexer.Position {
	return a.Location.End
}

//Literals or Raw values like: 1,2,3,4.6, "hello world", 'a' ...etc
type IntegerLiteralExpr struct {
	Value string
	Location
}

func (a IntegerLiteralExpr) INode() {
	//empty method implements Node interface
}
func (a IntegerLiteralExpr) StartPos() lexer.Position {
	return a.Location.Start
}
func (a IntegerLiteralExpr) EndPos() lexer.Position {
	return a.Location.End
}

type FloatLiteralExpr struct {
	Value string
	Location
}

func (a FloatLiteralExpr) INode() {
	//empty method implements Node interface
}
func (a FloatLiteralExpr) StartPos() lexer.Position {
	return a.Location.Start
}
func (a FloatLiteralExpr) EndPos() lexer.Position {
	return a.Location.End
}

type StringLiteralExpr struct {
	Value string
	Location
}

func (a StringLiteralExpr) INode() {
	//empty method implements Node interface
}
func (a StringLiteralExpr) StartPos() lexer.Position {
	return a.Location.Start
}
func (a StringLiteralExpr) EndPos() lexer.Position {
	return a.Location.End
}

type CharLiteralExpr struct {
	Value string
	Location
}

func (a CharLiteralExpr) INode() {
	//empty method implements Node interface
}
func (a CharLiteralExpr) StartPos() lexer.Position {
	return a.Location.Start
}
func (a CharLiteralExpr) EndPos() lexer.Position {
	return a.Location.End
}

type BooleanLiteralExpr struct {
	Value string
	Location
}

func (a BooleanLiteralExpr) INode() {
	//empty method implements Node interface
}
func (a BooleanLiteralExpr) StartPos() lexer.Position {
	return a.Location.Start
}
func (a BooleanLiteralExpr) EndPos() lexer.Position {
	return a.Location.End
}

type NullLiteralExpr struct {
	Value string
	Location
}

func (a NullLiteralExpr) INode() {
	//empty method implements Node interface
}
func (a NullLiteralExpr) StartPos() lexer.Position {
	return a.Location.Start
}
func (a NullLiteralExpr) EndPos() lexer.Position {
	return a.Location.End
}

type UnaryExpr struct {
	Operator lexer.Token
	Argument Node
	Location
}

func (a UnaryExpr) INode() {
	//empty method implements Node interface
}
func (a UnaryExpr) StartPos() lexer.Position {
	return a.Location.Start
}
func (a UnaryExpr) EndPos() lexer.Position {
	return a.Location.End
}

type BinaryExpr struct {
	Operator lexer.Token
	Left     Node
	Right    Node
	Location
}

func (a BinaryExpr) INode() {
	//empty method implements Node interface
}
func (a BinaryExpr) StartPos() lexer.Position {
	return a.Location.Start
}
func (a BinaryExpr) EndPos() lexer.Position {
	return a.Location.End
}

type VarAssignmentExpr struct {
	Assignee Node // Check later if we should use IdentifierExpr instead
	Value    Node
	Operator lexer.Token // Looks odd right? Well, we know the operator must be '='. But what about +=, -=, *= and so on?😀
	Location
}

func (a VarAssignmentExpr) INode() {
	//empty method implements Node interface
}
func (a VarAssignmentExpr) StartPos() lexer.Position {
	return a.Location.Start
}
func (a VarAssignmentExpr) EndPos() lexer.Position {
	return a.Location.End
}

type ArrayExpr struct {
	Values []Node
	Location
}

func (a ArrayExpr) INode() {
	//empty method implements Node interface
}
func (a ArrayExpr) StartPos() lexer.Position {
	return a.Location.Start
}
func (a ArrayExpr) EndPos() lexer.Position {
	return a.Location.End
}

type ArrayIndexAccess struct {
	Index      Node
	Arrayvalue Node
	Location
}

func (a ArrayIndexAccess) INode() {
	//empty method implements Node interface
}
func (a ArrayIndexAccess) StartPos() lexer.Position {
	return a.Location.Start
}
func (a ArrayIndexAccess) EndPos() lexer.Position {
	return a.Location.End
}

type StructLiteral struct {
	Name       IdentifierExpr
	Properties map[string]Node
	Location
}

func (a StructLiteral) INode() {
	//empty method implements Node interface
}
func (a StructLiteral) StartPos() lexer.Position {
	return a.Location.Start
}
func (a StructLiteral) EndPos() lexer.Position {
	return a.Location.End
}

type StructPropertyAccessExpr struct {
	Object   Node
	Property IdentifierExpr
	Location
}

func (a StructPropertyAccessExpr) INode() {
	//empty method implements Node interface
}
func (a StructPropertyAccessExpr) StartPos() lexer.Position {
	return a.Location.Start
}
func (a StructPropertyAccessExpr) EndPos() lexer.Position {
	return a.Location.End
}


type FunctionCallExpr struct {
	Name       IdentifierExpr
	Arguments  []Node
	Location
}

func (a FunctionCallExpr) INode() {
	//empty method implements Node interface
}
func (a FunctionCallExpr) StartPos() lexer.Position {
	return a.Location.Start
}
func (a FunctionCallExpr) EndPos() lexer.Position {
	return a.Location.End
}

type FunctionExpr struct {
	Params     []FunctionParam
	ReturnType DataType
	Block      BlockStmt
	Location
}

func (a FunctionExpr) INode() {
	//empty method implements Node interface
}

func (a FunctionExpr) StartPos() lexer.Position {
	return a.Location.Start
}

func (a FunctionExpr) EndPos() lexer.Position {
	return a.Location.End
}