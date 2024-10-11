package ast

import "walrus/lexer"

type DATA_TYPE string

type DataType interface {
	Type() DATA_TYPE
	StartPos() lexer.Position
	EndPos() lexer.Position
}

type IntegerType struct {
	TypeName DATA_TYPE
	Location
}

func (a IntegerType) Type() DATA_TYPE {
	return a.TypeName
}
func (a IntegerType) StartPos() lexer.Position {
	return a.Location.Start
}
func (a IntegerType) EndPos() lexer.Position {
	return a.Location.End
}

type FloatType struct {
	TypeName DATA_TYPE
	Location
}

func (a FloatType) Type() DATA_TYPE {
	return a.TypeName
}
func (a FloatType) StartPos() lexer.Position {
	return a.Location.Start
}
func (a FloatType) EndPos() lexer.Position {
	return a.Location.End
}

type StringType struct {
	TypeName DATA_TYPE
	Location
}

func (a StringType) Type() DATA_TYPE {
	return a.TypeName
}
func (a StringType) StartPos() lexer.Position {
	return a.Location.Start
}
func (a StringType) EndPos() lexer.Position {
	return a.Location.End
}

type CharType struct {
	TypeName DATA_TYPE
	Location
}

func (a CharType) Type() DATA_TYPE {
	return a.TypeName
}
func (a CharType) StartPos() lexer.Position {
	return a.Location.Start
}
func (a CharType) EndPos() lexer.Position {
	return a.Location.End
}

type BooleanType struct {
	TypeName DATA_TYPE
	Location
}

func (a BooleanType) Type() DATA_TYPE {
	return a.TypeName
}
func (a BooleanType) StartPos() lexer.Position {
	return a.Location.Start
}
func (a BooleanType) EndPos() lexer.Position {
	return a.Location.End
}

type NullType struct {
	TypeName DATA_TYPE
	Location
}

func (a NullType) Type() DATA_TYPE {
	return a.TypeName
}
func (a NullType) StartPos() lexer.Position {
	return a.Location.Start
}
func (a NullType) EndPos() lexer.Position {
	return a.Location.End
}

type VoidType struct {
	TypeName DATA_TYPE
	Location
}

func (a VoidType) Type() DATA_TYPE {
	return a.TypeName
}
func (a VoidType) StartPos() lexer.Position {
	return a.Location.Start
}
func (a VoidType) EndPos() lexer.Position {
	return a.Location.End
}

type ArrayType struct {
	TypeName  DATA_TYPE
	ArrayType DataType
	Location
}

func (a ArrayType) Type() DATA_TYPE {
	return a.TypeName
}
func (a ArrayType) StartPos() lexer.Position {
	return a.Location.Start
}
func (a ArrayType) EndPos() lexer.Position {
	return a.Location.End
}

type StructPropType struct {
	Prop		IdentifierExpr
	PropType	DataType
	IsPrivate	bool
}

type StructType struct {
	TypeName 	DATA_TYPE
	Properties	map[string]StructPropType
	Location
}
func (a StructType) Type() DATA_TYPE {
	return a.TypeName
}
func (a StructType) StartPos() lexer.Position {
	return a.Location.Start
}
func (a StructType) EndPos() lexer.Position {
	return a.Location.End
}

type FunctionTypeParam struct {
	Identifier 	IdentifierExpr
	Type       	DataType
	IsOptional	bool
	Location
}

type FunctionType struct {
	TypeName       	DATA_TYPE
	Parameters 		[]FunctionTypeParam
	ReturnType 		DataType
	Location
}
func (a FunctionType) Type() DATA_TYPE {
	return a.TypeName
}
func (a FunctionType) StartPos() lexer.Position {
	return a.Location.Start
}
func (a FunctionType) EndPos() lexer.Position {
	return a.Location.End
}

type UserDefinedType struct {
	TypeName	DATA_TYPE
	Location
}
func (a UserDefinedType) Type() DATA_TYPE {
	return a.TypeName
}
func (a UserDefinedType) StartPos() lexer.Position {
	return a.Location.Start
}
func (a UserDefinedType) EndPos() lexer.Position {
	return a.Location.End
}