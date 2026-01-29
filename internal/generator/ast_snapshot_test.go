package generator

import (
	"testing"

	"github.com/vishnunath-suresh/fin-project/internal/ast"
	"github.com/vishnunath-suresh/fin-project/internal/lexer"
	"github.com/vishnunath-suresh/fin-project/internal/parser"
)

func TestASTSnapshot_SimpleProgram(t *testing.T) {
	src := "set x 1\n" +
		"fn greet name\n" +
		"echo $name\n" +
		"end\n"

	l := lexer.New(src)
	tokens := parser.CollectTokens(l)
	p := parser.New(tokens)
	prog := p.ParseProgram()
	if errs := p.Errors(); len(errs) > 0 {
		t.Fatalf("parse errors: %v", errs)
	}

	got := ast.Format(prog)
	want := "Program @1:1\n" +
		"  SetStmt name=x @1:1\n" +
		"    value: NumberLit 1 @1:7\n" +
		"  FnDecl name=greet params=[name] @2:1\n" +
		"    body: EchoStmt @3:1\n" +
		"      value: IdentExpr name @3:6\n"

	if got != want {
		t.Fatalf("AST snapshot mismatch\nwant:\n%s\n\ngot:\n%s", want, got)
	}
}
