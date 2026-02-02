package generator

import (
	"testing"

	"github.com/vishnunathasuresh/fin-project/internal/ast"
)

func TestGenerate_TopLevelDeclEchoRun_PLACEHOLDER(t *testing.T) {
	g := NewBatchGenerator()
	prog := &ast.Program{Statements: []ast.Statement{
		// TODO(fin-v2): Generator not yet ported. Using DeclStmt to mirror v2; emit currently unsupported.
		&ast.DeclStmt{Names: []string{"x"}, Value: &ast.NumberLit{Value: "10", P: ast.Pos{Line: 1, Column: 8}}, P: ast.Pos{Line: 1, Column: 1}},
	}}
	_, err := g.Generate(prog)
	if err == nil {
		t.Fatalf("expected error until generator supports DeclStmt/run lowering")
	}
}

func TestGenerate_Assign_PLACEHOLDER(t *testing.T) {
	prog := &ast.Program{Statements: []ast.Statement{
		// TODO(fin-v2): Generator should handle DeclStmt then AssignStmt (Names slice). Currently unsupported.
		&ast.DeclStmt{Names: []string{"a"}, Value: &ast.NumberLit{Value: "1", P: ast.Pos{Line: 1, Column: 7}}, P: ast.Pos{Line: 1, Column: 1}},
		&ast.AssignStmt{Names: []string{"a"}, Value: &ast.NumberLit{Value: "2", P: ast.Pos{Line: 2, Column: 5}}, P: ast.Pos{Line: 2, Column: 1}},
	}}
	g := NewBatchGenerator()
	_, err := g.Generate(prog)
	if err == nil {
		t.Fatalf("expected error until generator supports DeclStmt/AssignStmt lowering")
	}
}

func TestGenerate_UnsupportedStmt(t *testing.T) {
	g := NewBatchGenerator()
	prog := &ast.Program{Statements: []ast.Statement{
		&ast.ReturnStmt{Value: &ast.NumberLit{Value: "1"}},
	}}

	_, err := g.Generate(prog)
	if err == nil {
		t.Fatalf("expected error for unsupported stmt")
	}

	if _, ok := err.(*GeneratorError); !ok {
		t.Fatalf("expected GeneratorError, got %T", err)
	}
}

func TestGenerate_DeclStmtNotYetSupported(t *testing.T) {
	g := NewBatchGenerator()
	prog := &ast.Program{Statements: []ast.Statement{
		&ast.DeclStmt{Names: []string{"a"}, Value: &ast.NumberLit{Value: "1", P: ast.Pos{Line: 1, Column: 6}}, P: ast.Pos{Line: 1, Column: 1}},
	}}

	_, err := g.Generate(prog)
	if err == nil {
		t.Fatalf("expected error for DeclStmt until generator is ported to v2")
	}
}

func TestGenerate_FunctionNotLifted(t *testing.T) {
	g := NewBatchGenerator()
	fn := &ast.FnDecl{Name: "x", Params: []ast.Param{}, Body: nil}

	if err := g.emitStmt(fn); err == nil {
		t.Fatalf("expected error for unlifted function")
	} else if _, ok := err.(*GeneratorError); !ok {
		t.Fatalf("expected GeneratorError, got %T", err)
	}
}

func TestGenerate_Call(t *testing.T) {
	t.Skip("fin-v2: generator not yet ported for FnDecl/Call lowering")
	g := NewBatchGenerator()
	prog := &ast.Program{Statements: []ast.Statement{
		&ast.FnDecl{
			Name:   "greet",
			Params: []ast.Param{{Name: "name", P: ast.Pos{Line: 1, Column: 12}}},
			Body:   []ast.Statement{}, // TODO(fin-v2): add run() lowering when available
		},
		&ast.CallStmt{Name: "greet", Args: []ast.Expr{&ast.StringLit{Value: "foo bar&baz"}}},
	}}
	if _, err := g.Generate(prog); err == nil {
		t.Fatalf("expected error until generator supports function/call lowering")
	}
}

func TestGenerate_Function(t *testing.T) {
	t.Skip("fin-v2: generator not yet ported for function body lowering")
	g := NewBatchGenerator()
	prog := &ast.Program{Statements: []ast.Statement{
		&ast.FnDecl{
			Name:   "greet",
			Params: []ast.Param{{Name: "name", P: ast.Pos{Line: 1, Column: 12}}},
			Body:   []ast.Statement{}, // TODO(fin-v2): add run()/assign lowering when generator is ported
		},
	}}
	if _, err := g.Generate(prog); err == nil {
		t.Fatalf("expected error until generator supports function lowering")
	}
}

func TestGenerate_IfElse(t *testing.T) {
	t.Skip("fin-v2: generator not yet ported for IfStmt lowering without EchoStmt")
	g := NewBatchGenerator()
	prog := &ast.Program{Statements: []ast.Statement{
		&ast.IfStmt{
			Cond: &ast.BoolLit{Value: true},
			Then: []ast.Statement{
				&ast.DeclStmt{Names: []string{"y"}, Value: &ast.NumberLit{Value: "1"}}, // placeholder
			},
			Else: []ast.Statement{
				&ast.DeclStmt{Names: []string{"n"}, Value: &ast.NumberLit{Value: "0"}}, // placeholder
			},
		},
	}}
	if _, err := g.Generate(prog); err == nil {
		t.Fatalf("expected error until generator supports IfStmt lowering in v2")
	}
}
