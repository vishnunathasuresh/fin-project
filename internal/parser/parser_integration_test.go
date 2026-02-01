package parser

import (
	"testing"

	"github.com/vishnunathasuresh/fin-project/internal/ast"
	"github.com/vishnunathasuresh/fin-project/internal/lexer"
)

func TestParseProgram_FullExample(t *testing.T) {
	src := `# sample fin v2 program
def greet(name: str) -> str:
    msg := "Hello " + name
    return msg

nums := [1, 2, 3]
for i .. 3
    x := i + 1

while true
    break

def a() -> int:
    if true
        y := 10
    return y
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
	if _, ok := prog.Statements[1].(*ast.DeclStmt); !ok {
		t.Fatalf("stmt1 not DeclStmt: %T", prog.Statements[1])
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
	src := "x := 1\n" +
		"if true\n" +
		"  y := x + 1\n" +
		"else\n" +
		"  y := x - 1\n"
	l := lexer.New(src)
	toks := CollectTokens(l)
	p := New(toks)
	prog := p.ParseProgram()
	if len(p.Errors()) != 0 {
		t.Fatalf("unexpected errors: %v", p.Errors())
	}
	out := ast.Format(prog)
	// Snapshot should contain DeclStmt nodes, not SetStmt
	if !contains(out, "DeclStmt") {
		t.Fatalf("snapshot missing DeclStmt node:\n%s", out)
	}
}

// contains is a helper to check if a string appears in another
func contains(s, substr string) bool {
	for i := 0; i < len(s)-len(substr)+1; i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func TestParseProgram_RecoveryThroughBadLine(t *testing.T) {
	src := `x := 1
if true
  y := 2
else
  z := 3

foo
`

	l := lexer.New(src)
	toks := CollectTokens(l)
	p := New(toks)
	prog := p.ParseProgram()

	if got := len(prog.Statements); got < 2 {
		t.Fatalf("stmt count = %d, want at least 2", got)
	}
	if _, ok := prog.Statements[0].(*ast.DeclStmt); !ok {
		t.Fatalf("stmt0 not DeclStmt: %T", prog.Statements[0])
	}
	if _, ok := prog.Statements[1].(*ast.IfStmt); !ok {
		t.Fatalf("stmt1 not IfStmt: %T", prog.Statements[1])
	}
}
