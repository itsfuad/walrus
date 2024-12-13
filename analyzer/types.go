package analyzer

import (
	"walrus/ast"
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

type TcValue interface {
	DType() builtins.TC_TYPE
}

type Int struct {
	DataType builtins.TC_TYPE
	BitSize  uint8
	IsSigned bool
}

func (t Int) DType() builtins.TC_TYPE {
	return t.DataType
}

type Float struct {
	DataType builtins.TC_TYPE
	BitSize  uint8
}

func (t Float) DType() builtins.TC_TYPE {
	return t.DataType
}

type Str struct {
	DataType builtins.TC_TYPE
}

func (t Str) DType() builtins.TC_TYPE {
	return t.DataType
}

type Bool struct {
	DataType builtins.TC_TYPE
}

func (t Bool) DType() builtins.TC_TYPE {
	return t.DataType
}

type Null struct {
	DataType builtins.TC_TYPE
}

func (t Null) DType() builtins.TC_TYPE {
	return t.DataType
}

type Void struct {
	DataType builtins.TC_TYPE
}

func (t Void) DType() builtins.TC_TYPE {
	return t.DataType
}

type Map struct {
	DataType  builtins.TC_TYPE
	KeyType   TcValue
	ValueType TcValue
}

func (t Map) DType() builtins.TC_TYPE {
	return t.DataType
}

type FnParam struct {
	Name       string
	IsOptional bool
	//DefaultValueType ValueTypeInterface
	Type TcValue
}

type Fn struct {
	DataType      builtins.TC_TYPE
	Params        []FnParam
	Returns       TcValue
	FunctionScope TypeEnvironment
}

func (t Fn) DType() builtins.TC_TYPE {
	return t.DataType
}

type ConditionBranch struct {
	DataType builtins.TC_TYPE
	Next     TcValue
	Returns  TcValue
}

type ConditionStmt struct {
	DataType builtins.TC_TYPE
	Branches []ConditionBranch
}

func (t ConditionStmt) DType() builtins.TC_TYPE {
	return t.DataType
}

type StructProperty struct {
	IsPrivate bool
	Type      TcValue
}

func (t StructProperty) DType() builtins.TC_TYPE {
	return t.Type.DType()
}

type StructMethod struct {
	IsPrivate bool
	Fn
}

func (t StructMethod) DType() builtins.TC_TYPE {
	return t.DataType
}

type Struct struct {
	DataType    builtins.TC_TYPE
	StructName  string
	StructScope TypeEnvironment
}

func (t Struct) DType() builtins.TC_TYPE {
	return t.DataType
}

type Array struct {
	DataType  builtins.TC_TYPE
	ArrayType TcValue
}

func (t Array) DType() builtins.TC_TYPE {
	return t.DataType
}

type UserDefined struct {
	DataType builtins.TC_TYPE
	TypeName string
	TypeDef  TcValue
}

func (t UserDefined) DType() builtins.TC_TYPE {
	return t.DataType
}

type ReturnType struct {
	DataType   builtins.TC_TYPE
	Expression TcValue
}

func (t ReturnType) DType() builtins.TC_TYPE {
	return t.DataType
}

type Block struct {
	DataType builtins.TC_TYPE
	Returns  TcValue
	Node     ast.Node
}

func (t Block) DType() builtins.TC_TYPE {
	return t.DataType
}

type InterfaceMethod struct {
	Name 	 	string
	Method		Fn
}

type Interface struct {
	DataType      builtins.TC_TYPE
	InterfaceName string
	Methods       []InterfaceMethod
}

func (t Interface) DType() builtins.TC_TYPE {
	return t.DataType
}

type Maybe struct {
	DataType  builtins.TC_TYPE
	MaybeType TcValue
}

func (t Maybe) DType() builtins.TC_TYPE {
	return t.DataType
}
