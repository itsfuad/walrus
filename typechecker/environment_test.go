package typechecker

import (
	"testing"
)

func TestNewTypeENV(t *testing.T) {
	env := NewTypeENV(nil, GLOBAL_SCOPE, "global", "/path/to/file")
	if env == nil {
		t.Fatal("Expected non-nil TypeEnvironment")
	}
	if env.scopeType != GLOBAL_SCOPE {
		t.Errorf("Expected scopeType to be GLOBAL_SCOPE, got %v", env.scopeType)
	}
	if env.scopeName != "global" {
		t.Errorf("Expected scopeName to be 'global', got %v", env.scopeName)
	}
	if env.filePath != "/path/to/file" {
		t.Errorf("Expected filePath to be '/path/to/file', got %v", env.filePath)
	}
}

func TestDeclareVar(t *testing.T) {
	env := NewTypeENV(nil, GLOBAL_SCOPE, "global", "/path/to/file")
	intType := Int{DataType: INT_TYPE}

	err := env.DeclareVar("x", intType, false, false)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
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
	env := NewTypeENV(nil, GLOBAL_SCOPE, "global", "/path/to/file")
	intType := Int{DataType: INT_TYPE}
	env.DeclareVar("x", intType, false, false)

	scope, err := env.ResolveVar("x")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if scope != env {
		t.Errorf("Expected scope to be the same as env")
	}

	_, err = env.ResolveVar("y")
	if err == nil {
		t.Fatalf("Expected error, got nil")
	}
}

func TestDeclareType(t *testing.T) {
	env := NewTypeENV(nil, GLOBAL_SCOPE, "global", "/path/to/file")
	structType := Struct{DataType: STRUCT_TYPE, StructName: "MyStruct"}

	err := env.DeclareType("MyStruct", structType)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if _, ok := env.types["MyStruct"]; !ok {
		t.Errorf("Expected type 'MyStruct' to be declared")
	}
}

func TestResolveType(t *testing.T) {
	env := NewTypeENV(nil, GLOBAL_SCOPE, "global", "/path/to/file")
	structType := Struct{DataType: STRUCT_TYPE, StructName: "MyStruct"}
	env.DeclareType("MyStruct", structType)

	scope, err := env.ResolveType("MyStruct")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if scope != env {
		t.Errorf("Expected scope to be the same as env")
	}

	_, err = env.ResolveType("UnknownType")
	if err == nil {
		t.Fatalf("Expected error, got nil")
	}
}

func TestResolveFunctionEnv(t *testing.T) {
	globalEnv := NewTypeENV(nil, GLOBAL_SCOPE, "global", "/path/to/file")
	funcEnv := NewTypeENV(globalEnv, FUNCTION_SCOPE, "function", "/path/to/file")

	resolvedEnv, err := funcEnv.ResolveFunctionEnv()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if resolvedEnv != funcEnv {
		t.Errorf("Expected resolvedEnv to be the same as funcEnv")
	}

	_, err = globalEnv.ResolveFunctionEnv()
	if err == nil {
		t.Fatalf("Expected error, got nil")
	}
}
