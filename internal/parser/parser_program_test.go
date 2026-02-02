package parser

import (
	"testing"

	"github.com/vishnunathasuresh/fin-project/internal/ast"
	"github.com/vishnunathasuresh/fin-project/internal/lexer"
)

func TestParseProgram_TopLevelLoop(t *testing.T) {
	tests := []struct {
		name          string
		src           string
		wantStmtCount int
		wantErrors    int
	}{
		{name: "empty file", src: "", wantStmtCount: 0, wantErrors: 0},
		{name: "only comments", src: "# hi\n# there\n", wantStmtCount: 0, wantErrors: 0},
		{name: "multiple blank lines", src: "\n\n\n", wantStmtCount: 0, wantErrors: 0},
		{name: "one call statement", src: "foo\n", wantStmtCount: 1, wantErrors: 0},
		{name: "syntax error recovery", src: "!\nbar\n", wantStmtCount: 1, wantErrors: 1},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.New(tt.src)
			toks := CollectTokens(l)
			p := New(toks)

			prog := p.ParseProgram()
			if prog == nil {
				t.Fatalf("ParseProgram returned nil")
			}

			if got := len(prog.Statements); got != tt.wantStmtCount {
				t.Fatalf("statement count = %d, want %d", got, tt.wantStmtCount)
			}

			if gotErrs := len(p.Errors()); gotErrs != tt.wantErrors {
				t.Fatalf("errors = %d, want %d", gotErrs, tt.wantErrors)
			}
		})
	}
}

func TestParseProgram_StressDeepNesting_Short(t *testing.T) {
	// TODO(fin-v2): extend once nested run()/call lowering is supported.
	src := `def a():
	    if true
	        while true
	            for i .. 3
	                x := 1
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

func TestParseProgram_StopsOnlyOnEOF(t *testing.T) {
	src := "foo\n# c\nbar\n"
	l := lexer.New(src)
	toks := CollectTokens(l)
	p := New(toks)

	prog := p.ParseProgram()
	if prog == nil {
		t.Fatalf("ParseProgram returned nil")
	}
	if len(prog.Statements) != 2 {
		t.Fatalf("statement count = %d, want 2", len(prog.Statements))
	}
	names := make([]string, 0, len(prog.Statements))
	for _, s := range prog.Statements {
		switch st := s.(type) {
		case *ast.CallStmt:
			names = append(names, st.Name)
		default:
			t.Fatalf("unexpected stmt type: %T", st)
		}
	}
	if len(names) != 2 || names[0] != "foo" || names[1] != "bar" {
		t.Fatalf("statement names = %v, want [foo bar]", names)
	}
	if len(p.Errors()) != 0 {
		t.Fatalf("expected no errors, got %d", len(p.Errors()))
	}

	// Ensure parser is at EOF without panic.
	if !p.isAtEnd() {
		t.Fatalf("parser not at EOF after ParseProgram")
	}
}

func TestParseProgram_ErrorRecovery_MissingEnd(t *testing.T) {
	src := "if true\n  x = 1\n"
	l := lexer.New(src)
	toks := CollectTokens(l)
	p := New(toks)
	prog := p.ParseProgram()
	// Should parse the if statement even without proper block termination
	if len(prog.Statements) < 1 {
		t.Fatalf("got %d statements, want at least 1", len(prog.Statements))
	}
	if len(p.Errors()) == 0 {
		t.Fatalf("expected errors for malformed if block")
	}
}
