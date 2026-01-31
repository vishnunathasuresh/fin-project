package ir

import (
	"testing"

	"github.com/vishnunath-suresh/fin-project/internal/ast"
)

func TestLowerSimpleFunction(t *testing.T) {
	// Create a simple AST function
	astProg := &ast.Program{
		Statements: []ast.Statement{
			&ast.FnDecl{
				Name:   "test",
				Params: []string{"a", "b"},
				Body: []ast.Statement{
					&ast.ReturnStmt{
						Value: &ast.BinaryExpr{
							Left:  &ast.IdentExpr{Name: "a"},
							Op:    "+",
							Right: &ast.IdentExpr{Name: "b"},
						},
					},
				},
				P: ast.Pos{Line: 1, Column: 1},
			},
		},
	}

	// Lower to IR
	irProg, err := Lower(astProg)
	if err != nil {
		t.Fatalf("Lower failed: %v", err)
	}

	// Verify function was created
	if len(irProg.Functions) != 1 {
		t.Fatalf("expected 1 function, got %d", len(irProg.Functions))
	}

	fn, ok := irProg.Functions["test"]
	if !ok {
		t.Fatal("function 'test' not found")
	}

	if fn.Name != "test" {
		t.Errorf("expected function name 'test', got '%s'", fn.Name)
	}

	if len(fn.Params) != 2 {
		t.Errorf("expected 2 parameters, got %d", len(fn.Params))
	}

	if len(fn.Body) != 1 {
		t.Errorf("expected 1 statement in body, got %d", len(fn.Body))
	}
}

func TestLowerIfStatement(t *testing.T) {
	astProg := &ast.Program{
		Statements: []ast.Statement{
			&ast.FnDecl{
				Name:   "test",
				Params: []string{"x"},
				Body: []ast.Statement{
					&ast.IfStmt{
						Cond: &ast.BinaryExpr{
							Left:  &ast.IdentExpr{Name: "x"},
							Op:    ">",
							Right: &ast.NumberLit{Value: "0"},
						},
						Then: []ast.Statement{
							&ast.ReturnStmt{Value: &ast.BoolLit{Value: true}},
						},
						Else: []ast.Statement{
							&ast.ReturnStmt{Value: &ast.BoolLit{Value: false}},
						},
					},
				},
			},
		},
	}

	irProg, err := Lower(astProg)
	if err != nil {
		t.Fatalf("Lower failed: %v", err)
	}

	fn := irProg.Functions["test"]
	if len(fn.Body) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(fn.Body))
	}

	ifStmt, ok := fn.Body[0].(*IfStmt)
	if !ok {
		t.Fatalf("expected IfStmt, got %T", fn.Body[0])
	}

	if len(ifStmt.Then) != 1 {
		t.Errorf("expected 1 then statement, got %d", len(ifStmt.Then))
	}

	if len(ifStmt.Else) != 1 {
		t.Errorf("expected 1 else statement, got %d", len(ifStmt.Else))
	}
}

func TestValidateIR(t *testing.T) {
	prog := &Program{
		Types:     make(map[string]*TypeDef),
		Functions: make(map[string]*Function),
	}

	prog.Functions["test"] = &Function{
		Name: "test",
		Params: []Param{
			{Name: "a", Type: &BasicType{Kind: "int"}},
			{Name: "b", Type: &BasicType{Kind: "int"}},
		},
		ReturnType: &BasicType{Kind: "int"},
		Body: []Stmt{
			&ReturnStmt{
				Value: &IntLit{Value: 42},
			},
		},
	}

	err := Validate(prog)
	if err != nil {
		t.Errorf("Validate failed: %v", err)
	}
}

func TestValidateDuplicateParam(t *testing.T) {
	prog := &Program{
		Types:     make(map[string]*TypeDef),
		Functions: make(map[string]*Function),
	}

	prog.Functions["test"] = &Function{
		Name: "test",
		Params: []Param{
			{Name: "a", Type: &BasicType{Kind: "int"}},
			{Name: "a", Type: &BasicType{Kind: "int"}}, // Duplicate
		},
		ReturnType: &BasicType{Kind: "int"},
		Body:       []Stmt{},
	}

	err := Validate(prog)
	if err == nil {
		t.Error("expected validation error for duplicate parameter")
	}
}
