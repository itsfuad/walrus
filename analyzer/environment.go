package analyzer

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
	CONDITIONAL_SCOPE
	LOOP_SCOPE
)

var typeDefinitions = map[string]TcValue{
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

var builtinValues = make(map[string]bool)

type TypeEnvironment struct {
	parent     *TypeEnvironment
	scopeType  SCOPE_TYPE
	scopeName  string
	variables  map[string]TcValue
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

func initVar(env *TypeEnvironment, name string, typeVar TcValue, isConst bool, isOptional bool) {
	err := env.declareVar(name, typeVar, isConst, isOptional)
	if err != nil {
		utils.RED.Println(err)
		os.Exit(-1)
	}
	builtinValues[name] = true
}

func NewTypeENV(parent *TypeEnvironment, scope SCOPE_TYPE, scopeName string, filePath string) *TypeEnvironment {
	return &TypeEnvironment{
		parent:     parent,
		scopeType:  scope,
		scopeName:  scopeName,
		filePath:   filePath,
		variables:  make(map[string]TcValue),
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

func (t *TypeEnvironment) resolveFunctionEnv() (*TypeEnvironment, error) {
	if t.scopeType == FUNCTION_SCOPE {
		return t, nil
	}
	if t.parent == nil {
		return nil, fmt.Errorf("function not found")
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

func (t *TypeEnvironment) declareVar(name string, typeVar TcValue, isConst bool, isOptional bool) error {

	if _, ok := typeDefinitions[name]; ok && name != "null" && name != "void" {
		return fmt.Errorf("type name '%s' cannot be used as variable name", name)
	}

	if ok, hasVal := builtinValues[name]; hasVal && ok {
		return fmt.Errorf("cannot declare builtin value '%s'", name)
	}

	//should not be declared
	if scope, err := t.resolveVar(name); err == nil && scope == t {
		return fmt.Errorf("'%s' is already declared in this scope", name)
	}

	t.variables[name] = typeVar
	t.constants[name] = isConst
	t.isOptional[name] = isOptional

	return nil
}

func declareType(name string, typeType TcValue) error {
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

	if _, ok := builtinValues[name]; ok {
		return true
	}

	return false
}

func getTypeDefinition(name string) (TcValue, error) {
	if typ, ok := typeDefinitions[name]; !ok {
		return nil, fmt.Errorf("unknown type '%s'", name)
	} else {
		if tp, ok := typ.(UserDefined); ok {
			if st, ok := tp.TypeDef.(Struct); ok {
				fmt.Printf("unwrapping type from: %s\n", st.StructName)
			}
		}
		return unwrapType(typ), nil
	}
}

func unwrapType(value TcValue) TcValue {
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
