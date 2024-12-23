package ast

import (
	"walrus/frontend/builtins"
	"walrus/position"
)

type DataType interface {
	Type() builtins.PARSER_TYPE
	StartPos() position.Coordinate
	EndPos() position.Coordinate
}

type IntegerType struct {
	TypeName builtins.PARSER_TYPE
	BitSize  uint8
	IsSigned bool
	position.Location
}

func (a IntegerType) Type() builtins.PARSER_TYPE {
	return a.TypeName
}
func (a IntegerType) StartPos() position.Coordinate {
	return a.Location.Start
}
func (a IntegerType) EndPos() position.Coordinate {
	return a.Location.End
}

type FloatType struct {
	TypeName builtins.PARSER_TYPE
	BitSize  uint8
	position.Location
}

func (a FloatType) Type() builtins.PARSER_TYPE {
	return a.TypeName
}
func (a FloatType) StartPos() position.Coordinate {
	return a.Location.Start
}
func (a FloatType) EndPos() position.Coordinate {
	return a.Location.End
}

type StringType struct {
	TypeName builtins.PARSER_TYPE
	position.Location
}

func (a StringType) Type() builtins.PARSER_TYPE {
	return a.TypeName
}
func (a StringType) StartPos() position.Coordinate {
	return a.Location.Start
}
func (a StringType) EndPos() position.Coordinate {
	return a.Location.End
}

type BooleanType struct {
	TypeName builtins.PARSER_TYPE
	position.Location
}

func (a BooleanType) Type() builtins.PARSER_TYPE {
	return a.TypeName
}
func (a BooleanType) StartPos() position.Coordinate {
	return a.Location.Start
}
func (a BooleanType) EndPos() position.Coordinate {
	return a.Location.End
}

type NullType struct {
	TypeName builtins.PARSER_TYPE
	position.Location
}

func (a NullType) Type() builtins.PARSER_TYPE {
	return a.TypeName
}
func (a NullType) StartPos() position.Coordinate {
	return a.Location.Start
}
func (a NullType) EndPos() position.Coordinate {
	return a.Location.End
}

type VoidType struct {
	TypeName builtins.PARSER_TYPE
	position.Location
}

func (a VoidType) Type() builtins.PARSER_TYPE {
	return a.TypeName
}
func (a VoidType) StartPos() position.Coordinate {
	return a.Location.Start
}
func (a VoidType) EndPos() position.Coordinate {
	return a.Location.End
}

type ArrayType struct {
	TypeName  builtins.PARSER_TYPE
	ArrayType DataType
	position.Location
}

func (a ArrayType) Type() builtins.PARSER_TYPE {
	return a.TypeName
}
func (a ArrayType) StartPos() position.Coordinate {
	return a.Location.Start
}
func (a ArrayType) EndPos() position.Coordinate {
	return a.Location.End
}

type StructPropType struct {
	Prop      IdentifierExpr
	PropType  DataType
	IsPrivate bool
}

type StructType struct {
	TypeName   builtins.PARSER_TYPE
	Properties []StructPropType
	position.Location
}

func (a StructType) Type() builtins.PARSER_TYPE {
	return a.TypeName
}
func (a StructType) StartPos() position.Coordinate {
	return a.Location.Start
}
func (a StructType) EndPos() position.Coordinate {
	return a.Location.End
}

type InterfaceMethod struct {
	Identifier IdentifierExpr
	FunctionType
}

type InterfaceType struct {
	TypeName builtins.PARSER_TYPE
	Methods  []InterfaceMethod
	position.Location
}

func (a InterfaceType) Type() builtins.PARSER_TYPE {
	return a.TypeName
}

func (a InterfaceType) StartPos() position.Coordinate {
	return a.Location.Start
}

func (a InterfaceType) EndPos() position.Coordinate {
	return a.Location.End
}

type FunctionTypeParam struct {
	Identifier IdentifierExpr
	Type       DataType
	IsOptional bool
	position.Location
}

type FunctionType struct {
	TypeName   builtins.PARSER_TYPE
	Parameters []FunctionTypeParam
	ReturnType DataType
	position.Location
}

func (a FunctionType) Type() builtins.PARSER_TYPE {
	return a.TypeName
}
func (a FunctionType) StartPos() position.Coordinate {
	return a.Location.Start
}
func (a FunctionType) EndPos() position.Coordinate {
	return a.Location.End
}

type MaybeType struct {
	TypeName  builtins.PARSER_TYPE
	MaybeType DataType
	position.Location
}

func (a MaybeType) Type() builtins.PARSER_TYPE {
	return a.TypeName
}

func (a MaybeType) StartPos() position.Coordinate {
	return a.Location.Start
}

func (a MaybeType) EndPos() position.Coordinate {
	return a.Location.End
}

type MapType struct {
	TypeName  builtins.PARSER_TYPE
	Map       IdentifierExpr
	KeyType   DataType
	ValueType DataType
	position.Location
}

func (a MapType) Type() builtins.PARSER_TYPE {
	return a.TypeName
}

func (a MapType) StartPos() position.Coordinate {
	return a.Location.Start
}

func (a MapType) EndPos() position.Coordinate {
	return a.Location.End
}

type UserDefinedType struct {
	TypeName  builtins.PARSER_TYPE
	AliasName string
	position.Location
}

func (a UserDefinedType) Type() builtins.PARSER_TYPE {
	return a.TypeName
}
func (a UserDefinedType) StartPos() position.Coordinate {
	return a.Location.Start
}
func (a UserDefinedType) EndPos() position.Coordinate {
	return a.Location.End
}
