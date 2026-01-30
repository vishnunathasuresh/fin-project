package generator

import (
	"errors"
	"testing"

	"github.com/vishnunath-suresh/fin-project/internal/ast"
	"github.com/vishnunath-suresh/fin-project/internal/lexer"
	"github.com/vishnunath-suresh/fin-project/internal/parser"
	"github.com/vishnunath-suresh/fin-project/internal/sema"
)

type goldenCase struct {
	name     string
	fin      string
	expected string
}

func TestGenerator_Golden(t *testing.T) {
	cases := []goldenCase{
		{
			name: "set_echo_call_fn",
			fin: "set x 1\n" +
				"greet \"Bob\"\n" +
				"fn greet name\n" +
				"    echo $name\n" +
				"end\n",
			expected: "@echo off\n" +
				"set x=1\n" +
				"call :fn_greet Bob\n" +
				"goto :eof\n" +
				":fn_greet\n" +
				"setlocal\n" +
				"set name=%1\n" +
				"set ret_greet_tmp_1=\n" +
				"    echo %name%\n" +
				":fn_ret_greet\n" +
				"endlocal & set fn_greet_ret=%ret_greet_tmp_1%\n" +
				"goto :eof\n",
		},
		{
			name: "control_flow_mix",
			fin: "set total 0\n" +
				"for i in 1..3\n" +
				"    echo $i\n" +
				"end\n" +
				"while false\n" +
				"    echo loop\n" +
				"end\n",
			expected: "@echo off\n" +
				"set total=0\n" +
				"set i=1\n" +
				":loop_continue_1\n" +
				"if %i% GTR 3 goto loop_break_1\n" +
				"    echo %i%\n" +
				"set /a i=%i%+1\n" +
				"goto loop_continue_1\n" +
				":loop_break_1\n" +
				":while_start_2\n" +
				"if not false goto while_end_2\n" +
				"echo %loop%\n" +
				"goto while_start_2\n" +
				":while_end_2\n",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			out := generateFromSource(t, tc.fin)
			if out != tc.expected {
				t.Fatalf("golden mismatch\nwant:\n%q\n\nhave:\n%q", tc.expected, out)
			}
		})
	}
}

func TestGenerator_Golden_Negative_SemaError(t *testing.T) {
	// Duplicate function should be caught by sema; generator should not run.
	prog := &ast.Program{Statements: []ast.Statement{
		&ast.FnDecl{Name: "foo", P: ast.Pos{Line: 1, Column: 1}},
		&ast.FnDecl{Name: "foo", P: ast.Pos{Line: 2, Column: 1}},
	}}

	a := sema.New()
	if err := a.Analyze(prog); err == nil {
		t.Fatalf("expected semantic error for duplicate function")
	} else {
		var df sema.DuplicateFunctionError
		if !errors.As(err, &df) {
			t.Fatalf("expected DuplicateFunctionError, got %T", err)
		}
	}
}

func generateFromSource(t *testing.T, src string) string {
	t.Helper()

	l := lexer.New(src)
	tokens := parser.CollectTokens(l)
	p := parser.New(tokens)
	prog := p.ParseProgram()
	if errs := p.Errors(); len(errs) > 0 {
		t.Fatalf("parse errors: %v", errs)
	}

	g := NewBatchGenerator()
	out, err := g.Generate(prog)
	if err != nil {
		t.Fatalf("generate error: %v", err)
	}
	return out
}

func generateFromSourceWithError(t *testing.T, src string) (string, error) {
	t.Helper()

	l := lexer.New(src)
	tokens := parser.CollectTokens(l)
	p := parser.New(tokens)
	prog := p.ParseProgram()
	if errs := p.Errors(); len(errs) > 0 {
		return "", errs[0]
	}

	// Run semantic analysis; expect errors to surface here for negative cases.
	a := sema.New()
	if err := a.Analyze(prog); err != nil {
		return "", err
	}

	g := NewBatchGenerator()
	return g.Generate(prog)
}

func generateDirect(prog *ast.Program) (string, error) {
	g := NewBatchGenerator()
	return g.Generate(prog)
}
