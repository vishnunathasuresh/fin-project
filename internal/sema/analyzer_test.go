package sema

import (
	"errors"
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

func TestAnalyze_ReservedSet(t *testing.T) {
	prog := &ast.Program{Statements: []ast.Statement{
		&ast.SetStmt{Name: "if", P: ast.Pos{Line: 1, Column: 1}},
	}}
	errs := Analyze(prog)
	if len(errs) == 0 {
		t.Fatalf("expected reserved name error for set")
	}
	var r ReservedNameError
	if !errors.As(errs[0], &r) {
		t.Fatalf("expected ReservedNameError, got %T", errs[0])
	}
}

func TestAnalyze_ReservedFnName(t *testing.T) {
	prog := &ast.Program{Statements: []ast.Statement{
		&ast.FnDecl{Name: "for", Params: nil, Body: nil, P: ast.Pos{Line: 2, Column: 3}},
	}}
	errs := Analyze(prog)
	if len(errs) == 0 {
		t.Fatalf("expected reserved name error for fn name")
	}
	var r ReservedNameError
	if !errors.As(errs[0], &r) {
		t.Fatalf("expected ReservedNameError, got %T", errs[0])
	}
}

func TestAnalyze_ReservedFnParam(t *testing.T) {
	prog := &ast.Program{Statements: []ast.Statement{
		&ast.FnDecl{Name: "valid", Params: []string{"while"}, Body: nil, P: ast.Pos{Line: 3, Column: 5}},
	}}
	errs := Analyze(prog)
	if len(errs) == 0 {
		t.Fatalf("expected reserved name error for fn param")
	}
	var r ReservedNameError
	if !errors.As(errs[0], &r) {
		t.Fatalf("expected ReservedNameError, got %T", errs[0])
	}
}

func TestAnalyze_UndefinedVariable(t *testing.T) {
	prog := &ast.Program{Statements: []ast.Statement{
		&ast.EchoStmt{Value: &ast.IdentExpr{Name: "x", P: ast.Pos{Line: 1, Column: 1}}, P: ast.Pos{Line: 1, Column: 1}},
	}}
	errs := Analyze(prog)
	if len(errs) == 0 {
		t.Fatalf("expected undefined variable error")
	}
	var u UndefinedVariableError
	if !errors.As(errs[0], &u) {
		t.Fatalf("expected UndefinedVariableError, got %T", errs[0])
	}
}

func TestAnalyze_DuplicateDefinition(t *testing.T) {
	prog := &ast.Program{Statements: []ast.Statement{
		&ast.SetStmt{Name: "a", Value: &ast.NumberLit{Value: "1", P: ast.Pos{Line: 1, Column: 10}}, P: ast.Pos{Line: 1, Column: 1}},
		&ast.SetStmt{Name: "a", Value: &ast.NumberLit{Value: "2", P: ast.Pos{Line: 2, Column: 10}}, P: ast.Pos{Line: 2, Column: 1}},
	}}
	errs := Analyze(prog)
	if len(errs) < 1 {
		t.Fatalf("expected duplicate definition error")
	}
	if err := errs[len(errs)-1]; err == nil {
		t.Fatalf("expected an error")
	}
}

func TestAnalyze_CallMissingFunction(t *testing.T) {
	prog := &ast.Program{Statements: []ast.Statement{
		&ast.CallStmt{Name: "missing", Args: nil, P: ast.Pos{Line: 1, Column: 1}},
	}}
	errs := Analyze(prog)
	if len(errs) == 0 {
		t.Fatalf("expected undefined function error")
	}
	var u UndefinedVariableError
	if !errors.As(errs[0], &u) {
		t.Fatalf("expected UndefinedVariableError, got %T", errs[0])
	}
}

func TestAnalyze_CallArityMismatch(t *testing.T) {
	prog := &ast.Program{Statements: []ast.Statement{
		&ast.FnDecl{Name: "foo", Params: []string{"a"}, Body: nil, P: ast.Pos{Line: 1, Column: 1}},
		&ast.CallStmt{Name: "foo", Args: []ast.Expr{}, P: ast.Pos{Line: 2, Column: 1}},
	}}
	errs := Analyze(prog)
	if len(errs) == 0 {
		t.Fatalf("expected arity error")
	}
	var ia InvalidArityError
	if !errors.As(errs[0], &ia) {
		t.Fatalf("expected InvalidArityError, got %T", errs[0])
	}
}

func TestAnalyze_CallOK(t *testing.T) {
	prog := &ast.Program{Statements: []ast.Statement{
		&ast.FnDecl{Name: "bar", Params: []string{"a", "b"}, Body: nil, P: ast.Pos{Line: 1, Column: 1}},
		&ast.CallStmt{Name: "bar", Args: []ast.Expr{&ast.NumberLit{Value: "1", P: ast.Pos{Line: 2, Column: 5}}, &ast.NumberLit{Value: "2", P: ast.Pos{Line: 2, Column: 8}}}, P: ast.Pos{Line: 2, Column: 1}},
	}}
	errs := Analyze(prog)
	if len(errs) != 0 {
		t.Fatalf("expected no errors, got %v", errs)
	}
}

func TestAnalyzeDefinitions_TracksFnScope(t *testing.T) {
	fn := &ast.FnDecl{Name: "foo", Params: []string{"p"}, Body: []ast.Statement{
		&ast.SetStmt{Name: "x", Value: &ast.NumberLit{Value: "1", P: ast.Pos{Line: 2, Column: 5}}, P: ast.Pos{Line: 2, Column: 1}},
	}, P: ast.Pos{Line: 1, Column: 1}}
	prog := &ast.Program{Statements: []ast.Statement{fn}}
	res := AnalyzeDefinitions(prog)
	if len(res.Errors) != 0 {
		t.Fatalf("unexpected errors: %v", res.Errors)
	}
	scope, ok := res.FuncScopes[fn]
	if !ok || scope == nil {
		t.Fatalf("expected function scope recorded")
	}
	if _, found := scope.Lookup("p"); !found {
		t.Fatalf("expected param p in function scope")
	}
	if _, found := scope.Lookup("x"); !found {
		t.Fatalf("expected set variable x in function scope")
	}
}

func TestAnalyzeDefinitions_TracksForScope(t *testing.T) {
	forStmt := &ast.ForStmt{Var: "i", Start: &ast.NumberLit{Value: "1", P: ast.Pos{Line: 1, Column: 10}}, End: &ast.NumberLit{Value: "3", P: ast.Pos{Line: 1, Column: 15}}, Body: []ast.Statement{
		&ast.SetStmt{Name: "j", Value: &ast.NumberLit{Value: "2", P: ast.Pos{Line: 2, Column: 5}}, P: ast.Pos{Line: 2, Column: 1}},
	}, P: ast.Pos{Line: 1, Column: 1}}
	prog := &ast.Program{Statements: []ast.Statement{forStmt}}
	res := AnalyzeDefinitions(prog)
	if len(res.Errors) != 0 {
		t.Fatalf("unexpected errors: %v", res.Errors)
	}
	scope, ok := res.ForScopes[forStmt]
	if !ok || scope == nil {
		t.Fatalf("expected for-loop scope recorded")
	}
	if _, found := scope.Lookup("i"); !found {
		t.Fatalf("expected loop var i in scope")
	}
	if _, found := scope.Lookup("j"); !found {
		t.Fatalf("expected set var j in loop scope")
	}
}

func TestAnalyze_NoShadowInFnParams(t *testing.T) {
	prog := &ast.Program{Statements: []ast.Statement{
		&ast.SetStmt{Name: "x", Value: &ast.NumberLit{Value: "1", P: ast.Pos{Line: 1, Column: 5}}, P: ast.Pos{Line: 1, Column: 1}},
		&ast.FnDecl{Name: "foo", Params: []string{"x"}, Body: nil, P: ast.Pos{Line: 2, Column: 1}},
	}}
	errs := Analyze(prog)
	if len(errs) == 0 {
		t.Fatalf("expected shadowing error for param x")
	}
	var sh ShadowingError
	if !errors.As(errs[0], &sh) {
		t.Fatalf("expected ShadowingError, got %T", errs[0])
	}
}

func TestAnalyze_NoShadowInNestedSet(t *testing.T) {
	prog := &ast.Program{Statements: []ast.Statement{
		&ast.SetStmt{Name: "y", Value: &ast.NumberLit{Value: "1", P: ast.Pos{Line: 1, Column: 5}}, P: ast.Pos{Line: 1, Column: 1}},
		&ast.IfStmt{Cond: &ast.BoolLit{Value: true, P: ast.Pos{Line: 2, Column: 4}}, Then: []ast.Statement{
			&ast.SetStmt{Name: "y", Value: &ast.NumberLit{Value: "2", P: ast.Pos{Line: 3, Column: 9}}, P: ast.Pos{Line: 3, Column: 1}},
		}, P: ast.Pos{Line: 2, Column: 1}},
	}}
	errs := Analyze(prog)
	if len(errs) == 0 {
		t.Fatalf("expected shadowing error for nested set y")
	}
	var sh ShadowingError
	if !errors.As(errs[0], &sh) {
		t.Fatalf("expected ShadowingError, got %T", errs[0])
	}
}
