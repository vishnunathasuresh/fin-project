package sema

import (
    "errors"
    "testing"

    "github.com/vishnunath-suresh/fin-project/internal/ast"
)

func TestIntegration_UndefinedVariable(t *testing.T) {
    prog := &ast.Program{Statements: []ast.Statement{
        &ast.EchoStmt{Value: &ast.IdentExpr{Name: "missing", P: ast.Pos{Line: 1, Column: 6}}, P: ast.Pos{Line: 1, Column: 1}},
    }}
    res := AnalyzeDefinitions(prog)
    if len(res.Errors) == 0 {
        t.Fatalf("expected undefined variable error")
    }
    var u UndefinedVariableError
    if !errors.As(res.Errors[0], &u) {
        t.Fatalf("expected UndefinedVariableError, got %T", res.Errors[0])
    }
    if u.P.Line != 1 || u.P.Column != 6 {
        t.Fatalf("expected position 1:6, got %d:%d", u.P.Line, u.P.Column)
    }
}

func TestIntegration_DuplicateFunction(t *testing.T) {
    prog := &ast.Program{Statements: []ast.Statement{
        &ast.FnDecl{Name: "foo", Params: nil, Body: nil, P: ast.Pos{Line: 1, Column: 1}},
        &ast.FnDecl{Name: "foo", Params: nil, Body: nil, P: ast.Pos{Line: 2, Column: 1}},
    }}
    res := AnalyzeDefinitions(prog)
    var df DuplicateFunctionError
    if !errors.As(res.Errors[0], &df) {
        t.Fatalf("expected DuplicateFunctionError, got %T", res.Errors[0])
    }
    if df.P.Line != 2 || df.P.Column != 1 {
        t.Fatalf("expected duplicate at 2:1, got %d:%d", df.P.Line, df.P.Column)
    }
}

func TestIntegration_InvalidArity(t *testing.T) {
    prog := &ast.Program{Statements: []ast.Statement{
        &ast.FnDecl{Name: "bar", Params: []string{"x"}, Body: nil, P: ast.Pos{Line: 1, Column: 1}},
        &ast.CallStmt{Name: "bar", Args: nil, P: ast.Pos{Line: 2, Column: 1}},
    }}
    res := AnalyzeDefinitions(prog)
    if len(res.Errors) == 0 {
        t.Fatalf("expected invalid arity error")
    }
    var ia InvalidArityError
    if !errors.As(res.Errors[0], &ia) {
        t.Fatalf("expected InvalidArityError, got %T", res.Errors[0])
    }
    if ia.P.Line != 2 || ia.P.Column != 1 {
        t.Fatalf("expected call position 2:1, got %d:%d", ia.P.Line, ia.P.Column)
    }
}

func TestIntegration_ShadowingAllowedAcrossFunctions(t *testing.T) {
    prog := &ast.Program{Statements: []ast.Statement{
        &ast.FnDecl{Name: "a", Params: []string{"x"}, Body: []ast.Statement{&ast.EchoStmt{Value: &ast.IdentExpr{Name: "x", P: ast.Pos{Line: 2, Column: 10}}, P: ast.Pos{Line: 2, Column: 5}}}, P: ast.Pos{Line: 1, Column: 1}},
        &ast.FnDecl{Name: "b", Params: []string{"x"}, Body: []ast.Statement{&ast.EchoStmt{Value: &ast.IdentExpr{Name: "x", P: ast.Pos{Line: 4, Column: 10}}, P: ast.Pos{Line: 4, Column: 5}}}, P: ast.Pos{Line: 3, Column: 1}},
    }}
    res := AnalyzeDefinitions(prog)
    if len(res.Errors) != 0 {
        t.Fatalf("expected no errors for distinct functions reusing param names, got %v", res.Errors)
    }
}

func TestIntegration_ReservedNameRejected(t *testing.T) {
    prog := &ast.Program{Statements: []ast.Statement{
        &ast.SetStmt{Name: "if", Value: &ast.NumberLit{Value: "1", P: ast.Pos{Line: 1, Column: 5}}, P: ast.Pos{Line: 1, Column: 1}},
    }}
    res := AnalyzeDefinitions(prog)
    if len(res.Errors) == 0 {
        t.Fatalf("expected reserved name error")
    }
    var r ReservedNameError
    if !errors.As(res.Errors[0], &r) {
        t.Fatalf("expected ReservedNameError, got %T", res.Errors[0])
    }
    if r.P.Line != 1 || r.P.Column != 1 {
        t.Fatalf("expected reserved position 1:1, got %d:%d", r.P.Line, r.P.Column)
    }
}
