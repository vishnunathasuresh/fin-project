package generator

import (
	"strings"
	"testing"

	"github.com/vishnunath-suresh/fin-project/internal/ast"
)

func TestLowerSetStmt_Scalar(t *testing.T) {
	ctx := NewContext()
	lowerSetStmt(ctx, &ast.SetStmt{Name: "x", Value: &ast.NumberLit{Value: "10"}})
	want := "set x=10\n"
	if ctx.String() != want {
		t.Fatalf("unexpected output:\n%s", ctx.String())
	}
}

func TestLowerSetStmt_List(t *testing.T) {
	ctx := NewContext()
	lowerSetStmt(ctx, &ast.SetStmt{Name: "nums", Value: &ast.ListLit{Elements: []ast.Expr{
		&ast.NumberLit{Value: "10"},
		&ast.NumberLit{Value: "20"},
	}}})
	want := strings.Join([]string{
		"set nums_0=10",
		"set nums_1=20",
		"set nums_len=2",
		"",
	}, "\n")
	if ctx.String() != want {
		t.Fatalf("unexpected output:\n%s", ctx.String())
	}
}

func TestLowerSetStmt_Map(t *testing.T) {
	ctx := NewContext()
	lowerSetStmt(ctx, &ast.SetStmt{Name: "user", Value: &ast.MapLit{Pairs: []ast.MapPair{
		{Key: "name", Value: &ast.StringLit{Value: "bob"}},
	}}})
	want := "set user_name=bob\n"
	if ctx.String() != want {
		t.Fatalf("unexpected output:\n%s", ctx.String())
	}
}

func TestLowerEchoStmt(t *testing.T) {
	ctx := NewContext()
	lowerEchoStmt(ctx, &ast.EchoStmt{Value: &ast.IdentExpr{Name: "name"}})
	want := "echo %name%\n"
	if ctx.String() != want {
		t.Fatalf("unexpected output:\n%s", ctx.String())
	}
}

func TestLowerRunStmt(t *testing.T) {
	ctx := NewContext()
	lowerRunStmt(ctx, &ast.RunStmt{Command: &ast.StringLit{Value: "git status"}})
	want := "git status\n"
	if ctx.String() != want {
		t.Fatalf("unexpected output:\n%s", ctx.String())
	}
}
