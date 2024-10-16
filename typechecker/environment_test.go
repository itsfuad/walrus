package typechecker

import (
	"testing"
)

const (
	FILE = "/path/to/file"
	EXPECTED_ERROR = "Expected error, got nil"
	EXPECTED_NO_ERROR = "Expected no error, got %v"
)

func TestNewTypeENV(t *testing.T) {
	env := NewTypeENV(nil, GLOBAL_SCOPE, "global", FILE)
	if env == nil {
		t.Fatal("Expected non-nil TypeEnvironment")
	}
	if env.scopeType != GLOBAL_SCOPE {
		t.Errorf("Expected scopeType to be GLOBAL_SCOPE, got %v", env.scopeType)
	}
	if env.scopeName != "global" {
		t.Errorf("Expected scopeName to be 'global', got %v", env.scopeName)
	}
	if env.filePath != FILE {
		t.Errorf("Expected filePath to be '/path/to/file', got %v", env.filePath)
	}
}

func TestDeclareVar(t *testing.T) {
	env := NewTypeENV(nil, GLOBAL_SCOPE, "global", FILE)
	intType := Int{DataType: INT_TYPE}

	err := env.DeclareVar("x", intType, false, false)
	if err != nil {
		t.Fatalf(EXPECTED_NO_ERROR, err)
	}

	if _, ok := env.variables["x"]; !ok {
		t.Errorf("Expected variable 'x' to be declared")
	}

	if env.constants["x"] {
		t.Errorf("Expected variable 'x' to be non-constant")
	}

	if env.isOptional["x"] {
		t.Errorf("Expected variable 'x' to be non-optional")
	}
}

func TestResolveVar(t *testing.T) {
	env := NewTypeENV(nil, GLOBAL_SCOPE, "global", FILE)
	intType := Int{DataType: INT_TYPE}
	env.DeclareVar("x", intType, false, false)

	scope, err := env.ResolveVar("x")
	if err != nil {
		t.Fatalf(EXPECTED_NO_ERROR, err)
	}

	if scope != env {
		t.Errorf("Expected scope to be the same as env")
	}

	_, err = env.ResolveVar("y")
	if err == nil {
		t.Fatalf("EXPECTED_ERROR")
	}
}

func TestDeclareType(t *testing.T) {
	env := NewTypeENV(nil, GLOBAL_SCOPE, "global", FILE)
	structType := Struct{DataType: STRUCT_TYPE, StructName: "MyStruct"}

	err := env.DeclareType("MyStruct", structType)
	if err != nil {
		t.Fatalf(EXPECTED_NO_ERROR, err)
	}

	if _, ok := env.types["MyStruct"]; !ok {
		t.Errorf("Expected type 'MyStruct' to be declared")
	}
}

func TestResolveType(t *testing.T) {
	env := NewTypeENV(nil, GLOBAL_SCOPE, "global", FILE)
	structType := Struct{DataType: STRUCT_TYPE, StructName: "MyStruct"}
	env.DeclareType("MyStruct", structType)

	scope, err := env.ResolveType("MyStruct")
	if err != nil {
		t.Fatalf(EXPECTED_NO_ERROR, err)
	}

	if scope != env {
		t.Errorf("Expected scope to be the same as env")
	}

	_, err = env.ResolveType("UnknownType")
	if err == nil {
		t.Fatalf("EXPECTED_ERROR")
	}
}

func TestResolveFunctionEnv(t *testing.T) {
	globalEnv := NewTypeENV(nil, GLOBAL_SCOPE, "global", FILE)
	funcEnv := NewTypeENV(globalEnv, FUNCTION_SCOPE, "function", FILE)

	resolvedEnv, err := funcEnv.ResolveFunctionEnv()
	if err != nil {
		t.Fatalf(EXPECTED_NO_ERROR, err)
	}

	if resolvedEnv != funcEnv {
		t.Errorf("Expected resolvedEnv to be the same as funcEnv")
	}

	_, err = globalEnv.ResolveFunctionEnv()
	if err == nil {
		t.Fatalf("EXPECTED_ERROR")
	}
}
