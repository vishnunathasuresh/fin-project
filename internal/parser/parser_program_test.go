package parser

import (
	"testing"

	"github.com/vishnunath-suresh/fin-project/internal/ast"
	"github.com/vishnunath-suresh/fin-project/internal/lexer"
)

func TestParseProgram_TopLevelLoop(t *testing.T) {
	tests := []struct {
		name          string
		src           string
		wantStmtCount int
		wantErrors    int
		wantFirstName string
	}{
		{name: "empty file", src: "", wantStmtCount: 0, wantErrors: 0},
		{name: "only comments", src: "# hi\n# there\n", wantStmtCount: 0, wantErrors: 0},
		{name: "multiple blank lines", src: "\n\n\n", wantStmtCount: 0, wantErrors: 0},
		{name: "one call statement", src: "echo\n", wantStmtCount: 1, wantErrors: 0, wantFirstName: "echo"},
		{name: "syntax error recovery", src: "!\nfoo\n", wantStmtCount: 1, wantErrors: 1, wantFirstName: "foo"},
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

			if tt.wantStmtCount > 0 {
				switch s := prog.Statements[0].(type) {
				case *ast.CallStmt:
					if s.Name != tt.wantFirstName {
						t.Fatalf("call name = %q, want %q", s.Name, tt.wantFirstName)
					}
				case *ast.EchoStmt:
					// echo keyword handled as EchoStmt; accept when expecting echo
					if tt.wantFirstName != "echo" {
						t.Fatalf("first stmt unexpected EchoStmt")
					}
				default:
					t.Fatalf("first stmt unexpected type: %T", s)
				}
			}

			if gotErrs := len(p.Errors()); gotErrs != tt.wantErrors {
				t.Fatalf("errors = %d, want %d", gotErrs, tt.wantErrors)
			}
		})
	}
}

func TestParseProgram_StopsOnlyOnEOF(t *testing.T) {
	src := "echo\n# c\nbar\n"
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
		case *ast.EchoStmt:
			names = append(names, "echo")
		default:
			t.Fatalf("unexpected stmt type: %T", st)
		}
	}
	if len(names) != 2 || names[0] != "echo" || names[1] != "bar" {
		t.Fatalf("statement names = %v, want [echo bar]", names)
	}
	if len(p.Errors()) != 0 {
		t.Fatalf("expected no errors, got %d", len(p.Errors()))
	}

	// Ensure parser is at EOF without panic.
	if !p.isAtEnd() {
		t.Fatalf("parser not at EOF after ParseProgram")
	}
}
