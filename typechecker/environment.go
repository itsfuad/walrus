package typechecker

import (
	"fmt"
	"walrus/builtins"
)

type VALUE_TYPE string

const (
	INT_TYPE      VALUE_TYPE = builtins.INT
	FLOAT_TYPE    VALUE_TYPE = builtins.FLOAT
	CHAR_TYPE     VALUE_TYPE = builtins.CHAR
	STRING_TYPE   VALUE_TYPE = builtins.STRING
	BOOLEAN_TYPE  VALUE_TYPE = builtins.BOOL
	NULL_TYPE     VALUE_TYPE = builtins.NULL
	VOID_TYPE     VALUE_TYPE = builtins.VOID
	FUNCTION_TYPE VALUE_TYPE = builtins.FUNCTION
	STRUCT_TYPE   VALUE_TYPE = builtins.STRUCT
	ARRAY_TYPE    VALUE_TYPE = builtins.ARRAY
)

type ValueTypeInterface interface {
	DType() VALUE_TYPE
}

type Int struct {
	DataType VALUE_TYPE
}

func (t Int) DType() VALUE_TYPE {
	return t.DataType
}

type Float struct {
	DataType VALUE_TYPE
}

func (t Float) DType() VALUE_TYPE {
	return t.DataType
}

type Chr struct {
	DataType VALUE_TYPE
}

func (t Chr) DType() VALUE_TYPE {
	return t.DataType
}

type Str struct {
	DataType VALUE_TYPE
}

func (t Str) DType() VALUE_TYPE {
	return t.DataType
}

type Bool struct {
	DataType VALUE_TYPE
}

func (t Bool) DType() VALUE_TYPE {
	return t.DataType
}

type Null struct {
	DataType VALUE_TYPE
}

func (t Null) DType() VALUE_TYPE {
	return t.DataType
}

type Void struct {
	DataType VALUE_TYPE
}

func (t Void) DType() VALUE_TYPE {
	return t.DataType
}

type Fn struct {
	DataType VALUE_TYPE
	Params   []ValueTypeInterface
	Returns  ValueTypeInterface
}

func (t Fn) DType() VALUE_TYPE {
	return t.DataType
}

type StructProperty struct {
	IsPrivate	bool
	Type		ValueTypeInterface
}

type Struct struct {
	DataType   VALUE_TYPE
	StructName string
	Elements   map[string]StructProperty
}

func (t Struct) DType() VALUE_TYPE {
	return t.DataType
}

type Array struct {
	DataType  VALUE_TYPE
	ArrayType ValueTypeInterface
}

func (t Array) DType() VALUE_TYPE {
	return t.DataType
}

type UserDefined struct {
	DataType	VALUE_TYPE
	TypeDef		ValueTypeInterface
}

func (t UserDefined) DType() VALUE_TYPE {
	return t.DataType
}

type TypeEnvironment struct {
	parent    	*TypeEnvironment
	variables 	map[string]ValueTypeInterface
	constants 	map[string]bool
	types		map[string]ValueTypeInterface
	filePath  	string
}

func NewTypeENV(parent *TypeEnvironment, filePath string) *TypeEnvironment {
	return &TypeEnvironment{
		parent:    	parent,
		filePath:  	filePath,
		variables: 	make(map[string]ValueTypeInterface),
		constants: 	make(map[string]bool),
		types: 		make(map[string]ValueTypeInterface),
	}
}

func (t *TypeEnvironment) ResolveVar(name string) (*TypeEnvironment, error) {
	if _, ok := t.variables[name]; ok {
		return t, nil
	}

	//check on the parent then
	if t.parent == nil {
		//no where is declared
		return nil, fmt.Errorf("'%s' is not declared in this scope", name)
	}

	return t.parent.ResolveVar(name)
}

func (t *TypeEnvironment) ResolveType(name string) (*TypeEnvironment, error) {
	if _, ok := t.types[name]; ok {
		return t, nil
	}
	if t.parent == nil {
		return nil, fmt.Errorf("type '%s' is not defined", name)
	}
	return t.parent.ResolveType(name)
}

func (t *TypeEnvironment) DeclareVar(name string, typeVar ValueTypeInterface, isConst bool) error {
	//should not be declared
	if _, err := t.ResolveVar(name); err == nil {
		return err
	}

	t.variables[name] = typeVar
	t.constants[name] = isConst

	return nil
}

func (t *TypeEnvironment) DeclareType(name string, typeType ValueTypeInterface) error {
	if _, err := t.ResolveType(name); err == nil {
		return err
	}
	t.types[name] = typeType
	return nil
}

func (t *TypeEnvironment) DeclareStruct(name string, typVar ValueTypeInterface) error {
	// will implement later
	return nil
}
