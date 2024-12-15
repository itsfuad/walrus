package analyzer

import (
	"walrus/builtins"
)

const (
	INT8_TYPE         builtins.TC_TYPE = builtins.INT8
	INT16_TYPE        builtins.TC_TYPE = builtins.INT16
	INT32_TYPE        builtins.TC_TYPE = builtins.INT32
	INT64_TYPE        builtins.TC_TYPE = builtins.INT64
	FLOAT32_TYPE      builtins.TC_TYPE = builtins.FLOAT32
	FLOAT64_TYPE      builtins.TC_TYPE = builtins.FLOAT64
	UINT8_TYPE        builtins.TC_TYPE = builtins.UINT8
	UINT16_TYPE       builtins.TC_TYPE = builtins.UINT16
	UINT32_TYPE       builtins.TC_TYPE = builtins.UINT32
	UINT64_TYPE       builtins.TC_TYPE = builtins.UINT64
	STRING_TYPE       builtins.TC_TYPE = builtins.STRING
	BYTE_TYPE         builtins.TC_TYPE = builtins.BYTE
	BOOLEAN_TYPE      builtins.TC_TYPE = builtins.BOOL
	NULL_TYPE         builtins.TC_TYPE = builtins.NULL
	VOID_TYPE         builtins.TC_TYPE = builtins.VOID
	FUNCTION_TYPE     builtins.TC_TYPE = builtins.FUNCTION
	STRUCT_TYPE       builtins.TC_TYPE = builtins.STRUCT
	INTERFACE_TYPE    builtins.TC_TYPE = builtins.INTERFACE
	ARRAY_TYPE        builtins.TC_TYPE = builtins.ARRAY
	MAP_TYPE          builtins.TC_TYPE = builtins.MAP
	MAYBE_TYPE        builtins.TC_TYPE = builtins.MAYBE
	USER_DEFINED_TYPE builtins.TC_TYPE = builtins.USER_DEFINED
	BLOCK_TYPE        builtins.TC_TYPE = "block"
	RETURN_TYPE       builtins.TC_TYPE = "return"
)

type AnalyzerNode interface {
	INodeResult()
}

type ExprType interface {
	DType() builtins.TC_TYPE
	INodeResult()
}

type Int struct {
	DataType builtins.TC_TYPE
	BitSize  uint8
	IsSigned bool
}

func (t Int) DType() builtins.TC_TYPE {
	return t.DataType
}

func (t Int) INodeResult() {} //empty function to implement interface

type Float struct {
	DataType builtins.TC_TYPE
	BitSize  uint8
}

func (t Float) DType() builtins.TC_TYPE {
	return t.DataType
}

func (t Float) INodeResult() {} //empty function to implement interface

type Str struct {
	DataType builtins.TC_TYPE
}

func (t Str) DType() builtins.TC_TYPE {
	return t.DataType
}

func (t Str) INodeResult() {} //empty function to implement interface

type Bool struct {
	DataType builtins.TC_TYPE
}

func (t Bool) DType() builtins.TC_TYPE {
	return t.DataType
}

func (t Bool) INodeResult() {} //empty function to implement interface

type Null struct {
	DataType builtins.TC_TYPE
}

func (t Null) DType() builtins.TC_TYPE {
	return t.DataType
}

func (t Null) INodeResult() {} //empty function to implement interface

type Void struct {
	DataType builtins.TC_TYPE
}

func (t Void) DType() builtins.TC_TYPE {
	return t.DataType
}

func (t Void) INodeResult() {} //empty function to implement interface

type Map struct {
	DataType  builtins.TC_TYPE
	KeyType   ExprType
	ValueType ExprType
}

func (t Map) DType() builtins.TC_TYPE {
	return t.DataType
}

func (t Map) INodeResult() {} //empty function to implement interface

type FnParam struct {
	Name       string
	IsOptional bool
	//DefaultValueType ValueTypeInterface
	Type ExprType
}

type Fn struct {
	DataType      builtins.TC_TYPE
	Params        []FnParam
	Returns       ExprType
	FunctionScope TypeEnvironment
}

func (t Fn) DType() builtins.TC_TYPE {
	return t.DataType
}

func (t Fn) INodeResult() {} //empty function to implement interface

type ConditionBranch struct {
	DataType builtins.TC_TYPE
	Next     ExprType
	Returns  ExprType
}

type ConditionStmt struct {
	DataType builtins.TC_TYPE
	Branches []ConditionBranch
}

func (t ConditionStmt) DType() builtins.TC_TYPE {
	return t.DataType
}

func (t ConditionStmt) INodeResult() {} //empty function to implement interface

type StructProperty struct {
	IsPrivate bool
	Type      ExprType
}

func (t StructProperty) DType() builtins.TC_TYPE {
	return t.Type.DType()
}

func (t StructProperty) INodeResult() {} //empty function to implement interface

type StructMethod struct {
	IsPrivate bool
	Fn
}

func (t StructMethod) DType() builtins.TC_TYPE {
	return t.DataType
}

func (t StructMethod) INodeResult() {} //empty function to implement interface

type Struct struct {
	DataType    builtins.TC_TYPE
	StructName  string
	StructScope TypeEnvironment
}

func (t Struct) DType() builtins.TC_TYPE {
	return t.DataType
}

func (t Struct) INodeResult() {} //empty function to implement interface

type Array struct {
	DataType  builtins.TC_TYPE
	ArrayType ExprType
}

func (t Array) DType() builtins.TC_TYPE {
	return t.DataType
}

func (t Array) INodeResult() {} //empty function to implement interface

type UserDefined struct {
	DataType builtins.TC_TYPE
	TypeName string
	TypeDef  ExprType
}

func (t UserDefined) DType() builtins.TC_TYPE {
	return t.DataType
}

func (t UserDefined) INodeResult() {} //empty function to implement interface

type ReturnType struct {
	DataType   builtins.TC_TYPE
	Expression ExprType
}

func (t ReturnType) DType() builtins.TC_TYPE {
	return t.DataType
}

func (t ReturnType) INodeResult() {} //empty function to implement interface

type InterfaceMethodType struct {
	Name   string
	Method Fn
}

type Interface struct {
	DataType      builtins.TC_TYPE
	InterfaceName string
	Methods       []InterfaceMethodType
}

func (t Interface) DType() builtins.TC_TYPE {
	return t.DataType
}

func (t Interface) INodeResult() {} //empty function to implement interface

type Maybe struct {
	DataType  builtins.TC_TYPE
	MaybeType ExprType
}

func (t Maybe) DType() builtins.TC_TYPE {
	return t.DataType
}

func (t Maybe) INodeResult() {} //empty function to implement interface
