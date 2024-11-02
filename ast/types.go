package ast

import (
	"walrus/builtins"
	"walrus/lexer"
)

type DataType interface {
	Type() builtins.DATA_TYPE
	StartPos() lexer.Position
	EndPos() lexer.Position
}

type IntegerType struct {
	TypeName builtins.DATA_TYPE
	BitSize  uint8
	IsSigned bool
	Location
}

func (a IntegerType) Type() builtins.DATA_TYPE {
	return a.TypeName
}
func (a IntegerType) StartPos() lexer.Position {
	return a.Location.Start
}
func (a IntegerType) EndPos() lexer.Position {
	return a.Location.End
}

type FloatType struct {
	TypeName builtins.DATA_TYPE
	BitSize  uint8
	Location
}

func (a FloatType) Type() builtins.DATA_TYPE {
	return a.TypeName
}
func (a FloatType) StartPos() lexer.Position {
	return a.Location.Start
}
func (a FloatType) EndPos() lexer.Position {
	return a.Location.End
}

type StringType struct {
	TypeName builtins.DATA_TYPE
	Location
}

func (a StringType) Type() builtins.DATA_TYPE {
	return a.TypeName
}
func (a StringType) StartPos() lexer.Position {
	return a.Location.Start
}
func (a StringType) EndPos() lexer.Position {
	return a.Location.End
}

type BooleanType struct {
	TypeName builtins.DATA_TYPE
	Location
}

func (a BooleanType) Type() builtins.DATA_TYPE {
	return a.TypeName
}
func (a BooleanType) StartPos() lexer.Position {
	return a.Location.Start
}
func (a BooleanType) EndPos() lexer.Position {
	return a.Location.End
}

type NullType struct {
	TypeName builtins.DATA_TYPE
	Location
}

func (a NullType) Type() builtins.DATA_TYPE {
	return a.TypeName
}
func (a NullType) StartPos() lexer.Position {
	return a.Location.Start
}
func (a NullType) EndPos() lexer.Position {
	return a.Location.End
}

type VoidType struct {
	TypeName builtins.DATA_TYPE
	Location
}

func (a VoidType) Type() builtins.DATA_TYPE {
	return a.TypeName
}
func (a VoidType) StartPos() lexer.Position {
	return a.Location.Start
}
func (a VoidType) EndPos() lexer.Position {
	return a.Location.End
}

type ArrayType struct {
	TypeName  builtins.DATA_TYPE
	ArrayType DataType
	Location
}

func (a ArrayType) Type() builtins.DATA_TYPE {
	return a.TypeName
}
func (a ArrayType) StartPos() lexer.Position {
	return a.Location.Start
}
func (a ArrayType) EndPos() lexer.Position {
	return a.Location.End
}

type StructPropType struct {
	Prop      IdentifierExpr
	PropType  DataType
	IsPrivate bool
}

type StructType struct {
	TypeName   builtins.DATA_TYPE
	Properties map[string]StructPropType
	Location
}

func (a StructType) Type() builtins.DATA_TYPE {
	return a.TypeName
}
func (a StructType) StartPos() lexer.Position {
	return a.Location.Start
}
func (a StructType) EndPos() lexer.Position {
	return a.Location.End
}

type InterfaceMethod struct {
	Identifier IdentifierExpr
	FunctionType
}

type InterfaceType struct {
	TypeName builtins.DATA_TYPE
	Methods  map[string]InterfaceMethod
	Location
}

func (a InterfaceType) Type() builtins.DATA_TYPE {
	return a.TypeName
}

func (a InterfaceType) StartPos() lexer.Position {
	return a.Location.Start
}

func (a InterfaceType) EndPos() lexer.Position {
	return a.Location.End
}

type FunctionTypeParam struct {
	Identifier IdentifierExpr
	Type       DataType
	IsOptional bool
	Location
}

type FunctionType struct {
	TypeName   builtins.DATA_TYPE
	Parameters []FunctionTypeParam
	ReturnType DataType
	Location
}

func (a FunctionType) Type() builtins.DATA_TYPE {
	return a.TypeName
}
func (a FunctionType) StartPos() lexer.Position {
	return a.Location.Start
}
func (a FunctionType) EndPos() lexer.Position {
	return a.Location.End
}

type MapType struct {
	TypeName builtins.DATA_TYPE
	KeyType  DataType
	ValueType DataType
	Location
}

func (a MapType) Type() builtins.DATA_TYPE {
	return a.TypeName
}

func (a MapType) StartPos() lexer.Position {
	return a.Location.Start
}

func (a MapType) EndPos() lexer.Position {
	return a.Location.End
}

type UserDefinedType struct {
	TypeName builtins.DATA_TYPE
	Location
}

func (a UserDefinedType) Type() builtins.DATA_TYPE {
	return a.TypeName
}
func (a UserDefinedType) StartPos() lexer.Position {
	return a.Location.Start
}
func (a UserDefinedType) EndPos() lexer.Position {
	return a.Location.End
}
