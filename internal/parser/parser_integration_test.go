package parser

import (
	"testing"

	"github.com/vishnunath-suresh/fin-project/internal/ast"
	"github.com/vishnunath-suresh/fin-project/internal/lexer"
)

func TestParseProgram_FullExample(t *testing.T) {
	src := `# sample fin program
fn greet name
    set msg "Hello $name"
    echo $msg
end

set nums [1,2,3]
for i in 1..3
    echo $i
end

while 1
    break
    continue
end
fn a
    if true
        while true
            for i in 1..3
                echo "x"
            end
        end
    end
end
`

	l := lexer.New(src)
	toks := CollectTokens(l)
	p := New(toks)
	prog := p.ParseProgram()

	if len(p.Errors()) != 0 {
		t.Fatalf("unexpected parse errors: %v", p.Errors())
	}

	if len(prog.Statements) != 5 {
		t.Fatalf("stmt count = %d, want 5", len(prog.Statements))
	}

	if _, ok := prog.Statements[0].(*ast.FnDecl); !ok {
		t.Fatalf("stmt0 not FnDecl: %T", prog.Statements[0])
	}
	if _, ok := prog.Statements[1].(*ast.SetStmt); !ok {
		t.Fatalf("stmt1 not SetStmt: %T", prog.Statements[1])
	}
	if _, ok := prog.Statements[2].(*ast.ForStmt); !ok {
		t.Fatalf("stmt2 not ForStmt: %T", prog.Statements[2])
	}
	if _, ok := prog.Statements[3].(*ast.WhileStmt); !ok {
		t.Fatalf("stmt3 not WhileStmt: %T", prog.Statements[3])
	}
	if _, ok := prog.Statements[4].(*ast.FnDecl); !ok {
		t.Fatalf("stmt4 not FnDecl: %T", prog.Statements[4])
	}
}

func TestParseProgram_StressDeepNesting(t *testing.T) {
	src := `fn a
    if true
        while true
            for i in 1..3
                echo "x"
            end
        end
    end
end
`
	l := lexer.New(src)
	toks := CollectTokens(l)
	p := New(toks)
	prog := p.ParseProgram()
	if len(p.Errors()) != 0 {
		t.Fatalf("unexpected errors: %v", p.Errors())
	}
	if len(prog.Statements) != 1 {
		t.Fatalf("stmt count = %d, want 1", len(prog.Statements))
	}
}

func TestParseProgram_StressLongExpression(t *testing.T) {
	src := `set x 1 + 2 * 3 == 7 && true || false
`
	l := lexer.New(src)
	toks := CollectTokens(l)
	p := New(toks)
	prog := p.ParseProgram()
	if len(p.Errors()) != 0 {
		t.Fatalf("unexpected errors: %v", p.Errors())
	}
	if len(prog.Statements) != 1 {
		t.Fatalf("stmt count = %d, want 1", len(prog.Statements))
	}
}

func TestParseProgram_StressRecovery(t *testing.T) {
	src := `set x
echo
fn test
    set a 1

`
	l := lexer.New(src)
	toks := CollectTokens(l)
	p := New(toks)
	prog := p.ParseProgram()
	if len(p.Errors()) == 0 {
		t.Fatalf("expected errors, got none")
	}
	if prog == nil || len(prog.Statements) == 0 {
		t.Fatalf("expected AST with statements, got nil/empty")
	}
}

func TestParseProgram_Snapshot(t *testing.T) {
	src := "set a 1\n" +
		"if exists \"f\"\n" +
		"    echo $a\n" +
		"else\n" +
		"    run \"cmd\"\n" +
		"end\n"
	l := lexer.New(src)
	toks := CollectTokens(l)
	p := New(toks)
	prog := p.ParseProgram()
	if len(p.Errors()) != 0 {
		t.Fatalf("unexpected errors: %v", p.Errors())
	}
	out := ast.Format(prog)
	want := "" +
		"Program @1:1\n" +
		"  SetStmt name=a @1:1\n" +
		"    value: NumberLit 1 @1:7\n" +
		"  IfStmt @2:1\n" +
		"    cond: ExistsCond @2:4\n" +
		"      path: StringLit \"f\" @2:11\n" +
		"    then:\n" +
		"      EchoStmt @3:5\n" +
		"        value: IdentExpr a @3:10\n" +
		"    else:\n" +
		"      RunStmt @5:5\n" +
		"        command: StringLit \"cmd\" @5:9\n"
	if out != want {
		t.Fatalf("snapshot mismatch:\nwant:\n%s\ngot:\n%s", want, out)
	}
}

func TestParseProgram_RecoveryThroughBadLine(t *testing.T) {
	src := `set a 1
if exists "file"
    set b 2
end
???
echo "after error"
`

	l := lexer.New(src)
	toks := CollectTokens(l)
	p := New(toks)
	prog := p.ParseProgram()

	if len(p.Errors()) == 0 {
		t.Fatalf("expected errors but got none")
	}
	if got := len(prog.Statements); got != 3 {
		t.Fatalf("stmt count = %d, want 3 (set, if, echo)", got)
	}
	if _, ok := prog.Statements[2].(*ast.EchoStmt); !ok {
		t.Fatalf("last stmt not EchoStmt: %T", prog.Statements[2])
	}
}
