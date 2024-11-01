package typechecker

import (
	"fmt"
	"walrus/ast"
	"walrus/builtins"
	"walrus/errgen"
)

const (
	INT8_TYPE         builtins.VALUE_TYPE = builtins.INT8
	INT16_TYPE        builtins.VALUE_TYPE = builtins.INT16
	INT32_TYPE        builtins.VALUE_TYPE = builtins.INT32
	INT64_TYPE        builtins.VALUE_TYPE = builtins.INT64
	FLOAT32_TYPE      builtins.VALUE_TYPE = builtins.FLOAT32
	FLOAT64_TYPE      builtins.VALUE_TYPE = builtins.FLOAT64
	UINT8_TYPE        builtins.VALUE_TYPE = builtins.UINT8
	UINT16_TYPE       builtins.VALUE_TYPE = builtins.UINT16
	UINT32_TYPE       builtins.VALUE_TYPE = builtins.UINT32
	UINT64_TYPE       builtins.VALUE_TYPE = builtins.UINT64
	STRING_TYPE       builtins.VALUE_TYPE = builtins.STRING
	BYTE_TYPE         builtins.VALUE_TYPE = builtins.BYTE
	BOOLEAN_TYPE      builtins.VALUE_TYPE = builtins.BOOL
	NULL_TYPE         builtins.VALUE_TYPE = builtins.NULL
	VOID_TYPE         builtins.VALUE_TYPE = builtins.VOID
	FUNCTION_TYPE     builtins.VALUE_TYPE = builtins.FUNCTION
	STRUCT_TYPE       builtins.VALUE_TYPE = builtins.STRUCT
	INTERFACE_TYPE    builtins.VALUE_TYPE = builtins.INTERFACE
	ARRAY_TYPE        builtins.VALUE_TYPE = builtins.ARRAY
	BLOCK_TYPE        builtins.VALUE_TYPE = "block"
	RETURN_TYPE       builtins.VALUE_TYPE = "return"
	USER_DEFINED_TYPE builtins.VALUE_TYPE = "user_defined"
)

type ValueTypeInterface interface {
	DType() builtins.VALUE_TYPE
}

type Int struct {
	DataType builtins.VALUE_TYPE
	BitSize  uint8
	IsSigned bool
}

func (t Int) DType() builtins.VALUE_TYPE {
	return t.DataType
}

type Float struct {
	DataType builtins.VALUE_TYPE
	BitSize  uint8
}

func (t Float) DType() builtins.VALUE_TYPE {
	return t.DataType
}

type Str struct {
	DataType builtins.VALUE_TYPE
}

func (t Str) DType() builtins.VALUE_TYPE {
	return t.DataType
}

type Bool struct {
	DataType builtins.VALUE_TYPE
}

func (t Bool) DType() builtins.VALUE_TYPE {
	return t.DataType
}

type Null struct {
	DataType builtins.VALUE_TYPE
}

func (t Null) DType() builtins.VALUE_TYPE {
	return t.DataType
}

type Void struct {
	DataType builtins.VALUE_TYPE
}

func (t Void) DType() builtins.VALUE_TYPE {
	return t.DataType
}

type FnParam struct {
	Name       string
	IsOptional bool
	//DefaultValueType ValueTypeInterface
	Type ValueTypeInterface
}

type Fn struct {
	DataType      builtins.VALUE_TYPE
	Params        []FnParam
	Returns       ValueTypeInterface
	FunctionScope TypeEnvironment
}

func (t Fn) DType() builtins.VALUE_TYPE {
	return t.DataType
}

type ConditionBranch struct {
	DataType builtins.VALUE_TYPE
	Next     ValueTypeInterface
	Returns  ValueTypeInterface
}

type ConditionStmt struct {
	DataType builtins.VALUE_TYPE
	Branches []ConditionBranch
}

func (t ConditionStmt) DType() builtins.VALUE_TYPE {
	return t.DataType
}

type StructProperty struct {
	IsPrivate bool
	Type      ValueTypeInterface
}

func (t StructProperty) DType() builtins.VALUE_TYPE {
	return t.Type.DType()
}

type StructMethod struct {
	IsPrivate bool
	Fn
}

func (t StructMethod) DType() builtins.VALUE_TYPE {
	return t.DataType
}

type Struct struct {
	DataType    builtins.VALUE_TYPE
	StructName  string
	StructScope TypeEnvironment
}

func (t Struct) DType() builtins.VALUE_TYPE {
	return t.DataType
}

type Array struct {
	DataType  builtins.VALUE_TYPE
	ArrayType ValueTypeInterface
}

func (t Array) DType() builtins.VALUE_TYPE {
	return t.DataType
}

type UserDefined struct {
	DataType builtins.VALUE_TYPE
	TypeName string
	TypeDef  ValueTypeInterface
}

func (t UserDefined) DType() builtins.VALUE_TYPE {
	return t.DataType
}

type ReturnType struct {
	DataType   builtins.VALUE_TYPE
	Expression ValueTypeInterface
}

func (t ReturnType) DType() builtins.VALUE_TYPE {
	return t.DataType
}

type Block struct {
	DataType builtins.VALUE_TYPE
	Returns  ValueTypeInterface
	Node     ast.Node
}

func (t Block) DType() builtins.VALUE_TYPE {
	return t.DataType
}

type Interface struct {
	DataType      builtins.VALUE_TYPE
	InterfaceName string
	Methods       map[string]Fn
}

func (t Interface) DType() builtins.VALUE_TYPE {
	return t.DataType
}

func makeNumericType(isInt bool, bitSize uint8, isSigned bool) builtins.VALUE_TYPE {
	rawDataType := ""
	if isInt {
		if isSigned {
			rawDataType = "i"
		} else {
			rawDataType = "u"
		}
	} else {
		rawDataType = "f"
	}

	rawDataType += fmt.Sprintf("%d", bitSize)

	return builtins.VALUE_TYPE(rawDataType)
}

// helper type initialization functions
func NewInt(bitSize uint8, isSigned bool) Int {
	return Int{DataType: makeNumericType(true, bitSize, isSigned), BitSize: bitSize, IsSigned: isSigned}
}

func NewFloat(bitSize uint8) Float {
	return Float{DataType: makeNumericType(false, bitSize, false), BitSize: bitSize}
}

func NewStr() Str {
	return Str{DataType: STRING_TYPE}
}

func NewBool() Bool {
	return Bool{DataType: BOOLEAN_TYPE}
}

func NewNull() Null {
	return Null{DataType: NULL_TYPE}
}

func NewVoid() Void {
	return Void{DataType: VOID_TYPE}
}

type SCOPE_TYPE int

const (
	GLOBAL_SCOPE SCOPE_TYPE = iota
	FUNCTION_SCOPE
	STRUCT_SCOPE
	CONDITIONAL_SCOPE
	LOOP_SCOPE
)

type TypeEnvironment struct {
	parent     *TypeEnvironment
	scopeType  SCOPE_TYPE
	scopeName  string
	variables  map[string]ValueTypeInterface
	constants  map[string]bool
	isOptional map[string]bool
	types      map[string]ValueTypeInterface
	interfaces map[string]Interface
	filePath   string
}

func NewTypeENV(parent *TypeEnvironment, scope SCOPE_TYPE, scopeName string, filePath string) *TypeEnvironment {
	return &TypeEnvironment{
		parent:     parent,
		scopeType:  scope,
		scopeName:  scopeName,
		filePath:   filePath,
		variables:  make(map[string]ValueTypeInterface),
		constants:  make(map[string]bool),
		isOptional: make(map[string]bool),
		types:      make(map[string]ValueTypeInterface),
		interfaces: make(map[string]Interface),
	}
}

func (t *TypeEnvironment) IsInStructScope() bool {
	if t.scopeType == STRUCT_SCOPE {
		return true
	}
	if t.parent == nil {
		return false
	}
	return t.parent.IsInStructScope()
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
		return nil, fmt.Errorf("'%s' was not declared in this scope", name)
	}

	return t.parent.ResolveVar(name)
}

func (t *TypeEnvironment) ResolveType(name string) (*TypeEnvironment, error) {
	if _, ok := t.types[name]; ok {
		return t, nil
	}
	if t.parent == nil {
		return nil, fmt.Errorf("type '%s' was not declared in this scope", name)
	}
	return t.parent.ResolveType(name)
}

// instead, find all the upper scopes and check if the type is declared in any of them
func (t *TypeEnvironment) GetTypeFromEnv(name string) (ValueTypeInterface, error) {
	if val, ok := t.types[name]; ok {
		return val.(UserDefined).TypeDef, nil
	}
	if t.parent == nil {
		return nil, fmt.Errorf("'%s' was not declared in this scope", name)
	}
	return t.parent.GetTypeFromEnv(name)
}

func (t *TypeEnvironment) DeclareVar(name string, typeVar ValueTypeInterface, isConst bool, isOptional bool) error {

	if _, ok := t.types[name]; ok {
		return fmt.Errorf("cannot declare variable with type '%s'", name)
	}

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
	fmt.Printf("Declared type %s, %T\n", name, typeType)
	return nil
}

func (t *TypeEnvironment) isDeclared(name string) bool {
	if _, ok := t.variables[name]; ok {
		return true
	}
	return false
}

func GetValueType(value ast.Node, t *TypeEnvironment) ValueTypeInterface {

	typ := CheckAST(value, t)

	typ, err := getValueTypeInterface(typ, t)
	if err != nil {
		//errgen.AddError(t.filePath, value.StartPos().Line, value.EndPos().Line, value.StartPos().Column, value.EndPos().Column, err.Error()).DisplayWithPanic()
		errgen.AddError(t.filePath, value.StartPos().Line, value.EndPos().Line, value.StartPos().Column, value.EndPos().Column, err.Error())
		return nil
	}

	return typ
}

func getValueTypeInterface(typ ValueTypeInterface, env *TypeEnvironment) (ValueTypeInterface, error) {
	switch t := typ.(type) {
	case UserDefined:
		//return getValueTypeInterface(t.TypeDef, env)
		//find the type in the env
		fmt.Printf("UserDefined type %s\n", t.TypeName)
		if val, err := env.GetTypeFromEnv(t.TypeName); err == nil {
			return val, nil
		} else {
			return nil, err
		}
	default:
		return t, nil
	}
}
