package ast

import (
	"walrus/frontend/lexer"
	"walrus/position"
)

// Any word which is not a keyword or literal
type IdentifierExpr struct {
	Name string
	position.Location
}

func (a IdentifierExpr) INode() {
	//empty method implements Node interface
}
func (a IdentifierExpr) StartPos() position.Coordinate {
	return a.Location.Start
}
func (a IdentifierExpr) EndPos() position.Coordinate {
	return a.Location.End
}

// Literals or Raw values like: 1,2,3,4.6, "hello world", 'a' ...etc
type IntegerLiteralExpr struct {
	Value    string
	BitSize  uint8
	IsSigned bool
	position.Location
}

func (a IntegerLiteralExpr) INode() {
	//empty method implements Node interface
}
func (a IntegerLiteralExpr) StartPos() position.Coordinate {
	return a.Location.Start
}
func (a IntegerLiteralExpr) EndPos() position.Coordinate {
	return a.Location.End
}

type FloatLiteralExpr struct {
	Value   string
	BitSize uint8
	position.Location
}

func (a FloatLiteralExpr) INode() {
	//empty method implements Node interface
}
func (a FloatLiteralExpr) StartPos() position.Coordinate {
	return a.Location.Start
}
func (a FloatLiteralExpr) EndPos() position.Coordinate {
	return a.Location.End
}

type StringLiteralExpr struct {
	Value string
	position.Location
}

func (a StringLiteralExpr) INode() {
	//empty method implements Node interface
}
func (a StringLiteralExpr) StartPos() position.Coordinate {
	return a.Location.Start
}
func (a StringLiteralExpr) EndPos() position.Coordinate {
	return a.Location.End
}

type ByteLiteralExpr struct {
	Value string
	position.Location
}

func (a ByteLiteralExpr) INode() {
	//empty method implements Node interface
}
func (a ByteLiteralExpr) StartPos() position.Coordinate {
	return a.Location.Start
}
func (a ByteLiteralExpr) EndPos() position.Coordinate {
	return a.Location.End
}

type BooleanLiteralExpr struct {
	Value string
	position.Location
}

func (a BooleanLiteralExpr) INode() {
	//empty method implements Node interface
}
func (a BooleanLiteralExpr) StartPos() position.Coordinate {
	return a.Location.Start
}
func (a BooleanLiteralExpr) EndPos() position.Coordinate {
	return a.Location.End
}

type NullLiteralExpr struct {
	Value string
	position.Location
}

func (a NullLiteralExpr) INode() {
	//empty method implements Node interface
}
func (a NullLiteralExpr) StartPos() position.Coordinate {
	return a.Location.Start
}
func (a NullLiteralExpr) EndPos() position.Coordinate {
	return a.Location.End
}

type MapProp struct {
	Key   Node
	Value Node
}

type MapLiteral struct {
	MapType
	Values []MapProp
	position.Location
}

func (a MapLiteral) INode() {
	//empty method implements Node interface
}
func (a MapLiteral) StartPos() position.Coordinate {
	return a.Location.Start
}
func (a MapLiteral) EndPos() position.Coordinate {
	return a.Location.End
}

type UnaryExpr struct {
	Operator lexer.Token
	Argument Node
	position.Location
}

func (a UnaryExpr) INode() {
	//empty method implements Node interface
}
func (a UnaryExpr) StartPos() position.Coordinate {
	return a.Location.Start
}
func (a UnaryExpr) EndPos() position.Coordinate {
	return a.Location.End
}

type TypeCastExpr struct {
	Expression Node
	ToCast     DataType
	position.Location
}

func (a TypeCastExpr) INode() {
	//empty method implements Node interface
}
func (a TypeCastExpr) StartPos() position.Coordinate {
	return a.Location.Start
}

func (a TypeCastExpr) EndPos() position.Coordinate {
	return a.Location.End
}

type TypeofExpr struct {
	Expression Node
	position.Location
}

func (a TypeofExpr) INode() {
	//empty method implements Node interface
}
func (a TypeofExpr) StartPos() position.Coordinate {
	return a.Location.Start
}

func (a TypeofExpr) EndPos() position.Coordinate {
	return a.Location.End
}

type BinaryExpr struct {
	Binop lexer.Token
	Left  Node
	Right Node
	position.Location
}

type IncrementalInterface interface {
	Arg() IdentifierExpr
	Op() lexer.Token
}

type PrefixExpr struct {
	OP       lexer.Token
	Argument IdentifierExpr
	position.Location
}

func (a PrefixExpr) INode() {
	//empty method implements Node interface
}
func (a PrefixExpr) StartPos() position.Coordinate {
	return a.Location.Start
}
func (a PrefixExpr) EndPos() position.Coordinate {
	return a.Location.End
}
func (a PrefixExpr) Arg() IdentifierExpr {
	return a.Argument
}
func (a PrefixExpr) Op() lexer.Token {
	return a.OP
}

type PostfixExpr struct {
	Operator lexer.Token
	Argument IdentifierExpr
	position.Location
}

func (a PostfixExpr) INode() {
	//empty method implements Node interface
}
func (a PostfixExpr) StartPos() position.Coordinate {
	return a.Location.Start
}
func (a PostfixExpr) EndPos() position.Coordinate {
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
func (a BinaryExpr) StartPos() position.Coordinate {
	return a.Location.Start
}
func (a BinaryExpr) EndPos() position.Coordinate {
	return a.Location.End
}

type VarAssignmentExpr struct {
	Assignee Node // Check later if we should use IdentifierExpr instead
	Value    Node
	Operator lexer.Token // Looks odd right? Well, we know the operator must be '='. But what about +=, -=, *= and so on?😀
	position.Location
}

func (a VarAssignmentExpr) INode() {
	//empty method implements Node interface
}
func (a VarAssignmentExpr) StartPos() position.Coordinate {
	return a.Location.Start
}
func (a VarAssignmentExpr) EndPos() position.Coordinate {
	return a.Location.End
}

type ArrayLiteral struct {
	Values []Node
	position.Location
}

func (a ArrayLiteral) INode() {
	//empty method implements Node interface
}
func (a ArrayLiteral) StartPos() position.Coordinate {
	return a.Location.Start
}
func (a ArrayLiteral) EndPos() position.Coordinate {
	return a.Location.End
}

type Indexable struct {
	Index     Node
	Container Node
	position.Location
}

func (a Indexable) INode() {
	//empty method implements Node interface
}
func (a Indexable) StartPos() position.Coordinate {
	return a.Location.Start
}
func (a Indexable) EndPos() position.Coordinate {
	return a.Location.End
}

type StructProp struct {
	Prop  IdentifierExpr
	Value Node
}

type StructLiteral struct {
	Identifier IdentifierExpr
	Properties []StructProp
	position.Location
}

func (a StructLiteral) INode() {
	//empty method implements Node interface
}
func (a StructLiteral) StartPos() position.Coordinate {
	return a.Location.Start
}
func (a StructLiteral) EndPos() position.Coordinate {
	return a.Location.End
}

type StructPropertyAccessExpr struct {
	Object   Node
	Property IdentifierExpr
	position.Location
}

func (a StructPropertyAccessExpr) INode() {
	//empty method implements Node interface
}
func (a StructPropertyAccessExpr) StartPos() position.Coordinate {
	return a.Location.Start
}
func (a StructPropertyAccessExpr) EndPos() position.Coordinate {
	return a.Location.End
}

type FunctionCallExpr struct {
	Caller    Node
	Arguments []Node
	position.Location
}

func (a FunctionCallExpr) INode() {
	//empty method implements Node interface
}
func (a FunctionCallExpr) StartPos() position.Coordinate {
	return a.Location.Start
}
func (a FunctionCallExpr) EndPos() position.Coordinate {
	return a.Location.End
}

type FunctionLiteral struct {
	Params     []FunctionParam
	Body       BlockStmt
	ReturnType DataType
	position.Location
}

func (a FunctionLiteral) INode() {
	//empty method implements Node interface
}

func (a FunctionLiteral) StartPos() position.Coordinate {
	return a.Location.Start
}

func (a FunctionLiteral) EndPos() position.Coordinate {
	return a.Location.End
}
