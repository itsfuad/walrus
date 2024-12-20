package typechecker

import (
	"fmt"
	"os"
	"walrus/utils"
)

type SCOPE_TYPE int

const (
	GLOBAL_SCOPE SCOPE_TYPE = iota
	FUNCTION_SCOPE
	STRUCT_SCOPE
	SAFE_SCOPE
	OTHERWISE_SCOPE
	CONDITIONAL_SCOPE
	LOOP_SCOPE
)

var typeDefinitions = map[string]Tc{
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
	string(BOOLEAN_TYPE): NewBool(),
	string(NULL_TYPE):    NewNull(),
	string(VOID_TYPE):    NewVoid(),
}

var builtinValues = make(map[string]bool)

type TypeEnvironment struct {
	parent     *TypeEnvironment
	scopeType  SCOPE_TYPE
	scopeName  string
	variables  map[string]Tc
	constants  map[string]bool
	isOptional map[string]bool
	interfaces map[string]Interface
	filePath   string
}

func ProgramEnv(filepath string) *TypeEnvironment {
	env := NewTypeENV(nil, GLOBAL_SCOPE, "global", filepath)
	initVar(env, "true", NewBool(), true, false)
	initVar(env, "false", NewBool(), true, false)
	initVar(env, "null", NewNull(), true, false)
	initVar(env, "PI", NewFloat(32), true, false)
	return env
}

func initVar(env *TypeEnvironment, name string, typeVar Tc, isConst bool, isOptional bool) {
	err := env.declareVar(name, typeVar, isConst, isOptional)
	if err != nil {
		utils.RED.Println(err)
		os.Exit(-1)
	}

	builtinValues[name] = true

	utils.BRIGHT_BROWN.Printf("Initialized builtin value '%s'\n", name)
}

func NewTypeENV(parent *TypeEnvironment, scope SCOPE_TYPE, scopeName string, filePath string) *TypeEnvironment {
	return &TypeEnvironment{
		parent:     parent,
		scopeType:  scope,
		scopeName:  scopeName,
		filePath:   filePath,
		variables:  make(map[string]Tc),
		constants:  make(map[string]bool),
		isOptional: make(map[string]bool),
		interfaces: make(map[string]Interface),
	}
}

func (t *TypeEnvironment) isInFunctionScope() bool {
	if t.scopeType == FUNCTION_SCOPE {
		return true
	}
	if t.parent == nil {
		return false
	}
	return t.parent.isInFunctionScope()
}

func (t *TypeEnvironment) isInStructScope() bool {
	if t.scopeType == STRUCT_SCOPE {
		return true
	}
	if t.parent == nil {
		return false
	}
	return t.parent.isInStructScope()
}

func (t *TypeEnvironment) resolveFunctionEnv() (*TypeEnvironment, error) {
	if t.scopeType == FUNCTION_SCOPE {
		return t, nil
	}
	if t.parent == nil {
		return nil, fmt.Errorf("no function found in this scope")
	}
	return t.parent.resolveFunctionEnv()
}

func (t *TypeEnvironment) resolveVar(name string) (*TypeEnvironment, error) {

	if t.isDeclared(name) {
		return t, nil
	}

	//check on the parent then
	if t.parent == nil {
		//no where is declared
		return nil, fmt.Errorf("'%s' was not declared in this scope", name)
	}

	return t.parent.resolveVar(name)
}

func (t *TypeEnvironment) declareVar(name string, typeVar Tc, isConst bool, isOptional bool) error {

	if _, ok := typeDefinitions[name]; ok && name != "null" && name != "void" {
		return fmt.Errorf("type name '%s' cannot be used as variable name", name)
	}

	if ok, hasVal := builtinValues[name]; hasVal && ok {
		return fmt.Errorf("cannot redeclare builtin value '%s'", name)
	}

	//should not be declared
	if scope, err := t.resolveVar(name); err == nil && scope == t {
		return fmt.Errorf("variable '%s' is already declared in this scope", name)
	}

	t.variables[name] = typeVar
	t.constants[name] = isConst
	t.isOptional[name] = isOptional

	return nil
}

func declareType(name string, typeType Tc) error {
	if _, ok := typeDefinitions[name]; ok {
		return fmt.Errorf("type '%s' is already defined", name)
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

func getTypeDefinition(name string) (Tc, error) {
	if typ, ok := typeDefinitions[name]; !ok {
		return nil, fmt.Errorf("unknown type '%s'", name)
	} else {
		return unwrapType(typ), nil
	}
}

func unwrapType(value Tc) Tc {
	switch t := value.(type) {
	case UserDefined:
		return unwrapType(t.TypeDef)
	default:
		return t
	}
}

func isTypeDefined(name string) bool {
	_, ok := typeDefinitions[name]
	return ok
}
