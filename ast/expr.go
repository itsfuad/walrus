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

type IncrementalInterface interface {
	Arg() IdentifierExpr
	Op() lexer.Token
}

type PrefixExpr struct {
	Operator lexer.Token
	Argument IdentifierExpr
	Location
}

func (a PrefixExpr) INode() {
	//empty method implements Node interface
}
func (a PrefixExpr) StartPos() lexer.Position {
	return a.Location.Start
}
func (a PrefixExpr) EndPos() lexer.Position {
	return a.Location.End
}
func (a PrefixExpr) Arg() IdentifierExpr {
	return a.Argument
}
func (a PrefixExpr) Op() lexer.Token {
	return a.Operator
}

type PostfixExpr struct {
	Operator lexer.Token
	Argument IdentifierExpr
	Location
}

func (a PostfixExpr) INode() {
	//empty method implements Node interface
}
func (a PostfixExpr) StartPos() lexer.Position {
	return a.Location.Start
}
func (a PostfixExpr) EndPos() lexer.Position {
	return a.Location.End
}
func (a PostfixExpr) Arg() IdentifierExpr {
	return a.Argument
}
func (a PostfixExpr) Op() lexer.Token {
	return a.Operator
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

type ArrayLiteral struct {
	Values []Node
	Location
}

func (a ArrayLiteral) INode() {
	//empty method implements Node interface
}
func (a ArrayLiteral) StartPos() lexer.Position {
	return a.Location.Start
}
func (a ArrayLiteral) EndPos() lexer.Position {
	return a.Location.End
}

type ArrayIndexAccess struct {
	Index Node
	Array Node
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
	Identifier IdentifierExpr
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
	Caller    Node
	Arguments []Node
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

type FunctionLiteral struct {
	Params     []FunctionParam
	Body       BlockStmt
	ReturnType DataType
	Location
}

func (a FunctionLiteral) INode() {
	//empty method implements Node interface
}

func (a FunctionLiteral) StartPos() lexer.Position {
	return a.Location.Start
}

func (a FunctionLiteral) EndPos() lexer.Position {
	return a.Location.End
}
