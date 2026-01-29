package generator

import (
	"testing"

	"github.com/vishnunath-suresh/fin-project/internal/lexer"
	"github.com/vishnunath-suresh/fin-project/internal/parser"
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
				"    echo %name%\n" +
				"endlocal\n" +
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
				"for /L %i in (1,1,3) do (\n" +
				"    echo %i%\n" +
				")\n" +
				":while_start_1\n" +
				"if not false goto while_end_1\n" +
				"echo %loop%\n" +
				"goto while_start_1\n" +
				":while_end_1\n",
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
