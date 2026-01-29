package sema

import (
	"errors"
	"fmt"
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

func TestIntegration_DeepNesting_NoPanic(t *testing.T) {
	depth := 200
	stmt := ast.Statement(&ast.WhileStmt{Cond: &ast.BoolLit{Value: true, P: ast.Pos{Line: 1, Column: 1}}, Body: nil, P: ast.Pos{Line: 1, Column: 1}})
	for i := 0; i < depth; i++ {
		stmt = &ast.WhileStmt{Cond: &ast.BoolLit{Value: true, P: ast.Pos{Line: i + 2, Column: 1}}, Body: []ast.Statement{stmt}, P: ast.Pos{Line: i + 2, Column: 1}}
	}
	prog := &ast.Program{Statements: []ast.Statement{stmt}}
	a := New()
	if err := a.Analyze(prog); err != nil {
		t.Fatalf("expected no errors for deep nesting, got %v", err)
	}
}

func TestIntegration_MixedLargeProgram_NoPanicAndAggregates(t *testing.T) {
	// many sets
	var stmts []ast.Statement
	for i := 0; i < 50; i++ {
		name := fmt.Sprintf("v%d", i)
		stmts = append(stmts, &ast.SetStmt{Name: name, Value: &ast.NumberLit{Value: "1", P: ast.Pos{Line: i + 1, Column: 5}}, P: ast.Pos{Line: i + 1, Column: 1}})
	}
	// many functions
	for i := 0; i < 10; i++ {
		fname := fmt.Sprintf("fn%d", i)
		stmts = append(stmts, &ast.FnDecl{Name: fname, Params: []string{"p"}, Body: []ast.Statement{
			&ast.EchoStmt{Value: &ast.IdentExpr{Name: "p", P: ast.Pos{Line: 200 + i, Column: 10}}, P: ast.Pos{Line: 200 + i, Column: 5}},
		}, P: ast.Pos{Line: 200 + i, Column: 1}})
	}
	// duplicate function and bad call + undefined call
	stmts = append(stmts,
		&ast.FnDecl{Name: "dup", Params: []string{"a"}, Body: nil, P: ast.Pos{Line: 300, Column: 1}},
		&ast.FnDecl{Name: "dup", Params: []string{"b"}, Body: nil, P: ast.Pos{Line: 301, Column: 1}},
		&ast.CallStmt{Name: "dup", Args: []ast.Expr{}, P: ast.Pos{Line: 302, Column: 1}},
		&ast.CallStmt{Name: "missingFn", Args: nil, P: ast.Pos{Line: 303, Column: 1}},
	)
	prog := &ast.Program{Statements: stmts}
	a := New()
	err := a.Analyze(prog)
	if err == nil {
		t.Fatalf("expected aggregated errors for mixed program")
	}
	// ensure specific error types were produced
	res := a.Result()
	if !containsErrorType[DuplicateFunctionError](res.Errors) {
		t.Fatalf("expected DuplicateFunctionError in mixed program")
	}
	if !containsErrorType[InvalidArityError](res.Errors) {
		t.Fatalf("expected InvalidArityError in mixed program")
	}
	if !containsErrorType[UndefinedVariableError](res.Errors) {
		t.Fatalf("expected UndefinedVariableError in mixed program")
	}
}

func containsErrorType[T error](errs []error) bool {
	for _, e := range errs {
		var tgt T
		if errors.As(e, &tgt) {
			return true
		}
	}
	return false
}
