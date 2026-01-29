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

func TestLowerIfStmt_Nested(t *testing.T) {
	ctx := NewContext()
	lowerIfStmt(ctx, &ast.IfStmt{
		Cond: &ast.BoolLit{Value: true},
		Then: []ast.Statement{
			&ast.IfStmt{
				Cond: &ast.BoolLit{Value: false},
				Then: []ast.Statement{
					&ast.EchoStmt{Value: &ast.StringLit{Value: "inner-then"}},
				},
				Else: []ast.Statement{
					&ast.EchoStmt{Value: &ast.StringLit{Value: "inner-else"}},
				},
			},
		},
		Else: []ast.Statement{
			&ast.EchoStmt{Value: &ast.StringLit{Value: "outer-else"}},
		},
	}, func(st ast.Statement) {
		switch s := st.(type) {
		case *ast.EchoStmt:
			lowerEchoStmt(ctx, s)
		case *ast.IfStmt:
			lowerIfStmt(ctx, s, func(n ast.Statement) {
				switch x := n.(type) {
				case *ast.EchoStmt:
					lowerEchoStmt(ctx, x)
				default:
					t.Fatalf("unexpected nested stmt %T", x)
				}
			})
		default:
			t.Fatalf("unexpected stmt type %T", s)
		}
	})

	want := strings.Join([]string{
		"if true (",
		"    if false (",
		"        echo inner-then",
		"    ) else (",
		"        echo inner-else",
		"    )",
		") else (",
		"    echo outer-else",
		")",
		"",
	}, "\n")

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

func TestLowerIfStmt_WithElse(t *testing.T) {
	ctx := NewContext()
	lowerIfStmt(ctx, &ast.IfStmt{
		Cond: &ast.BoolLit{Value: true},
		Then: []ast.Statement{
			&ast.EchoStmt{Value: &ast.StringLit{Value: "yes"}},
		},
		Else: []ast.Statement{
			&ast.EchoStmt{Value: &ast.StringLit{Value: "no"}},
		},
	}, func(st ast.Statement) {
		switch s := st.(type) {
		case *ast.EchoStmt:
			lowerEchoStmt(ctx, s)
		default:
			t.Fatalf("unexpected stmt type %T", s)
		}
	})

	want := strings.Join([]string{
		"if true (",
		"    echo yes",
		") else (",
		"    echo no",
		")",
		"",
	}, "\n")

	if ctx.String() != want {
		t.Fatalf("unexpected output:\n%s", ctx.String())
	}
}
