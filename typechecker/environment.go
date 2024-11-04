package typechecker

import (
	"fmt"
	"walrus/ast"
	"walrus/errgen"
)

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
