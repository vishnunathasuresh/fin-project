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

func TestParse_ForElse(t *testing.T) {
	src := "for i .. 3\n  x := i\nelse\n  y := 0\n"
	prog := parseProgram(t, src)
	if len(prog.Statements) != 1 {
		t.Fatalf("got %d statements, want 1", len(prog.Statements))
	}
	forStmt, ok := prog.Statements[0].(*ast.ForStmt)
	if !ok {
		t.Fatalf("stmt not ForStmt: %T", prog.Statements[0])
	}
	if len(forStmt.Body) != 1 {
		t.Fatalf("body len = %d, want 1", len(forStmt.Body))
	}
	if len(forStmt.Else) != 1 {
		t.Fatalf("else len = %d, want 1", len(forStmt.Else))
	}
}

func parseProgramWithParser(t *testing.T, src string) (*ast.Program, *Parser) {
	t.Helper()
	l := lexer.New(src)
	toks := CollectTokens(l)
	p := New(toks)
	return p.ParseProgram(), p
}

// ---- Declaration vs Assignment Tests ----

func TestParse_IfElifElse(t *testing.T) {
	src := "if a\n  x := 1\nelif b\n  x := 2\nelse\n  x := 3\n"
	prog := parseProgram(t, src)
	if len(prog.Statements) != 1 {
		t.Fatalf("got %d statements, want 1", len(prog.Statements))
	}
	ifStmt, ok := prog.Statements[0].(*ast.IfStmt)
	if !ok {
		t.Fatalf("stmt not IfStmt: %T", prog.Statements[0])
	}
	if len(ifStmt.Then) != 1 {
		t.Fatalf("then len = %d, want 1", len(ifStmt.Then))
	}
	if len(ifStmt.Else) != 1 {
		t.Fatalf("else len = %d, want 1 (elif as nested if)", len(ifStmt.Else))
	}
	elifStmt, ok := ifStmt.Else[0].(*ast.IfStmt)
	if !ok {
		t.Fatalf("elif node not IfStmt: %T", ifStmt.Else[0])
	}
	if len(elifStmt.Then) != 1 {
		t.Fatalf("elif then len = %d, want 1", len(elifStmt.Then))
	}
	if len(elifStmt.Else) != 1 {
		t.Fatalf("elif else len = %d, want 1", len(elifStmt.Else))
	}
}

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
	if len(decl.Names) != 1 || decl.Names[0] != "x" {
		t.Fatalf("decl names = %v, want [x]", decl.Names)
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
	if len(decl.Names) != 1 || decl.Names[0] != "name" {
		t.Fatalf("decl names = %v, want [name]", decl.Names)
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
		if len(decl.Names) != 1 || decl.Names[0] != name {
			t.Fatalf("decl %d names = %v, want [%s]", i, decl.Names, name)
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
	if len(assign.Names) != 1 || assign.Names[0] != "x" {
		t.Fatalf("assign names = %v, want [x]", assign.Names)
	}
}

func TestParse_TupleDeclVsTupleAssign(t *testing.T) {
	src := "(a, b) := foo()\n(a, b) = bar()\n"
	prog := parseProgram(t, src)
	if len(prog.Statements) != 2 {
		t.Fatalf("got %d statements, want 2", len(prog.Statements))
	}

	if decl, ok := prog.Statements[0].(*ast.DeclStmt); !ok {
		t.Fatalf("stmt0 not DeclStmt: %T", prog.Statements[0])
	} else {
		if want := []string{"a", "b"}; len(decl.Names) != 2 || decl.Names[0] != want[0] || decl.Names[1] != want[1] {
			t.Fatalf("decl names = %v, want %v", decl.Names, want)
		}
	}

	if assign, ok := prog.Statements[1].(*ast.AssignStmt); !ok {
		t.Fatalf("stmt1 not AssignStmt: %T", prog.Statements[1])
	} else {
		if want := []string{"a", "b"}; len(assign.Names) != 2 || assign.Names[0] != want[0] || assign.Names[1] != want[1] {
			t.Fatalf("assign names = %v, want %v", assign.Names, want)
		}
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
	if len(decl.Names) != 1 || decl.Names[0] != "x" {
		t.Fatalf("decl names = %v, want [x]", decl.Names)
	}

	assign, ok := prog.Statements[1].(*ast.AssignStmt)
	if !ok {
		t.Fatalf("stmt1 not AssignStmt: %T", prog.Statements[1])
	}
	if len(assign.Names) != 1 || assign.Names[0] != "x" {
		t.Fatalf("assign names = %v, want [x]", assign.Names)
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
	src := "if true\n  x := 1\nelse\n  x := 2\n"
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

// ---- Typed Function Declaration Tests ----

// TestParse_FnDecl_Simple tests: def add() -> int:
func TestParse_FnDecl_Simple(t *testing.T) {
	src := "def add() -> int:\n  return 0\n"
	prog := parseProgram(t, src)
	if len(prog.Statements) != 1 {
		t.Fatalf("got %d statements, want 1", len(prog.Statements))
	}
	fn, ok := prog.Statements[0].(*ast.FnDecl)
	if !ok {
		t.Fatalf("stmt not FnDecl: %T", prog.Statements[0])
	}
	if fn.Name != "add" {
		t.Fatalf("fn name = %q, want add", fn.Name)
	}
	if len(fn.Params) != 0 {
		t.Fatalf("fn params = %d, want 0", len(fn.Params))
	}
	if fn.Return == nil || fn.Return.Name != "int" {
		t.Fatalf("fn return type = %v, want int", fn.Return)
	}
}

// TestParse_FnDecl_SingleParam tests: def greet(name: str) -> str:
func TestParse_FnDecl_SingleParam(t *testing.T) {
	src := "def greet(name: str) -> str:\n  return name\n"
	prog := parseProgram(t, src)
	if len(prog.Statements) != 1 {
		t.Fatalf("got %d statements, want 1", len(prog.Statements))
	}
	fn, ok := prog.Statements[0].(*ast.FnDecl)
	if !ok {
		t.Fatalf("stmt not FnDecl: %T", prog.Statements[0])
	}
	if fn.Name != "greet" {
		t.Fatalf("fn name = %q, want greet", fn.Name)
	}
	if len(fn.Params) != 1 {
		t.Fatalf("fn params = %d, want 1", len(fn.Params))
	}
	if fn.Params[0].Name != "name" {
		t.Fatalf("param name = %q, want name", fn.Params[0].Name)
	}
	if fn.Params[0].Type == nil || fn.Params[0].Type.Name != "str" {
		t.Fatalf("param type = %v, want str", fn.Params[0].Type)
	}
	if fn.Return == nil || fn.Return.Name != "str" {
		t.Fatalf("return type = %v, want str", fn.Return)
	}
}

// TestParse_FnDecl_MultipleParams tests: def add(a: int, b: int) -> int:
func TestParse_FnDecl_MultipleParams(t *testing.T) {
	src := "def add(a: int, b: int) -> int:\n  return a + b\n"
	prog := parseProgram(t, src)
	if len(prog.Statements) != 1 {
		t.Fatalf("got %d statements, want 1", len(prog.Statements))
	}
	fn, ok := prog.Statements[0].(*ast.FnDecl)
	if !ok {
		t.Fatalf("stmt not FnDecl: %T", prog.Statements[0])
	}
	if fn.Name != "add" {
		t.Fatalf("fn name = %q, want add", fn.Name)
	}
	if len(fn.Params) != 2 {
		t.Fatalf("fn params = %d, want 2", len(fn.Params))
	}
	if fn.Params[0].Name != "a" || fn.Params[0].Type.Name != "int" {
		t.Fatalf("param 0 wrong: name=%q type=%v", fn.Params[0].Name, fn.Params[0].Type)
	}
	if fn.Params[1].Name != "b" || fn.Params[1].Type.Name != "int" {
		t.Fatalf("param 1 wrong: name=%q type=%v", fn.Params[1].Name, fn.Params[1].Type)
	}
	if fn.Return == nil || fn.Return.Name != "int" {
		t.Fatalf("return type = %v, want int", fn.Return)
	}
}

// TestParse_FnDecl_MixedParamTypes tests: def process(name: str, count: int) -> bool:
func TestParse_FnDecl_MixedParamTypes(t *testing.T) {
	src := "def process(name: str, count: int) -> bool:\n  return true\n"
	prog := parseProgram(t, src)
	if len(prog.Statements) != 1 {
		t.Fatalf("got %d statements, want 1", len(prog.Statements))
	}
	fn, ok := prog.Statements[0].(*ast.FnDecl)
	if !ok {
		t.Fatalf("stmt not FnDecl: %T", prog.Statements[0])
	}
	if len(fn.Params) != 2 {
		t.Fatalf("fn params = %d, want 2", len(fn.Params))
	}
	if fn.Params[0].Type.Name != "str" {
		t.Fatalf("param 0 type = %q, want str", fn.Params[0].Type.Name)
	}
	if fn.Params[1].Type.Name != "int" {
		t.Fatalf("param 1 type = %q, want int", fn.Params[1].Type.Name)
	}
	if fn.Return.Name != "bool" {
		t.Fatalf("return type = %q, want bool", fn.Return.Name)
	}
}

// TestParse_FnDecl_WithBody tests function with multiple statements
func TestParse_FnDecl_WithBody(t *testing.T) {
	src := "def add(a: int, b: int) -> int:\n  x := a + b\n  return x\n"
	prog := parseProgram(t, src)
	if len(prog.Statements) != 1 {
		t.Fatalf("got %d statements, want 1", len(prog.Statements))
	}
	fn, ok := prog.Statements[0].(*ast.FnDecl)
	if !ok {
		t.Fatalf("stmt not FnDecl: %T", prog.Statements[0])
	}
	if len(fn.Body) != 2 {
		t.Fatalf("fn body = %d statements, want 2", len(fn.Body))
	}
	if _, ok := fn.Body[0].(*ast.DeclStmt); !ok {
		t.Fatalf("body[0] not DeclStmt: %T", fn.Body[0])
	}
	if _, ok := fn.Body[1].(*ast.ReturnStmt); !ok {
		t.Fatalf("body[1] not ReturnStmt: %T", fn.Body[1])
	}
}

// TestParse_FnDecl_Negative_MissingParentheses tests: def add a: int -> int:
func TestParse_FnDecl_Negative_MissingParentheses(t *testing.T) {
	src := "def add a: int -> int:\n  return 0\n"
	_, p := parseProgramWithParser(t, src)
	if len(p.Errors()) == 0 {
		t.Fatalf("expected errors for missing parentheses, got none")
	}
}

// TestParse_FnDecl_Negative_MissingParamType tests: def add(a, b) -> int:
func TestParse_FnDecl_Negative_MissingParamType(t *testing.T) {
	src := "def add(a, b) -> int:\n  return 0\n"
	_, p := parseProgramWithParser(t, src)
	if len(p.Errors()) == 0 {
		t.Fatalf("expected errors for missing parameter type, got none")
	}
}

// TestParse_FnDecl_Negative_MissingReturnType tests: def add(a: int, b: int):
func TestParse_FnDecl_Negative_MissingReturnType(t *testing.T) {
	src := "def add(a: int, b: int):\n  return 0\n"
	_, p := parseProgramWithParser(t, src)
	if len(p.Errors()) == 0 {
		t.Fatalf("expected errors for missing return type, got none")
	}
}

// TestParse_FnDecl_Negative_MissingColon tests: def add(a: int) -> int\n
func TestParse_FnDecl_Negative_MissingColon(t *testing.T) {
	src := "def add(a: int) -> int\n  return 0\n"
	_, p := parseProgramWithParser(t, src)
	if len(p.Errors()) == 0 {
		t.Fatalf("expected errors for missing colon after return type, got none")
	}
}

// TestParse_FnDecl_Multiple tests multiple function declarations
func TestParse_FnDecl_Multiple(t *testing.T) {
	src := "def add(a: int, b: int) -> int:\n  return a + b\n\ndef sub(a: int, b: int) -> int:\n  return a - b\n"
	prog := parseProgram(t, src)
	if len(prog.Statements) != 2 {
		t.Fatalf("got %d statements, want 2", len(prog.Statements))
	}
	for i, name := range []string{"add", "sub"} {
		fn, ok := prog.Statements[i].(*ast.FnDecl)
		if !ok {
			t.Fatalf("stmt %d not FnDecl: %T", i, prog.Statements[i])
		}
		if fn.Name != name {
			t.Fatalf("fn %d name = %q, want %q", i, fn.Name, name)
		}
		if len(fn.Params) != 2 {
			t.Fatalf("fn %d params = %d, want 2", i, len(fn.Params))
		}
	}
}

// ---- Tuple Unpacking Tests ----

// TestParse_DeclStmt_TupleUnpacking tests: (x, y) := ...
func TestParse_DeclStmt_TupleUnpacking(t *testing.T) {
	src := "(x, y) := run()\n"
	prog := parseProgram(t, src)
	if len(prog.Statements) != 1 {
		t.Fatalf("got %d statements, want 1", len(prog.Statements))
	}
	decl, ok := prog.Statements[0].(*ast.DeclStmt)
	if !ok {
		t.Fatalf("stmt not DeclStmt: %T", prog.Statements[0])
	}
	if len(decl.Names) != 2 {
		t.Fatalf("decl names count = %d, want 2", len(decl.Names))
	}
	if decl.Names[0] != "x" || decl.Names[1] != "y" {
		t.Fatalf("decl names = %v, want [x y]", decl.Names)
	}
}

// TestParse_DeclStmt_TupleUnpacking_Three tests: (a, b, c) := ...
func TestParse_DeclStmt_TupleUnpacking_Three(t *testing.T) {
	src := "(out, err, code) := run()\n"
	prog := parseProgram(t, src)
	if len(prog.Statements) != 1 {
		t.Fatalf("got %d statements, want 1", len(prog.Statements))
	}
	decl, ok := prog.Statements[0].(*ast.DeclStmt)
	if !ok {
		t.Fatalf("stmt not DeclStmt: %T", prog.Statements[0])
	}
	if len(decl.Names) != 3 {
		t.Fatalf("decl names count = %d, want 3", len(decl.Names))
	}
	expected := []string{"out", "err", "code"}
	for i, name := range expected {
		if decl.Names[i] != name {
			t.Fatalf("decl.Names[%d] = %q, want %q", i, decl.Names[i], name)
		}
	}
}

// TestParse_AssignStmt_TupleUnpacking tests: (x, y) = ...
func TestParse_AssignStmt_TupleUnpacking(t *testing.T) {
	src := "(x, y) = values\n"
	prog := parseProgram(t, src)
	if len(prog.Statements) != 1 {
		t.Fatalf("got %d statements, want 1", len(prog.Statements))
	}
	assign, ok := prog.Statements[0].(*ast.AssignStmt)
	if !ok {
		t.Fatalf("stmt not AssignStmt: %T", prog.Statements[0])
	}
	if len(assign.Names) != 2 {
		t.Fatalf("assign names count = %d, want 2", len(assign.Names))
	}
	if assign.Names[0] != "x" || assign.Names[1] != "y" {
		t.Fatalf("assign names = %v, want [x y]", assign.Names)
	}
}

// ---- Function Call with Named Arguments Tests ----

// TestParse_CallExpr_NamedArgs tests: run(platform=bash, cmd=cmd)
func TestParse_CallExpr_NamedArgs(t *testing.T) {
	src := "(out, err) := run(platform=bash, cmd=cmd)\n"
	prog := parseProgram(t, src)
	if len(prog.Statements) != 1 {
		t.Fatalf("got %d statements, want 1", len(prog.Statements))
	}
	decl, ok := prog.Statements[0].(*ast.DeclStmt)
	if !ok {
		t.Fatalf("stmt not DeclStmt: %T", prog.Statements[0])
	}
	callExpr, ok := decl.Value.(*ast.CallExpr)
	if !ok {
		t.Fatalf("value not CallExpr: %T", decl.Value)
	}
	// Check named args
	if len(callExpr.NamedArgs) != 2 {
		t.Fatalf("call named args count = %d, want 2", len(callExpr.NamedArgs))
	}
	if callExpr.NamedArgs[0].Name != "platform" {
		t.Fatalf("first named arg name = %q, want platform", callExpr.NamedArgs[0].Name)
	}
	if callExpr.NamedArgs[1].Name != "cmd" {
		t.Fatalf("second named arg name = %q, want cmd", callExpr.NamedArgs[1].Name)
	}
}
