package sema

import (
	"errors"
	"testing"

	"github.com/vishnunathasuresh/fin-project/internal/ast"
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

func TestAnalyze_ReservedDecl(t *testing.T) {
	prog := &ast.Program{Statements: []ast.Statement{
		&ast.DeclStmt{Names: []string{"if"}, Value: &ast.NumberLit{Value: "1", P: ast.Pos{Line: 1, Column: 5}}, P: ast.Pos{Line: 1, Column: 1}},
	}}
	errs := Analyze(prog)
	if len(errs) == 0 {
		t.Fatalf("expected reserved name error for decl")
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
		&ast.FnDecl{Name: "valid", Params: []ast.Param{{Name: "while", P: ast.Pos{Line: 3, Column: 10}}}, Body: nil, P: ast.Pos{Line: 3, Column: 5}},
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
		&ast.AssignStmt{Names: []string{"x"}, Value: &ast.NumberLit{Value: "1", P: ast.Pos{Line: 1, Column: 5}}, P: ast.Pos{Line: 1, Column: 1}},
	}}
	err := New().Analyze(prog)
	if err == nil {
		t.Fatalf("expected undefined variable error")
	}
	var u UndefinedVariableError
	if !errors.As(err, &u) {
		t.Fatalf("expected UndefinedVariableError, got %T", err)
	}
}

func TestAnalyze_UndefinedVariable_PropertyBase(t *testing.T) {
	prog := &ast.Program{Statements: []ast.Statement{
		&ast.AssignStmt{Names: []string{"y"}, Value: &ast.PropertyExpr{Object: &ast.IdentExpr{Name: "obj", P: ast.Pos{Line: 1, Column: 6}}, Field: "f", P: ast.Pos{Line: 1, Column: 9}}, P: ast.Pos{Line: 1, Column: 1}},
	}}
	err := New().Analyze(prog)
	if err == nil {
		t.Fatalf("expected undefined variable error for property base")
	}
	var u UndefinedVariableError
	if !errors.As(err, &u) {
		t.Fatalf("expected UndefinedVariableError, got %T", err)
	}
}

func TestAnalyze_SkipReservedIdent(t *testing.T) {
	prog := &ast.Program{Statements: []ast.Statement{
		&ast.DeclStmt{Names: []string{"true"}, Value: &ast.NumberLit{Value: "1", P: ast.Pos{Line: 1, Column: 6}}, P: ast.Pos{Line: 1, Column: 1}},
	}}
	err := New().Analyze(prog)
	if err == nil {
		t.Fatalf("expected reserved name error for decl")
	}
	var r ReservedNameError
	if !errors.As(err, &r) {
		t.Fatalf("expected ReservedNameError, got %T", err)
	}
}

func TestAnalyze_DuplicateDefinition(t *testing.T) {
	prog := &ast.Program{Statements: []ast.Statement{
		&ast.DeclStmt{Names: []string{"a"}, Value: &ast.NumberLit{Value: "1", P: ast.Pos{Line: 1, Column: 10}}, P: ast.Pos{Line: 1, Column: 1}},
		&ast.DeclStmt{Names: []string{"a"}, Value: &ast.NumberLit{Value: "2", P: ast.Pos{Line: 2, Column: 10}}, P: ast.Pos{Line: 2, Column: 1}},
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
		&ast.FnDecl{Name: "foo", Params: []ast.Param{{Name: "a", P: ast.Pos{Line: 1, Column: 5}}}, Body: nil, P: ast.Pos{Line: 1, Column: 1}},
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
		&ast.FnDecl{Name: "bar", Params: []ast.Param{{Name: "a", P: ast.Pos{Line: 1, Column: 5}}, {Name: "b", P: ast.Pos{Line: 1, Column: 8}}}, Body: nil, P: ast.Pos{Line: 1, Column: 1}},
		&ast.CallStmt{Name: "bar", Args: []ast.Expr{&ast.NumberLit{Value: "1", P: ast.Pos{Line: 2, Column: 5}}, &ast.NumberLit{Value: "2", P: ast.Pos{Line: 2, Column: 8}}}, P: ast.Pos{Line: 2, Column: 1}},
	}}
	errs := Analyze(prog)
	if len(errs) != 0 {
		t.Fatalf("expected no errors, got %v", errs)
	}
}

func TestAnalyze_CallMixedCorrectAndWrong(t *testing.T) {
	prog := &ast.Program{Statements: []ast.Statement{
		&ast.FnDecl{Name: "f", Params: []ast.Param{{Name: "a", P: ast.Pos{Line: 1, Column: 5}}}, Body: nil, P: ast.Pos{Line: 1, Column: 1}},
		&ast.CallStmt{Name: "f", Args: []ast.Expr{&ast.NumberLit{Value: "1", P: ast.Pos{Line: 2, Column: 5}}}, P: ast.Pos{Line: 2, Column: 1}}, // ok
		&ast.CallStmt{Name: "f", Args: []ast.Expr{}, P: ast.Pos{Line: 3, Column: 1}},                                                           // wrong
	}}
	errs := Analyze(prog)
	if len(errs) != 1 {
		t.Fatalf("expected 1 error, got %v", errs)
	}
	var ia InvalidArityError
	if !errors.As(errs[0], &ia) {
		t.Fatalf("expected InvalidArityError, got %T", errs[0])
	}
}

func TestAnalyze_CallForwardReferenceArity(t *testing.T) {
	prog := &ast.Program{Statements: []ast.Statement{
		&ast.CallStmt{Name: "g", Args: []ast.Expr{}, P: ast.Pos{Line: 1, Column: 1}},
		&ast.FnDecl{Name: "g", Params: []ast.Param{}, Body: nil, P: ast.Pos{Line: 2, Column: 1}},
	}}
	errs := Analyze(prog)
	if len(errs) != 0 {
		t.Fatalf("expected no errors for forward reference with correct arity, got %v", errs)
	}
}

func TestAnalyze_CallUndefinedAndArityMixed(t *testing.T) {
	prog := &ast.Program{Statements: []ast.Statement{
		&ast.CallStmt{Name: "missing", Args: nil, P: ast.Pos{Line: 1, Column: 1}},
		&ast.FnDecl{Name: "h", Params: []ast.Param{{Name: "a", P: ast.Pos{Line: 2, Column: 5}}}, Body: nil, P: ast.Pos{Line: 2, Column: 1}},
		&ast.CallStmt{Name: "h", Args: []ast.Expr{}, P: ast.Pos{Line: 3, Column: 1}},
	}}
	errs := Analyze(prog)
	if len(errs) != 2 {
		t.Fatalf("expected 2 errors, got %v", errs)
	}
	var u UndefinedVariableError
	var ia InvalidArityError
	if !errors.As(errs[0], &u) || !errors.As(errs[1], &ia) {
		t.Fatalf("expected undefined then arity errors, got %v", errs)
	}
}

func TestAnalyze_CallZeroArgEdges(t *testing.T) {
	prog := &ast.Program{Statements: []ast.Statement{
		&ast.FnDecl{Name: "z", Params: []ast.Param{}, Body: nil, P: ast.Pos{Line: 1, Column: 1}},
		&ast.CallStmt{Name: "z", Args: []ast.Expr{&ast.NumberLit{Value: "1", P: ast.Pos{Line: 2, Column: 5}}}, P: ast.Pos{Line: 2, Column: 1}},
		&ast.FnDecl{Name: "u", Params: []ast.Param{{Name: "a", P: ast.Pos{Line: 3, Column: 5}}, {Name: "b", P: ast.Pos{Line: 3, Column: 8}}}, Body: nil, P: ast.Pos{Line: 3, Column: 1}},
		&ast.CallStmt{Name: "u", Args: []ast.Expr{}, P: ast.Pos{Line: 4, Column: 1}},
	}}
	errs := Analyze(prog)
	if len(errs) != 2 {
		t.Fatalf("expected 2 errors, got %v", errs)
	}
}

func TestAnalyzeDefinitionsWithLimit_ExceedsDepth(t *testing.T) {
	deep := &ast.Program{Statements: []ast.Statement{
		&ast.WhileStmt{Cond: &ast.BoolLit{Value: true, P: ast.Pos{Line: 1, Column: 1}}, Body: []ast.Statement{
			&ast.WhileStmt{Cond: &ast.BoolLit{Value: true, P: ast.Pos{Line: 2, Column: 1}}, Body: []ast.Statement{}, P: ast.Pos{Line: 2, Column: 1}},
		}, P: ast.Pos{Line: 1, Column: 1}},
	}}
	res := AnalyzeDefinitionsWithLimit(deep, 1)
	if len(res.Errors) == 0 {
		t.Fatalf("expected depth exceeded error")
	}
	var de DepthExceededError
	if !errors.As(res.Errors[0], &de) {
		t.Fatalf("expected DepthExceededError, got %T", res.Errors[0])
	}
}

func TestAnalyzeDefinitionsWithLimit_AllowsWithinLimit(t *testing.T) {
	prog := &ast.Program{Statements: []ast.Statement{
		&ast.WhileStmt{Cond: &ast.BoolLit{Value: true, P: ast.Pos{Line: 1, Column: 1}}, Body: []ast.Statement{}, P: ast.Pos{Line: 1, Column: 1}},
	}}
	res := AnalyzeDefinitionsWithLimit(prog, 2)
	if len(res.Errors) != 0 {
		t.Fatalf("expected no errors within depth limit, got %v", res.Errors)
	}
}

func TestAnalyzer_WrapperCollectsErrors(t *testing.T) {
	prog := &ast.Program{Statements: []ast.Statement{
		&ast.DeclStmt{Names: []string{"x"}, Value: &ast.NumberLit{Value: "1", P: ast.Pos{Line: 1, Column: 5}}, P: ast.Pos{Line: 1, Column: 1}},
		&ast.DeclStmt{Names: []string{"x"}, Value: &ast.NumberLit{Value: "2", P: ast.Pos{Line: 2, Column: 5}}, P: ast.Pos{Line: 2, Column: 1}},
	}}
	a := NewAnalyzer(prog, 0)
	a.Run()
	if len(a.Errors()) == 0 {
		t.Fatalf("expected errors collected by analyzer")
	}
}

func TestAnalyzer_PublicAPI_NoErrorsReturnsNil(t *testing.T) {
	prog := &ast.Program{Statements: []ast.Statement{
		&ast.FnDecl{Name: "f", Params: []ast.Param{}, Body: nil, P: ast.Pos{Line: 1, Column: 1}},
		&ast.CallStmt{Name: "f", Args: nil, P: ast.Pos{Line: 2, Column: 1}},
	}}
	a := New()
	err := a.Analyze(prog)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
}

func TestAnalyzer_PublicAPI_AggregatesErrors(t *testing.T) {
	prog := &ast.Program{Statements: []ast.Statement{
		&ast.DeclStmt{Names: []string{"x"}, Value: &ast.NumberLit{Value: "1", P: ast.Pos{Line: 1, Column: 5}}, P: ast.Pos{Line: 1, Column: 1}},
		&ast.DeclStmt{Names: []string{"x"}, Value: &ast.NumberLit{Value: "2", P: ast.Pos{Line: 2, Column: 5}}, P: ast.Pos{Line: 2, Column: 1}},
	}}
	a := New()
	err := a.Analyze(prog)
	if err == nil {
		t.Fatalf("expected aggregated error")
	}
}

func TestAnalyzeDefinitions_TracksFnScope(t *testing.T) {
	fn := &ast.FnDecl{Name: "foo", Params: []ast.Param{{Name: "p", P: ast.Pos{Line: 1, Column: 10}}}, Body: []ast.Statement{
		&ast.DeclStmt{Names: []string{"x"}, Value: &ast.NumberLit{Value: "1", P: ast.Pos{Line: 2, Column: 5}}, P: ast.Pos{Line: 2, Column: 1}},
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
		t.Fatalf("expected declared variable x in function scope")
	}
}

func TestAnalyzeDefinitions_TracksForScope(t *testing.T) {
	forStmt := &ast.ForStmt{Var: "i", Start: &ast.NumberLit{Value: "1", P: ast.Pos{Line: 1, Column: 10}}, End: &ast.NumberLit{Value: "3", P: ast.Pos{Line: 1, Column: 15}}, Body: []ast.Statement{
		&ast.DeclStmt{Names: []string{"j"}, Value: &ast.NumberLit{Value: "2", P: ast.Pos{Line: 2, Column: 5}}, P: ast.Pos{Line: 2, Column: 1}},
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
		t.Fatalf("expected declared var j in loop scope")
	}
}

func TestAnalyze_NoShadowInFnParams(t *testing.T) {
	prog := &ast.Program{Statements: []ast.Statement{
		&ast.DeclStmt{Names: []string{"x"}, Value: &ast.NumberLit{Value: "1", P: ast.Pos{Line: 1, Column: 5}}, P: ast.Pos{Line: 1, Column: 1}},
		&ast.FnDecl{Name: "foo", Params: []ast.Param{{Name: "x", P: ast.Pos{Line: 2, Column: 5}}}, Body: nil, P: ast.Pos{Line: 2, Column: 1}},
	}}
	errs := Analyze(prog)
	if len(errs) == 0 {
		t.Fatalf("expected shadowing error for param x")
	}
	var sh ShadowingError
	if !errors.As(errs[0], &sh) {
		t.Fatalf("expected ShadowingError, got %T", errs[0])
	}
	if sh.Def.Line != 1 || sh.Def.Column != 1 {
		t.Fatalf("expected def position 1:1, got %d:%d", sh.Def.Line, sh.Def.Column)
	}
}

func TestAnalyze_NoShadowInNestedDecl(t *testing.T) {
	prog := &ast.Program{Statements: []ast.Statement{
		&ast.DeclStmt{Names: []string{"y"}, Value: &ast.NumberLit{Value: "1", P: ast.Pos{Line: 1, Column: 5}}, P: ast.Pos{Line: 1, Column: 1}},
		&ast.IfStmt{Cond: &ast.BoolLit{Value: true, P: ast.Pos{Line: 2, Column: 4}}, Then: []ast.Statement{
			&ast.DeclStmt{Names: []string{"y"}, Value: &ast.NumberLit{Value: "2", P: ast.Pos{Line: 3, Column: 9}}, P: ast.Pos{Line: 3, Column: 1}},
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
	if sh.Def.Line != 1 || sh.Def.Column != 1 {
		t.Fatalf("expected def position 1:1, got %d:%d", sh.Def.Line, sh.Def.Column)
	}
}
