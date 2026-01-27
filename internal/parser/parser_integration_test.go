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
