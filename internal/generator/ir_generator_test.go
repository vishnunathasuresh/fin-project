package generator

import (
	"testing"

	"github.com/vishnunathasuresh/fin-project/internal/ir"
)

func TestIRBatchGenerator_SimpleFunction(t *testing.T) {
	prog := &ir.Program{
		Functions: map[string]*ir.Function{
			"main": {
				Name:   "main",
				Params: []ir.Param{},
				Body: []ir.Stmt{
					&ir.DeclStmt{
						Name: "a",
						Type: &ir.BasicType{Kind: "int"},
						Init: &ir.IntLit{Value: 42},
					},
				},
			},
		},
	}

	gen := NewIRBatchGenerator()
	output, err := gen.Generate(prog)
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	t.Logf("Generated output:\n%s", output)

	if output == "" {
		t.Fatal("expected non-empty output")
	}

	if !contains(output, "@echo off") {
		t.Error("expected @echo off in output")
	}

	if !contains(output, "setlocal EnableDelayedExpansion") {
		t.Error("expected setlocal EnableDelayedExpansion in output")
	}

	if !contains(output, "set /a a=42") {
		t.Error("expected set /a a=42 in output")
	}
}

func TestIRBatchGenerator_IfStatement(t *testing.T) {
	prog := &ir.Program{
		Functions: map[string]*ir.Function{
			"main": {
				Name:   "main",
				Params: []ir.Param{},
				Body: []ir.Stmt{
					&ir.DeclStmt{
						Name: "x",
						Type: &ir.BasicType{Kind: "int"},
						Init: &ir.IntLit{Value: 5},
					},
					&ir.IfStmt{
						Cond: &ir.BinaryOp{
							Left:  &ir.Ident{Name: "x", Type: &ir.BasicType{Kind: "int"}},
							Op:    ">",
							Right: &ir.IntLit{Value: 0},
							Type:  &ir.BasicType{Kind: "bool"},
						},
						Then: []ir.Stmt{
							&ir.AssignStmt{
								Name:  "x",
								Value: &ir.IntLit{Value: 10},
							},
						},
						Else: []ir.Stmt{},
					},
				},
			},
		},
	}

	gen := NewIRBatchGenerator()
	output, err := gen.Generate(prog)
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	if !contains(output, "if !x! > 0") {
		t.Error("expected if condition in output")
	}

	if !contains(output, "set /a x=10") {
		t.Error("expected assignment in if body")
	}
}

func TestIRBatchGenerator_ForLoop(t *testing.T) {
	prog := &ir.Program{
		Functions: map[string]*ir.Function{
			"main": {
				Name:   "main",
				Params: []ir.Param{},
				Body: []ir.Stmt{
					&ir.ForStmt{
						Var:   "i",
						Start: &ir.IntLit{Value: 1},
						End:   &ir.IntLit{Value: 3},
						Body: []ir.Stmt{
							&ir.DeclStmt{
								Name: "temp",
								Type: &ir.BasicType{Kind: "int"},
								Init: &ir.Ident{Name: "i", Type: &ir.BasicType{Kind: "int"}},
							},
						},
					},
				},
			},
		},
	}

	gen := NewIRBatchGenerator()
	output, err := gen.Generate(prog)
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	if !contains(output, "for /L %i in (1,1,3)") {
		t.Error("expected for loop in output")
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
