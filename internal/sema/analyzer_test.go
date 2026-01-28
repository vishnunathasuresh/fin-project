package sema

import (
	"testing"

	"github.com/vishnunath-suresh/fin-project/internal/ast"
)

func TestFunctionRegistry_DefineAndLookup(t *testing.T) {
	reg := NewFunctionRegistry()
	if err := reg.Define("foo", 2, ast.Pos{Line: 1, Column: 1}); err != nil {
		t.Fatalf("unexpected define error: %v", err)
	}
	if arity, ok := reg.Lookup("foo"); !ok || arity != 2 {
		t.Fatalf("lookup foo got ok=%v arity=%d, want ok=true arity=2", ok, arity)
	}
}

func TestFunctionRegistry_Duplicate(t *testing.T) {
	reg := NewFunctionRegistry()
	_ = reg.Define("foo", 1, ast.Pos{Line: 1, Column: 1})
	if err := reg.Define("foo", 1, ast.Pos{Line: 2, Column: 1}); err == nil {
		t.Fatalf("expected duplicate error, got nil")
	}
}

func TestFunctionRegistry_LookupMissing(t *testing.T) {
	reg := NewFunctionRegistry()
	if _, ok := reg.Lookup("missing"); ok {
		t.Fatalf("expected missing to be absent")
	}
}
