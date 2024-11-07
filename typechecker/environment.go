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

var typeDefinitions = map[string]ValueTypeInterface{
	string(INT8_TYPE):    NewInt(8, true),
	string(INT16_TYPE):   NewInt(16, true),
	string(INT32_TYPE):   NewInt(32, true),
	string(INT64_TYPE):   NewInt(64, true),
	string(UINT8_TYPE):   NewInt(8, false),
	string(UINT16_TYPE):  NewInt(16, false),
	string(UINT32_TYPE):  NewInt(32, false),
	string(UINT64_TYPE):  NewInt(64, false),
	string(BYTE_TYPE):    NewInt(8, false),
	string(STRING_TYPE):  NewStr(),
	string(FLOAT32_TYPE): NewFloat(32),
	string(FLOAT64_TYPE): NewFloat(64),
	string(NULL_TYPE):    NewNull(),
	string(VOID_TYPE):    NewVoid(),
}

type TypeEnvironment struct {
	parent     *TypeEnvironment
	scopeType  SCOPE_TYPE
	scopeName  string
	variables  map[string]ValueTypeInterface
	constants  map[string]bool
	isOptional map[string]bool
	interfaces map[string]Interface
	filePath   string
}

func ProgramEnv(filepath string) *TypeEnvironment {
	env := NewTypeENV(nil, GLOBAL_SCOPE, "global", filepath)
	env.DeclareVar("true", NewBool(), true, false)
	env.DeclareVar("false", NewBool(), true, false)
	env.DeclareVar("null", NewNull(), true, false)
	env.DeclareVar("PI", NewFloat(32), true, false)
	return env
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

func (t *TypeEnvironment) DeclareVar(name string, typeVar ValueTypeInterface, isConst bool, isOptional bool) error {

	if _, ok := typeDefinitions[name]; ok {
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

func DeclareType(name string, typeType ValueTypeInterface) error {
	if _, ok := typeDefinitions[name]; ok {
		return fmt.Errorf("type '%s' already defined", name)
	}
	typeDefinitions[name] = typeType
	return nil
}

func (t *TypeEnvironment) isDeclared(name string) bool {
	if _, ok := t.variables[name]; ok {
		return true
	}
	return false
}

func nodeType(value ast.Node, t *TypeEnvironment) ValueTypeInterface {

	typ := CheckAST(value, t)

	typ, err := unwrapType(typ)
	if err != nil {
		errgen.AddError(t.filePath, value.StartPos().Line, value.EndPos().Line, value.StartPos().Column, value.EndPos().Column, err.Error()).DisplayWithPanic()
	}

	return typ
}

func getTypeDefinition(name string) (ValueTypeInterface, error) {
	if typ, ok := typeDefinitions[name]; !ok {
		return nil, fmt.Errorf("unknown type '%s'", name)
	} else {
		switch t := typ.(type) {
		case UserDefined:
			return unwrapType(t.TypeDef)
		default:
			return typ, nil
		}
	}
}

func unwrapType(value ValueTypeInterface) (ValueTypeInterface, error) {
	switch t := value.(type) {
	case UserDefined:
		return unwrapType(t.TypeDef)
	default:
		return value, nil
	}
}

func isTypeDefined(name string) bool {
	_, ok := typeDefinitions[name]
	return ok
}
