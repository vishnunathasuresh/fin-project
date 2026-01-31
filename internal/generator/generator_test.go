package generator

import (
	"testing"

	"github.com/vishnunathasuresh/fin-project/internal/ast"
)

func TestGenerate_TopLevelSetEchoRun(t *testing.T) {
	g := NewBatchGenerator()
	prog := &ast.Program{Statements: []ast.Statement{
		&ast.SetStmt{Name: "x", Value: &ast.NumberLit{Value: "10"}},
		&ast.EchoStmt{Value: &ast.IdentExpr{Name: "x"}},
		&ast.RunStmt{Command: &ast.StringLit{Value: "git status"}},
	}}

	out, err := g.Generate(prog)
	if err != nil {
		t.Fatalf("Generate returned error: %v", err)
	}

	want := "@echo off\n" +
		"setlocal EnableDelayedExpansion\n" +
		"set x=10\n" +
		"echo !x!\n" +
		"git status\n" +
		"endlocal\n"
	if out != want {
		t.Fatalf("unexpected output:\n%s", out)
	}
}

func TestGenerate_Assign(t *testing.T) {
	prog := &ast.Program{Statements: []ast.Statement{
		&ast.SetStmt{Name: "a", Value: &ast.NumberLit{Value: "1", P: ast.Pos{Line: 1, Column: 7}}, P: ast.Pos{Line: 1, Column: 1}},
		&ast.AssignStmt{Name: "a", Value: &ast.NumberLit{Value: "2", P: ast.Pos{Line: 2, Column: 5}}, P: ast.Pos{Line: 2, Column: 1}},
	}}
	g := NewBatchGenerator()
	out, err := g.Generate(prog)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := "@echo off\nsetlocal EnableDelayedExpansion\nset a=1\nset a=2\nendlocal\n"
	if out != expected {
		t.Fatalf("unexpected output:\n%s", out)
	}
}

func TestGenerate_UnsupportedStmt(t *testing.T) {
	g := NewBatchGenerator()
	prog := &ast.Program{Statements: []ast.Statement{
		&ast.ReturnStmt{Value: &ast.NumberLit{Value: "1"}},
	}}

	_, err := g.Generate(prog)
	if err == nil {
		t.Fatalf("expected error for unsupported stmt")
	}

	if _, ok := err.(*GeneratorError); !ok {
		t.Fatalf("expected GeneratorError, got %T", err)
	}
}

func TestGenerate_FunctionNotLifted(t *testing.T) {
	g := NewBatchGenerator()
	fn := &ast.FnDecl{Name: "x"}

	if err := g.emitStmt(fn); err == nil {
		t.Fatalf("expected error for unlifted function")
	} else if _, ok := err.(*GeneratorError); !ok {
		t.Fatalf("expected GeneratorError, got %T", err)
	}
}

func TestGenerate_Call(t *testing.T) {
	g := NewBatchGenerator()
	prog := &ast.Program{Statements: []ast.Statement{
		&ast.FnDecl{
			Name:   "greet",
			Params: []string{"name"},
			Body: []ast.Statement{
				&ast.EchoStmt{Value: &ast.IdentExpr{Name: "name"}},
			},
		},
		&ast.CallStmt{Name: "greet", Args: []ast.Expr{&ast.StringLit{Value: "foo bar&baz"}}},
	}}

	out, err := g.Generate(prog)
	if err != nil {
		t.Fatalf("Generate returned error: %v", err)
	}

	want := "@echo off\n" +
		"setlocal EnableDelayedExpansion\n" +
		"call :fn_greet \"foo bar^&baz\"\n" +
		"goto :eof\n" +
		":fn_greet\n" +
		"setlocal EnableDelayedExpansion\n" +
		"set name=%1\n" +
		"set ret_greet_tmp_1=\n" +
		"    echo !name!\n" +
		":fn_ret_greet\n" +
		"endlocal & set fn_greet_ret=%ret_greet_tmp_1%\n" +
		"goto :eof\n" +
		"endlocal\n"

	if out != want {
		t.Fatalf("unexpected output:\n%s", out)
	}
}

func TestGenerate_Function(t *testing.T) {
	g := NewBatchGenerator()
	prog := &ast.Program{Statements: []ast.Statement{
		&ast.FnDecl{
			Name:   "greet",
			Params: []string{"name"},
			Body: []ast.Statement{
				&ast.EchoStmt{Value: &ast.StringLit{Value: "Hi"}},
				&ast.EchoStmt{Value: &ast.IdentExpr{Name: "name"}},
			},
		},
		// Top-level call should remain as-is (call lowering TBD), but function body must be emitted correctly.
	}}

	out, err := g.Generate(prog)
	if err != nil {
		t.Fatalf("Generate returned error: %v", err)
	}

	want := "@echo off\n" +
		"setlocal EnableDelayedExpansion\n" +
		"goto :eof\n" +
		":fn_greet\n" +
		"setlocal EnableDelayedExpansion\n" +
		"set name=%1\n" +
		"set ret_greet_tmp_1=\n" +
		"    echo Hi\n" +
		"    echo !name!\n" +
		":fn_ret_greet\n" +
		"endlocal & set fn_greet_ret=%ret_greet_tmp_1%\n" +
		"goto :eof\n" +
		"endlocal\n"

	if out != want {
		t.Fatalf("unexpected output:\n%s", out)
	}
}

func TestGenerate_IfElse(t *testing.T) {
	g := NewBatchGenerator()
	prog := &ast.Program{Statements: []ast.Statement{
		&ast.IfStmt{
			Cond: &ast.BoolLit{Value: true},
			Then: []ast.Statement{
				&ast.EchoStmt{Value: &ast.StringLit{Value: "yes"}},
			},
			Else: []ast.Statement{
				&ast.EchoStmt{Value: &ast.StringLit{Value: "no"}},
			},
		},
	}}

	out, err := g.Generate(prog)
	if err != nil {
		t.Fatalf("Generate returned error: %v", err)
	}

	want := "@echo off\n" +
		"setlocal EnableDelayedExpansion\n" +
		"if \"true\"==\"true\" (\n" +
		"    echo yes\n" +
		") else (\n" +
		"    echo no\n" +
		")\n" +
		"endlocal\n"

	if out != want {
		t.Fatalf("unexpected output:\n%s", out)
	}
}
