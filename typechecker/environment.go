package typechecker

import (
	"fmt"
	"walrus/ast"
	"walrus/builtins"
)

type VALUE_TYPE string

const (
	INT_TYPE      	VALUE_TYPE = builtins.INT
	FLOAT_TYPE    	VALUE_TYPE = builtins.FLOAT
	CHAR_TYPE     	VALUE_TYPE = builtins.BYTE
	STRING_TYPE   	VALUE_TYPE = builtins.STRING
	BOOLEAN_TYPE  	VALUE_TYPE = builtins.BOOL
	NULL_TYPE     	VALUE_TYPE = builtins.NULL
	VOID_TYPE     	VALUE_TYPE = builtins.VOID
	FUNCTION_TYPE 	VALUE_TYPE = builtins.FUNCTION
	STRUCT_TYPE   	VALUE_TYPE = builtins.STRUCT
	TRAIT_TYPE		VALUE_TYPE = builtins.TRAIT
	ARRAY_TYPE    	VALUE_TYPE = builtins.ARRAY
	BLOCK_TYPE    	VALUE_TYPE = "block"
	RETURN_TYPE   	VALUE_TYPE = "return"
	USER_DEFINED_TYPE VALUE_TYPE = "user_defined"
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

type FnParam struct {
	Name string
	IsOptional bool
	//DefaultValueType ValueTypeInterface
	Type ValueTypeInterface
}

type Fn struct {
	DataType      VALUE_TYPE
	Params        []FnParam
	Returns       ValueTypeInterface
	FunctionScope TypeEnvironment
}

func (t Fn) DType() VALUE_TYPE {
	return t.DataType
}

type ConditionBranch struct {
	DataType VALUE_TYPE
	Next     ValueTypeInterface
	Returns  ValueTypeInterface
}

type ConditionStmt struct {
	DataType VALUE_TYPE
	Branches []ConditionBranch
}

func (t ConditionStmt) DType() VALUE_TYPE {
	return t.DataType
}

type StructProperty struct {
	IsPrivate bool
	Type      ValueTypeInterface
}

type StructMethod struct {
	IsPrivate bool
	Fn
}

type Struct struct {
	DataType   		VALUE_TYPE
	StructName 		string
	Elements   		map[string]StructProperty
	Methods			map[string]StructMethod
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
	DataType VALUE_TYPE
	TypeDef  ValueTypeInterface
}

func (t UserDefined) DType() VALUE_TYPE {
	return t.DataType
}

type ReturnType struct {
	DataType   VALUE_TYPE
	Expression ValueTypeInterface
}

func (t ReturnType) DType() VALUE_TYPE {
	return t.DataType
}

type Block struct {
	DataType VALUE_TYPE
	Returns  ValueTypeInterface
	Node     ast.Node
}

func (t Block) DType() VALUE_TYPE {
	return t.DataType
}

type Trait struct {
	DataType VALUE_TYPE
	TraitName string
	Methods   map[string]Fn
}

func (t Trait) DType() VALUE_TYPE {
	return t.DataType
}

type SCOPE_TYPE int

const (
	GLOBAL_SCOPE SCOPE_TYPE = iota
	FUNCTION_SCOPE
	CONDITIONAL_SCOPE
	LOOP_SCOPE
)

type TypeEnvironment struct {
	parent    	*TypeEnvironment
	scopeType 	SCOPE_TYPE
	scopeName 	string
	variables 	map[string]ValueTypeInterface
	constants 	map[string]bool
	isOptional 	map[string]bool
	types    	map[string]ValueTypeInterface
	traits  	map[string]Trait
	filePath  	string
}

func NewTypeENV(parent *TypeEnvironment, scope SCOPE_TYPE, scopeName string, filePath string) *TypeEnvironment {
	return &TypeEnvironment{
		parent:    parent,
		scopeType: scope,
		scopeName: scopeName,
		filePath:  filePath,
		variables: make(map[string]ValueTypeInterface),
		constants: make(map[string]bool),
		isOptional: make(map[string]bool),
		types:     make(map[string]ValueTypeInterface),
		traits:    make(map[string]Trait),
	}
}

func (t *TypeEnvironment) ResolveFunctionEnv() (*TypeEnvironment, error) {
	if t.scopeType == FUNCTION_SCOPE {
		return t, nil
	}
	if t.parent == nil {
		return nil, fmt.Errorf("function not found")
	}
	return t.parent.ResolveFunctionEnv()
}

func (t *TypeEnvironment) ResolveVar(name string) (*TypeEnvironment, error) {
	if t.isDeclared(name) {
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

func (t *TypeEnvironment) ResolveStruct(name string) (*TypeEnvironment, error) {
	value, err := t.ResolveType(name)
	if err != nil {
		return nil, err
	}
	if _, ok := value.types[name].(UserDefined).TypeDef.(Struct); !ok {
		return nil, fmt.Errorf("'%s' is not a struct", name)
	}
	return value, nil
}

func (t *TypeEnvironment) DeclareVar(name string, typeVar ValueTypeInterface, isConst bool, isOptional bool) error {
	//should not be declared
	if scope, err := t.ResolveVar(name); err == nil && scope == t {
		return fmt.Errorf("'%s' is already declared in this scope", name)
	}

	t.variables[name] = typeVar
	t.constants[name] = isConst
	t.isOptional[name] = isOptional

	return nil
}

func (t *TypeEnvironment) DeclareType(name string, typeType ValueTypeInterface) error {
	if scope, err := t.ResolveType(name); err == nil && scope == t {
		return err
	}
	t.types[name] = typeType
	return nil
}

func (t *TypeEnvironment) isDeclared(name string) bool {
	if _, ok := t.variables[name]; ok {
		return true
	}
	return false
}

func (t *TypeEnvironment) DeclareTrait(name string, trait Trait) error {
	if _, ok := t.traits[name]; ok {
		return fmt.Errorf("trait '%s' is already declared", name)
	}

	t.traits[name] = trait

	return nil
}

func (t *TypeEnvironment) ResolveTrait(name string) (*TypeEnvironment, error) {
	if _, ok := t.traits[name]; ok {
		return t, nil
	}
	if t.parent == nil {
		return nil, fmt.Errorf("trait '%s' is not defined", name)
	}
	return t.parent.ResolveTrait(name)
}