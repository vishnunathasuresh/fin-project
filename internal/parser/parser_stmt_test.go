package parser

import (
	"testing"

	"github.com/vishnunathasuresh/fin-project/internal/ast"
	"github.com/vishnunathasuresh/fin-project/internal/lexer"
)

func parseProgram(t *testing.T, src string) *ast.Program {
	t.Helper()
	l := lexer.New(src)
	toks := CollectTokens(l)
	p := New(toks)
	return p.ParseProgram()
}

func parseProgramWithParser(t *testing.T, src string) (*ast.Program, *Parser) {
	t.Helper()
	l := lexer.New(src)
	toks := CollectTokens(l)
	p := New(toks)
	return p.ParseProgram(), p
}

// ---- Declaration vs Assignment Tests ----

// TestParse_DeclStmt_Simple parses "name := expr"
func TestParse_DeclStmt_Simple(t *testing.T) {
	src := "x := 10\n"
	prog := parseProgram(t, src)
	if len(prog.Statements) != 1 {
		t.Fatalf("got %d statements, want 1", len(prog.Statements))
	}
	decl, ok := prog.Statements[0].(*ast.DeclStmt)
	if !ok {
		t.Fatalf("stmt not DeclStmt: %T", prog.Statements[0])
	}
	if decl.Name != "x" {
		t.Fatalf("decl name = %q, want x", decl.Name)
	}
}

// TestParse_DeclStmt_WithString tests declaration with string value
func TestParse_DeclStmt_WithString(t *testing.T) {
	src := "name := \"fin\"\n"
	prog := parseProgram(t, src)
	if len(prog.Statements) != 1 {
		t.Fatalf("got %d statements, want 1", len(prog.Statements))
	}
	decl, ok := prog.Statements[0].(*ast.DeclStmt)
	if !ok {
		t.Fatalf("stmt not DeclStmt: %T", prog.Statements[0])
	}
	if decl.Name != "name" {
		t.Fatalf("decl name = %q, want name", decl.Name)
	}
}

// TestParse_DeclStmt_Multiple tests multiple declarations
func TestParse_DeclStmt_Multiple(t *testing.T) {
	src := "x := 1\ny := 2\nz := 3\n"
	prog := parseProgram(t, src)
	if len(prog.Statements) != 3 {
		t.Fatalf("got %d statements, want 3", len(prog.Statements))
	}
	for i, name := range []string{"x", "y", "z"} {
		decl, ok := prog.Statements[i].(*ast.DeclStmt)
		if !ok {
			t.Fatalf("stmt %d not DeclStmt: %T", i, prog.Statements[i])
		}
		if decl.Name != name {
			t.Fatalf("decl %d name = %q, want %q", i, decl.Name, name)
		}
	}
}

// TestParse_AssignStmt_Simple parses "name = expr"
func TestParse_AssignStmt_Simple(t *testing.T) {
	src := "x = 20\n"
	prog := parseProgram(t, src)
	if len(prog.Statements) != 1 {
		t.Fatalf("got %d statements, want 1", len(prog.Statements))
	}
	assign, ok := prog.Statements[0].(*ast.AssignStmt)
	if !ok {
		t.Fatalf("stmt not AssignStmt: %T", prog.Statements[0])
	}
	if assign.Name != "x" {
		t.Fatalf("assign name = %q, want x", assign.Name)
	}
}

// TestParse_DeclVsAssign tests the critical distinction between := and =
func TestParse_DeclVsAssign(t *testing.T) {
	src := "x := 1\nx = 2\n"
	prog := parseProgram(t, src)
	if len(prog.Statements) != 2 {
		t.Fatalf("got %d statements, want 2", len(prog.Statements))
	}

	decl, ok := prog.Statements[0].(*ast.DeclStmt)
	if !ok {
		t.Fatalf("stmt0 not DeclStmt: %T", prog.Statements[0])
	}
	if decl.Name != "x" {
		t.Fatalf("decl name = %q, want x", decl.Name)
	}

	assign, ok := prog.Statements[1].(*ast.AssignStmt)
	if !ok {
		t.Fatalf("stmt1 not AssignStmt: %T", prog.Statements[1])
	}
	if assign.Name != "x" {
		t.Fatalf("assign name = %q, want x", assign.Name)
	}
}

// TestParse_DeclWithExpression tests declaration with complex expression
func TestParse_DeclWithExpression(t *testing.T) {
	src := "result := 1 + 2\n"
	prog := parseProgram(t, src)
	if len(prog.Statements) != 1 {
		t.Fatalf("got %d statements, want 1", len(prog.Statements))
	}
	decl, ok := prog.Statements[0].(*ast.DeclStmt)
	if !ok {
		t.Fatalf("stmt not DeclStmt: %T", prog.Statements[0])
	}
	if decl.Value == nil {
		t.Fatalf("decl value is nil")
	}
}

// TestParse_AssignWithExpression tests assignment with complex expression
func TestParse_AssignWithExpression(t *testing.T) {
	src := "x = y + 1\n"
	prog := parseProgram(t, src)
	if len(prog.Statements) != 1 {
		t.Fatalf("got %d statements, want 1", len(prog.Statements))
	}
	assign, ok := prog.Statements[0].(*ast.AssignStmt)
	if !ok {
		t.Fatalf("stmt not AssignStmt: %T", prog.Statements[0])
	}
	if assign.Value == nil {
		t.Fatalf("assign value is nil")
	}
}

// TestParse_MixedDeclAndAssign tests mixing declarations and assignments
func TestParse_MixedDeclAndAssign(t *testing.T) {
	src := "x := 10\ny := 20\nx = x + 1\ny = y + 2\n"
	prog := parseProgram(t, src)
	if len(prog.Statements) != 4 {
		t.Fatalf("got %d statements, want 4", len(prog.Statements))
	}

	if _, ok := prog.Statements[0].(*ast.DeclStmt); !ok {
		t.Fatalf("stmt0 not DeclStmt: %T", prog.Statements[0])
	}
	if _, ok := prog.Statements[1].(*ast.DeclStmt); !ok {
		t.Fatalf("stmt1 not DeclStmt: %T", prog.Statements[1])
	}
	if _, ok := prog.Statements[2].(*ast.AssignStmt); !ok {
		t.Fatalf("stmt2 not AssignStmt: %T", prog.Statements[2])
	}
	if _, ok := prog.Statements[3].(*ast.AssignStmt); !ok {
		t.Fatalf("stmt3 not AssignStmt: %T", prog.Statements[3])
	}
}

// TestParse_Call_NoArgs tests function call without arguments
func TestParse_Call_NoArgs(t *testing.T) {
	src := "foo\n"
	prog := parseProgram(t, src)
	if len(prog.Statements) != 1 {
		t.Fatalf("got %d statements, want 1", len(prog.Statements))
	}
	call, ok := prog.Statements[0].(*ast.CallStmt)
	if !ok || call.Name != "foo" || len(call.Args) != 0 {
		t.Fatalf("call stmt wrong: %T name=%q args=%d", prog.Statements[0], call.Name, len(call.Args))
	}
}

// TestParse_Call_ManyArgs tests function call with multiple arguments
func TestParse_Call_ManyArgs(t *testing.T) {
	src := "foo 1 2 3\n"
	prog := parseProgram(t, src)
	if len(prog.Statements) != 1 {
		t.Fatalf("got %d statements, want 1", len(prog.Statements))
	}
	call, ok := prog.Statements[0].(*ast.CallStmt)
	if !ok || call.Name != "foo" || len(call.Args) != 3 {
		t.Fatalf("call stmt wrong: %T name=%q args=%d", prog.Statements[0], call.Name, len(call.Args))
	}
}

// TestParse_Return tests return statement
func TestParse_Return(t *testing.T) {
	src := "return 42\n"
	prog := parseProgram(t, src)
	if len(prog.Statements) != 1 {
		t.Fatalf("got %d statements, want 1", len(prog.Statements))
	}
	if _, ok := prog.Statements[0].(*ast.ReturnStmt); !ok {
		t.Fatalf("stmt not ReturnStmt: %T", prog.Statements[0])
	}
}

// TestParse_IfElse tests if/else control flow
func TestParse_IfElse(t *testing.T) {
	src := "if true\n  x = 1\nelse\n  x = 2\n"
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

// TestParse_For tests for loop
func TestParse_For(t *testing.T) {
	src := "for i .. 3\n  x = x + 1\n"
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

// TestParse_While tests while loop
func TestParse_While(t *testing.T) {
	src := "while true\n  x = x + 1\n"
	prog := parseProgram(t, src)
	if len(prog.Statements) != 1 {
		t.Fatalf("got %d statements, want 1", len(prog.Statements))
	}
	if _, ok := prog.Statements[0].(*ast.WhileStmt); !ok {
		t.Fatalf("stmt not WhileStmt: %T", prog.Statements[0])
	}
}

// TestParse_While_Nested tests nested while loops
func TestParse_While_Nested(t *testing.T) {
	src := "while true\n  while false\n    x = 1\n"
	prog := parseProgram(t, src)
	if len(prog.Statements) != 1 {
		t.Fatalf("got %d statements, want 1", len(prog.Statements))
	}
	outer, ok := prog.Statements[0].(*ast.WhileStmt)
	if !ok {
		t.Fatalf("stmt not WhileStmt: %T", prog.Statements[0])
	}
	if len(outer.Body) != 1 {
		t.Fatalf("outer body len = %d, want 1", len(outer.Body))
	}
	if _, ok := outer.Body[0].(*ast.WhileStmt); !ok {
		t.Fatalf("inner not WhileStmt: %T", outer.Body[0])
	}
}

// TestParse_BreakAndContinue tests break and continue
func TestParse_BreakAndContinue(t *testing.T) {
	src := "while true\n  break\n  continue\n"
	prog := parseProgram(t, src)
	if len(prog.Statements) != 1 {
		t.Fatalf("got %d statements, want 1", len(prog.Statements))
	}
	w, ok := prog.Statements[0].(*ast.WhileStmt)
	if !ok {
		t.Fatalf("stmt not WhileStmt: %T", prog.Statements[0])
	}
	if len(w.Body) != 2 {
		t.Fatalf("while body len = %d, want 2", len(w.Body))
	}
	if _, ok := w.Body[0].(*ast.BreakStmt); !ok {
		t.Fatalf("first not BreakStmt: %T", w.Body[0])
	}
	if _, ok := w.Body[1].(*ast.ContinueStmt); !ok {
		t.Fatalf("second not ContinueStmt: %T", w.Body[1])
	}
}
