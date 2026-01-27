package parser

import (
	"testing"

	"github.com/vishnunath-suresh/fin-project/internal/ast"
	"github.com/vishnunath-suresh/fin-project/internal/lexer"
)

func parseProgram(t *testing.T, src string) *ast.Program {
	t.Helper()
	l := lexer.New(src)
	toks := CollectTokens(l)
	p := New(toks)
	return p.ParseProgram()
}

func TestParse_SetEchoRunCallReturn(t *testing.T) {
	src := "set a 1\necho $a\nrun \"cmd\"\nfoo 1 2\nreturn 42\n"
	prog := parseProgram(t, src)
	if len(prog.Statements) != 5 {
		t.Fatalf("got %d statements, want 5", len(prog.Statements))
	}
	if setStmt, ok := prog.Statements[0].(*ast.SetStmt); !ok || setStmt.Name != "a" {
		t.Fatalf("stmt0 not set a: %T", prog.Statements[0])
	}
	if _, ok := prog.Statements[1].(*ast.EchoStmt); !ok {
		t.Fatalf("stmt1 not echo: %T", prog.Statements[1])
	}
	if _, ok := prog.Statements[2].(*ast.RunStmt); !ok {
		t.Fatalf("stmt2 not run: %T", prog.Statements[2])
	}
	if c, ok := prog.Statements[3].(*ast.CallStmt); !ok || c.Name != "foo" || len(c.Args) != 2 {
		t.Fatalf("stmt3 not call foo with 2 args: %T", prog.Statements[3])
	}
	if _, ok := prog.Statements[4].(*ast.ReturnStmt); !ok {
		t.Fatalf("stmt4 not return: %T", prog.Statements[4])
	}
}

func TestParse_IfElse(t *testing.T) {
	src := "if exists \"a\"\nfoo\nelse\nbar\nend\n"
	prog := parseProgram(t, src)
	if len(prog.Statements) != 1 {
		t.Fatalf("got %d statements, want 1", len(prog.Statements))
	}
	ifStmt, ok := prog.Statements[0].(*ast.IfStmt)
	if !ok {
		t.Fatalf("stmt not IfStmt: %T", prog.Statements[0])
	}
	if len(ifStmt.Then) != 1 || len(ifStmt.Else) != 1 {
		t.Fatalf("then/else sizes wrong: %d/%d", len(ifStmt.Then), len(ifStmt.Else))
	}
}

func TestParse_For(t *testing.T) {
	src := "for i in 1..3\nfoo\nend\n"
	prog := parseProgram(t, src)
	if len(prog.Statements) != 1 {
		t.Fatalf("got %d statements, want 1", len(prog.Statements))
	}
	f, ok := prog.Statements[0].(*ast.ForStmt)
	if !ok {
		t.Fatalf("stmt not ForStmt: %T", prog.Statements[0])
	}
	if f.Var != "i" {
		t.Fatalf("for var = %q, want i", f.Var)
	}
}

func TestParse_While(t *testing.T) {
	src := "while 1\nfoo\nend\n"
	prog := parseProgram(t, src)
	if len(prog.Statements) != 1 {
		t.Fatalf("got %d statements, want 1", len(prog.Statements))
	}
	if _, ok := prog.Statements[0].(*ast.WhileStmt); !ok {
		t.Fatalf("stmt not WhileStmt: %T", prog.Statements[0])
	}
}

func TestParse_Fn(t *testing.T) {
	src := "fn add x y\nreturn x + y\nend\n"
	prog := parseProgram(t, src)
	if len(prog.Statements) != 1 {
		t.Fatalf("got %d statements, want 1", len(prog.Statements))
	}
	fn, ok := prog.Statements[0].(*ast.FnDecl)
	if !ok {
		t.Fatalf("stmt not FnDecl: %T", prog.Statements[0])
	}
	if fn.Name != "add" || len(fn.Params) != 2 {
		t.Fatalf("fn name/params wrong: %q %d", fn.Name, len(fn.Params))
	}
	if len(fn.Body) != 1 {
		t.Fatalf("fn body size wrong: %d", len(fn.Body))
	}
}
