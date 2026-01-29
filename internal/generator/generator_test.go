package generator

import (
	"testing"

	"github.com/vishnunath-suresh/fin-project/internal/ast"
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
		"set x=10\n" +
		"echo %x%\n" +
		"git status\n"

	if out != want {
		t.Fatalf("unexpected output:\n%s", out)
	}
}
